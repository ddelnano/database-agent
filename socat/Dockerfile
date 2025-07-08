FROM alpine:latest

# Install required packages
RUN apk add --no-cache \
    bash \
    curl \
    socat \
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

# TODO(ddelnano): Switch this to a multi-stage build
COPY socat /app/

# Create non-root user
# RUN addgroup -g 1000 dagger && \
#     adduser -D -s /bin/bash -u 1000 -G dagger dagger

# USER dagger

WORKDIR /app

CMD ["/app/mcp-tcp-wrapper.sh"]
