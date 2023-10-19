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

# Copy sources
COPY third_parties /app/third_parties
COPY src /app/src
COPY Makefile /app/
COPY go.mod go.sum /app/

# Build
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

RUN go mod download && \
    go mod verify

RUN make gRPC
RUN make ${EXEC_TYPE}

# Run
CMD ./${EXEC_TYPE}
