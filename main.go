package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/trilopin/godinary/imagejob"
)

var Port string

func init() {

	Port = os.Getenv("GODINARY_PORT")
	if Port == "" {
		Port = "3002"
	}
	log.SetOutput(os.Stdout)
}

type Handler struct{}

// ServeHTTP overrides default mux. This replacement is needed in
// order to avoid doble slash removing in url
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	handler := new(Handler)
	http.HandleFunc("/v0.1/fetch/", imagejob.Fetch)
	fmt.Println("Listening on port", Port)
	log.Fatal(http.ListenAndServe(":"+Port, handler))
}
