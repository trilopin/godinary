package imagejob

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
)

func Middleware(domain string, allowedReferers []string, handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return raven.RecoveryHandler(logger(domainValidator(domain, refererValidator(allowedReferers, handler))))
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
