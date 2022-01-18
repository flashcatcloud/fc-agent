package integrations

import (
	"fmt"
	"reflect"

	v2 "github.com/grafana/agent/pkg/integrations/v2"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/grafana/agent/pkg/metrics"
	"github.com/grafana/agent/pkg/util"
	"github.com/prometheus/statsd_exporter/pkg/level"
	"github.com/weaveworks/common/server"
	"gopkg.in/yaml.v2"
)

type integrationsVersion int

const (
	IntegrationsVersion1 integrationsVersion = iota
	IntegrationsVersion2
)

// DefaultVersionedIntegrations is the default config for integrations.
var DefaultVersionedIntegrations = VersionedIntegrations{
	version:  IntegrationsVersion1,
	ConfigV1: &DefaultManagerConfig,
}

// VersionedIntegrations abstracts the subsystem configs for integrations v1
// and v2. VersionedIntegrations can only be unmarshaled as part of Load.
type VersionedIntegrations struct {
	version integrationsVersion
	raw     util.RawYAML

	ConfigV1 *ManagerConfig
	ConfigV2 *v2.SubsystemOptions
}

var (
	_ yaml.Unmarshaler = (*VersionedIntegrations)(nil)
	_ yaml.Marshaler   = (*VersionedIntegrations)(nil)
)

// UnmarshalYAML implements yaml.Unmarshaler. Full unmarshaling is deferred until
// setVersion is invoked.
func (c *VersionedIntegrations) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.ConfigV1 = nil
	c.ConfigV2 = nil
	return unmarshal(&c.raw)
}

// MarshalYAML implements yaml.Marshaler.
func (c VersionedIntegrations) MarshalYAML() (interface{}, error) {
	switch {
	case c.ConfigV1 != nil:
		return c.ConfigV1, nil
	case c.ConfigV2 != nil:
		return c.ConfigV2, nil
	default:
		return c.raw, nil
	}
}

// IsZero implements yaml.IsZeroer.
func (c VersionedIntegrations) IsZero() bool {
	switch {
	case c.ConfigV1 != nil:
		return reflect.ValueOf(*c.ConfigV1).IsZero()
	case c.ConfigV2 != nil:
		return reflect.ValueOf(*c.ConfigV2).IsZero()
	default:
		return len(c.raw) == 0
	}
}

// ApplyDefaults applies defaults to the subsystem based on globals.
func (c *VersionedIntegrations) ApplyDefaults(scfg *server.Config, mcfg *metrics.Config) error {
	if c.version != IntegrationsVersion2 {
		return c.ConfigV1.ApplyDefaults(scfg, mcfg)
	}
	return c.ConfigV2.ApplyDefaults(mcfg)
}

// SetVersion completes the deferred unmarshal and unmarshals the raw YAML into
// the subsystem config for version v.
func (c *VersionedIntegrations) SetVersion(v integrationsVersion, logger log.Logger) error {
	c.version = v

	switch c.version {
	case IntegrationsVersion1:
		cfg := DefaultManagerConfig
		c.ConfigV1 = &cfg
		err := yaml.UnmarshalStrict(c.raw, c.ConfigV1)
		// Node exporter has some post-processing that has to be done for migrations
		if err != nil {
			return err
		}
		if c.ConfigV1.Integrations.NodeExporter != nil {
			return c.ConfigV1.Integrations.NodeExporter.Config.PostProcessing()
		}
		return nil
	case IntegrationsVersion2:
		cfg := v2.DefaultSubsystemOptions
		c.ConfigV2 = &cfg
		result := yaml.UnmarshalStrict(c.raw, c.ConfigV2)
		return result
	default:
		panic(fmt.Sprintf("unknown integrations version %d", c.version))
	}
}

// IntegrationsGlobals is a global struct shared across integrations.
type IntegrationsGlobals = v2.Globals

// Integrations is an abstraction over both the v1 and v2 systems.
type Integrations interface {
	ApplyConfig(*VersionedIntegrations, IntegrationsGlobals) error
	WireAPI(*mux.Router)
	Stop()
}

// NewIntegrations creates a new subsystem. globals should be provided regardless
// of useV2. globals.SubsystemOptions will be automatically set if cfg.Version
// is set to IntegrationsVersion2.
func NewIntegrations(logger log.Logger, cfg *VersionedIntegrations, globals IntegrationsGlobals) (Integrations, error) {
	if cfg.version != IntegrationsVersion2 {
		instance, err := NewManager(*cfg.ConfigV1, logger, globals.Metrics.InstanceManager(), globals.Metrics.Validate)
		if err != nil {
			return nil, err
		}
		return &v1Integrations{Manager: instance}, nil
	}

	level.Warn(logger).Log("msg", "integrations-next is enabled. integrations-next is subject to change")

	globals.SubsystemOpts = *cfg.ConfigV2
	instance, err := v2.NewSubsystem(logger, globals)
	if err != nil {
		return nil, err
	}
	return &v2Integrations{Subsystem: instance}, nil
}

type v1Integrations struct{ *Manager }

func (s *v1Integrations) ApplyConfig(cfg *VersionedIntegrations, globals IntegrationsGlobals) error {
	return s.Manager.ApplyConfig(*cfg.ConfigV1)
}

type v2Integrations struct{ *v2.Subsystem }

func (s *v2Integrations) ApplyConfig(cfg *VersionedIntegrations, globals IntegrationsGlobals) error {
	globals.SubsystemOpts = *cfg.ConfigV2
	return s.Subsystem.ApplyConfig(globals)
}