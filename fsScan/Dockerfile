# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Commands:
# go build
# docker build -t my-fsscan-test .
# docker run -d --net fsscan-net -it --name mounttest1 --mount type=bind,source=/home/drossi/myTest,target=/home/drossi/myTest my-fsscan-test

# Start from the latest golang base image
# FROM golang:latest as builder
#FROM ubuntu:20.04
FROM golang:1.18
# Add Maintainer Info
LABEL maintainer="Dan"

# Set the Current Working Directory inside the container
WORKDIR /home
RUN mkdir drossi
WORKDIR /home/drossi
RUN mkdir myTest
WORKDIR /home/drossi/myTest
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
COPY fsScan/*.go .
COPY .env .
# COPY fsScan/pdftotext /usr/bin/

# *** Build the image
RUN go get scanfile.com/scanfile/proto
RUN go get scanfile.com/scanfile/types_pkg
RUN go get scanfile.com/scanfile/util_pkg

# Build the image
#RUN go build -o ./fsScan
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./fsScan .

# Copy the executable from the current directory to the Working Directory inside the container
# COPY fsScan/fsScan .

# Command to run the executable
CMD ["./fsScan"] 
