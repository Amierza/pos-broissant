package middleware

import (
	"fmt"
	"strings"

	"github.com/Amierza/pos-broissant/service"
	"github.com/golang-jwt/jwt/v5"
)

const (
	BearerPrefix = "Bearer "
	UserIDKey    = "user_id"
	RoleKey      = "role"
)

func ParseAuthorizationHeader(authHeader string, jwtService service.JWTService) (string, *jwt.Token, error) {
	if authHeader == "" {
		return "", nil, fmt.Errorf("Authorization header is missing")
	}

	if !strings.Contains(authHeader, "Bearer") {
		return "", nil, fmt.Errorf("Invalid token format")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	token, err := jwtService.ValidateToken(tokenString)
	if err != nil {
		return "", nil, err
	}

	if !token.Valid {
		return "", nil, fmt.Errorf("Token is invalid")
	}

	return tokenString, token, nil
}
