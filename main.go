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

		err := img.New(ps.ByName("info")[1:])
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = img.Download()
		if err != nil {
			http.Error(w, "Cannot download image", http.StatusInternalServerError)
			return
		}

		err = img.Process(w)
		if err != nil {
			http.Error(w, "Cannot process Image", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/"+img.TargetFormat)
	})

	fmt.Println("Listening on port", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), router))
}
