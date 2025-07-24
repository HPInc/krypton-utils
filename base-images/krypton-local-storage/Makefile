DOCKER_IMAGE_NAME=krypton-local-storage
FS_TEST_NAME=fs_storage_test

ifndef DOCKER_FILE
  DOCKER_FILE=Dockerfile
endif

ifndef TARGET_REPO
  TARGET_REPO=ghcr.io/hpinc/krypton
endif

lint:
	shellcheck -x *.sh

build: docker-image
docker-image:
	docker build -t $(DOCKER_IMAGE_NAME) -f $(DOCKER_FILE) .

test: stop
	./run_test.sh

stop:
	-docker stop $(FS_TEST_NAME)

clean:
	-docker rmi -f $(DOCKER_IMAGE_NAME)
	-docker rmi -f $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

publish: push

tag:
	docker tag $(DOCKER_IMAGE_NAME) $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

push: tag
	docker push $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

build_worker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		    go build -ldflags "-s -w" \
		    ./webhookworker.go

.SILENT:
.PHONY: all docker-image tag push clean publish