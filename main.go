package main

import (
	"fmt"
	jp "jaegerperf/jaegerperf"
)

func init() {
	err := jp.InitJobData()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	err := jp.StartHandler()
	if err != nil {
		panic(err)
	}
}
