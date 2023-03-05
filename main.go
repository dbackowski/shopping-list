package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os/exec"
)

var tmpl *template.Template

type Item struct {
	UUID string
	Name string
	Done bool
}

var items = []Item{
	{UUID: generateUUID(), Name: "Eggs x 10", Done: false},
}

func listItems(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, items)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var item Item

		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			w.WriteHeader(http.StatusProcessing)
			return
		}
		item.UUID = generateUUID()
		items = append(items, item)
		json.NewEncoder(w).Encode(item)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func generateUUID() string {
	newUUID, err := exec.Command("uuidgen").Output()

	if err != nil {
		log.Fatal(err)
	}

	return string(newUUID)
}

func main() {
	mux := http.NewServeMux()
	tmpl = template.Must(template.ParseFiles("templates/index.gohtml"))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", listItems)
	mux.HandleFunc("/create", addItem)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
