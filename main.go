package main

import (
	"encoding/json"
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
