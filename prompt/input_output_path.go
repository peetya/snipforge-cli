package prompt

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/peetya/snipforge-cli/data"
)

type OutputPathInputModel struct {
	Input            textinput.Model
	DetectedLanguage *data.Language
}

func NewOutputPathInputModel() OutputPathInputModel {
	m := OutputPathInputModel{}
	m.Init()
	return m
}

func (m *OutputPathInputModel) Init() {
	m.Input = textinput.New()
	m.Input.Validate = func(s string) error {
		if s == "" {
			return errors.New("output path cannot be empty")
		}

		return nil
	}
}

func (m *OutputPathInputModel) Update(pm *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return pm, tea.Quit

		case tea.KeyEnter:
			if pm.outputPathInputModel.Input.Err == nil {
				pm.navigator.Next(pm)
			}
		}
	}

	return pm, cmd
}

func (m *OutputPathInputModel) View() string {
	m.Input.Placeholder = m.placeholder()

	err := ""
	errStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	if m.Input.Err != nil {
		err = fmt.Sprintf(errStyle.Render("Error: %s"), m.Input.Err)
	}

	return fmt.Sprintf(
		"Where do you want to save the snippet? \n\n%s\n%s\n%s",
		m.Input.View(),
		err,
		"(enter to continue) (ctrl+c to quit)",
	)
}

func (m *OutputPathInputModel) placeholder() string {
	outputPath := "snippet.txt"
	if m.DetectedLanguage != nil && m.DetectedLanguage.PreferredFileName != "" {
		outputPath = m.DetectedLanguage.PreferredFileName
	}

	return fmt.Sprintf("e.g. %s", outputPath)
}
