package middleware

import (
	"context"
	"net/http"

	"github.com/Amierza/pos-broissant/dto"
	"github.com/Amierza/pos-broissant/service"
	"github.com/Amierza/pos-broissant/utils"
	"github.com/gin-gonic/gin"
)

func AuthorizeRoleIsAdmin(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString, _, err := ParseAuthorizationHeader(authHeader, jwtService)
		if err != nil {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		role, err := jwtService.GetRoleByToken(tokenString)
		if err != nil {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, res)
			return
		}

		if role != "admin" {
			res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_ROLE_NOT_ADMIN, nil)
			ctx.AbortWithStatusJSON(http.StatusForbidden, res)
			return
		}

		newCtx := context.WithValue(ctx.Request.Context(), RoleKey, role)
		ctx.Request = ctx.Request.WithContext(newCtx)
		ctx.Next()
	}
}
