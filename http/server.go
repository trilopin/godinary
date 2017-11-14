package http

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/trilopin/godinary/storage"
)

var (
	// SpecificThrotling is semaphore per domain
	SpecificThrotling map[string]chan struct{}
	// GlobalThrotling is global semaphore
	GlobalThrotling chan struct{}
)

// ServerOpts contains confif for Application
type ServerOpts struct {
	MaxRequest          int
	MaxRequestPerDomain int
	Port                string
	Domain              string
	SSLDir              string
	CDNTTL              string
	AllowedReferers     []string
	StorageDriver       storage.Driver
	FSBase              string
	GCEProject          string
	GSBucket            string
	GSCredentials       string
	APIAuth             map[string]string
}

// ------------------------------------
//             Main Server
// ------------------------------------

// Serve is the main func for start http server
func Serve(opts *ServerOpts) {
	var err error

	// semaphores control concurrent http client requests
	SpecificThrotling = make(map[string]chan struct{}, 20)
	GlobalThrotling = make(chan struct{}, opts.MaxRequest)

	mux := &Mux{
		Routes: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
	mux.Handle("/robots.txt", Middleware(RobotsTXT, opts))
	mux.Handle("/up", Up)
	mux.Handle("/image/fetch/", Middleware(Fetch(opts), opts))
	mux.Handle("/image/upload/", Middleware(Upload(opts), opts))
	mux.Handle("/v1_0/image/upload", AuthMiddleware(APIUpload(opts), opts))
	server := http.Server{
		Addr:    ":" + opts.Port,
		Handler: mux,
	}

	if opts.SSLDir == "" {
		log.Println("Listening on port", opts.Port)
		err = server.ListenAndServe()
	} else {
		log.Println("Listening with SSL on port", opts.Port)
		err = server.ListenAndServeTLS(opts.SSLDir+"server.pem", opts.SSLDir+"server.key")
	}

	if err != nil {
		log.Fatalln("ListenAndServe cannot start: ", err)
		raven.CaptureError(err, nil)
	}
}

// ------------------------------------
//             Mux
// ------------------------------------

// Mux is the custom Router needed in order to avoid URL cleaning
type Mux struct {
	Routes map[string]func(http.ResponseWriter, *http.Request)
}

// Handle adds new route with their handler
func (mux *Mux) Handle(route string, handler func(w http.ResponseWriter, r *http.Request)) {
	mux.Routes[route] = handler
}

// ServeHTTP manage custom url multiplexing avoiding path.clean in
// default go http mux.
func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, h := range mux.Routes {
		if strings.Index(r.URL.String(), key) == 0 {
			h(w, r)
		}
	}
}

// ------------------------------------
//             Middlewares
// ------------------------------------

// Middleware composes all middlewares
func Middleware(handler func(w http.ResponseWriter, r *http.Request), opts *ServerOpts) http.HandlerFunc {
	return raven.RecoveryHandler(logger(domainValidator(opts.Domain, refererValidator(opts.AllowedReferers, handler))))
}

// AuthMiddleware composes all middlewares + auth
func AuthMiddleware(handler func(w http.ResponseWriter, r *http.Request), opts *ServerOpts) http.HandlerFunc {
	return raven.RecoveryHandler(logger(domainValidator(opts.Domain, refererValidator(opts.AllowedReferers, auth(opts.APIAuth, handler)))))
}

// domainValidator is a middleware to check Host Header against configured domain
func domainValidator(domain string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if domain != "" && r.Host != domain {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// logger is a middleware which writes to standard output time and urls
func logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(time.Since(start), r.URL.Path)
	})
}

//refererValidator is a middleware to check Http-referer headers
func refererValidator(allowedReferers []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed := false
		httpReferer := r.Header.Get("Referer")
		if httpReferer != "" {
			info, _ := url.Parse(httpReferer)
			for _, domain := range allowedReferers {
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
		next.ServeHTTP(w, r)
	})
}

// auth validates API_KEY and API_SIGNATURE against shared API_SECRET
func auth(auth map[string]string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		APIKey := r.FormValue("apikey")
		signature := r.FormValue("signature")
		timestamp := r.FormValue("timestamp")
		if APIKey == "" || signature == "" {
			http.Error(w, "", http.StatusForbidden)
			return
		}
		storedSecret, ok := auth[APIKey]
		if !ok {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		ht := sha256.New()
		ht.Write([]byte(APIKey + timestamp + storedSecret))
		hash := hex.EncodeToString(ht.Sum(nil))
		if hash != signature {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
