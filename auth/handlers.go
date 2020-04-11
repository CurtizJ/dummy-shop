package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CurtizJ/dummy-shop/lib/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

// @Summary Sign up
// @Accept  json
// @Param user body User true "User with email and password"
// @Success 200
// @Router /signup [post]
func handlerSignUp(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errors.ReportAsJSON(w, "Cannot decode user", http.StatusBadRequest)
		return
	}

	exists, err := users.Exists(user.Email).Result()
	if err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	if exists == 1 {
		errors.ReportAsJSON(w, "User with email: "+user.Email+" already exists", http.StatusForbidden)
		return
	}

	if err := users.Set(user.Email, user.Password, 0).Err(); err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
	}
}

// @Summary Sign in
// @Accept  json
// @Param user body User true "User with email and password"
// @Success 200 {object} Token
// @Router /signin [post]
func handlerSignIn(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errors.ReportAsJSON(w, "Cannot decode user", http.StatusBadRequest)
		return
	}

	password, err := users.Get(user.Email).Result()
	if err == redis.Nil {
		errors.ReportAsJSON(w, "User not found", http.StatusNotAcceptable)
		return
	} else if err != nil {
		errors.ReportAsJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Password != password {
		errors.ReportAsJSON(w, "Incorrect password", http.StatusNotAcceptable)
		return
	}

	token, err := NewToken(user.Email)
	if err != nil {
		errors.ReportAsJSON(w, "Cannot create token", http.StatusInternalServerError)
		return
	}

	if err := sessions.Set(token.Refresh, user.Email, time.Duration(config.RefreshExpiration)*time.Second).Err(); err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(token)
}

// @Summary Verify
// @Description Verifies that access token is valid
// @Param Authorization header string true "Access jwt token to verify"
// @Success 200
// @Router /verify [post]
func handlerVerify(w http.ResponseWriter, r *http.Request) {
	jwtToken, ok := r.Header["Authorization"]
	if !ok {
		errors.ReportAsJSON(w, "Authorization requiered", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(jwtToken[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.AuthSecret), nil
	})

	if err != nil {
		errors.ReportAsJSON(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		json.NewEncoder(w).Encode(map[string]string{"email": claims["email"].(string)})
	} else {
		errors.ReportAsJSON(w, "Bad token", http.StatusUnauthorized)
	}
}

// @Summary Refresh
// @Description Returns new pair of access and refresh tokens
// @Param Authorization header string true "Refresh token"
// @Success 200
// @Router /refresh [post]
func handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, ok := r.Header["Authorization"]
	if !ok {
		errors.ReportAsJSON(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	email, err := sessions.Get(refreshToken[0]).Result()
	if err == redis.Nil {
		errors.ReportAsJSON(w, "Wrong or expired token", http.StatusUnauthorized)
		return
	} else if err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	if err := sessions.Del(refreshToken[0]).Err(); err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	token, err := NewToken(email)
	if err != nil {
		errors.ReportAsJSON(w, "Cannot create token", http.StatusInternalServerError)
		return
	}

	if err := sessions.Set(token.Refresh, email, time.Second*60).Err(); err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(token)
}
