//go:build !linux
// +build !linux

package cadvisor //nolint:golint

import (
	"github.com/grafana/agent/pkg/integrations/shared"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// NewIntegration creates a new cadvisor integration
func (c *Config) NewIntegration(logger log.Logger) (shared.Integration, error) {
	level.Warn(logger).Log("msg", "the cadvisor integration only works on linux; enabling it on other platforms will do nothing")
	return &shared.StubIntegration{}, nil
}
