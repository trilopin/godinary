package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/trilopin/godinary/imagejob"
	"github.com/trilopin/godinary/storage"
)

// Port exposed by http server
var Port string

// AllowedReferers is the list of hosts allowed to request images
var AllowedReferers []string

func init() {
	Port = os.Getenv("GODINARY_PORT")
	AllowedReferers = strings.Split(os.Getenv("GODINARY_ALLOW_HOSTS"), ",")

	if os.Getenv("GODINARY_STORAGE") == "gs" {
		storage.StorageDriver = storage.NewGoogleStorageDriver()
	} else {
		storage.StorageDriver = storage.NewFileDriver()
	}

	sort.Strings(AllowedReferers)
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
		"/concurrency": imagejob.Concurrency,
	}

	fmt.Println("Listening on port", Port)
	server.ListenAndServe()
}

type myHandler struct{}

// ServeHTTP manage custom url multiplexing avoiding path.clean in
// default go http mux.
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Manage authorization by htp Referer
	// List of authorized referers should provisioned via GODINARY_ALLOW_HOSTS
	// environment variable. Empty referer is always allowed because
	// development issues
	var (
		allowed     bool
		httpReferer string
	)

	httpReferer = r.Header.Get("Referer")
	for _, domain := range AllowedReferers {
		if domain == httpReferer {
			allowed = true
		}
	}

	if !allowed {
		http.Error(w, "Referer not allowed", http.StatusForbidden)
		return
	}
	// Manage route is is allowed
	for key, h := range mux {
		if strings.Index(r.URL.String(), key) == 0 {
			h(w, r)
		}
	}
}
