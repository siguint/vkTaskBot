# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
# RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY *.go ./

# Build the Go app
# RUN go mod init github.com/siguint/vktask
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
