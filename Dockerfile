# Use the official Golang image as the base image
FROM golang:1.18

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
