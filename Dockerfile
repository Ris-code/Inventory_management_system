# Use a smaller base image for the final stage
FROM golang:1.21 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /inventory-management

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the executable from the builder stage
COPY --from=builder /inventory-management .

# Copy the .env file into the final image
COPY --from=builder /app/.env ./

# Copy the template directory into the final image
COPY --from=builder /app/templates ./templates

# Copy the static directory into the final image
COPY --from=builder /app/static ./static

# Expose the port the application will run on
EXPOSE 8080

# Run the application
CMD ["./inventory-management"]
