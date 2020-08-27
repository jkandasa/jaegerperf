package main

import (
	"jaegerperf/cmd/handler"
	qapi "jaegerperf/pkg/api/query"
	ml "jaegerperf/pkg/model"
	"jaegerperf/pkg/util"

	cp "github.com/otiai10/copy"

	"go.uber.org/zap"
)

func init() {
	// init logger
	logger := util.GetLogger("debug", "console", false, 0)
	zap.ReplaceGlobals(logger)

	// load resources
	loadResources()

	// load tags
	err := qapi.LoadTags()
	if err != nil {
		zap.L().Error("Failed to init job data", zap.Error(err))
	}
}

func main() {
	err := handler.StartHandler()
	if err != nil {
		zap.L().Fatal("Fatal", zap.Error(err))
	}
}

// loadResources loads on first run
func loadResources() error {
	if !util.IsDirExists(ml.FileTemplatesGenerator) {
		err := cp.Copy(ml.FileResources, ml.FileBase)
		if err != nil {
			zap.L().Error("Failed to copy default resources", zap.Error(err))
		}
	}
	return nil
}
