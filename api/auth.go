package main

import (
	"net/http"

	"github.com/CurtizJ/dummy-shop/lib/errors"
)

const AUTH_ENDPOINT = "http://auth:8182/auth/verify"

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Header["Authorization"]
		if !ok {
			errors.ReportAsJSON(w, "Authorization requiered", http.StatusUnauthorized)
			return
		}

		request, _ := http.NewRequest("POST", AUTH_ENDPOINT, nil)
		request.Header.Add("Authorization", token[0])
		response, err := http.DefaultClient.Do(request)

		if err != nil {
			errors.ReportAsJSON(w, "Authorization service unavailable", http.StatusInternalServerError)
			return
		}

		if response.StatusCode == http.StatusUnauthorized {
			errors.ReportAsJSON(w, "Authorization failed", http.StatusUnauthorized)
			return
		} else if response.StatusCode != http.StatusOK {
			errors.ReportAsJSON(w, "Authorization service unavailable", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
