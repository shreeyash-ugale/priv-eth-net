package metrics

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Port      int
	Collector *Collector
}

// NewServer creates metrics HTTP server
func NewServer(port int, collector *Collector) *Server {
	return &Server{
		Port:      port,
		Collector: collector,
	}
}

// Start begins serving metrics
func (s *Server) Start() error {
	// Register collector
	prometheus.MustRegister(s.Collector)

	// Metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%d", s.Port)
	log.Printf("Metrics server listening on %s", addr)
	log.Printf("Visit http://localhost%s/metrics", addr)

	return http.ListenAndServe(addr, nil)
}
