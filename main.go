package main

import (
	jp "jaegerperf/jaegerperf"
	"time"
)

func main() {
	defer func() { jp.CloseTracers() }()
	jp.Execute(jp.GeneratorConfiguration{
		NumberOfDays:    2,
		SpansPerDay:     20,
		SpansChildDepth: 20,
		Tags:            nil,
		StartTime:       time.Now(),
	})
}
