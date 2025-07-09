package step

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"studyConductor/pkg"
	"syscall"
)

const TERMINAL_CMD = "gnome-terminal"

type BinaryStep struct {
	*AbstractStep
	cmd *exec.Cmd
}

func (s *BinaryStep) Activate(ctx context.Context) error {
	s.cmd = getCmd(s.config)
	err := s.cmd.Start()
	if err != nil {
		return err
	}
	if cv, exitst := s.config["oneshot"]; exitst && !cv.(bool) {
		s.state = StateRunning
	}
	return nil
}

func getCmd(cfg pkg.ConfigMap) *exec.Cmd {
	result := &exec.Cmd{}
	params := make([]string, 0)
	if cfg["BINARY_ARGS"] != nil {
		params = pkg.TypedSlice[string](cfg["BINARY_ARGS"].([]any))
	}
	cmdStr := cfg.GetString("BINARY_PATH")
	if terminal, hasKey := cfg["terminal"]; hasKey && terminal.(bool) {
		params = append([]string{"--", cmdStr}, params...)
		cmdStr = TERMINAL_CMD
	}
	if nohup, hasKey := cfg["nohup"]; hasKey && nohup.(bool) {
		params = append([]string{cmdStr}, params...)
		cmdStr = "nohup"
	}
	result = exec.Command(cmdStr, params...)
	if wd, hasWd := cfg["WORKING_DIR"]; hasWd {
		result.Dir = wd.(string)
	}
	return result
}

func (s *BinaryStep) Deactivate(ctx context.Context) error {
	err := syscall.Kill(s.cmd.Process.Pid, syscall.SIGTERM)
	if err != nil {
		return err
	}
	s.state = StatePending
	go func() {
		err = s.cmd.Wait()
		s.state = StateOff
	}()
	return nil
}

func (s *BinaryStep) Status(ctx context.Context) (*Status, error) {
	if cv, exist := s.config["oneshot"]; exist && cv.(bool) {
		s.state = StateOff
		return &Status{
			Label: SymbolPlayButton,
			Color: color.RGBA{G: uint8(255)},
		}, nil
	}
	if s.cmd == nil {
		pid, err := getPid(s.config["BINARY_PATH"].(string))
		if err != nil {
			return nil, err
		}
		// PID found, soft-attach process
		if pid > 0 {
			s.cmd = getCmd(s.config)
			s.cmd.Process, err = os.FindProcess(pid)
			if err != nil {
				return nil, err
			}
			s.state = StateRunning
		}
	}
	switch s.state {
	case StatePending:
		s.state = StatePending
		return &Status{
			Label: SymbolWorking,
			Color: color.RGBA{B: uint8(255)},
		}, nil
	case StateRunning:
		return &Status{
			Label: SymbolCheckmark,
			Color: color.RGBA{G: uint8(255)},
		}, nil
	case StatePaused:
		return &Status{
			Label: SymbolPause,
			Color: color.RGBA{R: uint8(255), G: uint8(255)},
		}, nil
	case StateOff:
		return &Status{
			Label: SymbolCross,
			Color: color.RGBA{R: uint8(255)},
		}, nil
	}
	return nil, fmt.Errorf("invalid state: %s", s.state)
}

func getPid(name string) (int, error) {
	procList, err := os.ReadDir("/proc")
	if err != nil {
		return 0, err
	}
	for _, pid := range procList {
		procName, _ := os.ReadFile("/proc/" + pid.Name() + "/cmdline")
		if strings.Contains(string(procName), name) {
			return strconv.Atoi(pid.Name())
		}
	}
	return 0, nil
}
