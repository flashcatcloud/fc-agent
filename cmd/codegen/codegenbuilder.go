package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	v1 "github.com/grafana/agent/pkg/integrations/v1"

	"github.com/grafana/agent/pkg/integrations/shared"
)

type codeGen struct {
}

func (c *codeGen) generateConfigMeta() []configMeta {
	configMetas := make([]configMeta, 0)
	for _, i := range v1.Configs {
		configMetas = append(configMetas, newConfigMeta(i.Config, i.DefaultConfig, i.Type))
	}
	return configMetas
}

func (c *codeGen) createV1Config() string {
	configs := c.generateConfigMeta()
	integrationTemplate, err := template.New("shared").Parse(`

type {{.Name}} struct {
  {{.ConfigStruct}} ` + "`yaml:\",omitempty,inline\"`" + `
  shared.Common ` + "`yaml:\",omitempty,inline\"`" + `
}

{{ if .DefaultConfig -}}
func (c *{{ .Name }}) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = {{ .PackageName }}.DefaultConfig
	type plain {{ .Name }}
	return unmarshal((*plain)(c))
}
{{ end -}}

`)
	v1Template, err := template.New("v1").Parse(`
type V1Integration struct {
  {{ range $index, $element := . -}}
	{{ $element.Name }} *{{$element.Name}} ` + "`yaml:\"{{ $element.PackageName }},omitempty\"`\n" +
		`{{ end -}}
   TestConfigs []shared.V1IntegrationConfig ` + "`yaml:\"-,omitempty\"`\n" + `
}

func (v *V1Integration) ActiveConfigs() []shared.V1IntegrationConfig {
    activeConfigs := make([]shared.V1IntegrationConfig,0)
	{{ range $index, $element := . -}}
	if v.{{ $element.Name }} != nil {
        activeConfigs = append(activeConfigs, newConfigWrapper(&v.{{ $element.Name}}.Config, v.{{ $element.Name}}.Common))
    }
	{{ end -}}
    for _, i := range v.TestConfigs {
        activeConfigs = append(activeConfigs, i)
    }
    return activeConfigs
}


type ConfigWrapper struct {
	cfg shared.Config
	cmn shared.Common
}

func (c *ConfigWrapper) Common() shared.Common {
	return c.cmn
}

func (c *ConfigWrapper) Config() shared.Config {
	return c.cfg
}

func newConfigWrapper(cfg shared.Config, cmn shared.Common) *ConfigWrapper {
	return &ConfigWrapper{
		cfg: cfg,
		cmn: cmn,
	}
}
`)
	if err != nil {
		panic(err)
	}
	v1ConfigBuilder := strings.Builder{}
	v1ConfigBuilder.WriteString("package v1 //nolint:golint\n")
	v1ConfigBuilder.WriteString(`
import (
"github.com/grafana/agent/pkg/integrations/shared"
`)
	for _, i := range configs {
		v1ConfigBuilder.WriteString(fmt.Sprintf("\"%s\"\n", i.PackagePath))
	}
	v1ConfigBuilder.WriteString(")\n")
	v1Buffer := bytes.Buffer{}
	err = v1Template.Execute(&v1Buffer, configs)
	if err != nil {
		panic(err)
	}

	v1ConfigBuilder.WriteString(v1Buffer.String())

	for _, cfg := range configs {
		bf := bytes.Buffer{}
		err = integrationTemplate.Execute(&bf, cfg)
		if err != nil {
			panic(err)
		}
		v1ConfigBuilder.WriteString(bf.String())
	}
	return v1ConfigBuilder.String()
}

