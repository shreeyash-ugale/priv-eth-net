# Forge

Spin up a tiny private Ethereum network, watch it come alive, and keep an eye on it with metrics. Forge gives you a simple Go helper to query node health, wire peers together, and export Prometheus metrics—paired with a Docker setup for three Geth nodes, Prometheus, and Grafana.

---

## What you get

- Three Geth nodes (PoS dev mode on node 1, followers on 2 and 3) via Docker
- A small Go CLI that can:
  - Show block height, peer count, and (simulated) mining status
  - Connect nodes as peers automatically (admin_addPeer)
  - Expose custom Prometheus metrics on a single endpoint
- Prometheus scraping both Geth and the custom exporter
- Grafana ready to dashboard your local chain

---

## Project layout

```
.
├─ main.go                # CLI entrypoint (status, connect peers, exporter)
├─ metrics/               # Custom Prometheus collector and HTTP server
├─ network/               # Node + PeerManager (RPC calls, peer connect)
├─ configs/
│  ├─ docker-compose.yml  # 3× Geth + Prometheus + Grafana
│  ├─ prometheus.yml      # Scrapes Geth and the custom exporter
│  ├─ genesis.json        # Chain genesis
│  ├─ config.yaml         # Beacon-like params for private net
│  └─ scripts/            # Geth startup scripts (PoS/dev + follower)
└─ ...
```

---

## Prerequisites

- Go 1.25+
- Docker (with Compose plugin)
- Windows, macOS, or Linux (commands below show PowerShell where it matters)

---

## Quick start (Docker)

1) From the repo root, start the stack:

```powershell
# Option A: run compose from the configs folder
cd configs
docker compose up -d

# Option B: run from repo root and point to the file
cd ..
docker compose -f configs/docker-compose.yml up -d
```

What this brings up:
- geth-node1: dev mode, produces a block every 5s
- geth-node2, geth-node3: follower nodes
- prometheus: scrapes each node on 6060 and the custom exporter on 9545
- grafana: default port from your env (see below)

2) Health check

```powershell
# Wait for all containers to be healthy
docker ps
```

3) Open the UIs

- Prometheus: http://localhost:9090 (or your mapped port)
- Grafana: http://localhost:3000 (or the port you set)

If Grafana asks for admin credentials, set them via env (see “Environment & ports”).

### Environment & ports

The compose file expects a few environment variables. If you don’t have an `.env`, create one next to `configs/docker-compose.yml` (or export them in your shell):

```dotenv
# Chain
NETWORK_ID=1337
HTTP_API=eth,net,web3,admin
HTTP_CORS_DOMAIN=*

# RPC ports (host:container). These are the RPC ports the Go CLI talks to.
NODE1_RPC_PORT=8545
NODE2_RPC_PORT=8546
NODE3_RPC_PORT=8547

# Metrics ports
# Prometheus (inside Docker) scrapes each Geth at geth-nodeX:6060.
# That means the process inside the container should listen on 6060.
# If you hit host-port conflicts, you can remove the metrics port mappings from compose;
# Prometheus doesn’t need them exposed on the host.
NODE1_METRICS_PORT=6060
NODE2_METRICS_PORT=6060
NODE3_METRICS_PORT=6060

# Prometheus & Grafana
PROMETHEUS_PORT=9090
GRAFANA_PORT=3000
GF_SECURITY_ADMIN_USER=admin
GF_SECURITY_ADMIN_PASSWORD=admin
```

Notes:
- The provided `prometheus.yml` scrapes `geth-node1:6060`, `geth-node2:6060`, `geth-node3:6060`, and the custom exporter at `host.docker.internal:9545`.
- If you want to expose Geth metrics on different host ports, you can, but keep the internal container port at 6060 for Prometheus, or adjust `prometheus.yml` accordingly.

---

## Using the CLI

Forge is a tiny Go program that talks to your nodes over HTTP/RPC. Defaults target localhost ports `8545`, `8546`, and `8547`.

Build it:

```powershell
go build
```

Run a quick status check:

```powershell
# From repo root
./priv-eth-net.exe -status
```

Watch the network (updates every 5s):

```powershell
./priv-eth-net.exe -watch
```

Connect peers (uses admin_addPeer and rewrites enode hosts to the container names so Docker DNS works):

```powershell
./priv-eth-net.exe -connect
```

Start the custom Prometheus exporter (default port 9545):

```powershell
./priv-eth-net.exe -exporter -port 9545
# Metrics will be at http://localhost:9545/metrics
```

CLI flags summary:

- `-status`   Show block height, peer count, mining indicator
- `-watch`    Keep printing status every 5 seconds
- `-connect`  Make each node peer with the others (admin_addPeer)
- `-exporter` Run Prometheus HTTP server
- `-port`     Exporter port (default 9545)

---

## How it works (short tour)

- `network.Node` wraps an RPC client and provides:
  - `GetBlockNumber`, `GetPeerCount` (via `net_peerCount`), `IsMining`
  - `GetNodeInfo` to fetch the enode and `admin_addPeer` to connect peers
- `network.PeerManager` holds a slice of nodes and:
  - Collects each node’s enode, rewrites the host to `geth-nodeN` for Docker DNS
  - Adds every node as a peer to every other node
  - Prints a clear status view
- `metrics.Collector` exports three gauges per node:
  - `ethereum_block_height` (height)
  - `ethereum_peer_count` (peers)
  - `ethereum_is_mining` (1/0; for dev PoS we infer via block production)
- `metrics.Server` is a minimal HTTP server that registers the collector and exposes `/metrics` and `/health`.

On the Docker side:
- `configs/scripts/start-geth-pos.sh` runs node1 in `--dev` mode with `--dev.period 5` (instant blocks), and nodes 2/3 as followers.
- `configs/prometheus.yml` scrapes Geth metrics from each container and the custom exporter on your host.

---

## Security and housekeeping

- Private keys and passwords are intentionally ignored by Git and excluded from Docker build contexts (`.gitignore` and `.dockerignore` are in place).
- Don’t commit contents of `configs/keystore*` or `configs/password.txt`.
- If you accidentally committed secrets in the past, rotate them and use `git rm --cached` plus a history rewrite tool (`git filter-repo`).

---

## Troubleshooting

- “Connection refused” from the CLI: make sure the containers are up and that RPC ports match what the CLI expects (`8545`, `8546`, `8547` by default). You can change the defaults in `main.go` or run the containers with those ports mapped.
- Prometheus can’t scrape node2/node3 metrics: ensure each Geth process listens on `6060` inside the container (that’s what `prometheus.yml` expects), or edit `prometheus.yml` to match whatever you configured.
- Grafana login: set `GF_SECURITY_ADMIN_USER` and `GF_SECURITY_ADMIN_PASSWORD` in your environment.

---

## Contributing

Open an issue or a PR if you spot something off or want to add a small feature. If you’re changing the Docker wiring, please note how it affects Prometheus scraping and the CLI defaults.

---

## License

MIT
