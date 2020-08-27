package generator

import (
	"errors"
	"fmt"
	"io"
	gml "jaegerperf/pkg/model/generator"
	jml "jaegerperf/pkg/model/job"
	"jaegerperf/pkg/util"
	"math/rand"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	ot "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var (
	gJob *jml.Job

	tags = []ot.Tag{
		{Key: "created by", Value: "golang qe perf automation"},
		{Key: "name1", Value: "foo"},
		{Key: "garbage string", Value: "sjdacdsakjcsadcds"},
	}

	logs = []otlog.Field{
		otlog.String("event", "soft error"),
		otlog.String("type", "cache timeout"),
		otlog.Int("waited.millis", 1500),
	}
)

// IsRunning job status
func IsRunning() bool {
	if gJob != nil {
		return gJob.IsRunning()
	}
	return false
}

// ExecuteSpansGenerator to dump data
func ExecuteSpansGenerator(jobID string, input gml.InputConfig) error {
	if IsRunning() {
		return errors.New("A generator job is in running state")
	}
	gJob = jml.New(jobID, jml.JobTypeSpansGenerator)
	report := gml.ExecutionReport{Input: input}
	// update job status on exit
	defer gJob.Update(false, &report)

	execStart := time.Now()

	err := execute(&report)
	if err != nil {
		report.Status.IsSuccess = false
		report.Status.Message = err.Error()
	} else {
		report.Status.IsSuccess = true
	}
	// update timetaken
	report.Status.StartTime = execStart
	report.Status.EndTime = time.Now()
	report.Status.TimeTaken = report.Status.EndTime.Sub(report.Status.StartTime).String()

	return err
}

func execute(report *gml.ExecutionReport) error {
	cfg := report.Input
	// reset tags and logs
	tags = []ot.Tag{
		{Key: "created by", Value: "golang qe perf automation"},
		{Key: "name1", Value: "foo"},
		{Key: "garbage string", Value: "sjdacdsakjcsadcds"},
	}

	logs = []otlog.Field{
		otlog.String("event", "soft error"),
		otlog.String("type", "cache timeout"),
		otlog.Int("waited.millis", 1500),
	}

	r := &jaegercfg.ReporterConfig{
		LogSpans: false,
	}
	if strings.ToLower(cfg.Endpoint.Type) == gml.EndpointTypeAgent {
		r.LocalAgentHostPort = cfg.Endpoint.URL
	} else if strings.ToLower(cfg.Endpoint.Type) == gml.EndpointTypeCollector {
		r = &jaegercfg.ReporterConfig{
			LogSpans:          false,
			CollectorEndpoint: "http://localhost:14268/api/traces",
		}
		if cfg.Endpoint.URL != "" {
			r.CollectorEndpoint = cfg.Endpoint.URL
		}
	} else {
		return fmt.Errorf("Invalid endpoint type: %s", cfg.Endpoint.Type)
	}

	if cfg.SpansConfig.ServiceName == "" {
		cfg.SpansConfig.ServiceName = "jaegerPerfTool_generated"
	}

	conf := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter:    r,
		ServiceName: cfg.SpansConfig.ServiceName,
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

	if cfg.SpansConfig.SpansPerSecond == 0 {
		return fmt.Errorf("Spans per second cannot be 0")
	}

	if cfg.SpansConfig.SpansPerSecond < cfg.SpansConfig.ChildDepth {
		return fmt.Errorf("SpansPerSecond[%d] can not be lesser than the ChildDepth[%d]", cfg.SpansConfig.SpansPerSecond, cfg.SpansConfig.ChildDepth)
	}

	// Update spans per second
	_rootSpans := cfg.SpansConfig.SpansPerSecond / (cfg.SpansConfig.ChildDepth + 1)
	cfg.SpansConfig.SpansPerSecond = _rootSpans * (cfg.SpansConfig.ChildDepth + 1)

	if cfg.SpansConfig.Tags != nil {
		for k, v := range cfg.SpansConfig.Tags {
			tags = append(tags, ot.Tag{Key: k, Value: v})
		}
	}

	if cfg.Mode == "" {
		return fmt.Errorf("Mode can not be empty")
	}

	switch strings.ToLower(cfg.Mode) {
	case "realtime":
		return generateRealtime(report, tracer)
	case "history":
		return generateHistory(report, tracer)
	default:
		return fmt.Errorf("Invalid mode:%s", cfg.Mode)
	}
}

