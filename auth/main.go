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
	"github.com/streadway/amqp"

	. "github.com/CurtizJ/dummy-shop/lib/config"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"

	swagger "github.com/swaggo/http-swagger"
)

const (
	USER_DATABASE          = 0
	SESSIONS_DATABASE      = 1
	CONFIRMATIONS_DATABASE = 2
)

var users *redis.Client
var sessions *redis.Client
var confirmations *redis.Client

var notifications amqp.Queue
var notificationsChannel *amqp.Channel

var config *Config

func main() {
	config = &Config{}
	err := NewConfigFromEnv(config)
	if err != nil {
		panic("Cannot create config: " + err.Error())
	}

	registerRedis()
	registerRabbitMQ()

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
	router.HandleFunc("/confirm", handlerConfirm)
}

func registerSwagger(router *mux.Router) {
	router.PathPrefix("/swagger").Handler(swagger.WrapHandler)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<p>Auth service</p><a href=/swagger/>Swagger</a>")
	})
}

func registerRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	users = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       USER_DATABASE,
	})

	sessions = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       SESSIONS_DATABASE,
	})

	confirmations = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       CONFIRMATIONS_DATABASE,
	})
}

func registerRabbitMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic("Cannot connect to RabbitMQ: " + err.Error())
	}

	notificationsChannel, err = conn.Channel()
	if err != nil {
		panic("Cannot connect to RabbitMQ: " + err.Error())
	}

	notifications, err = notificationsChannel.QueueDeclare(
		"notifications", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	if err != nil {
		panic("Cannot create queue in RabbitMQ: " + err.Error())
	}
}
