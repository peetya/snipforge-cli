package prompt

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peetya/snipforge-cli/data"
	"github.com/peetya/snipforge-cli/generator"
	"github.com/sashabaranov/go-openai"
	"time"
)

type generateCodeSnippetMsg struct {
	Result string
	Err    error
}

type generateCommand struct {
	PromptModel *Model
}

func (c *generateCommand) GenerateCodeSnippet() tea.Msg {
	if c.PromptModel.generateRequest.IsDryRun {
		time.Sleep(1 * time.Second)

		return generateCodeSnippetMsg{
			Result: "",
			Err:    nil,
		}
	}

	if c.PromptModel.generateRequest.IsQuiet {
		if err := c.validate(); err != nil {
			return generateCodeSnippetMsg{
				Result: "",
				Err:    err,
			}
		}

		c.setDefaultParams()
	} else {
		c.PromptModel.generateRequest.Goal = c.PromptModel.goalsInputModel.Input.Value()
		c.PromptModel.generateRequest.Language = c.PromptModel.languageInputModel.Input.Value()
		c.PromptModel.generateRequest.LanguageVersion = c.PromptModel.languageVersionInputModel.Input.Value()
		c.PromptModel.generateRequest.Output = c.PromptModel.outputPathInputModel.Input.Value()
		c.PromptModel.generateRequest.OpenAIKey = c.PromptModel.openaiApiKeyInputModel.Input.Value()
		c.PromptModel.generateRequest.OpenAIModel = ConvertItemToString(c.PromptModel.openaiModelInputModel.Input.SelectedItem())
	}

	res, tokenUsage, err := generator.GenerateCodeSnippet(c.PromptModel.generateRequest, c.PromptModel.detectedLanguage)

	c.PromptModel.result = res
	c.PromptModel.tokenUsage = tokenUsage

	return generateCodeSnippetMsg{
		Result: res,
		Err:    err,
	}
}

func (c *generateCommand) validate() error {
	if c.PromptModel.generateRequest.Goal == "" {
		return errors.New("missing mandatory parameter: goal")
	}

	if c.PromptModel.generateRequest.Language == "" {
		return errors.New("missing mandatory parameter: language")
	}

	if c.PromptModel.generateRequest.OpenAIKey == "" {
		return errors.New("missing mandatory parameter: openai key")
	}

	return nil
}

func (c *generateCommand) setDefaultParams() {
	if c.PromptModel.detectedLanguage == nil {
		c.PromptModel.detectedLanguage = data.DetectLanguage(c.PromptModel.generateRequest.Language)
	}

	if c.PromptModel.generateRequest.LanguageVersion == "" {
		version := "latest"
		if c.PromptModel.detectedLanguage != nil && c.PromptModel.detectedLanguage.PreferredVersion != "" {
			version = c.PromptModel.detectedLanguage.PreferredVersion
		}
		c.PromptModel.generateRequest.LanguageVersion = version
	}

	if c.PromptModel.generateRequest.Output == "" {
		outputPath := "snippet.txt"
		if c.PromptModel.detectedLanguage != nil && c.PromptModel.detectedLanguage.PreferredFileName != "" {
			outputPath = c.PromptModel.detectedLanguage.PreferredFileName
		}
		c.PromptModel.generateRequest.Output = outputPath
	}

	if c.PromptModel.generateRequest.OpenAIModel == "" {
		c.PromptModel.generateRequest.OpenAIModel = openai.GPT3Dot5Turbo
	}
}
