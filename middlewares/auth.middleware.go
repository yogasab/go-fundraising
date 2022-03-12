package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-fundraising/helper"
	"go-fundraising/service"
	"net/http"
	"strings"
)

func AuthorizeToken(jwtService service.JWTService, userService service.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//	Get auth header
		authHeader := ctx.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//	Split header
		jwtToken := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			jwtToken = arrayToken[1]
		}
		//	Validate token
		validatedJWTToken, err := jwtService.ValidateToken(jwtToken)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//	Get the payload from validated jwt token
		claim, ok := validatedJWTToken.Claims.(jwt.MapClaims)
		//	Check if it is valid
		if !ok || !validatedJWTToken.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//	Get user by user id
		userID := claim["user_id"].(float64)
		user, err := userService.GetUserByID(int(userID))
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//	Set user to request
		ctx.Set("user", user)
	}
}
