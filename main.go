package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":" + Port,
		Handler: &myHandler{},
	}

	mux = map[string]func(http.ResponseWriter, *http.Request){
		"/v0.1/fetch/": imagejob.Fetch,
	}

	fmt.Println("Listening on port", Port)
	server.ListenAndServe()
}

type myHandler struct{}

// ServeHTTP manage custom url multiplexing avoiding path.clean in
// default go http mux.
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, h := range mux {
		if strings.Index(r.URL.String(), key) == 0 {
			h(w, r)
		}
	}
}
