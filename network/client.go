package network

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Node struct {
	URL       string
	Client    *ethclient.Client
	RpcClient *rpc.Client
}

// Connect to Ethereum node
func Connect(url string) (*Node, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := ethclient.NewClient(rpcClient)

	return &Node{
		URL:       url,
		Client:    client,
		RpcClient: rpcClient,
	}, nil
}

// GetBlockNumber returns current block height
func (n *Node) GetBlockNumber(ctx context.Context) (uint64, error) {
	header, err := n.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

// GetPeerCount returns number of peers
func (n *Node) GetPeerCount(ctx context.Context) (int, error) {
	var result string
	err := n.RpcClient.CallContext(ctx, &result, "net_peerCount")
	if err != nil {
		return 0, err
	}

	count := new(big.Int)
	count.SetString(result[2:], 16)
	return int(count.Int64()), nil
}

// GetNodeInfo retrieves enode URL
func (n *Node) GetNodeInfo(ctx context.Context) (string, error) {
	var result map[string]interface{}
	err := n.RpcClient.CallContext(ctx, &result, "admin_nodeInfo")
	if err != nil {
		return "", err
	}

	if enode, ok := result["enode"].(string); ok {
		return enode, nil
	}
	return "", fmt.Errorf("enode not found")
}

// IsMining checks mining status
// Note: eth_mining API is deprecated in Geth 1.14+
// For PoS (dev mode), we check if blocks are being produced
func (n *Node) IsMining(ctx context.Context) (bool, error) {
	// Try the deprecated eth_mining first for backwards compatibility
	var result bool
	err := n.RpcClient.CallContext(ctx, &result, "eth_mining")
	if err == nil {
		return result, nil
	}

	// For PoS (dev mode) in modern Geth, check if node is actually producing blocks
	// by checking if block number is increasing (> 0 means block production has started)
	header, err := n.Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return false, nil
	}

	// If block number > 0, PoS dev mode is producing blocks
	return header.Number.Uint64() > 0, nil
}

func (n *Node) Close() {
	n.Client.Close()
}
