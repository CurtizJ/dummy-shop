package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CurtizJ/dummy-shop/lib/errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

func getStringHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

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

	err = users.HMSet(user.Email, map[string]interface{}{
		"password": getStringHash(user.Password),
		"verified": 0,
		"role":     RoleUser}).Err()

	if err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	if err := sendConfirmationLink(user.Email); err != nil {
		errors.ReportAsJSON(w, "Cannot send confirmation link", http.StatusInternalServerError)
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

	exists, err := users.Exists(user.Email).Result()
	if err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	if exists != 1 {
		errors.ReportAsJSON(w, "User not found", http.StatusUnauthorized)
		return
	}

	verified, err := users.HGet(user.Email, "verified").Int()
	if err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	if verified != 1 {
		errors.ReportAsJSON(w, "Email is not verified", http.StatusUnauthorized)
		return
	}

	password, err := users.HGet(user.Email, "password").Result()
	if err != nil {
		errors.ReportAsJSON(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	if getStringHash(user.Password) != password {
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

// @Summary Admin
// @Description Grant admin rights to user
// @Param Authorization header string true "Access token"
// @Accept  json
// @Param user body User true "User email to grant admin rights"
// @Success 200
// @Router /admin [post]
func handlerAdmin(w http.ResponseWriter, r *http.Request) {
	accessToken, ok := r.Header["Authorization"]
	if !ok {
		errors.ReportAsJSON(w, "Authorization required", http.StatusUnauthorized)
		return
	}

	response, err := VerifyAccessToken(accessToken[0])
	if err != nil || !response.Valid || response.Role != RoleAdmin {
		errors.ReportAsJSON(w, "Not enough rights", http.StatusUnauthorized)
		return
	}

	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		errors.ReportAsJSON(w, "Cannot decode user", http.StatusBadRequest)
	}

	email, exists := body["email"]
	if !exists {
		errors.ReportAsJSON(w, "Cannot decode user", http.StatusBadRequest)
	}

	if err := users.HSet(email, "role", RoleAdmin).Err(); err != nil {
		errors.ReportAsJSON(w, err.Error(), http.StatusInternalServerError)
	}
}

// Not a part of public api.
func handlerConfirm(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]
	email, err := confirmations.Get(code).Result()
	if err != nil {
		errors.ReportAsJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists, err := users.Exists(email).Result(); err != nil || exists != 1 {
		errors.ReportAsJSON(w, "Bad confirmation link", http.StatusBadRequest)
		return
	}

	if err := users.HSet(email, "verified", 1).Err(); err != nil {
		errors.ReportAsJSON(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"email": email, "verified": 1})
}
