package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CurtizJ/dummy-shop/errors"
	"github.com/CurtizJ/dummy-shop/items"
	"github.com/gorilla/mux"
)

const (
	DEFAULT_PAGE_LENGTH = 2
)

func handlerItem(w http.ResponseWriter, r *http.Request) {
	var item items.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	switch r.Method {
	case "POST":
		err = repo.Add(&item)
	case "PUT":
		err = repo.Update(&item)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	if err != nil {
		http.Error(w, err.Error(), getStatusFromError(err))
	}
}

func handlerItemId(w http.ResponseWriter, r *http.Request) {
	strId := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid id: %s", strId), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		handlerItemGet(id, w)
	case "DELETE":
		handlerItemDelete(id, w)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handlerItemGet(id uint64, w http.ResponseWriter) {
	item, err := repo.Get(id)
	if err != nil {
		http.Error(w, err.Error(), getStatusFromError(err))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerItemDelete(id uint64, w http.ResponseWriter) {
	if err := repo.Delete(id); err != nil {
		http.Error(w, err.Error(), getStatusFromError(err))
	}
}

func handlerList(w http.ResponseWriter, r *http.Request) {
	var list []items.Item
	var err error

	strPage, exists := r.URL.Query()["page"]
	if exists {
		page, err := strconv.ParseUint(strPage[0], 10, 64)
		if err != nil || page == 0 {
			http.Error(w, fmt.Sprintf("Invalid value for page: %s", strPage[0]), http.StatusBadRequest)
			return
		}
		list, err = repo.List(DEFAULT_PAGE_LENGTH, DEFAULT_PAGE_LENGTH*(page-1))
	} else {
		list, err = repo.ListAll()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getStatusFromError(err error) int {
	if appErr, ok := err.(errors.ApplicationError); ok {
		return appErr.GetHTTPStatus()
	}

	return http.StatusInternalServerError
}
