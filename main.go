package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	var port string

	for {
		fmt.Print("Ingrese el puerto que desea utilizar (ej:8000): ")
		fmt.Scanf("%s", &port)

		// Probar si el puerto está disponible
		ln, err := net.Listen("tcp", ":"+port)
		if err != nil {
			fmt.Printf("El puerto %s está en uso o no tienes permisos ❌\n", port)
			continue
		}

		// Puerto libre → cerramos el listener temporal y salimos del loop
		ln.Close()
		fmt.Printf("El servidor esta escuchando el puerto %s ✅\n", port)
		break
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method == http.MethodPost {
			hostname, _ := os.Hostname()
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"hostname": hostname, "port": port, "images": ""})
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Levantar servidor HTTP en el puerto elegido
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
