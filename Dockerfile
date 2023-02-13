# FROM ubuntu:22.04
FROM golang:1.19.3-alpine AS builder

# Requirements
# RUN apt-get update && apt-get install -y \
    # golang                 \
    # build-essential     \
    # && rm -rf /var/lib/apt/lists/*
RUN apk add --update \ 
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
RUN make

# Run
CMD ./santaclaus
