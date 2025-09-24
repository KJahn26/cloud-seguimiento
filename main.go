package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		json.NewEncoder(w).Encode(r.Body)
	})

	log.Fatal(http.ListenAndServe(":9097", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
