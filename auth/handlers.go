package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
)

func handlerSignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlerSignUp")
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// errors.ReportAsJSON(w, "Cannot decode user.", http.StatusBadRequest)
		http.Error(w, "Cannot decode user", http.StatusBadRequest)
		return
	}

	exists, err := users.Exists(user.Email).Result()
	if err != nil {
		http.Error(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	if exists == 1 {
		http.Error(w, "User with email: "+user.Email+" already exists", http.StatusForbidden)
		return
	}

	if err := users.Set(user.Email, user.Password, 0).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerSignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlerSignIn")
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// errors.ReportAsJSON(w, "Cannot decode user.", http.StatusBadRequest)
		http.Error(w, "Cannot decode user", http.StatusBadRequest)
		return
	}

	password, err := users.Get(user.Email).Result()
	if err == redis.Nil {
		http.Error(w, "User not found", http.StatusNotAcceptable)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(user.Email, password)

	if user.Password != password {
		http.Error(w, "Incorrect password", http.StatusNotAcceptable)
		return
	}

	token, err := NewToken(user.Email)
	if err != nil {
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}

	if err := sessions.Set(token.Refresh, user.Email, time.Second*60).Err(); err != nil {
		http.Error(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(token)
}

func handlerVerify(w http.ResponseWriter, r *http.Request) {
	jwtToken := r.Header["Authorization"][0]

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("kek"), nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Fprintf(w, "email: %s", claims["email"].(string))
	} else {
		http.Error(w, "Bad token", http.StatusUnauthorized)
	}
}

func handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header["Authorization"][0]
	email, err := sessions.Get(refreshToken).Result()
	if err == redis.Nil {
		http.Error(w, "Bad token", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	token, err := NewToken(email)
	if err != nil {
		http.Error(w, "Cannot create token", http.StatusInternalServerError)
		return
	}

	if err := sessions.Set(token.Refresh, email, time.Second*60).Err(); err != nil {
		http.Error(w, "Redis unavailable", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(token)
}
