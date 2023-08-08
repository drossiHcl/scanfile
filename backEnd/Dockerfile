# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Commands:
# go build
# docker build -t my-backend-test .
# docker run -d --net fsscan-net -it -p 50051:50051 --name mounttest --mount type=bind,source=/home/drossi/myTest,target=/home/drossi/myTest my-backend-test

# Start from the latest golang base image
# FROM golang:latest as builder
FROM ubuntu:20.04

ARG GRPC_SERVER_PORT

# Add Maintainer Info
LABEL maintainer="Dan"
#ENV GOPROXY=https://goproxy.io

# Set the Current Working Directory inside the container
WORKDIR /home
RUN mkdir drossi
WORKDIR /home/drossi
RUN mkdir myTest
WORKDIR /home/drossi/myTest
# RUN mkdir data
RUN ls
RUN pwd
RUN mkdir /app
WORKDIR /app

# Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -o backend .
# or
COPY ./backEnd/backEnd .
#COPY /home/drossi/myTest/data/local.env .

# EXPOSE 50051
EXPOSE $GRPC_SERVER_PORT

# Command to run the executable
CMD ["./backEnd"]
