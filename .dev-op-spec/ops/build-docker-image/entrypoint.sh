#!/bin/sh

cp -R tcp/swagger .tmp

path_to_working_dir=.dev-op-spec/ops/build-docker-image
cp $path_to_working_dir/Dockerfile .tmp

docker build -t ${DOCKER_REPO_NAME} .tmp