#!/bin/bash

set -o allexport
source ./env/unit_tests.env
source ./env/local_unit_tests.env
set +o allexport

make gRPC
make unit_tests

RETURN_VAL=0

function stop_docker() {
  	if [ "$1" == "--docker" ] ; then
		docker compose --env-file ./env/unit_tests.env stop
  	fi
}

if [ "$1" == "--docker" ] ; then
	docker exec -t maestro-mongo-1 mongosh -u $MONGO_USERNAME -p $MONGO_PASSWORD --eval "use santaclaus" --eval "db.dropDatabase()" ## remove database santaclaus to start on a blank canevas for tests
	docker compose --env-file ./env/unit_tests.env --profile launch down --volumes
	trap "echo \"Stopping docker container...\"; stop_docker $1; sleep 3; exit" SIGINT
	docker compose --env-file ./env/unit_tests.env --profile launch up --build
else
	if [ -x ./unit_tests ] ; then
		./unit_tests
	else
		echo "Unit test executable has not been found"
		RETURN_VAL=1
	fi
fi

stop_docker $1
exit $RETURN_VAL
