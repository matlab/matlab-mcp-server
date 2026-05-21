// Copyright 2026 The MathWorks, Inc.

package otel

type collectorConfig struct{}

func DefaultConfig() collectorConfig {
	return collectorConfig{}
}

func (collectorConfig) generateConfig() string {
	config := `
extensions:
  health_check:
    endpoint: "0.0.0.0:` + healthCheckPort + `"
receivers:
  otlp:
    protocols:
      http:
        endpoint: "0.0.0.0:4318"
exporters:
  file:
    path: "/tmp/telemetry/` + telemetryFileName + `"
    format: json
service:
  extensions: [health_check]
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [file]
`

	return config
}
