package imagejob

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
)

type Mux struct {
	Routes map[string]func(http.ResponseWriter, *http.Request)
}

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

func Middleware(handler func(w http.ResponseWriter, r *http.Request), opts *ServerOpts) http.HandlerFunc {
	return raven.RecoveryHandler(logger(domainValidator(opts.Domain, refererValidator(opts.AllowedReferers, handler))))
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
