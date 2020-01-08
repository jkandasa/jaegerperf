#!/bin/bash

# docker registry
DOCKER_ORG='jkandasa'
DOCKER_REPO="${DOCKER_ORG}/jaegerperf"

# tag version
TAG='1.0'

# build go project
go build ../main.go

# change permission
chmod +x ./main

# build docker image
docker build -t ${DOCKER_REPO}:${TAG} .

# push new image to docker hub
docker push ${DOCKER_REPO}:${TAG}

# remove build file
rm ./main -rf