#!/bin/bash

set -o allexport
source ./env/unit_tests.env
source ./env/local_unit_tests.env
set +o allexport

make gRPC
make unit_tests

if [ "$1" == "--docker" ] ; then
    docker compose --env-file ./env/unit_tests.env up --build &
    DOCKER_PID=$!
fi


function stop_docker() {
  if [ "$1" == "--docker" ] ; then
    kill $DOCKER_PID ## todo launch docker in background from docker arguments and stop with docker command
  fi
}

trap "echo \"Stopping docker container...\"; stop_docker $1; sleep 3; exit" SIGINT

./unit_tests
