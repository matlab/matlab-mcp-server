// Copyright 2026 The MathWorks, Inc.

package otel

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

// Unlike Collector, MockCollector keeps empty exports, which the real collector drops.
type MockCollector struct {
	server          *httptest.Server
	mu              sync.Mutex
	exports         []pmetric.Metrics
	connectionCount int
}

func StartMockCollector(t *testing.T) *MockCollector {
	t.Helper()

	mc := &MockCollector{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/metrics", mc.handleMetrics)

	mc.server = httptest.NewUnstartedServer(mux)
	mc.server.Config.ConnState = mc.trackConnection
	mc.server.Start()
	t.Cleanup(mc.server.Close)

	return mc
}

func (mc *MockCollector) Endpoint() string {
	return mc.server.URL
}

func (mc *MockCollector) ConnectionCount() int {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	return mc.connectionCount
}

func (mc *MockCollector) trackConnection(_ net.Conn, state http.ConnState) {
	if state != http.StateNew {
		return
	}

	mc.mu.Lock()
	mc.connectionCount++
	mc.mu.Unlock()
}

func (mc *MockCollector) Metrics() []pmetric.Metrics {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	result := make([]pmetric.Metrics, len(mc.exports))
	copy(result, mc.exports)
	return result
}

func (mc *MockCollector) AllMetricsByName(metricName string) []pmetric.Metric {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	var results []pmetric.Metric
	for _, export := range mc.exports {
		for i := range export.ResourceMetrics().Len() {
			rm := export.ResourceMetrics().At(i)
			for j := range rm.ScopeMetrics().Len() {
				sm := rm.ScopeMetrics().At(j)
				for k := range sm.Metrics().Len() {
					m := sm.Metrics().At(k)
					if m.Name() == metricName {
						results = append(results, m)
					}
				}
			}
		}
	}

	return results
}

func (mc *MockCollector) handleMetrics(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	unmarshaler := pmetric.ProtoUnmarshaler{}
	metrics, err := unmarshaler.UnmarshalMetrics(body)
	if err != nil {
		http.Error(w, "failed to unmarshal metrics: "+err.Error(), http.StatusBadRequest)
		return
	}

	mc.mu.Lock()
	mc.exports = append(mc.exports, metrics)
	mc.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}
