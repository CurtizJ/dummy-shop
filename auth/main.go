package main

// @title Small authorization service
// @version 1.0
// @description Small authorization service

// @host localhost:8182
// @BasePath /auth
import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/CurtizJ/dummy-shop/auth/docs"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"

	swagger "github.com/swaggo/http-swagger"
)

const (
	USER_DATABASE     = 0
	SESSIONS_DATABASE = 1
)

var users *redis.Client
var sessions *redis.Client
var config *ConfigEnv

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

	config = &ConfigEnv{}
	err := NewConfigFromEnv(config)
	if err != nil {
		panic("Cannot create config")
	}

	router := mux.NewRouter()
	registerSwagger(router)
	routerAuth := router.PathPrefix("/auth").Subrouter()
	registerHandlers(routerAuth)
	http.Handle("/", router)

	port, exists := os.LookupEnv("LISTEN_PORT")
	if !exists {
		port = ":8082"
	}

	http.ListenAndServe(port, nil)
}

func registerHandlers(router *mux.Router) {
	router.HandleFunc("/signup", handlerSignUp).Methods("POST")
	router.HandleFunc("/signin", handlerSignIn).Methods("POST")
	router.HandleFunc("/verify", handlerVerify).Methods("POST")
	router.HandleFunc("/refresh", handlerRefresh).Methods("POST")
}

func registerSwagger(router *mux.Router) {
	router.PathPrefix("/swagger").Handler(swagger.WrapHandler)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<p>Auth service</p><a href=/swagger/>Swagger</a>")
	})
}
