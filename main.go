package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
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
	if Port == "" {
		Port = "3002"
	}

	AllowedReferers = strings.Split(os.Getenv("GODINARY_ALLOW_HOSTS"), ",")

	if os.Getenv("GODINARY_STORAGE") == "gs" {
		storage.StorageDriver = storage.NewGoogleStorageDriver()
	} else {
		storage.StorageDriver = storage.NewFileDriver()
	}

	imagejob.MaxRequest, _ = strconv.Atoi(os.Getenv("GODINARY_MAX_REQUEST"))
	if imagejob.MaxRequest == 0 {
		imagejob.MaxRequest = 100
	}
	imagejob.MaxRequestPerDomain, _ = strconv.Atoi(os.Getenv("GODINARY_MAX_REQUEST_DOMAIN"))
	if imagejob.MaxRequestPerDomain == 0 {
		imagejob.MaxRequestPerDomain = 10
	}

	// globalSemaphore controls concurrent http client requests
	imagejob.SpecificThrotling = make(map[string]chan struct{}, 20)
	imagejob.GlobalThrotling = make(chan struct{}, imagejob.MaxRequest)

	sort.Strings(AllowedReferers)
	log.SetOutput(os.Stdout)
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":" + Port,
		Handler: &myHandler{},
	}

	mux = map[string]func(http.ResponseWriter, *http.Request){
		"/image/fetch/": imagejob.Fetch,
	}

	fmt.Println("Listening on port", Port)
	// server.ListenAndServe()
	err := server.ListenAndServeTLS("server.pem", "server.key")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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
	if httpReferer != "" {
		info, _ := url.Parse(httpReferer)
		for _, domain := range AllowedReferers {
			if domain == info.Host {
				allowed = true
			}
		}

		if !allowed {
			http.Error(w, "Referer not allowed", http.StatusForbidden)
			return
		}
	}
	// Manage route is is allowed
	for key, h := range mux {
		if strings.Index(r.URL.String(), key) == 0 {
			h(w, r)
		}
	}
}