func generateRealtime(report *gml.ExecutionReport, tracer ot.Tracer) error {
	cfg := report.Input
	if cfg.Realtime.Duration == "" {
		return fmt.Errorf("Duration can not be empty")
	}
	cfg.StartTime = time.Now()
	d, err := time.ParseDuration(cfg.Realtime.Duration)
	if err != nil {
		return err
	}
	endTime := cfg.StartTime.Add(d)
	ticker := time.NewTicker(1 * time.Second)
	spansSent := 0
	for {
		if endTime.Before(time.Now()) {
			break
		}
		_, ss := sendSpans(time.Now(), cfg.SpansConfig.SpansPerSecond, cfg.SpansConfig.ChildDepth, tracer)
		spansSent += ss
		<-ticker.C
	}
	ticker.Stop()

	// update spans sent
	report.Report.SpansSent = spansSent
	return nil
}

func generateHistory(report *gml.ExecutionReport, tracer ot.Tracer) error {
	cfg := report.Input
	if cfg.StartTime.IsZero() {
		cfg.StartTime = time.Now().Add(time.Duration(-2 * time.Hour))
	}

	if cfg.History.Days == 0 {
		return fmt.Errorf("Days can not be 0")
	}

	dDay := 24 * time.Hour

	if cfg.History.SpansPerDay == 0 {
		return fmt.Errorf("Spans per day can not be 0")
	}

	startTime := cfg.StartTime
	spansSent := 0
	for day := 1; day <= cfg.History.Days; day++ {
		totalSpans := cfg.History.SpansPerDay
		loopCount := totalSpans / cfg.SpansConfig.SpansPerSecond
		balanceCount := totalSpans % cfg.SpansConfig.SpansPerSecond
		ticker := time.NewTicker(1 * time.Second)
		sTime := startTime
		for count := 0; count < loopCount; count++ {
			var _spansSent int
			sTime, _spansSent = sendSpans(sTime, cfg.SpansConfig.SpansPerSecond, cfg.SpansConfig.ChildDepth, tracer)
			spansSent += _spansSent
			<-ticker.C
		}
		ticker.Stop()
		if balanceCount > 0 {
			var ss int
			sTime, ss = sendSpans(sTime, balanceCount, cfg.SpansConfig.ChildDepth, tracer)
			spansSent += ss
		}
		startTime = startTime.Add(time.Duration(-1 * dDay.Nanoseconds()))
	}

	// Update spans sent
	report.Report.SpansSent = spansSent
	return nil
}

func updateTags(s ot.Span, tags []ot.Tag) ot.Span {
	for _, t := range tags {
		s = s.SetTag(t.Key, t.Value)
	}
	return s
}

func sendSpans(startTime time.Time, spansCount, childDepth int, tracer ot.Tracer) (time.Time, int) {
	if spansCount == 0 {
		return startTime, 0
	}
	rootSpansCount := spansCount / (childDepth + 1)
	randMaxDuration := 950000 / rootSpansCount // 0.95 second(= 950000 microseconds) / rootSpansCount

	for count := 0; count < rootSpansCount; count++ {
		parentSpan := tracer.StartSpan("parent span", ot.StartTime(startTime))

		// child span
		var childSpan ot.Span

		// get span's context
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

		startTime = startTime.Add(rDuration)
		parentSpan.LogFields(logs...)
		updateTags(parentSpan, tags)
		parentSpan = parentSpan.SetTag("span type", "parent")
		parentSpan.FinishWithOptions(ot.FinishOptions{FinishTime: startTime})
	}
	// total spans sent
	totalSpansSent := rootSpansCount * (childDepth + 1)
	return startTime, totalSpansSent
}

// NewTracer return new jaeger tracer with given configuration
func NewTracer(cfg jaegercfg.Configuration) (opentracing.Tracer, error) {
	tracer, _, err := cfg.NewTracer(jaegercfg.Logger(util.NewCustomLogger("JAEGER_CLIENT", "info")))
	return tracer, err
}
