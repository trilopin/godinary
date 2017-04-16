# godinary
Image proxy with live resize &amp; tranformations


Install
```
go get github.com/trilopin/godinary
```



Tooling & Docker
- make build -> compiles and build docker image
- make get-deps -> retrieves dependencies
- make test -> launch tests
- make run -> start server


Configuration (via env vars)
```
- GODINARY_MAX_REQUEST: number of concurrent external requests (default 20)
- GODINARY_PORT: http server port (default 3002)
- GODINARY_ALLOW_HOSTS: list of referer hostnames allowed (blank is allways allowed)
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
