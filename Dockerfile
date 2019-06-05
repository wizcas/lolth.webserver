
#build stage
FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go install -v ./...

# final w/ scratch
FROM scratch
COPY --from=builder /go/bin/server.static /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./app"]
LABEL Name=lolth.server.static Version=0.0.1
EXPOSE 8080