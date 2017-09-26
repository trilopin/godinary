# godinary
Image proxy with live resize &amp; tranformations


Install
```
git clone https://github.com/trilopin/godinary
```



Docker flow:
- make build -> compiles and build docker image
- make run -> start server

Development flow:
- glide install
- GODINARY_FS_BASE=data GODINARY_ALLOW_HOSTS=<host>, GODINARY_SSL_DIR=./ go run main.go
- GODINARY_GS_CREDENTIALS=<credential>.json GODINARY_ALLOW_HOSTS=<host>, GODINARY_SSL_DIR=./ GODINARY_STORAGE=gs GODINARY_GCE_PROJECT=<gce project> GODINARY_GS_BUCKET=<gce bucket>  go run main.go

Configuration (via env vars)
```
- GODINARY_MAX_REQUEST: number of concurrent external requests (default 20)
- GODINARY_PORT: http server port (default 3002)
- GODINARY_ALLOW_HOSTS: list of referer hostnames allowed (blank is allways allowed)
- GODINARY_STORAGE: gs for google storage and fs for filesystem (default: "fs")
- GODINARY_FS_BASE: base dir for filesystem storage
- GODINARY_SENTRY_URL: sentry dsn for error tracking (default: "")
- GODINARY_RELEASE: commit hash for this release, used with sentry (default: "")
- GODINARY_SSL_DIR: SSL certs directory (default: "/app/")
- GODINARY_GS_PROJECT: GCE project for Google Storage
- GODINARY_GS_BUCKET: Bucket name in Google storage
- GODINARY_GS_CREDENTIALS: service account json file for Google storage
```


Use it
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

TODO:
- rate limiting
- log & better error handling
