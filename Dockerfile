# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8.3-stretch as builder

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ wget gcc libc6-dev make pkg-config ca-certificates git curl \
		libvips libvips-dev 

# setup go & glide
RUN curl https://glide.sh/get | sh
WORKDIR /go/src/github.com/trilopin/godinary/

# app
ENV SRC_DIR=/go/src/github.com/trilopin/godinary/
ADD . /go/src/github.com/trilopin/godinary/
RUN make get-deps
RUN go build -o godinary 
RUN mkdir /app && cp godinary /app/
 
ENTRYPOINT ["/app/godinary"]
