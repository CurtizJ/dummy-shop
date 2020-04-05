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
	"github.com/gorilla/mux"

	_ "github.com/CurtizJ/dummy-shop/api/docs"

	swagger "github.com/swaggo/http-swagger"
)

var repo repos.Repo

func main() {
	repo = repos.NewPostgresRepo()
	router := mux.NewRouter()
	registerHandlers(router)

	port, exists := os.LookupEnv("LISTEN_PORT")
	if !exists {
		port = ":8080"
	}
	http.ListenAndServe(port, nil)
}

func registerHandlers(router *mux.Router) {
	router.HandleFunc("/item", handlerItem).Methods("POST", "PUT")
	router.HandleFunc("/item/{id:[0-9]+}", handlerItemId).Methods("GET", "DELETE")
	router.HandleFunc("/items", handlerList).Methods("GET")

	router.PathPrefix("/swagger").Handler(swagger.WrapHandler)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<p>Shop API</p><a href=/swagger/>Swagger</a>")
	})

	http.Handle("/", router)
}
