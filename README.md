# godinary
Image proxy with live resize &amp; tranformations


Install
```
go get github.com/trilopin/godinary
```

Start server
```
godinary --port=8080
```


Use it
```
http://localhost:8080/v0.1/accountname/fetch/w_500/https%3A%2F%2Fphotos.roomorama-cache.com%2Fphotos%2Frooms%2F3001686%2F3001686_gallery.jpg
```

Parameters:
- type fetch -> last param is target URL
- w: max width
- h: max height
- c_limit: preserve height/width ratio
- f: format (jpg, jpeg, png, gif allowed)

TODO:
- remove julienschmidt/httprouter dependency 
- reduce/optimize resulting images
- concurrency :)
- dockerify
- log & better error handling
