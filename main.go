package main

import (
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

	log.Fatal(http.ListenAndServe(":8080", mux))
}
