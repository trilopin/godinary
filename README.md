# godinary
Image proxy with live resize &amp; tranformations


Install
```
go get github.com/trilopin/godinary
```

Start server
```
godinary
```


Docker
```
go build -o main .
docker build -t godinary .
```

Configuration (via env vars)
```
- GODINARY_MAX_REQUEST: number of concurrent external requests (default 20)
- GODINARY_PORT: http server port (default 3002)
```


Use it
```
http://localhost:3002/v0.1/fetch/w_500/https%3A%2F%2Fphotos.roomorama-cache.com%2Fphotos%2Frooms%2F3001686%2F3001686_gallery.jpg
```

Parameters:
- type fetch -> last param is target URL
- w: max width
- h: max height
- c_limit: preserve height/width ratio
- f: format (jpg, jpeg, png, gif allowed)

TODO:
- remove julienschmidt/httprouter dependency: need to solve weird double slash replacement in default go http mux
- reduce/optimize resulting images
- concurrency: global semaphore included, a semaphore per domain should be great
- log & better error handling
