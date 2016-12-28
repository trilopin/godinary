package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/trilopin/godinary/parser"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to Homepage!\n")
}

func Fetch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	img := parser.ImageJob{}
	_ = img.New(ps.ByName("info")[1:])
	_ = img.Download()
	_ = img.Process(w)
	w.Header().Set("Content-Type", "image/"+img.TargetFormat)
}

func StartServer(port int) {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/v0.1/:account/fetch/*info", Fetch)
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}
