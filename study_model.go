package main

import (
	"context"
	"studyConductor/pkg"
	"studyConductor/step"
)

//[]Step{
//	databaseStep,
//	brokerStep,
//	surveyServerStep,
//	latinSquirtStep,
//}

func NewStudyModel(config *pkg.Config) (*StudyModel, error) {
	var err error
	mdl := &StudyModel{
		Context:  context.Background(),
		steps:    make([]step.Step, len(config.Modules)),
		selected: make(map[int]struct{}),
	}
	for i, module := range config.Modules {
		mdl.steps[i], err = step.BuildStep(&module)
		if err != nil {
			return nil, err
		}
	}
	return mdl, nil
}
