DOCKER_IMAGE_NAME=krypton-cli
BIN=bin/krypton-cli.exe
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

ifndef GHCR
  GHCR=ghcr.io/hpinc/krypton
endif

ifndef TARGET
  TARGET=linux
  BIN=bin/krypton-cli
endif

# build info
GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILT_ON := $(shell hostname)
BUILD_DATE := $(shell date +%FT%T%z)

# scheduler protos
SCHEDULER_PROTOS_IMAGE=$(GHCR)/krypton-schedulerprotos
SCHEDULER_PROTOS_DIR=scheduler_protos
SCHEDULER_PROTOS_LOCAL=scheduler_protos_local

ifndef cmd
  cmd=enroll
endif
ifndef count
  count=1
endif
ifndef parallel
  parallel=-parallel
endif
ifndef verbose
  verbose=-verbose
endif
ifndef batch_size
  batch_size=100
endif

build: protos_refresh go_build

run: protos
	go run ./main.go

# extract protos from scheduler docker
protos: $(SCHEDULER_PROTOS_DIR)
$(SCHEDULER_PROTOS_DIR):
	docker pull $(SCHEDULER_PROTOS_IMAGE)
	mkdir -p $(SCHEDULER_PROTOS_DIR) && \
	docker create --name $(SCHEDULER_PROTOS_LOCAL) $(SCHEDULER_PROTOS_IMAGE) "" && \
	docker cp $(SCHEDULER_PROTOS_LOCAL):/protos $(SCHEDULER_PROTOS_DIR) && \
	mv $(SCHEDULER_PROTOS_DIR)/protos/* $(SCHEDULER_PROTOS_DIR) && \
	rm -rf $(SCHEDULER_PROTOS_DIR)/protos && \
	docker rm $(SCHEDULER_PROTOS_LOCAL)

# clean dirs as part of clean target
clean_scheduler_protos:
	-rm -rf $(SCHEDULER_PROTOS_DIR)
	-docker rm $(SCHEDULER_PROTOS_LOCAL)


vet:
	go vet ./...

gosec:
	gosec ./...

imports:
	goimports -w .


go_build:
	CGO_ENABLED=0 GOOS=$(TARGET) GOARCH=amd64 $(GOBUILD) \
	 -ldflags "-s -w -X main.gitCommitHash=$(GIT_COMMIT) -X main.builtAt=$(BUILD_DATE) \
	-X main.builtBy=$(USER) -X main.builtOn=$(BUILT_ON)" \
	-o $(BIN) main.go

install: build
	./install.sh

release: build
	./tools/make_release.sh

tidy:
	go mod tidy

test:
	go test ./...

lint: protos_refresh
	./tools/run_linter.sh

docker: protos_refresh
	docker build -t $(DOCKER_IMAGE_NAME) .

tag:
	docker tag $(DOCKER_IMAGE_NAME) $(GHCR)/$(DOCKER_IMAGE_NAME)

publish: tag
	docker push $(GHCR)/$(DOCKER_IMAGE_NAME)

protos_refresh: clean_scheduler_protos protos

clean: $(GOCLEAN) clean_scheduler_protos

help:
	-@go run ./...
.PHONY: all vet imports build test docker-image tidy help tag publish clean clean_scheduler_protos protos
