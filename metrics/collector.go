package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shreeyash-ugale/priv-eth-net/network"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	manager     *network.PeerManager
	blockHeight *prometheus.Desc
	peerCount   *prometheus.Desc
	isMining    *prometheus.Desc
	mu          sync.Mutex
}

// NewCollector creates Prometheus collector for Ethereum metrics
func NewCollector(manager *network.PeerManager) *Collector {
	return &Collector{
		manager: manager,
		blockHeight: prometheus.NewDesc(
			"ethereum_block_height",
			"Current block height",
			[]string{"node", "url"},
			nil,
		),
		peerCount: prometheus.NewDesc(
			"ethereum_peer_count",
			"Number of connected peers",
			[]string{"node", "url"},
			nil,
		),
		isMining: prometheus.NewDesc(
			"ethereum_is_mining",
			"Mining status (1=mining, 0=not mining)",
			[]string{"node", "url"},
			nil,
		),
	}
}

// Describe implements prometheus.Collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.blockHeight
	ch <- c.peerCount
	ch <- c.isMining
}

// Collect implements prometheus.Collector
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i, node := range c.manager.Nodes {
		nodeName := fmt.Sprintf("node-%d", i+1)

		// Collect block height
		if blockNum, err := node.GetBlockNumber(ctx); err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.blockHeight,
				prometheus.GaugeValue,
				float64(blockNum),
				nodeName,
				node.URL,
			)
		}

		// Collect peer count
		if peerCount, err := node.GetPeerCount(ctx); err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.peerCount,
				prometheus.GaugeValue,
				float64(peerCount),
				nodeName,
				node.URL,
			)
		}

		// Collect mining status
		if mining, err := node.IsMining(ctx); err == nil {
			value := 0.0
			if mining {
				value = 1.0
			}
			ch <- prometheus.MustNewConstMetric(
				c.isMining,
				prometheus.GaugeValue,
				value,
				nodeName,
				node.URL,
			)
		}
	}
}
