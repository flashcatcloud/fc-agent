package config

import (
	"bytes"
	"fmt"
	_ "github.com/grafana/agent/pkg/integrations/install"
	"github.com/grafana/agent/pkg/integrations/node_exporter"
	"github.com/grafana/agent/pkg/integrations/redis_exporter"
	v2 "github.com/grafana/agent/pkg/integrations/v2"
	"github.com/grafana/agent/pkg/integrations/windows_exporter"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigMaker(t *testing.T) {
	configStr := `wal_directory: /tmp/wal`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	fullpath := filepath.Join(tDir, "metrics-1.yml")
	err = ioutil.WriteFile(fullpath, []byte(configStr), fs.ModePerm)
	assert.Nil(t, err)
	fileFS := fmt.Sprintf("file://%s", tDir)
	loaderCfg := LoaderConfig{
		Sources:       nil,
		TemplatePaths: []string{fileFS},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	assert.NoError(t, err)
	configs, err := cmf.processMetric()
	assert.Nil(t, err)
	assert.NotNil(t, configs)
	assert.Equal(t, configs.WALDir, "/tmp/wal")
}

func TestConfigMakerWithFakeFiles(t *testing.T) {
	configStr := `wal_directory: /tmp/wal`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	fullpath := filepath.Join(tDir, "metrics-1.yml")
	err = ioutil.WriteFile(fullpath, []byte(configStr), fs.ModePerm)
	assert.Nil(t, err)

	fullpath = filepath.Join(tDir, "fake.yml")
	err = ioutil.WriteFile(fullpath, []byte(configStr), fs.ModePerm)
	assert.Nil(t, err)

	fileFS := fmt.Sprintf("file://%s", tDir)
	loaderCfg := LoaderConfig{
		Sources:       nil,
		TemplatePaths: []string{fileFS},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	assert.NoError(t, err)
	configs, err := cmf.processMetric()
	assert.Nil(t, err)
	assert.NotNil(t, configs)
	assert.Equal(t, configs.WALDir, "/tmp/wal")
}

func TestConfigMakerWithMultipleMetrics(t *testing.T) {
	configStr := `wal_directory: /tmp/wal`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	fullpath := filepath.Join(tDir, "metrics-1.yml")
	err = ioutil.WriteFile(fullpath, []byte(configStr), fs.ModePerm)
	assert.Nil(t, err)

	fullpath = filepath.Join(tDir, "metrics-2.yml")
	err = ioutil.WriteFile(fullpath, []byte(configStr), fs.ModePerm)
	assert.Nil(t, err)

	fileFS := fmt.Sprintf("file://%s", tDir)
	loaderCfg := LoaderConfig{
		Sources:       nil,
		TemplatePaths: []string{fileFS},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	assert.Nil(t, err)
	_, err = cmf.processMetric()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "multiple metrics configuration found")
}

func TestConfigMakerWithMetricsAndInstances(t *testing.T) {
	configStr := `wal_directory: /tmp/wal`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	writeFile(t, tDir, "metrics-1.yml", configStr)
	writeFile(t, tDir, "metrics_instances-1.yml", "name: t1")
	writeFile(t, tDir, "metrics_instances-2.yml", "name: t2")
	fileFS := fmt.Sprintf("file://%s", tDir)
	loaderCfg := LoaderConfig{
		Sources:       nil,
		TemplatePaths: []string{fileFS},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	assert.Nil(t, err)
	cfg := &Config{}
	err = cmf.ProcessConfigs(cfg)
	assert.Nil(t, err)
	assert.Len(t, cfg.Metrics.Configs, 2)
}

func TestConfigMakerWithExporter(t *testing.T) {
	configStr := `
windows_exporter:
  enabled_collectors: one,two,three
  instance: testinstance
`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	writeFile(t, tDir, "exporters-1.yml", configStr)
	fileFS := fmt.Sprintf("file://%s", tDir)
	loaderCfg := LoaderConfig{
		Sources:       nil,
		TemplatePaths: []string{fileFS},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	assert.Nil(t, err)
	configs, err := cmf.processExporters()
	assert.Len(t, configs, 1)
	wincfg, _ := configs[0].(v2.UpgradedConfig).LegacyConfig()

	assert.True(t, wincfg.(*windows_exporter.Config).EnabledCollectors == "one,two,three")
}

func TestConfigMakerWithExporterWithTemplate(t *testing.T) {
	configStr := `
windows_exporter:
  enabled_collectors: {{ (datasource "vars").value }}
  instance: testinstance
`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	writeFile(t, tDir, "vars.yaml", "value: banana")
	fullpath := filepath.Join(tDir, "vars.yaml")
	writeFile(t, tDir, "exporters-1.yml", configStr)
	fileFS := fmt.Sprintf("file://%s", tDir)
	loaderCfg := LoaderConfig{
		Sources: []Datasource{{
			Name: "vars",
			URL:  fmt.Sprintf("file://%s", fullpath),
		}},
		TemplatePaths: []string{fileFS},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	assert.Nil(t, err)
	configs, err := cmf.processExporters()
	assert.Len(t, configs, 1)
	wincfg, _ := configs[0].(v2.UpgradedConfig).LegacyConfig()
	assert.True(t, wincfg.(*windows_exporter.Config).EnabledCollectors == "banana")
}

func TestConfigMakerWithMultipleExporter(t *testing.T) {
	configStr := `
windows_exporter:
  enabled_collectors: one,two,three
  instance: testinstance
node_exporter:
  enabled: false
`
	tDir, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	writeFile(t, tDir, "exporters-1.yml", configStr)
	fileFS := fmt.Sprintf("file://%s", tDir)
	loaderCfg := LoaderConfig{
		Sources:       nil,
		TemplatePaths: []string{fileFS},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	assert.Nil(t, err)
	configs, err := cmf.processExporters()
	assert.Len(t, configs, 2)
	for _, cfg := range configs {
		switch v := cfg.(type) {
		default:
			fmt.Printf("unexpected type %T", v)
		case v2.UpgradedConfig:
			oldConfig, _ := v.LegacyConfig()
			switch oc := oldConfig.(type) {
			case *windows_exporter.Config:
				assert.True(t, "one,two,three" == oc.EnabledCollectors)
			case *node_exporter.Config:
				assert.NotNil(t, v)
			}
		}
	}
	_ = os.RemoveAll(tDir)
}

func TestLoadingFromS3(t *testing.T) {
	configStr := `
windows_exporter:
  enabled_collectors: one,two,three
  instance: testinstance
`
	backend := s3mem.New()
	faker := gofakes3.New(backend)

	srv := httptest.NewServer(faker.Server())
	backend.CreateBucket("mybucket")
	t.Cleanup(srv.Close)
	_, err := backend.PutObject(
		"mybucket",
		"exporters-1.yml",
		map[string]string{"Content-Type": "application/yaml"},
		bytes.NewBufferString(configStr),
		int64(len(configStr)),
	)

	u, err := url.Parse(srv.URL)
	os.Setenv("AWS_ANON", "true")
	defer os.Unsetenv("AWS_ANON")

	s3Url := "s3://mybucket/?region=us-east-1&disableSSL=true&s3ForcePathStyle=true&endpoint=" + u.Host
	assert.NoError(t, err)
	loaderCfg := LoaderConfig{
		Sources:       nil,
		TemplatePaths: []string{s3Url},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	cfg, err := cmf.processExporters()
	assert.NoError(t, err)
	assert.Len(t, cfg, 1)
	oc, _ := cfg[0].(v2.UpgradedConfig).LegacyConfig()
	winCfg := oc.(*windows_exporter.Config)
	assert.True(t, winCfg.EnabledCollectors == "one,two,three")
}

func TestLoadingFromS3LoadingVarsLocally(t *testing.T) {
	configStr := `
windows_exporter:
  enabled_collectors: {{ (datasource "vars").value }}
  instance: testinstance
`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	writeFile(t, tDir, "vars.yaml", "value: banana")

	backend := s3mem.New()
	faker := gofakes3.New(backend)

	srv := httptest.NewServer(faker.Server())
	backend.CreateBucket("mybucket")
	t.Cleanup(srv.Close)
	_, err = backend.PutObject(
		"mybucket",
		"exporters-1.yml",
		map[string]string{"Content-Type": "application/yaml"},
		bytes.NewBufferString(configStr),
		int64(len(configStr)),
	)

	u, err := url.Parse(srv.URL)
	os.Setenv("AWS_ANON", "true")
	defer os.Unsetenv("AWS_ANON")
	fullpath := filepath.Join(tDir, "vars.yaml")
	s3Url := "s3://mybucket/?region=us-east-1&disableSSL=true&s3ForcePathStyle=true&endpoint=" + u.Host
	assert.NoError(t, err)
	loaderCfg := LoaderConfig{
		Sources: []Datasource{{
			Name: "vars",
			URL:  fmt.Sprintf("file://%s", fullpath),
		}},
		TemplatePaths: []string{s3Url},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	cfg, err := cmf.processExporters()
	assert.NoError(t, err)
	assert.Len(t, cfg, 1)
	oc, _ := cfg[0].(v2.UpgradedConfig).LegacyConfig()
	winCfg := oc.(*windows_exporter.Config)
	assert.True(t, winCfg.EnabledCollectors == "banana")
}

func TestLoadingFromS3LoadingVarsLocallyWithRange(t *testing.T) {
	configStr := `
windows_exporter:
  enabled_collectors: banana
  instance: testinstance
  autoscrape:
    metric_relabel_configs: {{ range (datasource "vars").value }}
    - source_labels: [__address__]
      target_label: {{ . }}
      replacement: "{{ . }}-value"
    {{ end }}
`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	writeFile(t, tDir, "vars.yaml", "value: [banana,apple,pear]")

	backend := s3mem.New()
	faker := gofakes3.New(backend)

	srv := httptest.NewServer(faker.Server())
	backend.CreateBucket("mybucket")
	t.Cleanup(srv.Close)
	_, err = backend.PutObject(
		"mybucket",
		"exporters-1.yml",
		map[string]string{"Content-Type": "application/yaml"},
		bytes.NewBufferString(configStr),
		int64(len(configStr)),
	)

	u, err := url.Parse(srv.URL)
	os.Setenv("AWS_ANON", "true")
	defer os.Unsetenv("AWS_ANON")
	fullpath := filepath.Join(tDir, "vars.yaml")
	s3Url := "s3://mybucket/?region=us-east-1&disableSSL=true&s3ForcePathStyle=true&endpoint=" + u.Host
	assert.NoError(t, err)
	loaderCfg := LoaderConfig{
		Sources: []Datasource{{
			Name: "vars",
			URL:  fmt.Sprintf("file://%s", fullpath),
		}},
		TemplatePaths: []string{s3Url},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	cfg, err := cmf.processExporters()
	assert.NoError(t, err)
	assert.Len(t, cfg, 1)
	oc := cfg[0].(*v2.ConfigShim)
	assert.Len(t, oc.Common.Autoscrape.MetricRelabelConfigs, 3)
	foundApple := 0
	foundPear := 0
	foundBanana := 0
	for _, rc := range oc.Common.Autoscrape.MetricRelabelConfigs {
		if rc.TargetLabel == "apple" {
			foundApple++
		}
		if rc.TargetLabel == "pear" {
			foundPear++
		}
		if rc.TargetLabel == "banana" {
			foundBanana++
		}
	}
	assert.True(t, (foundPear+foundApple+foundBanana) == 3)

}

func TestMultiplex(t *testing.T) {
	configStr := `
redis_exporter_configs:
- redis_addr: localhost:6379
  autoscrape:
    metric_relabel_configs: 
    - source_labels: [__address__]
      target_label: "banana"
      replacement: "apple"
- redis_addr: localhost:6380
`
	tDir, err := os.MkdirTemp("", "")
	defer os.RemoveAll(tDir)
	assert.Nil(t, err)
	writeFile(t, tDir, "exporters-1.yml", configStr)
	fileFS := fmt.Sprintf("file://%s", tDir)
	loaderCfg := LoaderConfig{
		Sources:       nil,
		TemplatePaths: []string{fileFS},
	}
	cmf, err := NewConfigLoader(loaderCfg)
	assert.Nil(t, err)
	configs, err := cmf.processExporters()
	assert.Nil(t, err)
	assert.Len(t, configs, 2)
	found := 0
	for _, c := range configs {
		cs, _ := c.(*v2.ConfigShim)
		rc := cs.Orig.(*redis_exporter.Config)
		if rc.RedisAddr == "localhost:6379" {
			assert.Len(t, cs.Common.Autoscrape.MetricRelabelConfigs, 1)
			found++
		}
		if rc.RedisAddr == "localhost:6380" {
			found++
		}
	}

	assert.True(t, found == 2)
}

func writeFile(t *testing.T, directory string, path string, contents string) {
	fullpath := filepath.Join(directory, path)
	err := ioutil.WriteFile(fullpath, []byte(contents), fs.ModePerm)
	assert.Nil(t, err)
}
