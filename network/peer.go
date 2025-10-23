package network

import (
	"context"
	"fmt"
	"strings"
)

type PeerManager struct {
	Nodes []*Node
}

// NewPeerManager creates manager for multiple nodes
func NewPeerManager(urls []string) (*PeerManager, error) {
	nodes := make([]*Node, 0, len(urls))

	for _, url := range urls {
		node, err := Connect(url)
		if err != nil {
			fmt.Printf("Warning: failed to connect to %s: %v\n", url, err)
			continue
		}
		nodes = append(nodes, node)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes connected")
	}

	return &PeerManager{Nodes: nodes}, nil
}

// ConnectPeers connects all nodes to each other
func (pm *PeerManager) ConnectPeers(ctx context.Context) error {
	fmt.Println("Connecting peers...")

	// Get enode URLs from all nodes
	enodes := make([]string, len(pm.Nodes))
	for i, node := range pm.Nodes {
		enode, err := node.GetNodeInfo(ctx)
		if err != nil {
			return fmt.Errorf("failed to get enode from node %d: %w", i+1, err)
		}

		// Replace localhost/IP with container name for Docker networking
		enode = replaceHostInEnode(enode, i+1)
		enodes[i] = enode
		fmt.Printf("Node %d enode: %s\n", i+1, enode)
	}

	// Connect each node to all others
	for i, node := range pm.Nodes {
		for j, enode := range enodes {
			if i != j {
				var result bool
				err := node.RpcClient.CallContext(ctx, &result, "admin_addPeer", enode)
				if err != nil {
					fmt.Printf("Warning: failed to add peer from node%d to node%d: %v\n", i+1, j+1, err)
				} else {
					fmt.Printf("✓ Connected node%d -> node%d\n", i+1, j+1)
				}
			}
		}
	}

	return nil
}

// replaceHostInEnode replaces IP with Docker container name
func replaceHostInEnode(enode string, nodeNum int) string {
	// enode://pubkey@ip:port
	parts := strings.Split(enode, "@")
	if len(parts) != 2 {
		return enode
	}

	portParts := strings.Split(parts[1], ":")
	if len(portParts) != 2 {
		return enode
	}

	containerName := fmt.Sprintf("geth-node%d", nodeNum)
	return fmt.Sprintf("%s@%s:%s", parts[0], containerName, portParts[1])
}

// GetNetworkStatus returns status of all nodes
func (pm *PeerManager) GetNetworkStatus(ctx context.Context) {
	fmt.Println("\n=== Network Status ===")

	for i, node := range pm.Nodes {
		fmt.Printf("\nNode %d (%s):\n", i+1, node.URL)

		blockNum, err := node.GetBlockNumber(ctx)
		if err != nil {
			fmt.Printf("  Block: Error - %v\n", err)
		} else {
			fmt.Printf("  Block: %d\n", blockNum)
		}

		peerCount, err := node.GetPeerCount(ctx)
		if err != nil {
			fmt.Printf("  Peers: Error - %v\n", err)
		} else {
			fmt.Printf("  Peers: %d\n", peerCount)
		}

		mining, err := node.IsMining(ctx)
		if err != nil {
			fmt.Printf("  Mining: Error - %v\n", err)
		} else {
			fmt.Printf("  Mining: %v\n", mining)
		}
	}
	fmt.Println()
}

func (pm *PeerManager) Close() {
	for _, node := range pm.Nodes {
		node.Close()
	}
}
