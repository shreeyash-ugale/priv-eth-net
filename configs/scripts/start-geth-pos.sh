#!/bin/sh
set -e

# Initialize genesis if not already done
if [ ! -f /root/.ethereum/geth/chaindata/LOCK ]; then
    echo "Initializing genesis block for PoS..."
    geth init --datadir /root/.ethereum /genesis.json
fi

if [ "$ENABLE_MINING" = "true" ]; then
    echo "Starting Geth node 1 in dev mode (PoS simulation with instant mining)..."
    exec geth --datadir /root/.ethereum \
        --http --http.addr 0.0.0.0 --http.port "$RPC_PORT" \
        --http.api "$HTTP_API" --http.corsdomain "$HTTP_CORS_DOMAIN" \
        --ws --ws.addr 0.0.0.0 --ws.port 8546 --ws.api "$HTTP_API" --ws.origins "*" \
        --metrics --metrics.addr 0.0.0.0 --metrics.port "$METRICS_PORT" \
        --authrpc.addr 0.0.0.0 --authrpc.port 8551 --authrpc.vhosts "*" \
        --nodiscover \
        --dev --dev.period 5 \
        --verbosity 3
else
    echo "Starting Geth follower node in PoS mode..."
    exec geth --datadir /root/.ethereum \
        --networkid "$NETWORK_ID" \
        --http --http.addr 0.0.0.0 --http.port "$RPC_PORT" \
        --http.api "$HTTP_API" --http.corsdomain "$HTTP_CORS_DOMAIN" \
        --ws --ws.addr 0.0.0.0 --ws.port 8546 --ws.api "$HTTP_API" --ws.origins "*" \
        --metrics --metrics.addr 0.0.0.0 --metrics.port "$METRICS_PORT" \
        --authrpc.addr 0.0.0.0 --authrpc.port 8551 --authrpc.vhosts "*" \
        --nodiscover \
        --syncmode full \
        --verbosity 3
fi
