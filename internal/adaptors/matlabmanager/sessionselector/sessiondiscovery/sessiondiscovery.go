// Copyright 2026 The MathWorks, Inc.

package sessiondiscovery

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strconv"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient/embeddedconnector"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
)

const sessionDetailsFileName = "sessionDetails.json"

var ErrInvalidSessionDetails = errors.New("invalid session details")

type AppDataDirGetter interface {
	AppDataDir() (string, error)
}

type OSLayer interface {
	ReadFile(filePath string) ([]byte, error)
}

type sessionDetailsJSON struct {
	Port        json.Number `json:"port"`
	Certificate string      `json:"certificate"`
	APIKey      string      `json:"apiKey"`
	PID         json.Number `json:"pid"`
}

type SessionDiscoverer struct {
	appDataDirGetter AppDataDirGetter
	osLayer          OSLayer
}

func New(appDataDirGetter AppDataDirGetter, osLayer OSLayer) *SessionDiscoverer {
	return &SessionDiscoverer{
		appDataDirGetter: appDataDirGetter,
		osLayer:          osLayer,
	}
}

func (d *SessionDiscoverer) FromSessionDetails(logger entities.Logger, sessionDetails []byte) (embeddedconnector.ConnectionDetails, error) {
	var zeroValue embeddedconnector.ConnectionDetails

	var details sessionDetailsJSON
	if err := json.Unmarshal(sessionDetails, &details); err != nil {
		return zeroValue, err
	}

	port := details.Port.String()
	portAsInt, err := strconv.Atoi(port)
	if err != nil {
		logger.Debug("Failed to parse port as int")
		return zeroValue, ErrInvalidSessionDetails
	}

	if portAsInt < 1 || portAsInt > 65535 {
		logger.Debug("Port is out of range")
		return zeroValue, ErrInvalidSessionDetails
	}

	if details.APIKey == "" {
		logger.Debug("Invalid empty API Key")
		return zeroValue, ErrInvalidSessionDetails
	}

	certificatePEM, err := d.osLayer.ReadFile(details.Certificate)
	if err != nil {
		return zeroValue, err
	}

	if len(certificatePEM) == 0 {
		logger.Debug("Invalid empty certificate PEM")
		return zeroValue, ErrInvalidSessionDetails
	}

	return embeddedconnector.ConnectionDetails{
		Host:           "localhost",
		Port:           port,
		APIKey:         details.APIKey,
		CertificatePEM: certificatePEM,
	}, nil
}

func (d *SessionDiscoverer) DiscoverSessions(logger entities.Logger) []embeddedconnector.ConnectionDetails {
	appDataDir, err := d.appDataDirGetter.AppDataDir()
	if err != nil {
		logger.WithError(err).Debug("Failed to determine app data directory for session discovery")
		return nil
	}

	// Hardcoding v1 for now, if we end up having multiple version, we'll need version based handlers
	sessionFilePath := filepath.Join(appDataDir, "v1", sessionDetailsFileName)
	data, err := d.osLayer.ReadFile(sessionFilePath)
	if err != nil {
		logger.WithError(err).Debug("No shared MATLAB session file found")
		return nil
	}

	connectionDetails, err := d.FromSessionDetails(logger, data)
	if err != nil {
		logger.WithError(err).Debug("Failed to get connection details from session details")
		return nil
	}

	return []embeddedconnector.ConnectionDetails{connectionDetails}
}
