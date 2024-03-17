# This Dockerfile is used to build the container image for the msvc-orders microservice.
# It starts with the Alpine Linux base image, installs the necessary ca-certificates,
# copies the compiled application binary from the builder stage to the /app directory,
# sets the command to run the application, and exposes port 8000 for incoming connections.
# Build stage
FROM golang:alpine AS builder

# Install git
RUN apk add --no-cache git
WORKDIR /tmp
RUN git init
ENV GOPRIVATE=github.com/SOAT1StackGoLang/
ENV GO111MODULE=on 
ARG GITHUB_ACCESS_TOKEN
RUN git config --global url."https://$GITHUB_ACCESS_TOKEN:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# Set working directory
WORKDIR /go/src/app

ADD ./ .
RUN ls -alth
RUN go get -d -v ./...
## Run Swag
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go install github.com/swaggo/swag/cmd/swag
RUN swag init -g helpers.go -o ./docs/ -d ./internal/transport/routes -ot go
# Build the application
RUN go build -o /go/bin/app -v cmd/server/*.go
#RUN ls -alth cmd/migrations/files
RUN go build -o /go/bin/migs -v cmd/migrations/*.go

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/bin/migs /migs
CMD /migs; /app
EXPOSE 8000