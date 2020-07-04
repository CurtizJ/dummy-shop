package main

// @title Shop API
// @version 1.0
// @description Educational online-shop API

// @host localhost:8181
// @BasePath /
import (
	"fmt"
	"net/http"
	"os"

	"github.com/CurtizJ/dummy-shop/api/repos"
	"github.com/CurtizJ/dummy-shop/lib/pb"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	_ "github.com/CurtizJ/dummy-shop/api/docs"

	swagger "github.com/swaggo/http-swagger"
)

var repo repos.Repo
var grpcClient pb.VerificationClient

func main() {
	repo = repos.NewPostgresRepo()
	router := mux.NewRouter()
	registerHandlers(router)

	port, exists := os.LookupEnv("LISTEN_PORT")
	if !exists {
		port = ":8081"
	}

	conn, err := grpc.Dial(os.Getenv("GRPC_ADDRESS"), grpc.WithInsecure())
	if err != nil {
		panic("Cannot create grpc connection. Error: " + err.Error())
	}
	defer conn.Close()
	grpcClient = pb.NewVerificationClient(conn)

	http.ListenAndServe(port, nil)
}

func registerHandlers(router *mux.Router) {
	routerSecured := router.NewRoute().Subrouter()
	routerSecured.Use(getMiddleware(RoleAdmin))

	routerSecured.HandleFunc("/item", handlerItem).Methods("POST", "PUT")
	routerSecured.HandleFunc("/item/{id:[0-9]+}", handlerItemId).Methods("DELETE")

	// router.Use(getMiddleware(RoleUser))
	router.HandleFunc("/item/{id:[0-9]+}", handlerItemId).Methods("GET")
	router.HandleFunc("/items", handlerList).Methods("GET")

	router.PathPrefix("/swagger").Handler(swagger.WrapHandler)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<p>Shop API</p><a href=/swagger/>Swagger</a>")
	})

	http.Handle("/", router)
}
