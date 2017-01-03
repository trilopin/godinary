# https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/
FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD main /
CMD ["/main"]