func (c *codeGen) createV2Config() string {
	configs := c.generateConfigMeta()

	integrationTemplate, err := template.New("shared").Parse(`

type {{.Name}} struct {
  {{.ConfigStruct}} ` + "`yaml:\",omitempty,inline\"`" + `
  Cmn common.MetricsConfig  ` + "`yaml:\",inline\"`" + `
}

{{ if .DefaultConfig -}}
func (c *{{ .Name }}) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Config = {{ .PackageName }}.DefaultConfig
	type plain {{ .Name }}
	return unmarshal((*plain)(c))
}
{{ end -}}

`)
	// Type: 1 = Singleton, 2 = Multiplex
	v2template, err := template.New("v1").Parse(`
type Integrations struct {
  {{ range $index, $element := . -}}
	{{ if eq .Type 0 -}}
		{{ $element.Name }} *{{$element.Name}} ` + "`yaml:\"{{ $element.PackageName }},omitempty\"`" + `
    {{ end -}}
    {{ if eq .Type  1 -}}
       {{ $element.Name }}Configs []*{{$element.Name}} ` + "`yaml:\"{{ $element.PackageName }}_configs,omitempty\"`" + `
    {{ end -}}
    {{ end -}}
   TestConfigs []Config  ` + "`yaml:\"-,omitempty\"`\n" + `

}

func (v *Integrations) ActiveConfigs() []Config {
    activeConfigs := make([]Config,0)
	{{ range $index, $element := . -}}
    {{ if eq .Type  0 -}}
	if v.{{ $element.Name }} != nil {
        activeConfigs = append(activeConfigs, newConfigWrapper(v.{{ $element.Name}}, v.{{ $element.Name}}.Cmn, v.{{ $element.Name}}.NewIntegration, v.{{ $element.Name}}.InstanceKey))
    }
    {{ end -}}
	{{ if eq .Type  1 -}}
	for _, i := range v.{{ $element.Name}}Configs {
		activeConfigs = append(activeConfigs, newConfigWrapper(i, i.Cmn, i.NewIntegration, i.InstanceKey))
	}
    {{ end -}}
	{{ end -}}
    for _, i := range v.TestConfigs {
        activeConfigs = append(activeConfigs, i)
    }
    return activeConfigs
}

type configWrapper struct {
	cfg                shared.Config
	cmn                common.MetricsConfig
	configInstanceFunc configInstance
	newInstanceFunc    newIntegration
}

func (c *configWrapper) ApplyDefaults(globals Globals) error {
	c.cmn.ApplyDefaults(globals.SubsystemOpts.Metrics.Autoscrape)
	if id, err := c.Identifier(globals); err == nil {
		c.cmn.InstanceKey = &id
	}
	return nil
}

func (c *configWrapper) Identifier(globals Globals) (string, error) {
	if c.cmn.InstanceKey != nil {
		return *c.cmn.InstanceKey, nil
	}
	return c.configInstanceFunc(globals.AgentIdentifier)
}

func (c *configWrapper) NewIntegration(logger log.Logger, globals Globals) (Integration, error) {
	return newIntegrationFromV1(c, logger, globals, c.newInstanceFunc)
}

func (c *configWrapper) Cfg() Config {
	return c
}

func (c *configWrapper) Name() string {
	return c.cfg.Name()
}

func (c *configWrapper) Common() common.MetricsConfig {
	return c.cmn
}

type newIntegration func(l log.Logger) (shared.Integration, error)

type configInstance func(agentKey string) (string, error)

func newConfigWrapper(cfg shared.Config, cmn common.MetricsConfig, ni newIntegration, ci configInstance) *configWrapper {
	return &configWrapper{
		cfg:                cfg,
		cmn:                cmn,
		configInstanceFunc: ci,
		newInstanceFunc:    ni,
	}
}
`)
	if err != nil {
		panic(err)
	}
	v2ConfigBuilder := strings.Builder{}
	v2ConfigBuilder.WriteString("package v2 //nolint:golint\n")
	v2ConfigBuilder.WriteString(`
import (
"context"
	"errors"
	"github.com/go-kit/log"
"github.com/grafana/agent/pkg/integrations/shared"
	"github.com/grafana/agent/pkg/integrations/v2/common"

`)
	for _, i := range configs {
		v2ConfigBuilder.WriteString(fmt.Sprintf("\"%s\"\n", i.PackagePath))
	}
	v2ConfigBuilder.WriteString(")\n")
	v1Buffer := bytes.Buffer{}
	err = v2template.Execute(&v1Buffer, configs)
	if err != nil {
		panic(err)
	}

	v2ConfigBuilder.WriteString(v1Buffer.String())

	for _, cfg := range configs {
		bf := bytes.Buffer{}
		err = integrationTemplate.Execute(&bf, cfg)
		if err != nil {
			panic(err)
		}
		v2ConfigBuilder.WriteString(bf.String())
	}
	v2ConfigBuilder.WriteString(`
func newIntegrationFromV1(c IntegrationConfig, logger log.Logger, globals Globals, newInt func(l log.Logger) (shared.Integration, error)) (Integration, error) {

	v1Integration, err := newInt(logger)
	if err != nil {
		return nil, err
	}

	id, err := c.Cfg().Identifier(globals)
	if err != nil {
		return nil, err
	}

	// Generate our handler. Original integrations didn't accept a prefix, and
	// just assumed that they would be wired to /metrics somewhere.
	handler, err := v1Integration.MetricsHandler()
	if err != nil {
		return nil, fmt.Errorf("generating http handler: %w", err)
	} else if handler == nil {
		handler = http.NotFoundHandler()
	}

	// Generate targets. Original integrations used a static set of targets,
	// so this mapping can always be generated just once.
	//
	// Targets are generated from the result of ScrapeConfigs(), which returns a
	// tuple of job name and relative metrics path.
	//
	// Job names were prefixed at the subsystem level with integrations/, so we
	// will retain that behavior here.
	v1ScrapeConfigs := v1Integration.ScrapeConfigs()
	targets := make([]handlerTarget, 0, len(v1ScrapeConfigs))
	for _, sc := range v1ScrapeConfigs {
		targets = append(targets, handlerTarget{
			MetricsPath: sc.MetricsPath,
			Labels: model.LabelSet{
				model.JobLabel: model.LabelValue("integrations/" + sc.JobName),
			},
		})
	}

	// Convert he run function. Original integrations sometimes returned
	// ctx.Err() on exit. This isn't recommended anymore, but we need to hide the
	// error if it happens, since the error was previously ignored.
	runFunc := func(ctx context.Context) error {
		err := v1Integration.Run(ctx)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.Canceled) && ctx.Err() != nil:
			// Hide error that no longer happens in newer integrations.
			return nil
		default:
			return err
		}
	}

	// Aggregate our converted settings into a v2 integration.
	return &metricsHandlerIntegration{
		integrationName: c.Cfg().Name(),
		instanceID:      id,

		common:  c.Common(),
		globals: globals,
		handler: handler,
		targets: targets,

		runFunc: runFunc,
	}, nil
}
`)
	return v2ConfigBuilder.String()
}

type configMeta struct {
	Name          string
	ConfigStruct  string
	DefaultConfig string
	PackageName   string
	Type          shared.Type
	IsNativeV2    bool
	PackagePath   string
}

func newConfigMeta(c interface{}, defaultConfig interface{}, p shared.Type) configMeta {
	path := reflect.TypeOf(c).PkgPath()
	configType := fmt.Sprintf("%T", c)
	packageName := strings.ReplaceAll(configType, ".Config", "")
	name := strings.ReplaceAll(configType, ".Config", "")
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.Title(name)
	name = strings.ReplaceAll(name, " ", "")
	var dc string
	if defaultConfig != nil {
		dc = fmt.Sprintf("%T", defaultConfig)
	}

	return configMeta{
		Name:          name,
		ConfigStruct:  configType,
		DefaultConfig: dc,
		PackageName:   packageName,
		Type:          p,
		IsNativeV2:    false,
		PackagePath:   path,
	}
}
