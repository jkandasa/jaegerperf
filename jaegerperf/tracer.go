package jaegerperf

import (
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

var tracers = make(map[string]opentracing.Tracer)

// NewTracer return new jaeger tracer with given configuration
func NewTracer(name string, cfg jaegercfg.Configuration) opentracing.Tracer {
	tracer, _, err := cfg.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil
	}
	tracers[name] = tracer
	return tracer
}

// GetTracer returns a tracer by name
func GetTracer(name string) opentracing.Tracer {
	return tracers[name]
}

// CloseTracers all the tracers
func CloseTracers() {
	for k, v := range tracers {
		v.(io.Closer).Close()
		delete(tracers, k)
	}
}

// CloseTracer the specified tracer
func CloseTracer(name string) {
	t, ok := tracers[name]
	if ok {
		t.(io.Closer).Close()
		delete(tracers, name)
	}
}
