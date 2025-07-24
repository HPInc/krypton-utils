DOCKER_IMAGE_NAME=krypton-local-sqs

ifndef DOCKER_FILE
  DOCKER_FILE=Dockerfile
endif

ifndef TARGET_REPO
  TARGET_REPO=ghcr.io/hpinc/krypton
endif

build: docker-image
docker-image:
	docker build -t $(DOCKER_IMAGE_NAME) -f $(DOCKER_FILE) .

clean:
	-docker rmi -f $(DOCKER_IMAGE_NAME)
	-docker rmi -f $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

publish: push

tag:
	docker tag $(DOCKER_IMAGE_NAME) $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

push: tag
	docker push $(TARGET_REPO)/$(DOCKER_IMAGE_NAME)

.SILENT:
.PHONY: all docker-image tag push clean publish