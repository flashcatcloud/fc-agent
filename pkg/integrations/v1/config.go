package v1

import (
	"github.com/grafana/agent/pkg/integrations/shared"
	"github.com/grafana/agent/pkg/integrations/v1/agent"
	"github.com/grafana/agent/pkg/integrations/v1/cadvisor"
	"github.com/grafana/agent/pkg/integrations/v1/consul_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/dnsmasq_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/elasticsearch_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/github_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/kafka_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/memcached_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/mongodb_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/mysqld_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/node_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/postgres_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/process_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/redis_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/statsd_exporter"
	"github.com/grafana/agent/pkg/integrations/v1/windows_exporter"
)

type V1Integration struct {
  Agent *Agent `yaml:"agent,omitempty"`
Cadvisor *Cadvisor `yaml:"cadvisor,omitempty"`
ConsulExporter *ConsulExporter `yaml:"consul_exporter,omitempty"`
DnsmasqExporter *DnsmasqExporter `yaml:"dnsmasq_exporter,omitempty"`
ElasticsearchExporter *ElasticsearchExporter `yaml:"elasticsearch_exporter,omitempty"`
GithubExporter *GithubExporter `yaml:"github_exporter,omitempty"`
KafkaExporter *KafkaExporter `yaml:"kafka_exporter,omitempty"`
MemcachedExporter *MemcachedExporter `yaml:"memcached_exporter,omitempty"`
MongodbExporter *MongodbExporter `yaml:"mongodb_exporter,omitempty"`
MysqldExporter *MysqldExporter `yaml:"mysqld_exporter,omitempty"`
NodeExporter *NodeExporter `yaml:"node_exporter,omitempty"`
PostgresExporter *PostgresExporter `yaml:"postgres_exporter,omitempty"`
ProcessExporter *ProcessExporter `yaml:"process_exporter,omitempty"`
RedisExporter *RedisExporter `yaml:"redis_exporter,omitempty"`
StatsdExporter *StatsdExporter `yaml:"statsd_exporter,omitempty"`
WindowsExporter *WindowsExporter `yaml:"windows_exporter,omitempty"`
TestConfigs []shared.V1IntegrationConfig `yaml:"-,omitempty"`

}

func (v *V1Integration) ActiveConfigs() []shared.V1IntegrationConfig {
    activeConfigs := make([]shared.V1IntegrationConfig,0)
	if v.Agent != nil {
        activeConfigs = append(activeConfigs, v.Agent)
    }
	if v.Cadvisor != nil {
        activeConfigs = append(activeConfigs, v.Cadvisor)
    }
	if v.ConsulExporter != nil {
        activeConfigs = append(activeConfigs, v.ConsulExporter)
    }
	if v.DnsmasqExporter != nil {
        activeConfigs = append(activeConfigs, v.DnsmasqExporter)
    }
	if v.ElasticsearchExporter != nil {
        activeConfigs = append(activeConfigs, v.ElasticsearchExporter)
    }
	if v.GithubExporter != nil {
        activeConfigs = append(activeConfigs, v.GithubExporter)
    }
	if v.KafkaExporter != nil {
        activeConfigs = append(activeConfigs, v.KafkaExporter)
    }
	if v.MemcachedExporter != nil {
        activeConfigs = append(activeConfigs, v.MemcachedExporter)
    }
	if v.MongodbExporter != nil {
        activeConfigs = append(activeConfigs, v.MongodbExporter)
    }
	if v.MysqldExporter != nil {
        activeConfigs = append(activeConfigs, v.MysqldExporter)
    }
	if v.NodeExporter != nil {
        activeConfigs = append(activeConfigs, v.NodeExporter)
    }
	if v.PostgresExporter != nil {
        activeConfigs = append(activeConfigs, v.PostgresExporter)
    }
	if v.ProcessExporter != nil {
        activeConfigs = append(activeConfigs, v.ProcessExporter)
    }
	if v.RedisExporter != nil {
        activeConfigs = append(activeConfigs, v.RedisExporter)
    }
	if v.StatsdExporter != nil {
        activeConfigs = append(activeConfigs, v.StatsdExporter)
    }
	if v.WindowsExporter != nil {
        activeConfigs = append(activeConfigs, v.WindowsExporter)
    }
	for _, i := range v.TestConfigs {
        activeConfigs = append(activeConfigs, i)
    }
    return activeConfigs
}


