package middleware

import (
	"context"
	"net/http"

	"github.com/Amierza/pos-broissant/dto"
	"github.com/Amierza/pos-broissant/service"
	"github.com/Amierza/pos-broissant/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString, _, err := ParseAuthorizationHeader(authHeader, jwtService)
		if err != nil {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userID, err := jwtService.GetUserIDByToken(tokenString)
		if err != nil {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		newCtx := context.WithValue(ctx.Request.Context(), UserIDKey, userID)
		ctx.Request = ctx.Request.WithContext(newCtx)
		ctx.Next()
	}
}
