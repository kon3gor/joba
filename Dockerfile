# Use an official Go image as the build stage
FROM golang:1.23 as builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o joba cmd/main.go

# Use a lightweight image for the final container
FROM alpine:latest

# Copy the Go binary from the builder stage
COPY --from=builder /app/joba /joba
COPY --from=builder /app/config.yaml /config.yaml

# Command to run the Go binary
ENTRYPOINT ["/joba"]
