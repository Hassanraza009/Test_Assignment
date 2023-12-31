package rest

import (
	"net/http"
	"strings"

	"test/models"
	"test/service"

	"github.com/gin-gonic/gin"
)

type MiddleWare interface {
	ValidateToken() gin.HandlerFunc
}

type middleWare struct {
	jwtService service.JWTService
}

// NewAuthController return new instance of authcontroller
func NewMiddleWare(jwtService service.JWTService) MiddleWare {
	return &middleWare{
		jwtService: jwtService,
	}
}

const keyCurrentEmail = "__current_email"

func (m *middleWare) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Get the token from the Authorization header with the "Bearer" prefix
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, NewStandardResponse(false, 1001, "Token not found", nil))
			return
		}

		// Extract the token from the "Bearer " prefix
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, NewStandardResponse(false, 1001, "Invalid Bearer token format", nil))
			return
		}

		token := splitToken[1]

		decodedClaims, err := m.jwtService.VerifyLoginToken(token)
		if ctx == nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, NewStandardResponse(false, models.INVALID_TOKEN, models.INVALID_TOKEN_MESSAGE, nil))
			return
		}
		if err != nil {
			if standardError, ok := err.(*service.StandardError); ok {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, NewStandardResponse(false, standardError.Code, standardError.Message, nil))
				return
			} else {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, NewStandardResponse(false, models.INTERNAL_SERVER_ERROR, "Error conversion failed", nil))
				return
			}
		}

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewStandardResponse(false, models.INTERNAL_SERVER_ERROR, "Error fetching user data", nil))
			return
		}

		ctx.Set(keyCurrentEmail, decodedClaims.Email)
		ctx.Next()
	}
}
