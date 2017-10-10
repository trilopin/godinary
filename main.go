package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trilopin/godinary/imagejob"
	"github.com/trilopin/godinary/storage"
)

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

	if viper.GetString("sentry_url") != "" {
		raven.SetDSN(viper.GetString("sentry_url"))
		if viper.GetString("release") != "" {
			raven.SetRelease(viper.GetString("release"))
		}
		raven.CapturePanic(func() {
			// do all of the scary things here
		}, nil)
	}

	if viper.GetString("storage") == "gs" {
		storage.StorageDriver = storage.NewGoogleStorageDriver()
	} else {
		storage.StorageDriver = storage.NewFileDriver()
	}

	log.SetOutput(os.Stdout)
}

func main() {
	var err error

	opts := &imagejob.ServerOpts{
		Port:                viper.GetString("port"),
		Domain:              viper.GetString("domain"),
		AllowedReferers:     strings.Split(viper.GetString("allow_hosts"), ","),
		MaxRequest:          viper.GetInt("max_request"),
		MaxRequestPerDomain: viper.GetInt("max_request_domain"),
	}

	// semaphores control concurrent http client requests
	imagejob.SpecificThrotling = make(map[string]chan struct{}, 20)
	imagejob.GlobalThrotling = make(chan struct{}, opts.MaxRequest)

	mux := &imagejob.Mux{
		Routes: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
	mux.Handle("/robots.txt", imagejob.Middleware(imagejob.RobotsTXT, opts))
	mux.Handle("/up", imagejob.Middleware(imagejob.Up, opts))
	mux.Handle("/image/fetch/", imagejob.Middleware(imagejob.Fetch(opts), opts))
	server := http.Server{
		Addr:    ":" + opts.Port,
		Handler: mux,
	}

	if SSLDir := viper.GetString("ssl_dir"); SSLDir == "" {
		fmt.Println("Listening on port", opts.Port)
		err = server.ListenAndServe()
	} else {
		fmt.Println("Listening with SSL on port", opts.Port)
		err = server.ListenAndServeTLS(SSLDir+"server.pem", SSLDir+"server.key")
	}

	if err != nil {
		log.Fatal("ListenAndServe cannot start: ", err)
		raven.CaptureError(err, nil)
	}
}
