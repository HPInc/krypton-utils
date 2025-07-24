# Krypton utilities
This repository contains common utilities consumed by various Krypton micro-services. This includes common docker images used:
- To build the source code
- A minimal Alpine Linux based image used to run various microservices
- Minimal base images for dependencies such as PostgreSQL, Redis etc.


## Base images
There is a folder called ```base-images``` at the root of the repository. This folder includes dockerfiles for all base images required for the DSTS service including:
1. ```krypton-go-builder``` - An Alpine Linux based docker image that can be used to build the Krypton Go micro-services. It includes a working Go environment and the protoc compiler and other dependencies required for building and running unit tests for the service.
2. ```krypton-go-base``` - A minimal Alpine Linux docker image that is used to run the Krypton services.
3. ```postgres``` - A docker image for the PostgreSQL server that is used as a database for various Krypton services including the DSTS. This is provided for local development and testing purposes. You can use a managed database service in the cloud or spin up a PostgreSQL instance on a VM in production environments.
4. ```redis``` - A docker image for Redis server which is used by various Krypton services for caching purposes. This is provided for local development and testing purposes. You can use a managed caching service in the cloud or spin up a Redis instance on a VM in production environments.

**NOTE:** These docker images are published to the GHCR (Github Container Registry) docker repository to be consumed by other Krypton services in their respective Github repositories.
