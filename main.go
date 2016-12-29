package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/trilopin/godinary/imagejob"
)

func main() {
	port := flag.Int("port", 3001, "Port to listen to")
	flag.Parse()
	router := httprouter.New()

	router.GET("/v0.1/fetch/*info", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		img := imagejob.ImageJob{}
		_ = img.New(ps.ByName("info")[1:])
		_ = img.Download()
		_ = img.Process(w)
		w.Header().Set("Content-Type", "image/"+img.TargetFormat)
	})

	fmt.Println("Listening on port", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), router))
}
