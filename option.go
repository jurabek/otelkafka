package otelkafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net"
	"strings"
)

const defaultTracerName = "go.opentelemetry.io/contrib/instrumentation/github.com/confluentinc/confluent-kafka-go/otelkafka"

type config struct {
	TracerProvider   trace.TracerProvider
	Propagators      propagation.TextMapPropagator
	Tracer           trace.Tracer
	consumerGroupID  string
	bootstrapServers string
}

// newConfig returns a config with all Options set.
func newConfig(opts ...Option) config {
	cfg := config{
		Propagators:    otel.GetTextMapPropagator(),
		TracerProvider: otel.GetTracerProvider(),
	}
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	cfg.Tracer = cfg.TracerProvider.Tracer(
		defaultTracerName,
		trace.WithInstrumentationVersion(Version()),
	)

	return cfg
}

// Option interface used for setting optional config properties.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (fn optionFunc) apply(c *config) {
	fn(c)
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.TracerProvider = provider
		}
	})
}

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(cfg *config) {
		if propagators != nil {
			cfg.Propagators = propagators
		}
	})
}

// WithConfig extracts the config information from kafka.ConfigMap for the client
func WithConfig(cg *kafka.ConfigMap) Option {
	return optionFunc(func(cfg *config) {
		if groupID, err := cg.Get("group.id", ""); err == nil {
			cfg.consumerGroupID = groupID.(string)
		}
		if bs, err := cg.Get("bootstrap.servers", ""); err == nil && bs != "" {
			for _, addr := range strings.Split(bs.(string), ",") {
				host, _, err := net.SplitHostPort(addr)
				if err == nil {
					cfg.bootstrapServers = host
					return
				}
			}
		}
	})
}
