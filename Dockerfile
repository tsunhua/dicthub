# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
# golan dockerhub: https://hub.docker.com/_/golang
FROM golang:1.12.7-alpine3.10  AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum main ./

######## Start a new stage from scratch #######
FROM alpine:3.10

# Solved: https://github.com/aws/aws-sdk-go/issues/2322
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8888
# Export port 9090 for pprof
EXPOSE 9090

# Command to run the executable
CMD ["./main"]
