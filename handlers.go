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
	DEFAULT_PAGE_LENGTH = 20
)

func handlerItem(w http.ResponseWriter, r *http.Request) {
	var item items.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "POST":
		handlerItemPost(&item, w)
	case "PUT":
		handlerItemPut(&item, w)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Add item
// @Description Add single item
// @Tags items
// @Accept  json
// @Param item body items.Item true "Item to add"
// @Success 200
// @Router /item [post]
func handlerItemPost(item *items.Item, w http.ResponseWriter) {
	if err := repo.Add(item); err != nil {
		http.Error(w, err.Error(), getStatusFromError(err))
	}
}

// @Summary Update item
// @Description Update single item
// @Tags items
// @Accept  json
// @Param item body items.Item true "New item to update. Item with updated id should be already added."
// @Success 200
// @Router /item [put]
func handlerItemPut(item *items.Item, w http.ResponseWriter) {
	if err := repo.Update(item); err != nil {
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

// @Summary Get item
// @Description Get single item by id
// @Tags items
// @Produce json
// @Param id path int true "Item id"
// @Success 200 {object} items.Item
// @Router /item/{id} [get]
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

// @Summary Delete item
// @Description Delete single item by id
// @Tags items
// @Param id path int true "Item id"
// @Success 200
// @Router /item/{id} [delete]
func handlerItemDelete(id uint64, w http.ResponseWriter) {
	if err := repo.Delete(id); err != nil {
		http.Error(w, err.Error(), getStatusFromError(err))
	}
}

// @Summary List items
// @Description List all items with optional pagination
// @Tags items
// @Produce json
// @Param page query int false "Return items from this page. If not specified, return all items"
// @Success 200 {array} items.Item
// @Router /items [get]
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
