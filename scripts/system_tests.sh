#!/bin/bash

ARG_BUILD_SERVICE=false # --build-service
ARG_BUILD_TESTS=false   # --build-tests
ARG_RUN_SERVICE=false   # --run-service
ARG_RUN_TESTS=false     # --run-tests
ARG_STOP=false          # --stop

CURRENT_FILE_DIR="$(dirname $0)/"

usage()
{
    echo "Usage: $0 [--help] [--build-service] [--build-tests] [--run-service] [--run-tests] [--stop]"
    echo "\t--help: Prints this message"
    echo "\t--build-service: Build santaclaus"
    echo "\t--build-tests: Build system tests"
    echo "\t--run-service: Run santaclaus"
    echo "\t--run-tests: Run system tests"
    echo "\t--stop: Stop service"
    exit 0
}

check_exit_failure()
{
    EXIT_STATUS=$?
    if [ $EXIT_STATUS -ne 0 ]; then
        echo -e "\033[31m$1\033[0m" 1>&2
        exit $EXIT_STATUS
    fi
}

for arg in "$@"; do
    case $arg in
        --help)
            usage
        ;;
        --build-service)
            ARG_BUILD_SERVICE=true
        ;;
        --build-tests)
            ARG_BUILD_TESTS=true
        ;;
        --run-service)
            ARG_RUN_SERVICE=true
        ;;
        --run-tests)
            ARG_RUN_TESTS=true
        ;;
        --stop)
            ARG_STOP=true
        ;;
        *)
            echo "Invalid option: $arg" >&2
            exit 1
        ;;
    esac
done

if $ARG_BUILD_SERVICE; then
    docker compose --profile launch --env-file env/santaclaus.env build
    check_exit_failure "Failed to build santaclaus"
fi

if $ARG_BUILD_TESTS; then
    echo "Not implemented yet" >&2
    # TODO: Build system tests (and check exit status)
fi

if $ARG_RUN_SERVICE; then
    docker compose --profile launch --env-file env/santaclaus.env up
    check_exit_failure "Failed to run santaclaus"
fi

if $ARG_RUN_TESTS; then
    echo "Not implemented yet" >&2
    # set -o allexport
    # source $CURRENT_FILE_DIR/../env/system_tests.env
    # set +o allexport

    # exec ./santaclaus
    # check_exit_failure "System tests failed"
    # TODO: Run system tests (and check exit status) when build will be done
fi

if $ARG_STOP; then
    docker compose --profile launch --env-file env/santaclaus.env down
    check_exit_failure "Failed to stop santaclaus"
fi

exit 0