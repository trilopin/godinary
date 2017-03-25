
get-deps:
	go get github.com/disintegration/imaging
	go get github.com/stretchr/testify/assert
	go get github.com/chai2010/webp

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main
	docker build -t godinary .

test:
	cd imagejob && go test -v .

run:
	docker run -p 3002:3002 -ti godinary
