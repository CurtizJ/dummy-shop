package main

import (
	"context"
	"fmt"

	"github.com/CurtizJ/dummy-shop/lib/pb"
	"github.com/dgrijalva/jwt-go"
)

type VerificationServer struct {
}

func VerifyAccessToken(accessToken string) (*pb.VerificationResponse, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.AuthSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		role, _ := users.HGet(email, "role").Int()
		return &pb.VerificationResponse{
			UserEmail: email,
			Valid:     true,
			Role:      int32(role),
		}, nil
	}

	return &pb.VerificationResponse{UserEmail: "", Valid: false, Role: RoleUnknown}, nil
}

func (server *VerificationServer) Verify(ctx context.Context, request *pb.VerificationRequest) (*pb.VerificationResponse, error) {
	return VerifyAccessToken(request.AccessToken)
}
