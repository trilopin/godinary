package main

import (
	"flag"
	"log"
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

	log.SetOutput(os.Stdout)
}

func main() {
	opts := &imagejob.ServerOpts{
		Port:                viper.GetString("port"),
		Domain:              viper.GetString("domain"),
		AllowedReferers:     strings.Split(viper.GetString("allow_hosts"), ","),
		MaxRequest:          viper.GetInt("max_request"),
		MaxRequestPerDomain: viper.GetInt("max_request_domain"),
	}

	if viper.GetString("storage") == "gs" {
		opts.GCEProject = viper.GetString("gce_project")
		if opts.GCEProject == "" {
			panic("GODINARY_GCE_PROJECT should be setted")
		}
		opts.GSBucket = viper.GetString("gs_bucket")
		if opts.GSBucket == "" {
			panic("GODINARY_GS_BUCKET should be setted")
		}
		opts.GSCredencials = viper.GetString("gs_credentials")
		if opts.GSCredencials == "" {
			panic("GODINARY_GS_CREDENTIUALS shold be setted")
		}
		opts.StorageDriver = storage.NewGoogleStorageDriver(opts.GCEProject, opts.GSBucket, opts.GSCredencials)
	} else {
		opts.FSBase = viper.GetString("fs_base")
		if opts.FSBase == "" {
			panic("GODINARY_FS_BASE should be setted")
		}
		opts.StorageDriver = storage.NewFileDriver(opts.FSBase)
	}

	imagejob.Serve(opts)
}
