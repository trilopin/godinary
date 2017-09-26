
get-deps:
	glide install

build:
	docker build -t godinary:latest .

local-test:
	cd imagejob && go test -v .
	cd storage && go test -v .

local-certs:
	openssl genrsa -out server.key 2048 && openssl ecparam -genkey -name secp384r1 -out server.key && openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650 -subj /C=US/ST=City/L=City/O=company/OU=SSLServers/CN=localhost/emailAddress=me@example.com

test-docker-image:
	docker run -p 3002:3002 --env-file .env --entrypoint sh -ti godinary:latest

run:
	docker run -p 3002:3002 --env-file .env \
	       -v $$PWD/data/:/data/ \
		   -v $$PWD/server.key:/app/server.key \
		   -v $$PWD/server.pem:/app/server.pem \
		   -ti godinary:latest
