package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type JWTService interface {
	GenerateToken(UserID int64) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	return &jwtService{
		secretKey: secretKey,
		issuer:    "user",
	}
}

func (s *jwtService) GenerateToken(UserID int64) (string, error) {
	claims := jwtCustomClaim{
		UserID: UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	plainToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := plainToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	// Parse token
	parsedToken, err := jwt.Parse(encodedToken, func(jwtToken *jwt.Token) (interface{}, error) {
		// Check sign in method algorithm
		_, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid Token")
		}
		// Return Secret Key
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return parsedToken, err
	}
	return parsedToken, nil
}
