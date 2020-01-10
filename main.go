package main

import (
	jp "jaegerperf/jaegerperf"
)

func main() {
	err := jp.StartHandler()
	if err != nil {
		panic(err)
	}
}

