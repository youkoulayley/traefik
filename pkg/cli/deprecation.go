package cli

import (
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type rawConfiguration map[string]interface{}

// deprecationNotice prints warns and hints if deprecated/removed static option are in use,
// and returns whether at least one of these options is incompatible with the current version.
func (r *rawConfiguration) deprecationNotice(logger zerolog.Logger) bool {
	if r == nil {
		return false
	}

	marshal, err := json.Marshal(r)
	if err != nil {
		log.Error().Err(err).Send()
	}
	config := &configuration{}

	err = json.Unmarshal(marshal, config)
	if err != nil {
		log.Error().Err(err).Send()
	}

	return config.deprecationNotice(logger)
}

// configuration holds the static configuration removed/deprecated options.
type configuration struct {
	Experimental *experimental `json:"experimental,omitempty" toml:"experimental,omitempty" yaml:"experimental,omitempty"`
	Pilot        *interface{}  `json:"pilot,omitempty" toml:"pilot,omitempty" yaml:"pilot,omitempty"`
	Providers    *providers    `json:"providers,omitempty" toml:"providers,omitempty" yaml:"providers,omitempty"`
	Tracing      *tracing      `json:"tracing,omitempty" toml:"tracing,omitempty" yaml:"tracing,omitempty"`
}

func (c *configuration) deprecationNotice(logger zerolog.Logger) bool {
	if c == nil {
		return false
	}

	var incompatible bool
	if c.Pilot != nil {
		incompatible = true
		logger.Error().Msg("Pilot configuration has been removed in v3, please remove all Marathon-related static configuration for Traefik to start." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	return incompatible ||
		c.Experimental.deprecationNotice(logger) ||
		c.Providers.deprecationNotice(logger) ||
		c.Tracing.deprecationNotice(logger)
}

type providers struct {
	Docker        *docker            `json:"docker,omitempty" toml:"docker,omitempty" yaml:"docker,omitempty"`
	Consul        *hashicorpProvider `json:"consul,omitempty" toml:"consul,omitempty" yaml:"consul,omitempty"`
	ConsulCatalog *hashicorpProvider `json:"consulCatalog,omitempty" toml:"consulCatalog,omitempty" yaml:"consulCatalog,omitempty"`
	Nomad         *hashicorpProvider `json:"nomad,omitempty" toml:"nomad,omitempty" yaml:"nomad,omitempty"`
	Marathon      *interface{}       `json:"marathon,omitempty" toml:"marathon,omitempty" yaml:"marathon,omitempty"`
	Rancher       *interface{}       `json:"rancher,omitempty" toml:"rancher,omitempty" yaml:"rancher,omitempty"`
}

type hashicorpProvider struct {
	Namespace *string `json:"namespace,omitempty" toml:"namespace,omitempty" yaml:"namespace,omitempty"`
}

func (p *providers) deprecationNotice(logger zerolog.Logger) bool {
	if p == nil {
		return false
	}

	var incompatible bool
	if p.Consul != nil && p.Consul.Namespace != nil {
		incompatible = true
		logger.Error().Msg("Consul provider `namespace` option has been removed, please use the `namespaces` option instead." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if p.ConsulCatalog != nil && p.ConsulCatalog.Namespace != nil {
		incompatible = true
		logger.Error().Msg("ConsulCatalog provider `namespace` option has been removed, please use the `namespaces` option instead." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if p.Nomad != nil && p.Nomad.Namespace != nil {
		incompatible = true
		logger.Error().Msg("ConsulCatalog provider `namespace` option has been removed, please use the `namespaces` option instead." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if p.Marathon != nil {
		incompatible = true
		logger.Error().Msg("Marathon provider has been removed in v3, please remove all Marathon-related static configuration for Traefik to start." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if p.Rancher != nil {
		incompatible = true
		logger.Error().Msg("Rancher provider has been removed in v3, please remove all Rancher-related static configuration for Traefik to start." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	return incompatible || p.Docker.deprecationNotice(logger)
}

type docker struct {
	SwarmMode *bool `json:"swarmMode,omitempty" toml:"swarmMode,omitempty" yaml:"swarmMode,omitempty"`
}

func (d *docker) deprecationNotice(logger zerolog.Logger) bool {
	if d == nil {
		return false
	}

	var incompatible bool

	if d.SwarmMode != nil {
		incompatible = true
		logger.Error().Msg("Docker provider `swarmMode` option has been removed in v3, please use the Swarm Provider instead." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	return incompatible
}

type experimental struct {
	HTTP3 *bool `json:"http3,omitempty" toml:"http3,omitempty" yaml:"http3,omitempty"`
}

func (e *experimental) deprecationNotice(logger zerolog.Logger) bool {
	if e == nil {
		return false
	}

	if e.HTTP3 != nil {
		logger.Warn().Msg("HTTP3 is not an experimental feature in v3 and the associated enablement option will be remove in a future major release." +
			"This option has now no effect, but we recommend to stop using it to avoid hassle considering future major release upgrade." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	// As long as HTTP3 is still a field of the experimental static configuration,
	// there will be no parsing error, so we just want to warn here.
	return false
}

type tracing struct {
	SpanNameLimit *int         `json:"spanNameLimit,omitempty" toml:"spanNameLimit,omitempty" yaml:"spanNameLimit,omitempty" export:"true"`
	Jaeger        *interface{} `json:"jaeger,omitempty" toml:"jaeger,omitempty" yaml:"jaeger,omitempty"`
	Zipkin        *interface{} `json:"zipkin,omitempty" toml:"zipkin,omitempty" yaml:"zipkin,omitempty"`
	Datadog       *interface{} `json:"datadog,omitempty" toml:"datadog,omitempty" yaml:"datadog,omitempty"`
	Instana       *interface{} `json:"instana,omitempty" toml:"instana,omitempty" yaml:"instana,omitempty"`
	Haystack      *interface{} `json:"haystack,omitempty" toml:"haystack,omitempty" yaml:"haystack,omitempty"`
	Elastic       *interface{} `json:"elastic,omitempty" toml:"elastic,omitempty" yaml:"elastic,omitempty"`
}

func (t *tracing) deprecationNotice(logger zerolog.Logger) bool {
	if t == nil {
		return false
	}
	var incompatible bool
	if t.SpanNameLimit != nil {
		incompatible = true
		logger.Error().Msg("SpanNameLimit option for Tracing has been removed in v3, as Span names are now of a fixed length." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if t.Jaeger != nil {
		incompatible = true
		logger.Error().Msg("Jaeger Tracing backend has been removed in v3, please remove all Jaeger-related Tracing static configuration for Traefik to start." +
			"In v3, Open Telemetry replaces specific tracing backend implementations, and an collector/exporter can be used to export metrics in a vendor specific format." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if t.Zipkin != nil {
		incompatible = true
		logger.Error().Msg("Zipkin Tracing backend has been removed in v3, please remove all Jaeger-related Tracing static configuration for Traefik to start." +
			"In v3, Open Telemetry replaces specific tracing backend implementations, and an collector/exporter can be used to export metrics in a vendor specific format." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if t.Datadog != nil {
		incompatible = true
		logger.Error().Msg("Datadog Tracing backend has been removed in v3, please remove all Jaeger-related Tracing static configuration for Traefik to start." +
			"In v3, Open Telemetry replaces specific tracing backend implementations, and an collector/exporter can be used to export metrics in a vendor specific format." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if t.Instana != nil {
		incompatible = true
		logger.Error().Msg("Instana Tracing backend has been removed in v3, please remove all Jaeger-related Tracing static configuration for Traefik to start." +
			"In v3, Open Telemetry replaces specific tracing backend implementations, and an collector/exporter can be used to export metrics in a vendor specific format." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if t.Haystack != nil {
		incompatible = true
		logger.Error().Msg("Haystack Tracing backend has been removed in v3, please remove all Haystack-related Tracing static configuration for Traefik to start." +
			"In v3, Open Telemetry replaces specific tracing backend implementations, and an collector/exporter can be used to export metrics in a vendor specific format." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	if t.Elastic != nil {
		incompatible = true
		logger.Error().Msg("Elastic Tracing backend has been removed in v3, please remove all Elastic-related Tracing static configuration for Traefik to start." +
			"In v3, Open Telemetry replaces specific tracing backend implementations, and an collector/exporter can be used to export metrics in a vendor specific format." +
			"For more information please read the migration guide: https://doc.traefik.io/traefik/v3.0/migration/v2-to-v3/#")
	}

	return incompatible
}
