package prompt

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/peetya/snipforge-cli/data"
	"math/rand"
)

type languageInputModel struct {
	Input textinput.Model
}

func newLanguageInputModel() languageInputModel {
	m := languageInputModel{}
	m.Init()
	return m
}

func (m *languageInputModel) Init() {
	m.Input = textinput.New()
	m.Input.Placeholder = m.placeholder()
	m.Input.Validate = func(s string) error {
		if s == "" {
			return errors.New("language cannot be empty")
		}

		return nil
	}
}

func (m *languageInputModel) Update(pm *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return pm, tea.Quit

		case tea.KeyEnter:
			if pm.languageInputModel.Input.Err == nil {
				pm.detectedLanguage = data.DetectLanguage(pm.languageInputModel.Input.Value())
				pm.languageVersionInputModel.LanguageValue = pm.languageInputModel.Input.Value()
				pm.languageVersionInputModel.DetectedLanguage = pm.detectedLanguage
				pm.outputPathInputModel.DetectedLanguage = pm.detectedLanguage
				pm.navigator.Next(pm)
			}
		}
	}

	return pm, cmd
}

func (m *languageInputModel) View() string {
	err := ""
	errStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	if m.Input.Err != nil {
		err = fmt.Sprintf(errStyle.Render("Error: %s"), m.Input.Err)
	}
	return fmt.Sprintf(
		"Which programming or tooling language do you want to use?\n\n%s\n%s\n%s",
		m.Input.View(),
		err,
		"(enter to continue) (ctrl+c to quit)",
	)
}

func (m *languageInputModel) placeholder() string {
	n := rand.Intn(len(data.Languages))
	return fmt.Sprintf("e.g. %s", data.Languages[n].Names[0])
}
