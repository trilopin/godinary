package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trilopin/godinary/imagejob"
	"github.com/trilopin/godinary/storage"
)

// Port exposed by http server
var Port string

// AllowedReferers is the list of hosts allowed to request images
var AllowedReferers []string

func setupConfig() {
	// flags setup
	flag.String("domain", "", "Domain to validate with Host header, it will deny any other request (if port is not standard must be passed as host:port)")
	flag.String("sentry_url", "", "Sentry DSN for error tracking")
	flag.String("release", "", "Release hash to notify sentry")
	flag.String("allow_hosts", "", "Domains authorized to ask godinary separated by commas (A comma at the end allows empty referers)")
	flag.String("port", "3002", "Port where the https server listen")
	flag.String("ssl_dir", "", "Path to directory with server.key and server.pem SSL files")
	flag.Int("max_request", 100, "Maximum number of simultaneous downloads")
	flag.Int("max_request_domain", 10, "Maximum number of simultaneous downloads per domain")
	flag.String("cdn_ttl", "604800", "Number of seconds images wil be cached in CDN")
	flag.String("storage", "fs", "Storage type: 'gs' for google storage or 'fs' for filesystem")
	flag.String("fs_base", "", "FS option: Base dir for filesystem storage")
	flag.String("gce_project", "", "GS option: Sentry DSN for error tracking")
	flag.String("gs_bucket", "", "GS option: Bucket name")
	flag.String("gs_credentials", "", "GS option: Path to service account file with Google Storage credentials")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// env setup
	viper.AutomaticEnv()
	viper.SetEnvPrefix("godinary")
	flag.VisitAll(func(f *flag.Flag) {
		viper.BindEnv(f.Name)
	})
}

func init() {

	setupConfig()
	Port = viper.GetString("port")

	if viper.GetString("sentry_url") != "" {
		raven.SetDSN(viper.GetString("sentry_url"))
		if viper.GetString("release") != "" {
			raven.SetRelease(viper.GetString("release"))
		}
		raven.CapturePanic(func() {
			// do all of the scary things here
		}, nil)
	}

	AllowedReferers = strings.Split(viper.GetString("allow_hosts"), ",")

	if viper.GetString("storage") == "gs" {
		storage.StorageDriver = storage.NewGoogleStorageDriver()
	} else {
		storage.StorageDriver = storage.NewFileDriver()
	}

	imagejob.MaxRequest = viper.GetInt("max_request")
	imagejob.MaxRequestPerDomain = viper.GetInt("max_request_domain")

	// globalSemaphore controls concurrent http client requests
	imagejob.SpecificThrotling = make(map[string]chan struct{}, 20)
	imagejob.GlobalThrotling = make(chan struct{}, imagejob.MaxRequest)

	sort.Strings(AllowedReferers)
	log.SetOutput(os.Stdout)
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	var err error

	server := http.Server{
		Addr:    ":" + Port,
		Handler: &myHandler{},
	}

	mux = map[string]func(http.ResponseWriter, *http.Request){
		"/image/fetch/": raven.RecoveryHandler(imagejob.Fetch),
		"/robots.txt": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "User-Agent: *")
			fmt.Fprintln(w, "Disallow: /")
		},
		"/up": func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "up")
		},
	}

	fmt.Println("Listening with SSL on port", Port)
	if SSLDir := viper.GetString("ssl_dir"); SSLDir == "" {
		err = server.ListenAndServe()
	} else {
		err = server.ListenAndServeTLS(SSLDir+"server.pem", SSLDir+"server.key")
	}
	if err != nil {
		raven.CaptureError(err, nil)
		log.Fatal("ListenAndServe cannot start: ", err)
	}
}

type myHandler struct{}

// ServeHTTP manage custom url multiplexing avoiding path.clean in
// default go http mux.
func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Manage authorization by http Referer
	// List of authorized referers should provisioned via GODINARY_ALLOW_HOSTS
	// environment variable. Empty referer is always allowed because
	// development issues
	var (
		allowed     bool
		httpReferer string
	)

	httpReferer = r.Header.Get("Referer")
	domain := viper.GetString("domain")
	if domain != "" && r.Host != domain {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if httpReferer != "" {
		info, _ := url.Parse(httpReferer)
		for _, domain := range AllowedReferers {
			if domain != "" && strings.HasSuffix(info.Host, domain) {
				allowed = true
				break
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
