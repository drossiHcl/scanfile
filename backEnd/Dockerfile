# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Commands:
# go build
# docker build -t my-backend-test .
# docker run -d --net fsscan-net -it -p 50051:50051 --name mounttest --mount type=bind,source=/home/drossi/myTest,target=/home/drossi/myTest my-backend-test

# Start from the latest golang base image
# FROM golang:latest as builder FROM ubuntu:20.04
FROM golang:1.18

ARG GRPC_SERVER_PORT

# Add Maintainer Info
LABEL maintainer="Dan"
#ENV GOPROXY=https://goproxy.io

# Set the Current Working Directory inside the container
WORKDIR /home
RUN mkdir daniele
WORKDIR /home/daniele
RUN mkdir Daniele
WORKDIR /home/daniele/Daniele
RUN mkdir myData
WORKDIR /home/daniele/Daniele/myData
RUN mkdir scanfile
WORKDIR /home/daniele/Daniele
RUN mkdir scanfile
WORKDIR /home/daniele/Daniele/scanfile
RUN mkdir data
RUN mkdir log
RUN ls
RUN pwd
RUN mkdir /app
WORKDIR /app

RUN mkdir util_pkg
COPY util_pkg/* util_pkg
RUN mkdir types_pkg
COPY types_pkg/* types_pkg
RUN mkdir proto
COPY proto/* proto/

# Set the GOPROXY environment variable
#RUN export GOPROXY=https://goproxy.io,direct
#ENV GOPROXY=direct
#ENV GOPROXY=https://goproxy.io,direct
# *** Copy go.mod and go.sum
#COPY /home/drossi/myTest/provaGo/scanfile/go.mod .
#COPY /home/drossi/myTest/provaGo/scanfile/go.sum .
# starting from upper dir: docker build -t my-fsscan-test -f fsScan/Dockerfile .

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
COPY backEnd/*.go .
# COPY .env .
# COPY fsScan/pdftotext /usr/bin/

# *** Build the image
RUN go get scanfile.com/scanfile/proto
RUN go get scanfile.com/scanfile/types_pkg
RUN go get scanfile.com/scanfile/util_pkg

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o backEnd .
# or COPY ./backEnd/backEnd .
#COPY /home/drossi/myTest/data/local.env .

# EXPOSE 50051
EXPOSE $GRPC_SERVER_PORT

# Command to run the executable
CMD ["./backEnd"]
