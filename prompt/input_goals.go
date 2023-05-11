package prompt

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
)

type goalsInputModel struct {
	Input textarea.Model
}

func newGoalsInputModel() goalsInputModel {
	m := goalsInputModel{}
	m.Init()
	return m
}

func (m *goalsInputModel) Init() {
	m.Input = textarea.New()
	m.Input.Placeholder = m.placeholder()
	m.Input.SetWidth(100)
	m.Input.SetHeight(5)
}

func (m *goalsInputModel) Update(pm *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Input.Value() == "" {
			m.Input.Err = errors.New("goals cannot be empty")
		} else {
			m.Input.Err = nil
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			return pm, tea.Quit

		case tea.KeyCtrlS:
			if m.Input.Err == nil {
				pm.navigator.Next(pm)
			}
		}
	}

	return pm, cmd
}

func (m *goalsInputModel) View() string {
	err := ""
	errStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	if m.Input.Err != nil {
		err = fmt.Sprintf(errStyle.Render("Error: %s"), m.Input.Err)
	}

	return fmt.Sprintf(
		"What are your goals?\n\n%s\n%s\n%s",
		m.Input.View(),
		err,
		"(ctrl+s to continue) (ctrl+c to quit)",
	)
}

func (m *goalsInputModel) placeholder() string {
	placeholders := []string{
		"A Hello World script",
		"A script that returns the next prime number",
		"A script that returns the current time",
	}

	n := rand.Intn(len(placeholders))
	return fmt.Sprintf("e.g. %s", placeholders[n])
}
