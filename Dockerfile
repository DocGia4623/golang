# Use the official Golang image as the base image
FROM golang:1.23.4-alpine

# Set the environment variables for Go to build a binary compatible with your container's OS and architecture
ENV GOARCH=amd64
ENV GOOS=linux

# Set the Current Working Directory inside the container
WORKDIR ./app

# Copy go mod and sum files
COPY go.mod go.sum ./ 

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN apk update && apk add --no-cache git

# Build the Go app inside the container
RUN go build -o app .

# Ensure the 'app' file is executable
RUN chmod +x ./app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
