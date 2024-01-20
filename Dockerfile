# This Dockerfile is used to build the container image for the msvc-payments microservice.
# It starts with the Alpine Linux base image, installs the necessary ca-certificates,
# copies the compiled application binary from the builder stage to the /app directory,
# sets the command to run the application, and exposes port 8000 for incoming connections.
# Build stage
FROM golang:alpine AS builder

# Install git
RUN apk add --no-cache git

# Set working directory
WORKDIR /go/src/app

# Copy source code to working directory
ADD ./ .

# List files in the working directory (debugging)
RUN ls -alth

# Download dependencies
RUN go get -d -v ./...

# Build the application
RUN go build -o /go/bin/app -v cmd/server/*.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
CMD /app
EXPOSE 8000