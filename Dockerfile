# Build stage
FROM golang:1.21 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the CA certificate
COPY ca.pem ./

# Copy the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /inventory-management

# Final stage
FROM alpine:latest

# Install necessary dependencies (optional, if required)
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

# Set the working directory
WORKDIR /app

# Copy the executable from the builder stage
COPY --from=builder /inventory-management .

# Copy the .env file into the image
COPY --from=builder /app/.env ./ 

# Copy the CA certificate into the image (optional)
COPY --from=builder /app/ca.pem ./ 

# Copy the templates directory into the image
COPY --from=builder /app/templates ./templates

# Copy the static directory into the image
COPY --from=builder /app/static ./static

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./inventory-management"]
