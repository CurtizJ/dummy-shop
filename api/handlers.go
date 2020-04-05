package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CurtizJ/dummy-shop/lib/errors"
	"github.com/CurtizJ/dummy-shop/api/items"
	"github.com/gorilla/mux"
)

func handlerItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var item items.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		errors.ReportAsJSON(w, err.Error(), http.StatusBadRequest)
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
		errors.ReportErrorAsJSON(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(item); err != nil {
		errors.ReportErrorAsJSON(w, err)
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
		errors.ReportErrorAsJSON(w, err)
		return
	}
	if err := json.NewEncoder(w).Encode(item); err != nil {
		errors.ReportErrorAsJSON(w, err)
	}
}

func handlerItemId(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	strId := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		errors.ReportAsJSON(w, fmt.Sprintf("Invalid id: %s", strId), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		handlerItemGet(id, w)
	case "DELETE":
		handlerItemDelete(id, w)
	default:
		errors.ReportAsJSON(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
		errors.ReportErrorAsJSON(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(item); err != nil {
		errors.ReportErrorAsJSON(w, err)
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
		errors.ReportErrorAsJSON(w, err)
	}
}

// @Summary List items
// @Description List all items with optional pagination
// @Tags items
// @Produce json
// @Param length query int false "Return at most 'length' items."
// @Param offset query int false "Skip first 'offset' items. Must be specified with 'length'"
// @Success 200 {array} items.Item
// @Router /items [get]
func handlerList(w http.ResponseWriter, r *http.Request) {
	var list []items.Item
	var err error

	strLength, existsLimit := r.URL.Query()["length"]
	strOffset, existsOffset := r.URL.Query()["offset"]

	if existsOffset && !existsLimit {
		errors.ReportAsJSON(w, "Offset must be specified only with length", http.StatusBadRequest)
		return
	}

	if existsLimit {
		length, err := strconv.ParseUint(strLength[0], 10, 64)
		if err != nil {
			errors.ReportAsJSON(w, fmt.Sprintf("Invalid value for length: %s", strLength[0]), http.StatusBadRequest)
			return
		}

		offset := uint64(0)
		if existsOffset {
			offset, err = strconv.ParseUint(strOffset[0], 10, 64)
			if err != nil {
				errors.ReportAsJSON(w, fmt.Sprintf("Invalid value for offset: %s", strLength[0]), http.StatusBadRequest)
				return
			}
		}

		list, err = repo.List(length, offset)
	} else {
		list, err = repo.ListAll()
	}

	if err != nil {
		errors.ReportErrorAsJSON(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		errors.ReportErrorAsJSON(w, err)
	}
}
