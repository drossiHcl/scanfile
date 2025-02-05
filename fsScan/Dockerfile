# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Commands:
# go build
# docker build -t my-fsscan-test .
# or
# docker build -t my-fsscan-test -f fsScan/Dockerfile .


# docker run -d --net fsscan-net -it -p 50051:50051 --name mounttest --mount type=bind,source=/home/drossi/myTest,target=/home/drossi/myTest my-myfsscan-test
# or
# deploy with k8s:
# envsubst < my-fsscan-dep-copilot.yaml | kubectl apply -f -


# Start from the latest golang base image
# FROM golang:latest as builder
#FROM ubuntu:20.04
FROM golang:1.18
# Add Maintainer Info
LABEL maintainer="Dan"

# Set the Current Working Directory inside the container
RUN mkdir /app
WORKDIR /app
RUN mkdir myData
WORKDIR /app/myData
RUN mkdir scanfile
WORKDIR /app/myData/scanfile
RUN mkdir log

WORKDIR /app
RUN mkdir data
COPY data/* data

RUN mkdir util_pkg
COPY util_pkg/* util_pkg
RUN mkdir types_pkg
COPY types_pkg/* types_pkg
RUN mkdir proto
COPY proto/* proto/


RUN apt-get update && apt-get install -y apt-transport-https
RUN apt-get -y update
COPY go.mod .
COPY go.sum .
# *** Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# dependencies required for code.sajari.com/docconv
RUN apt-get install -y poppler-utils wv unrtf tidy
RUN go get github.com/JalfResi/justext
RUN go get code.sajari.com/docconv

# *** Copy sources
COPY fsScan/*.go .
# COPY .env .
# COPY fsScan/pdftotext /usr/bin/

# *** Build the image
RUN go get scanfile.com/scanfile/proto
RUN go get scanfile.com/scanfile/types_pkg
RUN go get scanfile.com/scanfile/util_pkg

# Build the image
#RUN go build -o ./fsScan
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./fsScan .

# Command to run the executable
CMD ["./fsScan"] 
