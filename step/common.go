package step

import (
	"context"
	"image/color"
	"studyConductor/pkg"
	"sync/atomic"
)

const (
	SymbolCross      = "❌"
	SymbolCheckmark  = "✓"
	SymbolPause      = "⏸"
	SymbolWorking    = "⏳"
	SymbolPlayButton = "▶"
)

var Conf *pkg.Config

type Status struct {
	Label string
	Color color.Color
}

type Step interface {
	Activate(ctx context.Context) error
	Deactivate(ctx context.Context) error
	Status(ctx context.Context) (*Status, error)
	Label() string
	State() StepState
}

type AbstractStep struct {
	label     string
	config    pkg.ConfigMap
	state     StepState
	workingMx atomic.Bool
}

type StepState string

const (
	StateOff     StepState = "off"
	StatePending StepState = "pending"
	StatePaused  StepState = "paused"
	StateRunning StepState = "running"
)

func (s *AbstractStep) State() StepState {
	return s.state
}

func (s *AbstractStep) Label() string {
	return s.label
}

func (s *AbstractStep) String() string {
	var state string
	switch s.state {
	case StateRunning:
		state = "Stop"
	case StatePending:
		state = "Wait..."
	case StateOff:
	default:
		state = "Start"
	}
	return state + " " + s.label
}

func BuildStep(module *pkg.Module) (Step, error) {
	var s Step
	aStep := &AbstractStep{
		label:     module.Name,
		config:    module.Configuration,
		workingMx: atomic.Bool{},
	}
	switch module.Type {
	case pkg.ModuleTypeBinary:
		aStep.state = StateOff
		s = &BinaryStep{aStep, nil}
	case pkg.ModuleTypeDocker:
		s = &ContainerStep{aStep, nil, nil}
	}
	return s, nil
}
