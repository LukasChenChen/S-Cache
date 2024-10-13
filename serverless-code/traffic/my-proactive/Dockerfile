FROM golang:1.19-bullseye as builder
MAINTAINER Chen Chen <lukas.c.chen@gmail.com>

WORKDIR /app
COPY *.go ./ 
COPY go.mod ./
#COPY ./config/* ./config/
COPY ./config/* .
RUN CGO_ENABLED=0 go build -o scheduler
RUN chmod a+x scheduler
CMD ["./scheduler"]

