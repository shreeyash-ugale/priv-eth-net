#!/bin/sh
set -e

# Copy keystore if provided  
if [ -d /keystore ] && [ "$(ls -A /keystore)" ]; then
    echo "Copying keystore files..."
    mkdir -p /root/.ethereum/keystore
    cp -n /keystore/* /root/.ethereum/keystore/ 2>/dev/null || true
fi

# Initialize genesis if not already done
if [ ! -f /root/.ethereum/geth/chaindata/LOCK ]; then
    echo "Initializing genesis block..."
    geth init --datadir /root/.ethereum /genesis.json
fi

if [ "$ENABLE_MINING" = "true" ]; then
    echo "Starting Geth with Clique sealing (Geth 1.13.x)..."
    exec geth --datadir /root/.ethereum \
        --networkid "$NETWORK_ID" \
        --http --http.addr 0.0.0.0 --http.port "$RPC_PORT" \
        --http.api "$HTTP_API" --http.corsdomain "$HTTP_CORS_DOMAIN" \
        --metrics --metrics.addr 0.0.0.0 --metrics.port "$METRICS_PORT" \
        --nodiscover \
        --mine \
        --miner.etherbase "$ETHERBASE_ADDRESS" \
        --unlock "$ETHERBASE_ADDRESS" \
        --password /password.txt \
        --allow-insecure-unlock
else
    echo "Starting Geth node (non-sealing)..."
    exec geth --datadir /root/.ethereum \
        --networkid "$NETWORK_ID" \
        --http --http.addr 0.0.0.0 --http.port "$RPC_PORT" \
        --http.api "$HTTP_API" --http.corsdomain "$HTTP_CORS_DOMAIN" \
        --metrics --metrics.addr 0.0.0.0 --metrics.port "$METRICS_PORT" \
        --nodiscover
fi
