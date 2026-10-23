// Copyright 2026 The MathWorks, Inc.

package otel_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/matlab/matlab-mcp-server/tests/testutils/otel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func buildExport(configure func(metrics pmetric.Metrics)) []byte {
	md := pmetric.NewMetrics()
	configure(md)
	marshaler := pmetric.ProtoMarshaler{}
	data, _ := marshaler.MarshalMetrics(md)
	return data
}

func addCounter(sm pmetric.ScopeMetrics, name, unit string, value int64) {
	m := sm.Metrics().AppendEmpty()
	m.SetName(name)
	m.SetUnit(unit)
	m.SetEmptySum().SetIsMonotonic(true)
	m.Sum().DataPoints().AppendEmpty().SetIntValue(value)
}

func postToCollector(t *testing.T, collector *otel.MockCollector, data []byte) {
	t.Helper()
	resp, err := http.Post(collector.Endpoint()+"/v1/metrics", "application/x-protobuf", bytes.NewReader(data)) //nolint:noctx // Test helper
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, resp.Body.Close())
}

func TestMockCollector_Metrics_HappyPath(t *testing.T) {
	// Arrange
	collector := otel.StartMockCollector(t)
	data := buildExport(func(md pmetric.Metrics) {
		sm := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()
		addCounter(sm, "server.starts", "{start}", 1)
	})

	// Act
	postToCollector(t, collector, data)

	// Assert
	exports := collector.Metrics()
	require.Len(t, exports, 1)
	assert.Equal(t, 1, exports[0].DataPointCount())
}

func TestMockCollector_Metrics_EmptyExport(t *testing.T) {
	// Arrange
	collector := otel.StartMockCollector(t)
	data := buildExport(func(md pmetric.Metrics) {
		md.ResourceMetrics().AppendEmpty().Resource().Attributes().PutStr("service.name", "matlab-mcp-server")
	})

	// Act
	postToCollector(t, collector, data)

	// Assert
	exports := collector.Metrics()
	require.Len(t, exports, 1)
	assert.Zero(t, exports[0].DataPointCount())
}

func TestMockCollector_ConnectionCount_RejectedRequest(t *testing.T) {
	// Arrange
	collector := otel.StartMockCollector(t)

	// Act
	resp, err := http.Post(collector.Endpoint()+"/not-metrics", "application/x-protobuf", bytes.NewReader(nil)) //nolint:noctx // Test request
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, 1, collector.ConnectionCount())
	assert.Empty(t, collector.Metrics())
}

func TestMockCollector_AllMetricsByName_HappyPath(t *testing.T) {
	// Arrange
	collector := otel.StartMockCollector(t)
	for i := range 3 {
		data := buildExport(func(md pmetric.Metrics) {
			sm := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()
			addCounter(sm, "server.starts", "{start}", int64(i+1))
			addCounter(sm, "server.client_connections", "{connection}", 1)
		})
		postToCollector(t, collector, data)
	}

	// Act
	metrics := collector.AllMetricsByName("server.starts")

	// Assert
	assert.Len(t, metrics, 3)
}
