# FROM ubuntu:22.04
FROM golang:1.19.3-alpine AS base

ARG EXEC_TYPE=santaclaus
ENV EXEC_TYPE=${EXEC_TYPE}
## EXEC_TYPE can be either "santaclaus" or "unit_tests"

# Requirements
RUN apk add --update --no-cache \ 
        make \ 
        build-base

WORKDIR /app

# Copy sources
COPY third_parties /app/third_parties
COPY src /app/src
COPY Makefile /app/
COPY go.mod go.sum /app/

# Build
RUN go mod download && \
    go mod verify

RUN make ${EXEC_TYPE}

# Run
CMD ./${EXEC_TYPE}