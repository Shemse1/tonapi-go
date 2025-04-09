// Code generated by ogen, DO NOT EDIT.

package tonapi

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/otelogen"
)

var (
	// Allocate option closure once.
	clientSpanKind = trace.WithSpanKind(trace.SpanKindClient)
)

type (
	optionFunc[C any] func(*C)
	otelOptionFunc    func(*otelConfig)
)

type otelConfig struct {
	TracerProvider trace.TracerProvider
	Tracer         trace.Tracer
	MeterProvider  metric.MeterProvider
	Meter          metric.Meter
}

func (cfg *otelConfig) initOTEL() {
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	if cfg.MeterProvider == nil {
		cfg.MeterProvider = otel.GetMeterProvider()
	}
	cfg.Tracer = cfg.TracerProvider.Tracer(otelogen.Name,
		trace.WithInstrumentationVersion(otelogen.SemVersion()),
	)
	cfg.Meter = cfg.MeterProvider.Meter(otelogen.Name,
		metric.WithInstrumentationVersion(otelogen.SemVersion()),
	)
}

type clientConfig struct {
	otelConfig
	Client ht.Client
}

// ClientOption is client config option.
type ClientOption interface {
	applyClient(*clientConfig)
}

var _ ClientOption = (optionFunc[clientConfig])(nil)

func (o optionFunc[C]) applyClient(c *C) {
	o(c)
}

var _ ClientOption = (otelOptionFunc)(nil)

func (o otelOptionFunc) applyClient(c *clientConfig) {
	o(&c.otelConfig)
}

func newClientConfig(opts ...ClientOption) clientConfig {
	cfg := clientConfig{
		Client: http.DefaultClient,
	}
	for _, opt := range opts {
		opt.applyClient(&cfg)
	}
	cfg.initOTEL()
	return cfg
}

type baseClient struct {
	cfg      clientConfig
	requests metric.Int64Counter
	errors   metric.Int64Counter
	duration metric.Float64Histogram
}

func (cfg clientConfig) baseClient() (c baseClient, err error) {
	c = baseClient{cfg: cfg}
	if c.requests, err = otelogen.ClientRequestCountCounter(c.cfg.Meter); err != nil {
		return c, err
	}
	if c.errors, err = otelogen.ClientErrorsCountCounter(c.cfg.Meter); err != nil {
		return c, err
	}
	if c.duration, err = otelogen.ClientDurationHistogram(c.cfg.Meter); err != nil {
		return c, err
	}
	return c, nil
}

// Option is config option.
type Option interface {
	ClientOption
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
//
// If none is specified, the global provider is used.
func WithTracerProvider(provider trace.TracerProvider) Option {
	return otelOptionFunc(func(cfg *otelConfig) {
		if provider != nil {
			cfg.TracerProvider = provider
		}
	})
}

// WithMeterProvider specifies a meter provider to use for creating a meter.
//
// If none is specified, the otel.GetMeterProvider() is used.
func WithMeterProvider(provider metric.MeterProvider) Option {
	return otelOptionFunc(func(cfg *otelConfig) {
		if provider != nil {
			cfg.MeterProvider = provider
		}
	})
}

// WithClient specifies http client to use.
func WithClient(client ht.Client) ClientOption {
	return optionFunc[clientConfig](func(cfg *clientConfig) {
		if client != nil {
			cfg.Client = client
		}
	})
}
