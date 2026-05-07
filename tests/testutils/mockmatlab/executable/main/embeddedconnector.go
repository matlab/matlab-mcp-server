// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/tests/testutils/mockmatlab/mockruntime"
)

func startConnectorServer(ctx context.Context, sessionDir, apiKey string, tlsCfg *tls.Config, recorder *mockruntime.EventRecorder) error {
	ctx, cancelExit := context.WithCancel(ctx)
	defer cancelExit()

	mux := http.NewServeMux()
	mux.HandleFunc("/messageservice/json/secure", makeHandler(apiKey, handleEval(cancelExit, recorder)))
	mux.HandleFunc("/messageservice/json/state", makeHandler(apiKey, handleState))

	listener, err := tls.Listen("tcp", "127.0.0.1:0", tlsCfg)
	if err != nil {
		return fmt.Errorf("failed to start TLS listener: %w", err)
	}

	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		closeListenerWithLog(listener)
		return fmt.Errorf("failed to parse listener address: %w", err)
	}

	portFile := filepath.Join(sessionDir, "connector.securePort")
	if err := os.WriteFile(portFile, []byte(port), 0o600); err != nil {
		closeListenerWithLog(listener)
		return fmt.Errorf("failed to write port file: %w", err)
	}

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	return server.Shutdown(context.Background())
}

func makeHandler(apiKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
			return
		}
		if r.Header.Get("mwapikey") != apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func handleEval(cancelExit context.CancelFunc, recorder *mockruntime.EventRecorder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}

		var request embeddedconnector.ConnectorPayload
		if err := json.Unmarshal(body, &request); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		var response embeddedconnector.ConnectorPayload

		if len(request.Messages.Eval) > 0 {
			code := request.Messages.Eval[0].Code
			recorder.RecordEval(code)
			response = embeddedconnector.ConnectorPayload{
				Messages: embeddedconnector.ConnectorMessage{
					EvalResponse: []embeddedconnector.EvalResponseMessage{
						{
							IsError:     false,
							ResponseStr: code,
						},
					},
				},
			}
			if mockruntime.ShouldExitFromEvalCode(code) {
				recorder.RecordExitRequested()
				defer cancelExit()
			}
		}

		if len(request.Messages.FEval) > 0 {
			recorder.RecordFeval(request.Messages.FEval[0].Function)
			response = embeddedconnector.ConnectorPayload{
				Messages: embeddedconnector.ConnectorMessage{
					FevalResponse: []embeddedconnector.FevalResponseMessage{
						{
							IsError:       false,
							MessageFaults: nil,
							Results:       []interface{}{request.Messages.FEval[0].Function},
						},
					},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}

func closeListenerWithLog(listener net.Listener) {
	if err := listener.Close(); err != nil {
		log.Printf("failed to close listener: %v", err)
	}
}

func handleState(w http.ResponseWriter, _ *http.Request) {
	response := embeddedconnector.ConnectorPayload{
		Messages: embeddedconnector.ConnectorMessage{
			PingResponse: []embeddedconnector.PingResponseMessage{
				{
					MessageFaults: nil,
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
