package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peetya/snipforge-cli/util"
	"time"
)

type saveCodeSnippetMsg struct {
	Err error
}

type saveCommand struct {
	PromptModel *Model
}

func (c *saveCommand) SaveCodeSnippet() tea.Msg {
	if c.PromptModel.generateRequest.IsDryRun {
		time.Sleep(1 * time.Second)
		return saveCodeSnippetMsg{
			Err: nil,
		}
	}

	return saveCodeSnippetMsg{
		Err: util.SaveSnippet(c.PromptModel.result, c.PromptModel.generateRequest.Output),
	}
}
