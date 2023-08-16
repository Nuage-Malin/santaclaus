#!/bin/bash

set -o allexport
source ./env/unit_tests.env
source ./env/local_unit_tests.env
set +o allexport

make gRPC
make unit_tests

function stop_docker() {
  	if [ "$1" == "--docker" ] ; then
		docker compose --env-file ./env/unit_tests.env stop
  	fi
}

if [ "$1" == "--docker" ] ; then
	docker compose --env-file ./env/unit_tests.env --profile launch down --volumes
	trap "echo \"Stopping docker container...\"; stop_docker $1; sleep 3; exit" SIGINT
	docker compose --env-file ./env/unit_tests.env --profile launch up --build
else
	./unit_tests
fi

stop_docker $1