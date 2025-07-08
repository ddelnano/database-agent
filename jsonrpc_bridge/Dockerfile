# Build stage
FROM golang:1.23.6-alpine AS builder

# Set working directory
WORKDIR /build

# Copy source code
COPY main.go ./

# Initialize go module
RUN go mod init jsonrpc_bridge

# Build the binary with CGO disabled
RUN CGO_ENABLED=0 GOOS=linux go build -a -o jsonrpc_bridge main.go

# Runtime stage
FROM alpine:latest

# Install required packages
RUN apk add --no-cache \
    bash \
    curl \
    docker \
    ca-certificates

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/

# Install dagger
RUN curl -L https://dl.dagger.io/dagger/install.sh | sh \
    && mv bin/dagger /usr/local/bin/

# Create app directory
RUN mkdir -p /app

# Copy the wrapper script
COPY mcp-tcp-wrapper.sh /app/mcp-tcp-wrapper.sh
RUN chmod +x /app/mcp-tcp-wrapper.sh

# Copy the compiled binary from builder stage
COPY --from=builder /build/jsonrpc_bridge /app/jsonrpc_bridge
RUN chmod +x /app/jsonrpc_bridge

WORKDIR /app

CMD ["/app/mcp-tcp-wrapper.sh"]
