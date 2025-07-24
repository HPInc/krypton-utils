GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILT_ON := $(shell hostname)
BUILD_DATE := $(shell date +%FT%T%z)

BINARY_NAME=jwtserver
BINARY_DIR=bin

## TARGET_REPO can be overridden to push to another
## eg: TARGET_REPO=docker.io/krypton-images make publish for eg.
ifndef TARGET_REPO
  TARGET_REPO=ghcr.io/hpinc/krypton
endif

ifndef DOCKER_FILE
  DOCKER_FILE=Dockerfile
endif

pull:

clean-base:
	-docker rmi -f $(DOCKER_IMAGE_NAME)
	-docker rmi -f $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

publish: push

tag:
	docker tag $(DOCKER_IMAGE_NAME) $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

push: tag
	docker push $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

tidy:
	go mod tidy

vet:
	go vet ./...

.SILENT:
.PHONY: all docker-image tag push clean-base publish
