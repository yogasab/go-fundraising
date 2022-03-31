package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go-fundraising/dto"
	"go-fundraising/service"
	"net/http"
)

type AuthenticationHandler interface {
	LoginIndex(ctx *gin.Context)
	LoginStore(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type authenticationHandler struct {
	userService service.UserService
}

func NewAuthenticationHandler(userService service.UserService) AuthenticationHandler {
	return &authenticationHandler{userService: userService}
}

func (h *authenticationHandler) LoginIndex(ctx *gin.Context) {
	ctx.HTML(http.StatusFound, "login.html", nil)
}

func (h *authenticationHandler) LoginStore(ctx *gin.Context) {
	var request dto.LoginRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	user, err := h.userService.LoginUser(request)
	if err != nil || user.Role != "admin" {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	session := sessions.Default(ctx)
	session.Set("userID", user.ID)
	session.Set("username", user.Name)
	session.Save()

	ctx.Redirect(http.StatusFound, "/users")
}

func (h *authenticationHandler) Logout(ctx *gin.Context) {
	sessions := sessions.Default(ctx)
	sessions.Clear()
	sessions.Save()
	ctx.Redirect(http.StatusFound, "/auth/login")
}
