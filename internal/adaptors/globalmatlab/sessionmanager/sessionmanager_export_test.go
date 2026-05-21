// Copyright 2026 The MathWorks, Inc.

package sessionmanager

import "time"

func (s *SessionManager) SetDiscoveryRetryInterval(interval time.Duration) {
	s.discoveryRetryInterval = interval
}
