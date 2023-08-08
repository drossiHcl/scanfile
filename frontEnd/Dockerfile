# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Commands:
# go build
# docker build -t my-fsscan-test .
# docker run -d --net fsscan-net -it --name mounttest1 --mount type=bind,source=/home/drossi/myTest,target=/home/drossi/myTest my-fsscan-test

# Start from the latest golang base image
# FROM golang:latest as builder
FROM ubuntu:20.04
FROM golang:latest

ARG HTTP_FRONTEND_PORT

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
RUN mkdir templates
COPY ./frontEnd/templates/* templates/

# Copy the source from the current directory to the Working Directory inside the container
COPY ./frontEnd/frontEnd .
# or Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -o backend .
#COPY ./go.mod .
#RUN go mod init go.mod
#RUN go mod tidy

# EXPOSE 8082
EXPOSE $HTTP_FRONTEND_PORT

# Command to run the executable
CMD ["./frontEnd"]