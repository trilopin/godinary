package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
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

func main() {
	router := httprouter.New()

	router.GET("/v0.1/fetch/*info", imagejob.Fetch)

	fmt.Println("Listening on port", Port)
	log.Fatal(http.ListenAndServe(":"+Port, router))
}
