DOCKER_IMAGE_NAME=krypton-jwtserver
TEST_SERVER=jwtserver_test
PORT=9090

build: docker-image
docker-image:
	docker build -t $(DOCKER_IMAGE_NAME) -f $(DOCKER_FILE) .

# Build the binaries for the JWT service.
build-binaries:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) \
	-ldflags "-X main.gitCommitHash=$(GIT_COMMIT) -X main.builtAt=$(BUILD_DATE) \
	-X main.builtBy=$(USER) -X main.builtOn=$(BUILT_ON)" \
	-o $(BINARY_DIR)/$(BINARY_NAME) -v

run:
	go run ./...

imports:
	goimports -w .

stop:
	-docker stop $(TEST_SERVER)

test: stop build privateKey.pem
	docker run --rm -itd \
	-p$(PORT):$(PORT) \
	-v$(CURDIR)/privateKey.pem:/go/bin/privateKey.pem \
	--name $(TEST_SERVER) $(DOCKER_IMAGE_NAME)
	
	curl -s --write-out "%{http_code}" localhost:$(PORT)/api/v1/token | grep 200

privateKey.pem:
	openssl genrsa -out privateKey.pem 2048

publicKey.pem:
	openssl rsa -in privateKey.pem -pubout -out publicKey.pem	

# assumes a local private key
local_test: test
	curl -s localhost:$(PORT)/token

clean: clean_local clean-base

clean_local:
	-rm publicKey.pem privateKey.pem
	-docker stop $(TEST_SERVER)

.PHONY: run vet imports stop test local_test clean clean_local

include ./common.mk
