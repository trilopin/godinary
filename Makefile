
get-deps:
	glide install

build-test:
	docker build --build-arg RUNTESTS=1 -t godinary:latest .

build:
	docker build -t godinary:latest .

build-dev:
	docker build -t godinary:dev -f Dockerfile.dev .

test:
	go test --cover github.com/trilopin/godinary/http github.com/trilopin/godinary/storage github.com/trilopin/godinary/image

local-certs:
	openssl genrsa -out server.key 2048 && openssl ecparam -genkey -name secp384r1 -out server.key && openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650 -subj /C=US/ST=City/L=City/O=company/OU=SSLServers/CN=localhost/emailAddress=me@example.com

test-docker-image:
	docker run -p 3002:3002 --env-file .env --entrypoint sh -ti godinary:latest

run:
	docker run --rm -p 3000:3000 --env-file .env \
	       -v $$PWD/:/go/src/github.com/trilopin/godinary/ \
		   -ti godinary:dev

sh-dev:
	docker run --rm -p 3000:3000 --env-file .env \
	       -v $$PWD/:/go/src/github.com/trilopin/godinary/ \
		   -ti godinary:dev bash