type Agent struct {
  agent.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *Agent) Cfg() shared.Config {
	return &c.Config
}

func (c *Agent) Cmn() shared.Common {
	return c.Common
}




type Cadvisor struct {
  cadvisor.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *Cadvisor) Cfg() shared.Config {
	return &c.Config
}

func (c *Cadvisor) Cmn() shared.Common {
	return c.Common
}


func (c *Cadvisor) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = cadvisor.DefaultConfig
	type plain Cadvisor
	return unmarshal((*plain)(c))
}



type ConsulExporter struct {
  consul_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *ConsulExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *ConsulExporter) Cmn() shared.Common {
	return c.Common
}


func (c *ConsulExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = consul_exporter.DefaultConfig
	type plain ConsulExporter
	return unmarshal((*plain)(c))
}



type DnsmasqExporter struct {
  dnsmasq_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *DnsmasqExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *DnsmasqExporter) Cmn() shared.Common {
	return c.Common
}


func (c *DnsmasqExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = dnsmasq_exporter.DefaultConfig
	type plain DnsmasqExporter
	return unmarshal((*plain)(c))
}



type ElasticsearchExporter struct {
  elasticsearch_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *ElasticsearchExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *ElasticsearchExporter) Cmn() shared.Common {
	return c.Common
}


func (c *ElasticsearchExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = elasticsearch_exporter.DefaultConfig
	type plain ElasticsearchExporter
	return unmarshal((*plain)(c))
}



type GithubExporter struct {
  github_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *GithubExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *GithubExporter) Cmn() shared.Common {
	return c.Common
}


func (c *GithubExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = github_exporter.DefaultConfig
	type plain GithubExporter
	return unmarshal((*plain)(c))
}



type KafkaExporter struct {
  kafka_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *KafkaExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *KafkaExporter) Cmn() shared.Common {
	return c.Common
}


func (c *KafkaExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = kafka_exporter.DefaultConfig
	type plain KafkaExporter
	return unmarshal((*plain)(c))
}



type MemcachedExporter struct {
  memcached_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *MemcachedExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *MemcachedExporter) Cmn() shared.Common {
	return c.Common
}


func (c *MemcachedExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = memcached_exporter.DefaultConfig
	type plain MemcachedExporter
	return unmarshal((*plain)(c))
}



type MongodbExporter struct {
  mongodb_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *MongodbExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *MongodbExporter) Cmn() shared.Common {
	return c.Common
}




type MysqldExporter struct {
  mysqld_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *MysqldExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *MysqldExporter) Cmn() shared.Common {
	return c.Common
}


func (c *MysqldExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = mysqld_exporter.DefaultConfig
	type plain MysqldExporter
	return unmarshal((*plain)(c))
}



type NodeExporter struct {
  node_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *NodeExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *NodeExporter) Cmn() shared.Common {
	return c.Common
}


func (c *NodeExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = node_exporter.DefaultConfig
	type plain NodeExporter
	return unmarshal((*plain)(c))
}



type PostgresExporter struct {
  postgres_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *PostgresExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *PostgresExporter) Cmn() shared.Common {
	return c.Common
}




type ProcessExporter struct {
  process_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *ProcessExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *ProcessExporter) Cmn() shared.Common {
	return c.Common
}


func (c *ProcessExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = process_exporter.DefaultConfig
	type plain ProcessExporter
	return unmarshal((*plain)(c))
}



type RedisExporter struct {
  redis_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *RedisExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *RedisExporter) Cmn() shared.Common {
	return c.Common
}


func (c *RedisExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = redis_exporter.DefaultConfig
	type plain RedisExporter
	return unmarshal((*plain)(c))
}



type StatsdExporter struct {
  statsd_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *StatsdExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *StatsdExporter) Cmn() shared.Common {
	return c.Common
}


func (c *StatsdExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = statsd_exporter.DefaultConfig
	type plain StatsdExporter
	return unmarshal((*plain)(c))
}



type WindowsExporter struct {
  windows_exporter.Config `yaml:",omitempty,inline"`
  shared.Common `yaml:",omitempty,inline"`
}

func (c *WindowsExporter) Cfg() shared.Config {
	return &c.Config
}

func (c *WindowsExporter) Cmn() shared.Common {
	return c.Common
}


func (c *WindowsExporter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = windows_exporter.DefaultConfig
	type plain WindowsExporter
	return unmarshal((*plain)(c))
}

