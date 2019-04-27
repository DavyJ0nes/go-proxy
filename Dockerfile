# -- Builder Image
FROM golang:1.12.4-alpine3.9 As Builder

ARG projectPath=github.com/davyj0nes/go-proxy

RUN apk --no-cache add ca-certificates

COPY . /go/src/$projectPath
WORKDIR /go/src/$projectPath

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo --installsuffix netgo -o proxy .

# -- Main Image
FROM alpine:3.9

LABEL MAINTAINER="davyj0nes <davyrogerjones@gmail.com>"

RUN adduser -D -s /bin/sh app

# Copy SSL Certs
COPY --from=Builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder image
COPY --from=Builder /go/src/github.com/davyj0nes/go-proxy /bin/proxy
RUN chmod a+x /bin/proxy

# Ensure not runnning as root user
USER app

ENTRYPOINT ["proxy"]