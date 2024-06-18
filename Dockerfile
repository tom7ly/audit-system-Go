# syntax=docker/dockerfile:1

# Stage 1: Build the application and generate the Ent schema
FROM golang:1.22.4 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Generate Ent schema
RUN go run entgo.io/ent/cmd/ent generate ./ent/schema

# Build the Go app
RUN go build -o main ./cmd/main.go

# Stage 2: Run the application
FROM golang:1.22.4

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/main .

# Copy the ent directory from the builder stage (if needed)
COPY --from=builder /app/ent ./ent

# Copy wait-for-it.sh script
COPY wait-for-it.sh /app/wait-for-it.sh

# Make wait-for-it.sh executable
RUN chmod +x /app/wait-for-it.sh

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./wait-for-it.sh", "db:5432", "--", "tail", "-f", "/dev/null"]
