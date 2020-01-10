#!/bin/bash

# docker registry
DOCKER_ORG='quay.io/jkandasa'
DOCKER_REPO="${DOCKER_ORG}/jaegerperf"

# alpine golang builder
GOLANG_BUILDER_TAG="1.13.5-alpine3.11"

# tag version
TAG='1.0'

# build go project
# go build ../main.go
podman run --rm -v \
    "$PWD"/..:/usr/src/jaegerperf -w /usr/src/jaegerperf \
    golang:${GOLANG_BUILDER_TAG} \
    go build -o docker/appbin -v

# change permission
chmod +x ./appbin

# build docker image
podman build -t ${DOCKER_REPO}:${TAG} .

# remove bin file
rm ./appbin -rf

# push new image to docker hub
podman push ${DOCKER_REPO}:${TAG}

# remove build file
rm ./main -rf
