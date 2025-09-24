package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	var port string
	var addr string

	for {
		fmt.Print("Ingrese el puerto que desea utilizar (ej:8000): ")
		fmt.Scan(&port)

		// Probar si el puerto está disponible
		ln, err := net.Listen("tcp", ":"+addr)
		if err != nil {
			fmt.Printf("El puerto %s está en uso o no tienes permisos ❌\n", port)
			continue
		}

		// Puerto libre → cerramos el listener temporal y salimos del loop
		ln.Close()
		fmt.Printf("El servidor esta escuchando el puerto %s ✅\n", port)
		break
	}

	// Levantar servidor HTTP en el puerto elegido
	log.Fatal(http.ListenAndServe(":"+addr, nil))
}
