// Copyright 2026 The MathWorks, Inc.

package otel

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/matlab/matlab-mcp-server/internal/entities"
	otelapi "go.opentelemetry.io/otel"
)

const RedactedEndpoint = "[redacted]"

func NewErrorHandler(logger entities.Logger, endpoint string) ErrorHandler {
	return &loggerErrorHandler{
		logger:          logger,
		endpointPattern: endpointPattern(endpoint),
	}
}

func SetErrorHandler(handler ErrorHandler) {
	otelapi.SetErrorHandler(handler)
}

type loggerErrorHandler struct {
	logger          entities.Logger
	endpointPattern *regexp.Regexp
	warnOnce        sync.Once
}

func (h *loggerErrorHandler) Handle(err error) {
	h.warnOnce.Do(func() {
		h.logger.Warn("Telemetry error; metrics may not be reaching the collector. Further telemetry errors are logged at debug level.")
	})
	h.logger.WithError(h.redact(err)).Debug("Telemetry error")
}

func (h *loggerErrorHandler) redact(err error) error {
	if h.endpointPattern == nil {
		return err
	}
	return errors.New(h.endpointPattern.ReplaceAllLiteralString(err.Error(), RedactedEndpoint))
}

func endpointPattern(endpoint string) *regexp.Regexp {
	if endpoint == "" {
		return nil
	}

	alternatives := []string{regexp.QuoteMeta(endpoint)}
	if parsed, err := url.Parse(endpoint); err == nil && parsed.Hostname() != "" {
		alternatives = append(alternatives, regexp.QuoteMeta(parsed.Hostname()))
	}

	return regexp.MustCompile("(?i)" + strings.Join(alternatives, "|"))
}
