#!/usr/bin/env bash

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
export REPO=kasvith/kache
export TAG=`if [ "$TRAVIS_BRANCH" == "master" ]; then echo "latest"; else echo $TRAVIS_TAG ; fi`
echo Tagging $REPO with $TAG
docker build -f docker/Dockerfile -t $REPO:$TAG .
echo Pushing image to Docker Hub
docker push $REPO
echo Done