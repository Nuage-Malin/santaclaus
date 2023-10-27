# FROM ubuntu:22.04
FROM golang:1.19.3-alpine

ARG EXEC_TYPE=santaclaus
ENV EXEC_TYPE=${EXEC_TYPE}
## EXEC_TYPE can be either "santaclaus" or "unit_tests"

# Requirements
RUN apk add --update --no-cache \
        make \
        build-base  \
        protobuf-dev

WORKDIR /app

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

COPY go.mod go.sum /app/
RUN go mod download && \
    go mod verify

COPY Makefile /app/
COPY third_parties /app/third_parties
RUN make gRPC

COPY src /app/src
RUN make ${EXEC_TYPE}

# Run
CMD ./${EXEC_TYPE}
