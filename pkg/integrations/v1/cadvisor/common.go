package cadvisor

import (
	"time"
)

const name = "cadvisor"

// DefaultConfig holds the default settings for the cadvisor integration
var DefaultConfig Config = Config{
	// Common cadvisor shared defaults
	StoreContainerLabels: true,
	ResctrlInterval:      0,

	StorageDuration: 2 * time.Minute,

	// Containerd shared defaults
	Containerd:          "/run/containerd/containerd.sock",
	ContainerdNamespace: "k8s.io",

	// Docker shared defaults
	Docker:        "unix:///var/run/docker.sock",
	DockerTLS:     false,
	DockerTLSCert: "cert.pem",
	DockerTLSKey:  "key.pem",
	DockerTLSCA:   "ca.pem",

	// Raw shared defaults
	DockerOnly: false,
}

// Config controls cadvisor
type Config struct {
	// Common cadvisor shared options
	// StoreContainerLabels converts container labels and environment variables into labels on prometheus metrics for each container. If false, then only metrics exported are container name, first alias, and image name.
	StoreContainerLabels bool `yaml:"store_container_labels,omitempty"`

	// AllowlistedContainerLabels list of container labels to be converted to labels on prometheus metrics for each container. store_container_labels must be set to false for this to take effect.
	AllowlistedContainerLabels []string `yaml:"allowlisted_container_labels,omitempty"`

	// EnvMetadataAllowlist list of environment variable keys matched with specified prefix that needs to be collected for containers, only support containerd and docker runtime for now.
	EnvMetadataAllowlist []string `yaml:"env_metadata_allowlist,omitempty"`

	// RawCgroupPrefixAllowlist list of cgroup path prefix that needs to be collected even when -docker_only is specified.
	RawCgroupPrefixAllowlist []string `yaml:"raw_cgroup_prefix_allowlist,omitempty"`

	// PerfEventsConfig path to a JSON file containing configuration of perf events to measure. Empty value disabled perf events measuring.
	PerfEventsConfig string `yaml:"perf_events_config,omitempty"`

	// ResctrlInterval resctrl mon groups updating interval. Zero value disables updating mon groups.
	ResctrlInterval int `yaml:"resctrl_interval,omitempty"`

	// DisableMetrics list of `metrics` to be disabled.
	DisabledMetrics []string `yaml:"disabled_metrics,omitempty"`

	// EnableMetrics list of `metrics` to be enabled. If set, overrides 'disable_metrics'.
	EnabledMetrics []string `yaml:"enabled_metrics,omitempty"`

	// StorageDuration length of time to keep data stored in memory (Default: 2m)
	StorageDuration time.Duration `yaml:"storage_duration,omitempty"`

	// Containerd shared options
	// Containerd containerd endpoint
	Containerd string `yaml:"containerd,omitempty"`

	// ContainerdNamespace containerd namespace
	ContainerdNamespace string `yaml:"containerd_namespace,omitempty"`

	// Docker shared options
	// Docker docker endpoint
	Docker string `yaml:"docker,omitempty"`

	// DockerTLS use TLS to connect to docker
	DockerTLS bool `yaml:"docker_tls,omitempty"`

	// DockerTLSCert path to client certificate
	DockerTLSCert string `yaml:"docker_tls_cert,omitempty"`

	// DockerTLSKey path to private key
	DockerTLSKey string `yaml:"docker_tls_key,omitempty"`

	// DockerTLSCA path to trusted CA
	DockerTLSCA string `yaml:"docker_tls_ca,omitempty"`

	// Raw shared options
	// DockerOnly only report docker containers in addition to root stats
	DockerOnly bool `yaml:"docker_only,omitempty"`
}

// Name returns the name of the integration that this shared represents.
func (c *Config) Name() string {
	return name
}

// InstanceKey returns the agentKey
func (c *Config) InstanceKey(agentKey string) (string, error) {
	return agentKey, nil
}
