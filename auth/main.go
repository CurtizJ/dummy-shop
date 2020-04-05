package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

const (
	USER_DATABASE     = 0
	SESSIONS_DATABASE = 1
)

var users *redis.Client
var sessions *redis.Client

func main() {
	users = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       USER_DATABASE,
	})

	sessions = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       SESSIONS_DATABASE,
	})

	router := mux.NewRouter().PathPrefix("/auth").Subrouter()
	registerHandlers(router)

	port, exists := os.LookupEnv("LISTEN_PORT")
	if !exists {
		port = ":8081"
	}

	fmt.Println("listening on port: ", port)

	http.ListenAndServe(port, nil)

	// port, exists := os.LookupEnv("LISTEN_PORT")
	// if !exists {
	// 	port = ":8081"
	// }
	// http.ListenAndServe(port, nil)

	// token, err := CreateToken(uint64(228))
	// fmt.Println(token, err)

	// token1, err := jwt.Parse(token.Access, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	// 	}

	// 	return []byte("kek"), nil
	// })

	// fmt.Println(token1.Valid)
	// fmt.Println(token1.Claims)
}

func registerHandlers(router *mux.Router) {
	router.HandleFunc("/signup", handlerSignUp).Methods("POST")
	router.HandleFunc("/signin", handlerSignIn).Methods("POST")
	router.HandleFunc("/verify", handlerVerify).Methods("POST")
	router.HandleFunc("/refresh", handlerRefresh).Methods("POST")

	http.Handle("/", router)
}
