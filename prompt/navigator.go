package prompt

import (
	"errors"
	"fmt"
	"github.com/peetya/snipforge-cli/data"
	"log"
)

type NavigationStep int

const (
	Goals NavigationStep = iota
	Language
	LanguageVersion
	OutputPath
	OpenaiApiKey
	OpenaiModel
	Loading
	End
)

type Navigator struct {
	CurrentStep NavigationStep
}

func NewNavigator(pm *Model) *Navigator {
	n := &Navigator{CurrentStep: Goals}

	n.Current(pm)
	pm.goalsInputModel.Input.Focus()

	return n
}

func (n *Navigator) Next(pm *Model) NavigationStep {
	n.CurrentStep++
	return n.Current(pm)
}

func (n *Navigator) Current(pm *Model) NavigationStep {
	if n.CurrentStep == Goals {
		pm.goalsInputModel.Input.Focus()

		if pm.generateRequest.Goal != "" {
			pm.goalsInputModel.Input.SetValue(pm.generateRequest.Goal)
			return n.Next(pm)
		}
	}

	if n.CurrentStep == Language {
		pm.goalsInputModel.Input.Blur()
		pm.languageInputModel.Input.Focus()

		if pm.generateRequest.Language != "" {
			pm.languageInputModel.Input.SetValue(pm.generateRequest.Language)
			pm.detectedLanguage = data.DetectLanguage(pm.generateRequest.Language)
			pm.languageVersionInputModel.LanguageValue = pm.generateRequest.Language
			pm.languageVersionInputModel.DetectedLanguage = pm.detectedLanguage
			pm.outputPathInputModel.DetectedLanguage = pm.detectedLanguage
			return n.Next(pm)
		}
	}

	if n.CurrentStep == LanguageVersion {
		pm.languageInputModel.Input.Blur()
		pm.languageVersionInputModel.Input.Focus()

		if pm.generateRequest.LanguageVersion != "" {
			pm.languageVersionInputModel.Input.SetValue(pm.generateRequest.LanguageVersion)
			return n.Next(pm)
		}
	}

	if n.CurrentStep == OutputPath {
		pm.languageVersionInputModel.Input.Blur()
		pm.outputPathInputModel.Input.Focus()

		if pm.generateRequest.Output != "" {
			pm.outputPathInputModel.Input.SetValue(pm.generateRequest.Output)
			return n.Next(pm)
		}
	}

	if n.CurrentStep == OpenaiApiKey {
		pm.outputPathInputModel.Input.Blur()
		pm.openaiApiKeyInputModel.Input.Focus()

		if pm.generateRequest.OpenAIKey != "" {
			pm.openaiApiKeyInputModel.Input.SetValue(pm.generateRequest.OpenAIKey)
			return n.Next(pm)
		}
	}

	if n.CurrentStep == OpenaiModel {
		pm.openaiApiKeyInputModel.Input.Blur()

		if pm.generateRequest.OpenAIModel != "" {
			for i, item := range Items {
				if ConvertItemToString(item) == pm.generateRequest.OpenAIModel {
					pm.openaiModelInputModel.Input.Select(i)
					return n.Next(pm)
				}
			}

			log.Fatal(errors.New(fmt.Sprintf(
				"invalid OpenAI model: %s. Please use one of the following: %s",
				pm.generateRequest.OpenAIModel,
				Items,
			)))
		}
	}

	return n.CurrentStep
}
