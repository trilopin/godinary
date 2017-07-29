# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM marcbachmann/libvips:latest
MAINTAINER tomas@aparicio.me


# Go version to use
ENV GOLANG_VERSION 1.8.1

# gcc for cgo
RUN apt-get update && apt-get install -y \
    gcc curl git libc6-dev make ca-certificates \
    --no-install-recommends \
  && rm -rf /var/lib/apt/lists/*

ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz

RUN curl -fsSL --insecure "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

# Fetch the latest version of the package
RUN go get -u github.com/trilopin/godinary

# Run the outyet command by default when the container starts.
ENTRYPOINT ["/go/bin/godinary"]
