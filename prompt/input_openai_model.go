package prompt

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sashabaranov/go-openai"
	"io"
	"time"
)

var Items = []list.Item{
	OpenaiModelItem(openai.GPT432K0314),
	OpenaiModelItem(openai.GPT432K),
	OpenaiModelItem(openai.GPT40314),
	OpenaiModelItem(openai.GPT4),
	OpenaiModelItem(openai.GPT3Dot5Turbo0301),
	OpenaiModelItem(openai.GPT3Dot5Turbo),
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                              { return 1 }
func (d itemDelegate) Spacing() int                             { return 0 }
func (d itemDelegate) Update(ms tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(OpenaiModelItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	if index == m.Index() {
		str = " > " + str
	} else {
		str = "   " + str
	}

	fmt.Fprint(w, str)
}

type OpenaiModelInputModel struct {
	Input list.Model
}

func NewOpenaiModelInputModel() OpenaiModelInputModel {
	m := OpenaiModelInputModel{}
	m.Init()
	return m
}

func (m *OpenaiModelInputModel) Init() {
	m.Input = list.New(Items, itemDelegate{}, 100, 13)
	m.Input.SetShowTitle(false)
	m.Input.ResetSelected()
	m.Input.Select(5) // Select GPT3Dot5Turbo by default
}

func (m *OpenaiModelInputModel) Update(pm *Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)

	m.Input, cmd = m.Input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return pm, tea.Quit

		case tea.KeyEnter:
			pm.spinnerStart = time.Now()
			pm.navigator.Next(pm)

			return pm, pm.getExecutionCmd()
		}
	}

	return pm, cmd
}

func (m *OpenaiModelInputModel) View() string {
	return fmt.Sprintf(
		"What OpenAI model do you want to use?\n\n%s\n\n%s",
		m.Input.View(),
		"(enter to continue) (ctrl+c to quit)",
	)
}
