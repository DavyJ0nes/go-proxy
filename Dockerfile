# --  SET ARGS
ARG ORG_NAME=davyj0nes
ARG REPO_NAME=go-proxy
ARG MAINTAINER_NAME=DavyJ0nes
ARG APP_NAME=proxy
ARG MAIN_PATH=main.go

# -- BUILDER IMAGE
FROM golang:1.12.4-alpine3.9 As Builder

ARG ORG_NAME
ARG REPO_NAME
ARG MAINTAINER_NAME
ARG APP_NAME
ARG MAIN_PATH

ENV GO111MODULE=on

# install git
RUN apk --no-cache add git

WORKDIR /go/src/github.com/${ORG_NAME}/${REPO_NAME}

# set up dependencies
 COPY ./go.mod go.mod
 COPY ./go.sum go.sum
 RUN go mod vendor

# copy rest of the package code
COPY . .

# build the statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo --installsuffix netgo -o $APP_NAME $MAIN_PATH

# -- MAIN IMAGE
FROM alpine:3.9

ARG ORG_NAME
ARG REPO_NAME
ARG MAINTAINER_NAME
ARG APP_NAME
ARG MAIN_PATH

LABEL MAINTAINER=${MAINTAINER_NAME}
LABEL OWNER=${MAINTAINER_NAME}

# copy binary from builder image
COPY --from=Builder /go/src/github.com/${ORG_NAME}/${REPO_NAME}/${APP_NAME} /bin/proxy
RUN chmod +x /bin/proxy

# ensure not runnning as root user
RUN adduser -D -s /bin/sh dockmaster
USER dockmaster

EXPOSE 8081
ENTRYPOINT ["proxy"]
