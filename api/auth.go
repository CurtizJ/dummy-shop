package main

import (
	"context"
	"net/http"

	"github.com/CurtizJ/dummy-shop/lib/errors"
	"github.com/CurtizJ/dummy-shop/lib/pb"
	"github.com/gorilla/mux"
)

const (
	RoleUnknown int32 = iota
	RoleUser
	RoleAdmin
)

func getMiddleware(allowedRole int32) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := r.Header["Authorization"]
			if !ok {
				errors.ReportAsJSON(w, "Authorization requiered", http.StatusUnauthorized)
				return
			}

			response, err := grpcClient.Verify(context.Background(), &pb.VerificationRequest{
				AccessToken: token[0],
			})

			if err != nil || !response.Valid || response.Role < allowedRole {
				errors.ReportAsJSON(w, "Authorization failed", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
