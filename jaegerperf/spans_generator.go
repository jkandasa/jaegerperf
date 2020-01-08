package jaegerperf

import (
	"fmt"
	"math/rand"
	"time"

	ot "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var spansPerMinute = 500

// GeneratorConfiguration to send spans
type GeneratorConfiguration struct {
	NumberOfDays    int
	SpansPerDay     int
	SpansChildDepth int
	Tags            map[string]interface{}
	StartTime       time.Time
}

// Execute spans generator
func Execute(cfg GeneratorConfiguration) {
	conf := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: false,
			//CollectorEndpoint: "http://hello.com:4455",
		},
		ServiceName: "service-1",
	}
	tracer := NewTracer("tracer_1", conf)
	if cfg.StartTime.IsZero() {
		cfg.StartTime = time.Now()
	}

	if cfg.SpansPerDay == 0 {
		cfg.SpansPerDay = 10
	}

	if cfg.NumberOfDays == 0 {
		cfg.NumberOfDays = 1
	}

	startTime := cfg.StartTime
	for day := 0; day < cfg.NumberOfDays; day++ {
		totalSpans := cfg.SpansPerDay
		loopCount := totalSpans / spansPerMinute
		var spansCount int
		balanceCount := totalSpans
		if loopCount > 0 {
			spansCount = totalSpans / loopCount
			balanceCount = totalSpans % loopCount
		}
		fmt.Printf("day:%d, loopCount:%d, spansCount:%d, balanceCount:%d\n",
			day, loopCount, spansCount, balanceCount)
		ticker := time.NewTicker(1 * time.Second)
		for count := 0; count < loopCount; count++ {
			sendSpans(startTime, spansCount, cfg.SpansChildDepth, tracer)
			<-ticker.C
		}
		ticker.Stop()
		if balanceCount > 0 {
			sendSpans(startTime, balanceCount, cfg.SpansChildDepth, tracer)
		}
		startTime = startTime.Add(time.Hour * 24)
	}
}

func sendSpans(startTime time.Time, spansCount, childDepth int, tracer ot.Tracer) {
	if spansCount == 0 {
		return
	}
	spansDone := 0
	for {
		if spansDone >= spansCount {
			break
		}
		parentSpan := tracer.StartSpan("test tracer", ot.StartTime(startTime))
		var childSpan ot.Span
		cSpan := func() ot.SpanContext {
			if childSpan != nil {
				return childSpan.Context()
			}
			return parentSpan.Context()
		}
		depth := 0
		for depth < childDepth {
			childSpan = ot.StartSpan(
				"child_span_"+string(depth),
				ot.ChildOf(cSpan()))
			depth++
		}

		spansDone = spansDone + childDepth + 1
		parentSpan.Finish()
		startTime = startTime.Add(time.Millisecond * time.Duration(rand.Intn(10000)+500))
	}

}
