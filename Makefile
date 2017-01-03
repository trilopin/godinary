
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main
	docker build -t godinary .

run:
	docker run -p 3002:3002 -ti godinary
