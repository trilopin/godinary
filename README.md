# godinary
Image proxy with live resize &amp; tranformations


### Install
```
git clone https://github.com/trilopin/godinary
```


### Docker flow:
- make build -> compiles and build docker image
- make run -> start server

### Development flow:
- mkdir data && cp .env.example .env
- make build-dev
- make run

### Configuration
Variables can be passed as arguments or as env vars (uppercase and with GODINARY_ prefix)
```
$ godinary -h
Usage of godinary:
      --allow_hosts string       Domains authorized to ask godinary separated by commas (A comma at the end allows empty referers)
      --cdn_ttl string           Number of seconds images wil be cached in CDN (default "604800")
      --domain string            Domain to validate with Host header, it will deny any other request (if port is not standard must be passed as host:port)
      --fs_base string           FS option: Base dir for filesystem storage
      --gce_project string       GS option: Sentry DSN for error tracking
      --gs_bucket string         GS option: Bucket name
      --gs_credentials string    GS option: Path to service account file with Google Storage credentials
      --max_request int          Maximum number of simultaneous downloads (default 100)
      --max_request_domain int   Maximum number of simultaneous downloads per domain (default 10)
      --port string              Port where the https server listen (default "3002")
      --release string           Release hash to notify sentry
      --sentry_url string        Sentry DSN for error tracking
      --ssl_dir string           Path to directory with server.key and server.pem SSL files (default "/app/")
      --storage string           Storage type: 'gs' for google storage or 'fs' for filesystem (default "fs")
```


### Use it
```
http://localhost:3002/image/fetch/w_500/https%3A%2F%2Fphotos.roomorama-cache.com%2Fphotos%2Frooms%2F3001686%2F3001686_gallery.jpg
```

Parameters:
- type fetch -> last param is target URL
- w: max width
- h: max height
- c: crop type (scale, fit and limit allowed)
- f: format (jpg, jpeg, png, gif, webp and auto allowed)
- q: quality (75 by default)

### TODO
- rate limiting
- log & better error handling
