# Stage 1: Build
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go mod and install deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build the binary
RUN go build -o server ./cmd/api

# Stage 2: Run
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]
