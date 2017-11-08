package main

import (
	"flag"
	"log"
	"os"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trilopin/godinary/http"
	"github.com/trilopin/godinary/storage"
)

func setupConfig() {
	// flags setup
	flag.String("domain", "", "Domain to validate with Host header, it will deny any other request (if port is not standard must be passed as host:port)")
	flag.String("auth", "", "List of apikey,apisecret (one per line) allowed to use API")
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

	log.SetOutput(os.Stdout)
}

func main() {
	var err error

	opts := &http.ServerOpts{
		Port:                viper.GetString("port"),
		Domain:              viper.GetString("domain"),
		AllowedReferers:     strings.Split(viper.GetString("allow_hosts"), ","),
		MaxRequest:          viper.GetInt("max_request"),
		MaxRequestPerDomain: viper.GetInt("max_request_domain"),
		SSLDir:              viper.GetString("ssl_dir"),
		CDNTTL:              viper.GetString("cdn_ttl"),
	}
	opts.APIAuth = make(map[string]string)
	auth := viper.GetString("auth")
	if auth != "" {
		parts := strings.Split(auth, ",")
		if len(parts) == 2 {
			opts.APIAuth[parts[0]] = parts[1]
		} else {
			log.Fatalln("Invalid auth ", parts)
		}
	}

	if viper.GetString("storage") == "gs" {
		opts.GCEProject = viper.GetString("gce_project")
		opts.GSBucket = viper.GetString("gs_bucket")
		opts.GSCredentials = viper.GetString("gs_credentials")
		if opts.GCEProject == "" {
			log.Fatalln("GoogleStorage project should be setted")
		}
		if opts.GSBucket == "" {
			log.Fatalln("GoogleStorage bucket should be setted")
		}
		if opts.GSCredentials == "" {
			log.Fatalln("GoogleStorage Credentials shold be setted")
		}
		opts.StorageDriver = &storage.GoogleStorageDriver{
			BucketName:  opts.GSBucket,
			ProjectName: opts.GCEProject,
			Credentials: opts.GSCredentials,
		}
	} else {
		opts.FSBase = viper.GetString("fs_base")
		if opts.FSBase == "" {
			log.Fatalln("filesystem base path should be setted: %v", err)
		}
		opts.StorageDriver = storage.NewFileDriver(opts.FSBase)
	}

	http.Serve(opts)
}
