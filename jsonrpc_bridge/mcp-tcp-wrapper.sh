#!/bin/bash

# TCP wrapper script for persistent dagger MCP server
# This script starts a persistent dagger MCP process and bridges TCP connections to it

set -e

# Default port if not specified
PORT=${MCP_PORT:-8080}
USE_LOCAL_DAGGER=${USE_LOCAL_DAGGER:-false}

# Create named pipes for communication
PIPE_DIR="/tmp/mcp-pipes"
mkdir -p "$PIPE_DIR"
IN_PIPE="$PIPE_DIR/in"
OUT_PIPE="$PIPE_DIR/out"

# Clean up pipes on exit
cleanup() {
    echo "Cleaning up..."
    rm -rf "$PIPE_DIR"
    if [ -n "$DAGGER_PID" ]; then
        kill "$DAGGER_PID" 2>/dev/null || true
    fi
    if [ -n "$SOCAT_PID" ]; then
        kill "$SOCAT_PID" 2>/dev/null || true
    fi
}
trap cleanup EXIT

# Create named pipes
mkfifo "$IN_PIPE" "$OUT_PIPE"

exec 3<>"$IN_PIPE"
exec 4<>"$OUT_PIPE"

if [ "$USE_LOCAL_DAGGER" = false ] && [ -z "$DAGGER_ENGINE_POD_NAME" ]; then
  # Get the dagger engine pod name
  DAGGER_ENGINE_POD_NAME=$(kubectl get pod \
      --selector=name=dagger-dagger-helm-engine --namespace=dagger \
      --output=jsonpath='{.items[0].metadata.name}')
else
  echo "Using DAGGER_ENGINE_POD_NAME: $DAGGER_ENGINE_POD_NAME"
fi

if [ "$USE_LOCAL_DAGGER" = true ]; then
    echo "Using local dagger instance"
    unset "$DAGGER_ENGINE_POD_NAME"
else
    echo "Using remote dagger instance in pod: $DAGGER_ENGINE_POD_NAME"
    export _EXPERIMENTAL_DAGGER_RUNNER_HOST="kube-pod://$DAGGER_ENGINE_POD_NAME?namespace=dagger"
fi

# Start the persistent dagger MCP server
echo "Starting persistent dagger MCP server..."
dagger mcp --allow-llm all -m https://github.com/ddelnano/database-agent < "$IN_PIPE" > "$OUT_PIPE" &
DAGGER_PID=$!

# Give the dagger process a moment to start
sleep 2

# Check if dagger process is still running
if ! kill -0 "$DAGGER_PID" 2>/dev/null; then
    echo "Error: dagger MCP server failed to start"
    exit 1
fi

echo "Dagger MCP server started with PID: $DAGGER_PID"

# Start jsonrpc_bridge to bridge TCP connections to the named pipes
echo "Starting JSON RPC bridge on port $PORT..."
/app/jsonrpc_bridge &
SOCAT_PID=$!

echo "TCP bridge started with PID: $SOCAT_PID"
echo "MCP server is ready on port $PORT"

# Wait for either process to exit
wait
