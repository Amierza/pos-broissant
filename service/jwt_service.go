package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	JWTService interface {
		GenerateToken(userID, role string) (string, string, error)
		ValidateToken(tokenString string) (*jwt.Token, error)
		GetUserIDByToken(accessToken string) (string, error)
		GetRoleByToken(accessToken string) (string, error)
	}

	jwtCustomClaim struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
		jwt.RegisteredClaims
	}

	jwtService struct {
		secretKey string
		ctx       context.Context
		issuer    string
	}
)

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "Template"
	}
	return secretKey
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "Template",
	}
}

func (j *jwtService) GenerateToken(userID, role string) (string, string, error) {
	accessClaims := jwtCustomClaim{
		userID,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 120)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshClaims := jwtCustomClaim{
		userID,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 3600 * 24 * 7)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return accessTokenString, refreshTokenString, nil
}

func (j *jwtService) ParseToken(t_ *jwt.Token) (any, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t_.Header["alg"])
	}

	return []byte(j.secretKey), nil
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, j.ParseToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j *jwtService) GetUserIDByToken(accessToken string) (string, error) {
	token, err := j.ValidateToken(accessToken)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	userID := fmt.Sprintf("%v", claims["user_id"])
	return userID, nil
}

func (j *jwtService) GetRoleByToken(accessToken string) (string, error) {
	token, err := j.ValidateToken(accessToken)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	role := fmt.Sprintf("%v", claims["role"])
	return role, nil
}
