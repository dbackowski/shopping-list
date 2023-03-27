package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/google/uuid"
)

type Item struct {
	UUID string
	Name string
	Done bool
}

type withId struct {
	id string
}

var items = []Item{}

func Serve(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var id string

	p := r.URL.Path
	switch {
	case match(p, "/static/([^/]+[css|js]$)"):
		h = get(serveStaticFiles)
	case match(p, "/alive"):
		h = get(alive)
	case match(p, "/"):
		h = get(serveIndex)
	case match(p, "/items"):
		h = get(listItems)
	case match(p, "/items/create"):
		h = post(addItem)
	case match(p, "/items/update/([^/]+)", &id):
		h = put(withId{id}.updateItem)
	case match(p, "/items/delete/([^/]+)", &id):
		h = delete(withId{id}.deleteItem)
	default:
		http.NotFound(w, r)
		return
	}
	h.ServeHTTP(w, r)
}

func match(path, pattern string, vars ...*string) bool {
	regex := regexp.MustCompile("^" + pattern + "$")
	matches := regex.FindStringSubmatch(path)

	if len(matches) <= 0 {
		return false
	}

	if len(vars) > 0 {
		for i, match := range matches[1:] {
			*vars[i] = match
		}
	}

	return true
}

func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "GET")
}

func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "POST")
}

func put(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "PUT")
}

func delete(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "DELETE")
}

func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Allow", method)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}

func findItemIndexByUuid(uuid string) (int, error) {
	for i, v := range items {
		if v.UUID == uuid {
			return i, nil
		}
	}

	return -1, fmt.Errorf("UUID: %s not found in items", uuid)
}

func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	var path = "." + r.URL.Path

	_, error := os.Stat(path)

	if os.IsNotExist(error) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, path)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func alive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "{\"alive\": true}")
}

func listItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	item.UUID = generateUUID()
	items = append(items, item)
	json.NewEncoder(w).Encode(item)
}

func (h withId) updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(h.id) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	item_index, err := findItemIndexByUuid(h.id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var item = items[item_index]

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	items[item_index] = item

	json.NewEncoder(w).Encode(items)
}

func (h withId) deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(h.id) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	item_index, err := findItemIndexByUuid(h.id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	items = append(items[:item_index], items[item_index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func generateUUID() string {
	id := uuid.New()
	return id.String()
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(Serve)))
}
