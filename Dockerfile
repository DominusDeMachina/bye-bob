# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Install required packages for building
RUN apk add --no-cache git build-base

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Install templ tool for template compilation
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy the source code
COPY . .

# Run templ generate before building to compile templates
RUN templ generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o byebob ./cmd/server

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/byebob .
COPY --from=builder /app/.env .
COPY --from=builder /app/static ./static

# Expose the application port
EXPOSE 3000

# Run the application
CMD ["./byebob"] 