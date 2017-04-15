
get-deps:
	go get gopkg.in/h2non/bimg.v1
	go get github.com/stretchr/testify/assert

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main
	docker build -t godinary .

test:
	cd imagejob && go test -v .

run:
	docker run -p 3002:3002 --env-file .env -ti godinary

local:
	go build -a -v -o main
	GODINARY_FS_BASE=/Users/jpeso/godinary_data/ ./main
