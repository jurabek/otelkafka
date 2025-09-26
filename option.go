package otelkafka

import (
	"net"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const defaultTracerName = "go.opentelemetry.io/contrib/instrumentation/github.com/confluentinc/confluent-kafka-go/otelkafka"

type config struct {
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider

	Propagators      propagation.TextMapPropagator
	Tracer           trace.Tracer
	consumerGroupID  string
	bootstrapServers string

	attributeInjectFunc func(msg *kafka.Message) []attribute.KeyValue
	messageCarrierFunc  func(*kafka.Message) propagation.TextMapCarrier
}

// newConfig returns a config with all Options set.
func newConfig(opts ...Option) config {
	cfg := config{
		Propagators:    otel.GetTextMapPropagator(),
		TracerProvider: otel.GetTracerProvider(),
		MeterProvider:  otel.GetMeterProvider(),
		messageCarrierFunc: func(msg *kafka.Message) propagation.TextMapCarrier {
			return NewMessageCarrier(msg)
		},
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

func WithMeterProvider(meterProvider metric.MeterProvider) Option {
	return optionFunc(func(cfg *config) {
		if meterProvider != nil {
			cfg.MeterProvider = meterProvider
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

func WithMessageCarrier(carrier func(*kafka.Message) propagation.TextMapCarrier) Option {
	return optionFunc(func(cfg *config) {
		cfg.messageCarrierFunc = carrier
	})
}

func WithCustomAttributeInjector(fn func(msg *kafka.Message) []attribute.KeyValue) Option {
	return optionFunc(func(cfg *config) {
		cfg.attributeInjectFunc = fn
	})
}

// withConfig extracts the config information from kafka.ConfigMap for the client
func withConfig(cg *kafka.ConfigMap) Option {
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
