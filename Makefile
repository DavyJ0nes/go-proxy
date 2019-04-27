.DEFAULT_TARGET=help
.PHONY: all
all: help

# VARIABLES
USERNAME = davyj0nes
APP_NAME = proxy

GO_VERSION ?= 1.12.4
GO_PROJECT_PATH ?= github.com/davyj0nes/go-proxy
GO_FILES = $(shell go list ./... | grep -v /vendor/)

APP_PORT = 8080
LOCAL_PORT = 8080

VERSION = 0.1.0
COMMIT = $(shell git rev-parse HEAD | cut -c 1-6)
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

DOCKER_RUN_CMD = docker run -it --rm --name ${APP_NAME} ${USERNAME}/${APP_NAME}:${VERSION} "\$$@"

# COMMANDS

## image: builds a docker image for the application
.PHONY: image
image:
	$(call blue, "# Building Docker Image...")
	@docker build --no-cache --label APP_VERSION=${VERSION} --label BUILT_ON=${BUILD_TIME} --label GIT_HASH=${COMMIT} -t ${USERNAME}/${APP_NAME}:${VERSION} .
	@docker tag ${USERNAME}/${APP_NAME}:${VERSION} ${USERNAME}/${APP_NAME}:latest
	@$(MAKE) clean

## publish: pushes the tagged docker image to docker hub
.PHONY: publish
publish: image
	$(call blue, "# Publishing Docker Image...")
	@docker push docker.io/${USERNAME}/${APP_NAME}:${VERSION}

## run: runs the application locally
.PHONY: run
run:
	$(call blue, "# Running App...")
	@docker run -it --rm -v "$(GOPATH)":/go -v "$(CURDIR)":/go/src/app -p ${LOCAL_PORT}:${APP_PORT} -w /go/src/app golang:${GO_VERSION} go run cmd/main.go

## run_image: builds and runs the docker image locally
.PHONY: run_image
run_image: image
	$(call blue, "# Running Docker Image Locally...")
	@docker run -it --rm --name ${APP_NAME} -p ${LOCAL_PORT}:${APP_PORT} ${USERNAME}/${APP_NAME}:${VERSION}

## test: run test suitde for application
.PHONY: test
test:
	$(call blue, "# Testing Golang Code...")
	@docker run --rm -it -v "$(GOPATH):/go" -v "$(CURDIR)":/go/src/app -w /go/src/app golang:${GO_VERSION} sh -c 'go test -v -race ${GO_FILES}' 

## clean: remove binary from non release directory
.PHONY: clean
clean: 
	@rm -f ${APP_NAME} 

## help: Show this help message
.PHONY: help
help: Makefile
	@echo "${APP_NAME} - v${VERSION}"
	@echo
	@echo " Choose a command run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^## //p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

# FUNCTIONS
define blue
	@tput setaf 4
	@echo $1
	@tput sgr0
endef
