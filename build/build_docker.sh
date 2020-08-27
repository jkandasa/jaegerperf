#!/bin/bash

# docker registry
DOCKER_ORG='quay.io/jkandasa'
DOCKER_REPO="${DOCKER_ORG}/jaegerperf"

# alpine golang builder image tag
GOLANG_BUILDER_TAG="1.15.0-alpine3.12"

# tag version
TAG='1.3'

# build go project
# go build ../main.go
podman run --rm -v \
    "$PWD"/..:/usr/src/jaegerperf -w /usr/src/jaegerperf \
    golang:${GOLANG_BUILDER_TAG} \
    go build -v -o build/jaegerperf cmd/main.go

# change permission
chmod +x ./jaegerperf

# copy default templates
cp ../resources ./resources -r

# copy UI files
cp ../web/build ./web -r

# build image
podman build -t ${DOCKER_REPO}:${TAG} .

# push image to registry
podman push ${DOCKER_REPO}:${TAG}

# remove application bin file and UI files
rm ./jaegerperf -rf
rm ./resources -rf
rm ./web -rf
