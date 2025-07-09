package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"studyConductor/pkg"
	"studyConductor/step"
)

type StudyModel struct {
	Context         context.Context
	steps           []step.Step
	cursor          int
	selected        map[int]struct{}
	gestureDetected bool
}

func (m *StudyModel) Init() tea.Cmd {
	return nil
}

func (m *StudyModel) ToggleGesture() error {
	m.gestureDetected = !m.gestureDetected
	gesture := make(map[string]bool)
	gesture["g"] = m.gestureDetected
	gestureJson, err := json.Marshal(gesture)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8398/override-gestures", bytes.NewReader(gestureJson))
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}

func (m *StudyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ctx := context.Background()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "g":
			if err := m.ToggleGesture(); err != nil {
				panic(any(err))
			}
			return m, nil
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "w":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "s":
			if m.cursor < len(m.steps)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.steps[m.cursor].State() == step.StateOff || m.steps[m.cursor].State() == "" {
				go func() {
					err := m.steps[m.cursor].Activate(ctx)
					if err != nil {
						panic(any(err))
					}
				}()
				m.selected[m.cursor] = struct{}{}
			} else {
				go func() {
					_ = m.steps[m.cursor].Deactivate(ctx)
				}()
				delete(m.selected, m.cursor)
			}
		}
	}

	return m, nil
}

func (m *StudyModel) View() string {
	s := "Steps to prepare for the study\n\n"

	for i, stp := range m.steps {
		cursor := " "
		if m.cursor == i {
			cursor = "⇨"
		}

		check := " "
		//		if _, ok := m.selected[i]; ok {
		status, err := stp.Status(m.Context)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		check = status.Label
		//		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, check, stp.Label())
	}
	s += "\nGesture: "
	if m.gestureDetected {
		s += color.HiGreenString("[Active]") + "\n"
	} else {
		s += "[None]\n"
	}
	s += "\nPress q to quit.\n"
	return s
}

func initModel(config *pkg.Config) (*StudyModel, error) {
	return NewStudyModel(config)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := pkg.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config")
	}
	usbDrive, err := pkg.GetStorage(cfg)
	if err != nil {
		log.Fatal("Error loading USB drive:", err)
	}
	step.Conf = cfg
	step.Conf.Study.Storage["path"] = usbDrive
	mdl, err := initModel(cfg)
	if err != nil {
		log.Fatal("Error initializing study model")
	}
	p := tea.NewProgram(mdl)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
