package jaegerperf

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	ot "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"gopkg.in/yaml.v2"
)

var spansPerMinute = 500

var tags = []ot.Tag{
	ot.Tag{Key: "created by", Value: "golang qe perf automation"},
	ot.Tag{Key: "name1", Value: "foo"},
	ot.Tag{Key: "garbage string", Value: "sjdacdsakjcsadcds"},
}

var logs = []otlog.Field{
	otlog.String("event", "soft error"),
	otlog.String("type", "cache timeout"),
	otlog.Int("waited.millis", 1500),
}

// GeneratorConfiguration to send spans
type GeneratorConfiguration struct {
	Target       string                 `yaml:"target"`
	Endpoint     string                 `yaml:"endpoint"`
	ServiceName  string                 `yaml:"serviceName"`
	NumberOfDays int                    `yaml:"numberOfDays"`
	SpansPerDay  int                    `yaml:"spansPerDay"`
	ChildDepth   int                    `yaml:"childDepth"`
	Tags         map[string]interface{} `yaml:"tags"`
	StartTime    time.Time              `yaml:"startTime"`
}

// NewTracer return new jaeger tracer with given configuration
func NewTracer(cfg jaegercfg.Configuration) (opentracing.Tracer, error) {
	tracer, _, err := cfg.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	return tracer, err
}

// ExecuteSpansGenerator to dump data
func ExecuteSpansGenerator(config string) error {
	cfg := GeneratorConfiguration{}
	err := yaml.Unmarshal([]byte(config), &cfg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	dDay := 24 * time.Hour

	if cfg.NumberOfDays > 0 {
		cfg.StartTime = time.Now().Add(time.Duration(-1 * dDay.Nanoseconds() * int64(cfg.NumberOfDays)))
	}
	if cfg.Tags != nil {
		for k, v := range cfg.Tags {
			tags = append(tags, ot.Tag{Key: k, Value: v})
		}
	}
	err = execute(cfg)
	return err
}

func execute(cfg GeneratorConfiguration) error {
	r := &jaegercfg.ReporterConfig{
		LogSpans: false,
	}
	if cfg.Endpoint != "" {
		r.LocalAgentHostPort = cfg.Endpoint
	}
	if strings.EqualFold("collector", cfg.Target) {
		r = &jaegercfg.ReporterConfig{
			LogSpans:          false,
			CollectorEndpoint: "http://localhost:14268/api/traces",
		}
		if cfg.Endpoint != "" {
			r.CollectorEndpoint = cfg.Endpoint
		}
	}

	conf := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter:    r,
		ServiceName: cfg.ServiceName,
	}
	tracer, err := NewTracer(conf)
	if err != nil {
		return err
	}
	defer func() {
		if tracer != nil {
			tracer.(io.Closer).Close()
		}
	}()
	if cfg.StartTime.IsZero() {
		cfg.StartTime = time.Now().Add(-1 * time.Hour)
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
		ticker := time.NewTicker(1 * time.Second)
		for count := 0; count < loopCount; count++ {
			sendSpans(startTime, spansCount, cfg.ChildDepth, tracer)
			<-ticker.C
		}
		ticker.Stop()
		if balanceCount > 0 {
			sendSpans(startTime, balanceCount, cfg.ChildDepth, tracer)
		}
		startTime = startTime.Add(time.Hour * 24)
	}
	return nil
}

func updateTags(s ot.Span, tags []ot.Tag) ot.Span {
	for _, t := range tags {
		s = s.SetTag(t.Key, t.Value)
	}
	return s
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
		parentSpan := tracer.StartSpan("parent span", ot.StartTime(startTime))
		var childSpan ot.Span
		cSpan := func() ot.SpanContext {
			if childSpan != nil {
				return childSpan.Context()
			}
			return parentSpan.Context()
		}
		depth := 1
		sTime := startTime
		rand.Seed(time.Now().UnixNano())
		rDuration := time.Duration(1000000 * (rand.Intn(60000) + 500))
		var rcDuration time.Duration
		if childDepth > 0 {
			rcDuration = time.Duration(rDuration.Nanoseconds() / int64(childDepth))
		}
		for depth <= childDepth {
			childSpan = tracer.StartSpan(
				fmt.Sprintf("child_span_%d", depth),
				ot.StartTime(sTime),
				ot.ChildOf(cSpan()))
			sTime = sTime.Add(rcDuration)
			updateTags(childSpan, tags)
			childSpan = childSpan.SetTag("span type", fmt.Sprintf("child:%d", depth))
			childSpan.FinishWithOptions(ot.FinishOptions{FinishTime: startTime.Add(rDuration)})
			time.Sleep(2 * time.Millisecond)
			depth++
		}

		spansDone = spansDone + childDepth + 1
		startTime = startTime.Add(rDuration)
		parentSpan.LogFields(logs...)
		updateTags(parentSpan, tags)
		parentSpan = parentSpan.SetTag("span type", "parent")
		parentSpan.FinishWithOptions(ot.FinishOptions{FinishTime: startTime})
	}
}
