package prompt

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peetya/snipforge-cli/data"
)

type languageVersionInputModel struct {
	Input            textinput.Model
	DetectedLanguage *data.Language
	LanguageValue    string
}

func newLanguageVersionInputModel() languageVersionInputModel {
	m := languageVersionInputModel{}
	m.Init()
	return m
}

func (m *languageVersionInputModel) Init() {
	m.Input = textinput.New()
}

func (m *languageVersionInputModel) Update(pm *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return pm, tea.Quit

		case tea.KeyEnter:
			pm.navigator.Next(pm)
		}
	}

	return pm, cmd
}

func (m *languageVersionInputModel) View() string {
	lang := m.LanguageValue

	if m.DetectedLanguage != nil {
		lang = m.DetectedLanguage.Names[0]
	}

	m.Input.Placeholder = m.placeholder()

	return fmt.Sprintf(
		"Which version of %s do you want to use? (optional)\n\n%s\n\n%s",
		lang,
		m.Input.View(),
		"(enter to continue) (ctrl+c to quit)",
	)
}

func (m *languageVersionInputModel) placeholder() string {
	version := "latest"
	if m.DetectedLanguage != nil && m.DetectedLanguage.PreferredVersion != "" {
		version = m.DetectedLanguage.PreferredVersion
	}

	return fmt.Sprintf("e.g. %s", version)
}
