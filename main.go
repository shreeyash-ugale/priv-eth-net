package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/shreeyash-ugale/priv-eth-net/metrics"
	"github.com/shreeyash-ugale/priv-eth-net/network"
)

func main() {
	connect := flag.Bool("connect", false, "Connect all peers")
	status := flag.Bool("status", false, "Show network status")
	exporter := flag.Bool("exporter", false, "Start Prometheus exporter")
	watch := flag.Bool("watch", false, "Watch network status (updates every 5s)")
	port := flag.Int("port", 9545, "Metrics exporter port")
	flag.Parse()

	// Use 127.0.0.1 for Windows IPv4 compatibility
	nodeURLs := []string{
		"http://127.0.0.1:8545",
		"http://127.0.0.1:8546",
		"http://127.0.0.1:8547",
	}

	fmt.Println("🔗 Connecting to Ethereum nodes...")
	manager, err := network.NewPeerManager(nodeURLs)
	if err != nil {
		log.Fatalf("❌ Failed to create peer manager: %v", err)
	}
	defer manager.Close()

	fmt.Printf("✅ Connected to %d nodes\n", len(manager.Nodes))

	ctx := context.Background()

	if *connect {
		if err := manager.ConnectPeers(ctx); err != nil {
			log.Fatalf("❌ Failed to connect peers: %v", err)
		}
		time.Sleep(2 * time.Second)
		manager.GetNetworkStatus(ctx)
	}

	if *status {
		manager.GetNetworkStatus(ctx)
	}

	if *watch {
		fmt.Println("👁️  Watching network status (Ctrl+C to stop)...")
		for {
			manager.GetNetworkStatus(ctx)
			time.Sleep(5 * time.Second)
		}
	}

	if *exporter {
		collector := metrics.NewCollector(manager)
		server := metrics.NewServer(*port, collector)

		fmt.Printf("\n📊 Starting Prometheus exporter on port %d...\n", *port)
		fmt.Printf("📈 Metrics: http://localhost:%d/metrics\n", *port)
		fmt.Println("Press Ctrl+C to stop")

		if err := server.Start(); err != nil {
			log.Fatalf("❌ Failed to start exporter: %v", err)
		}
	}

	if !*connect && !*status && !*exporter && !*watch {
		manager.GetNetworkStatus(ctx)
		fmt.Println("\n📚 Usage:")
		fmt.Println("  -connect   Connect all nodes as peers")
		fmt.Println("  -status    Show network status")
		fmt.Println("  -watch     Watch network status (updates every 5s)")
		fmt.Println("  -exporter  Start Prometheus metrics exporter")
		fmt.Println("  -port      Exporter port (default: 9545)")
		fmt.Println("\nExample: priv-eth-net.exe -status")
	}
}
