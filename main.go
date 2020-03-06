package main

import (
	"net/http"

	"github.com/CurtizJ/dummy-shop/repos"
	"github.com/gorilla/mux"
)

var repo repos.Repo

func main() {
	repo = repos.NewPostgresRepo()
	router := mux.NewRouter()
	registerHandlers(router)
	http.ListenAndServe(":8181", nil)
}

func registerHandlers(r *mux.Router) {
	r.HandleFunc("/item", handlerItem).Methods("POST", "PUT")
	r.HandleFunc("/item/{id:[0-9]+}", handlerItemId).Methods("GET", "DELETE")
	r.HandleFunc("/items", handlerList).Methods("GET")

	http.Handle("/", r)
}
