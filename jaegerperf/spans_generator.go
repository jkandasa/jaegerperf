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
)

const spansGenerator = "SpansGenerator"

var (
	gJob = Job{js: JobStatus{}, jobType: spansGenerator}

	tags = []ot.Tag{
		ot.Tag{Key: "created by", Value: "golang qe perf automation"},
		ot.Tag{Key: "name1", Value: "foo"},
		ot.Tag{Key: "garbage string", Value: "sjdacdsakjcsadcds"},
	}

	logs = []otlog.Field{
		otlog.String("event", "soft error"),
		otlog.String("type", "cache timeout"),
		otlog.Int("waited.millis", 1500),
	}
)

// GeneratorConfiguration to send spans
type GeneratorConfiguration struct {
	Target            string                 `yaml:"target" json:"target"`
	Endpoint          string                 `yaml:"endpoint" json:"endpoint"`
	ServiceName       string                 `yaml:"serviceName" json:"serviceName"`
	Mode              string                 `yaml:"mode" json:"mode"`
	ExecutionDuration string                 `yaml:"executionDuration" json:"executionDuration"`
	NumberOfDays      int                    `yaml:"numberOfDays" json:"numberOfDays"`
	SpansPerDay       int                    `yaml:"spansPerDay" json:"spansPerDay"`
	SpansPerSecond    int                    `yaml:"spansPerSecond" json:"spansPerSecond"`
	ChildDepth        int                    `yaml:"childDepth" json:"childDepth"`
	Tags              map[string]interface{} `yaml:"tags" json:"tags"`
	StartTime         time.Time              `yaml:"startTime" json:"startTime"`
}

// IsGeneratorRunning job status
func IsGeneratorRunning() bool {
	return gJob.IsRunning()
}

// NewTracer return new jaeger tracer with given configuration
func NewTracer(cfg jaegercfg.Configuration) (opentracing.Tracer, error) {
	tracer, _, err := cfg.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	return tracer, err
}

// ExecuteSpansGenerator to dump data
func ExecuteSpansGenerator(jobID string, cfg GeneratorConfiguration) error {
	gJob.SetStatus(true, jobID, cfg)
	defer gJob.SetCompleted()

	// update job configuration
	gJob.SetStatus(true, jobID, cfg)
	err := execute(&cfg)
	if err != nil {
		// store job result
	}
	return err
}

func execute(cfg *GeneratorConfiguration) error {
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

	if cfg.ServiceName == "" {
		cfg.ServiceName = "jaegerPerfTool_generated"
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

	if cfg.SpansPerSecond == 0 {
		cfg.SpansPerSecond = 500
	}

	if cfg.Tags != nil {
		for k, v := range cfg.Tags {
			tags = append(tags, ot.Tag{Key: k, Value: v})
		}
	}

	if cfg.Mode == "" {
		cfg.Mode = "realtime"
	}

	switch strings.ToLower(cfg.Mode) {
	case "realtime":
		return generateRealtime(cfg, tracer)
	case "history":
		return generateHistory(cfg, tracer)
	default:
		return fmt.Errorf("Invalid mode:%s", cfg.Mode)
	}
}

func generateRealtime(cfg *GeneratorConfiguration, tracer ot.Tracer) error {
	if cfg.ExecutionDuration == "" {
		cfg.ExecutionDuration = "5m"
	}
	cfg.StartTime = time.Now()
	d, err := time.ParseDuration(cfg.ExecutionDuration)
	if err != nil {
		return err
	}
	endTime := cfg.StartTime.Add(d)
	ticker := time.NewTicker(1 * time.Second)
	for {
		if endTime.Before(time.Now()) {
			break
		}
		sendSpans(time.Now(), cfg.SpansPerSecond, cfg.ChildDepth, tracer)
		<-ticker.C
	}
	ticker.Stop()

	return nil
}

func generateHistory(cfg *GeneratorConfiguration, tracer ot.Tracer) error {
	if cfg.StartTime.IsZero() {
		cfg.StartTime = time.Now().Add(time.Duration(-2 * time.Hour))
	}

	if cfg.NumberOfDays == 0 {
		cfg.NumberOfDays = 1
	}

	// limit maximum number of days in 2 years
	// There is no reason to limit the days
	if cfg.NumberOfDays > 730 {
		cfg.NumberOfDays = 730
	}

	dDay := 24 * time.Hour

	//if cfg.NumberOfDays > 1 {
	//	cfg.StartTime = cfg.StartTime.Add(time.Duration(-1 * dDay.Nanoseconds() * int64(cfg.NumberOfDays)))
	//}

	if cfg.SpansPerDay == 0 {
		cfg.SpansPerDay = 10
	}

	startTime := cfg.StartTime
	for day := 1; day <= cfg.NumberOfDays; day++ {
		totalSpans := cfg.SpansPerDay
		loopCount := totalSpans / cfg.SpansPerSecond
		balanceCount := totalSpans % cfg.SpansPerSecond
		ticker := time.NewTicker(1 * time.Second)
		sTime := startTime
		for count := 0; count < loopCount; count++ {
			sTime = sendSpans(sTime, cfg.SpansPerSecond, cfg.ChildDepth, tracer)
			<-ticker.C
		}
		ticker.Stop()
		if balanceCount > 0 {
			sTime = sendSpans(sTime, balanceCount, cfg.ChildDepth, tracer)
		}
		startTime = startTime.Add(time.Duration(-1 * dDay.Nanoseconds()))
	}
	return nil
}

func updateTags(s ot.Span, tags []ot.Tag) ot.Span {
	for _, t := range tags {
		s = s.SetTag(t.Key, t.Value)
	}
	return s
}

func sendSpans(startTime time.Time, spansCount, childDepth int, tracer ot.Tracer) time.Time {
	if spansCount == 0 {
		return startTime
	}
	parentSpansCount := spansCount / (childDepth + 1)
	randMaxDuration := 950000 / parentSpansCount // 0.95 second(= 950000 microseconds) / parentSpansCount
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
		rDuration := time.Duration(1000 * (rand.Intn(randMaxDuration)))
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
			time.Sleep(2 * time.Nanosecond)
			depth++
		}

		spansDone = spansDone + childDepth + 1
		startTime = startTime.Add(rDuration)
		parentSpan.LogFields(logs...)
		updateTags(parentSpan, tags)
		parentSpan = parentSpan.SetTag("span type", "parent")
		parentSpan.FinishWithOptions(ot.FinishOptions{FinishTime: startTime})
	}
	return startTime
}
