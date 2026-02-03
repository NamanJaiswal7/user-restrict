FROM golang:1.23-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy source code FIRST so go mod tidy can see imports
COPY . .

# Create dummy go.sum if it doesn't exist yet
RUN touch go.sum

# Download and tidy dependencies
RUN go mod tidy
RUN go mod download

# Build
RUN go build -o /user-restrict cmd/server/main.go

# Expose port
EXPOSE 8080

# Run
CMD ["/user-restrict"]
