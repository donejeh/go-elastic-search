# Use Go official image
FROM golang:1.24-alpine

# Set working directory inside the container
WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the app
RUN go build -o app .

# Expose the port
EXPOSE 8080

# Run the app
CMD ["./app"]
