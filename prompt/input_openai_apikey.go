package prompt

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/peetya/snipforge-cli/data"
)

type openaiApiKeyInputModel struct {
	Input            textinput.Model
	DetectedLanguage *data.Language
}

func newOpenaiApiKeyInputModel() openaiApiKeyInputModel {
	m := openaiApiKeyInputModel{}
	m.Init()
	return m
}

func (m *openaiApiKeyInputModel) Init() {
	m.Input = textinput.New()
	m.Input.EchoMode = textinput.EchoPassword
	m.Input.EchoCharacter = '•'
	m.Input.Validate = func(s string) error {
		if s == "" {
			return errors.New("api key cannot be empty")
		}

		return nil
	}
}

func (m *openaiApiKeyInputModel) Update(pm *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return pm, tea.Quit

		case tea.KeyCtrlT:
			if pm.openaiApiKeyInputModel.Input.EchoMode == textinput.EchoNormal {
				pm.openaiApiKeyInputModel.Input.EchoMode = textinput.EchoPassword
			} else {
				pm.openaiApiKeyInputModel.Input.EchoMode = textinput.EchoNormal
			}

		case tea.KeyEnter:
			if m.Input.Err == nil {
				pm.navigator.Next(pm)
			}
		}
	}

	return pm, cmd
}

func (m *openaiApiKeyInputModel) View() string {
	err := ""
	errStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	if m.Input.Err != nil {
		err = fmt.Sprintf(errStyle.Render("Error: %s"), m.Input.Err)
	}
	return fmt.Sprintf(
		"What is your OpenAI API key? (More info: https://platform.openai.com/account/api-keys) \n\n%s\n%s\n%s",
		m.Input.View(),
		err,
		"(enter to continue) (ctrl+t to toggle visibility) (ctrl+c to quit)",
	)
}
