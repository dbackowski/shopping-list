package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var tmpl *template.Template

type Item struct {
	UUID string
	Name string
	Done bool
}

type withId struct {
	id string
}

var items = []Item{
	{UUID: generateUUID(), Name: "Eggs x 10", Done: false},
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	var id string

	p := r.URL.Path
	switch {
	case match(p, "/static/([^/]+[css|js]$)"):
		h = get(serveStaticFiles)
	case match(p, "/"):
		h = get(listItems)
	case match(p, "/alive"):
		h = get(alive)
	case match(p, "/create"):
		h = post(addItem)
	case match(p, "/update/([^/]+)", &id):
		h = put(withId{id}.updateItem)
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

func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	var path = "." + r.URL.Path

	_, error := os.Stat(path)

	if os.IsNotExist(error) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, path)
}

func alive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "{\"alive\": true}")
}

func listItems(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, items)
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
	id := strings.TrimPrefix(r.URL.Path, "/update/")

	if len(id) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	for i, v := range items {
		if v.UUID == id {
			var item = items[i]

			if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			items[i] = item
		}
	}

	json.NewEncoder(w).Encode(items)
}

func generateUUID() string {
	newUUID, err := exec.Command("uuidgen").Output()

	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSuffix(string(newUUID), "\n")
}

func main() {
	tmpl = template.Must(template.ParseFiles("templates/index.gohtml"))
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(Serve)))
}
