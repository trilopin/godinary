package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/trilopin/godinary/importer"
	"github.com/trilopin/godinary/storage"
)

func setupConfig() {
	// flags setup
	flag.String("storage", "fs", "Storage type: 'gs' for google storage or 'fs' for filesystem")
	flag.String("fs_base", "", "FS option: Base dir for filesystem storage")
	flag.String("gce_project", "", "GS option: Sentry DSN for error tracking")
	flag.String("gs_bucket", "", "GS option: Bucket name")
	flag.String("gs_credentials", "", "GS option: Path to service account file with Google Storage credentials")
	flag.String("cloudinary_userspace", "", "Cloudinary User Space")
	flag.String("cloudinary_apikey", "", "Cloudinary API Key")
	flag.String("cloudinary_apisecret", "", "Cloudinary API Secret")
	flag.String("memprofile", "", "write memory profile to `file`")

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

	log.SetOutput(os.Stdout)
}

func main() {
	var sd storage.Driver
	var err error

	if viper.GetString("storage") == "gs" {
		GCEProject := viper.GetString("gce_project")
		GSBucket := viper.GetString("gs_bucket")
		CSCredentials := viper.GetString("gs_credentials")
		if GCEProject == "" {
			log.Fatalln("GoogleStorage project should be setted")
		}
		if GSBucket == "" {
			log.Fatalln("GoogleStorage bucket should be setted")
		}
		if CSCredentials == "" {
			log.Fatalln("GoogleStorage Credentials shold be setted")
		}
		sd = &storage.GoogleStorageDriver{
			BucketName:  GSBucket,
			ProjectName: GCEProject,
			Credentials: CSCredentials,
		}
		if err != nil {
			log.Fatalf("can not create GoogleStorage Driver: %v", err)
		}
	} else {
		FSBase := viper.GetString("fs_base")
		if FSBase == "" {
			log.Fatalln("filesystem base path should be setted")
		}
		sd = storage.NewFileDriver(FSBase)
	}

	fmt.Println("Import")
	ci := &importer.CloudinaryImporter{
		UserSpace: viper.GetString("cloudinary_userspace"),
		APIKey:    viper.GetString("cloudinary_apikey"),
		APISecret: viper.GetString("cloudinary_apisecret"),
	}

	defer func() {
		memProfile := viper.GetString("memprofile")
		if memProfile != "" {
			f, err := os.Create(memProfile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			f.Close()
		}
	}()
	ci.Import(sd)
}
