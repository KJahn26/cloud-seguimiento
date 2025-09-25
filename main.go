package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Respuesta struct {
	Hostname string      `json:"hostname"`
	Port     string      `json:"port"`
	Images   []ImageData `json:"images"`
}

type ImageData struct {
	Filename string `json:"filename"`
	Base64   string `json:"base64"`
}

func main() {
	var port string

	for {

		//pedimos al usuario ingresar el puerto a utilizar
		fmt.Print("Ingrese el puerto que desea utilizar (ej:8000): ")
		//lo guardamos en port
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

	// respuesta del servidor a solicitud
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		//metodo get para enviarle un json
		if r.Method == http.MethodGet {
			hostname, _ := os.Hostname()
			w.Header().Set("Content-type", "application/json")
			respuesta := Respuesta{
				Hostname: hostname,
				Port:     port,
				Images:   cargarImagenesAleatorias("images", 4),
			}
			json.NewEncoder(w).Encode(respuesta)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	//cargamos la pagina
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template/index.html")
	})

	// Levantar servidor HTTP en el puerto elegido
	log.Fatal(http.ListenAndServe(":"+port, nil))

	//cargamos imagensita de la nube
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("template"))))

}

// cargamos las imagenes de la ruta
func cargarImagenesAleatorias(ruta string, cantidad int) []ImageData {
	archivos, err := ioutil.ReadDir(ruta)
	if err != nil {
		log.Fatal("Error leyendo carpeta:", err)
	}

	var imagenesValidas []string

	//descartamos los archivos que no nos sirven y tomamos las validas
	for _, archivo := range archivos {
		if !archivo.IsDir() {
			ext := strings.ToLower(filepath.Ext(archivo.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				imagenesValidas = append(imagenesValidas, archivo.Name())
			}
		}
	}

	//Mezclamos aleatoriamente
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(imagenesValidas), func(i, j int) {
		imagenesValidas[i], imagenesValidas[j] = imagenesValidas[j], imagenesValidas[i]
	})

	//ajusta la cantidad para que no haya desbordes ni errores
	if cantidad > len(imagenesValidas) {
		cantidad = len(imagenesValidas)
	}

	var resultado []ImageData

	// Codificamos cada imagen seleccionada en Base64
	for i := 0; i < cantidad; i++ {
		rutaCompleta := filepath.Join(ruta, imagenesValidas[i])
		contenido, err := ioutil.ReadFile(rutaCompleta)
		if err != nil {
			continue
		}
		//el tipo MIME según extensión
		mime := "image/jpeg"
		if strings.HasSuffix(imagenesValidas[i], ".png") {
			mime = "image/png"
		}

		//codificamos y Agregamos la imagen al resultado
		base64Str := fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(contenido))
		resultado = append(resultado, ImageData{
			Filename: imagenesValidas[i],
			Base64:   base64Str,
		})
	}

	return resultado
}
