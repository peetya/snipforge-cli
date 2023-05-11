package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peetya/snipforge-cli/util"
	"time"
)

type prepareMsg struct {
	Err error
}

type prepareCommand struct {
	PromptModel *Model
}

func (c *prepareCommand) Prepare() tea.Msg {
	if c.PromptModel.generateRequest.IsDryRun {
		time.Sleep(1 * time.Second)
		return prepareMsg{
			Err: nil,
		}
	}

	return prepareMsg{
		Err: util.PrepareOutputFolderPath(c.PromptModel.generateRequest.Output),
	}
}
