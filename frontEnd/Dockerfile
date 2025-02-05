# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Commands:
# go build
# docker build -t my-frontend-test .
# or
# docker build -t my-frontend-test -f frontEnd/Dockerfile .

# docker run -d --net fsscan-net -it -p 50051:50051 --name mounttest --mount type=bind,source=/home/drossi/myTest,target=/home/drossi/myTest my-frontend-test
# or
# deploy with k8s:
# envsubst < my-frontend-dep.yaml | kubectl apply -f -


# Start from the latest golang base image
# FROM golang:latest as builder
FROM ubuntu:20.04
FROM golang:1.18

ARG HTTP_FRONTEND_PORT

# Add Maintainer Info
LABEL maintainer="Dan"

# Set the Current Working Directory inside the container
RUN mkdir /app
WORKDIR /app
RUN mkdir templates
COPY ./frontEnd/templates/* templates/
RUN mkdir proto
COPY proto/* proto/

# Copy the source from the current directory to the Working Directory inside the container COPY ./frontEnd/frontEnd .
# or Build the Go app
RUN apt-get update && apt-get install -y apt-transport-https
RUN apt-get -y update
COPY go.mod .
COPY go.sum .
# *** Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# *** Copy sources
COPY frontEnd/*.go .

# *** Build the image
RUN go get scanfile.com/scanfile/proto

# Build the image
#RUN go build -o ./fsScan
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./frontEnd .

# EXPOSE 8082
EXPOSE $HTTP_FRONTEND_PORT

# Command to run the executable
CMD ["./frontEnd"]