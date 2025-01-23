# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with static linking
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /xrd2crd ./cmd/xrd2crd

# Final stage using distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the binary from builder
COPY --from=builder /xrd2crd /usr/local/bin/xrd2crd

# Use nonroot user
USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/xrd2crd"] 