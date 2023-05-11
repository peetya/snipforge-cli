package prompt

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peetya/snipforge-cli/data"
	"github.com/peetya/snipforge-cli/generator"
	"github.com/peetya/snipforge-cli/model"
	"log"
	"os"
	"time"
)

type Model struct {
	navigator                   *Navigator
	spinnerStart                time.Time
	isPreparationStarted        bool
	isPreparationDone           bool
	isCodeSnippetGenerationDone bool
	isCodeSnippetSavingDone     bool
	detectedLanguage            *data.Language

	goalsInputModel           goalsInputModel
	languageInputModel        languageInputModel
	languageVersionInputModel languageVersionInputModel
	outputPathInputModel      OutputPathInputModel
	openaiApiKeyInputModel    openaiApiKeyInputModel
	openaiModelInputModel     OpenaiModelInputModel
	spinnerLoading            spinner.Model

	generateRequest *model.GenerateRequest
	result          string
	tokenUsage      generator.TokenUsage
}

func InitializeModel(req *model.GenerateRequest) *Model {
	m := &Model{
		isPreparationStarted:        false,
		isPreparationDone:           false,
		isCodeSnippetGenerationDone: false,
		isCodeSnippetSavingDone:     false,

		goalsInputModel:           newGoalsInputModel(),
		languageInputModel:        newLanguageInputModel(),
		languageVersionInputModel: newLanguageVersionInputModel(),
		outputPathInputModel:      NewOutputPathInputModel(),
		openaiApiKeyInputModel:    newOpenaiApiKeyInputModel(),
		openaiModelInputModel:     NewOpenaiModelInputModel(),
		spinnerLoading:            spinner.New(spinner.WithSpinner(spinner.Dot)),

		generateRequest: req,
	}

	m.navigator = NewNavigator(m)

	return m
}

func (pm *Model) Init() tea.Cmd {
	if pm.generateRequest.IsQuiet {
		pm.navigator.CurrentStep = Loading
	}

	return pm.spinnerLoading.Tick
}

func (pm *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if pm.navigator.CurrentStep == Loading && !pm.isPreparationStarted {
		var spinnerCmd tea.Cmd
		pm.spinnerLoading, spinnerCmd = pm.spinnerLoading.Update(msg)
		pm.isPreparationStarted = true
		return pm, tea.Batch(spinnerCmd, pm.getExecutionCmd())
	}

	switch msg := msg.(type) {
	case prepareMsg:
		if msg.Err != nil {
			log.Fatal(msg.Err)
		}

		pm.isPreparationDone = true

		return pm, nil

	case generateCodeSnippetMsg:
		if msg.Err != nil {
			log.Fatal(msg.Err)
		}

		pm.isCodeSnippetGenerationDone = true

		return pm, nil

	case saveCodeSnippetMsg:
		if msg.Err != nil {
			log.Fatal(msg.Err)
		}

		pm.isCodeSnippetSavingDone = true

		if pm.generateRequest.IsStdout {
			fmt.Println(pm.result)
			os.Exit(0)
		}

		return pm, tea.Quit

	case tea.KeyMsg:
		switch pm.navigator.CurrentStep {
		case Goals:
			return pm.goalsInputModel.Update(pm, msg)

		case Language:
			return pm.languageInputModel.Update(pm, msg)

		case LanguageVersion:
			return pm.languageVersionInputModel.Update(pm, msg)

		case OutputPath:
			return pm.outputPathInputModel.Update(pm, msg)

		case OpenaiApiKey:
			return pm.openaiApiKeyInputModel.Update(pm, msg)

		case OpenaiModel:
			return pm.openaiModelInputModel.Update(pm, msg)

		case Loading, End:
			switch msg.Type {
			case tea.KeyCtrlC:
				return pm, tea.Quit
			}
		}
	default:
		var spinnerCmd tea.Cmd
		pm.spinnerLoading, spinnerCmd = pm.spinnerLoading.Update(msg)
		return pm, spinnerCmd
	}
	return pm, nil
}

func (pm *Model) View() string {
	view := `Welcome to SnipForge's interactive mode! In this mode, you will be prompted to provide the necessary information for generating a code snippet.

  More info: https://github.com/peetya/snipforge-cli

`
	switch pm.navigator.CurrentStep {
	case Goals:
		view += pm.goalsInputModel.View()

	case Language:
		view += pm.languageInputModel.View()

	case LanguageVersion:
		view += pm.languageVersionInputModel.View()

	case OutputPath:
		view += pm.outputPathInputModel.View()

	case OpenaiApiKey:
		view += pm.openaiApiKeyInputModel.View()

	case OpenaiModel:
		view += pm.openaiModelInputModel.View()

	case Loading, End:
		prepTxt := fmt.Sprintf("%s Preparing output folder...", pm.spinnerLoading.View())
		if pm.isPreparationDone {
			prepTxt = "✅  Output folder is ready!"
		}
		genTxt := fmt.Sprintf("%s Generating code snippet...", pm.spinnerLoading.View())
		if pm.isCodeSnippetGenerationDone {
			genTxt = "✅  Code snippet is generated!"
		}
		saveTxt := fmt.Sprintf("%s Saving code snippet...", pm.spinnerLoading.View())
		if pm.isCodeSnippetSavingDone {
			saveTxt = "✅  Code snippet is saved!"
		}

		if pm.generateRequest.IsStdout {
			prepTxt = "✅  Output folder generation is skipped!"
			saveTxt = "✅  Code snippet saving is skipped!"
		}

		view += fmt.Sprintf(
			"%s\n%s\n%s\n\n",
			prepTxt,
			genTxt,
			saveTxt,
		)

		if pm.isPreparationDone && pm.isCodeSnippetGenerationDone && pm.isCodeSnippetSavingDone {
			if pm.generateRequest.IsDryRun {
				view += fmt.Sprintf("Dry run is completed! Your request passed all checks. No code snippet is generated.\n")
			} else {
				maxTokenTxt := ""

				if pm.generateRequest.OpenAIMaxTokens != 0 {
					maxTokenTxt = fmt.Sprintf("/%d", pm.generateRequest.OpenAIMaxTokens)
				}

				view += fmt.Sprintf("The snippet is successfully generated and saved to %s\n\n", pm.generateRequest.Output)
				view += fmt.Sprintf("Model used: %s\n", pm.generateRequest.OpenAIModel)
				view += fmt.Sprintf("Tokens used: %d (prompt: %d, completion: %d%s)\n", pm.tokenUsage.TotalTokens, pm.tokenUsage.PromptTokens, pm.tokenUsage.CompletionTokens, maxTokenTxt)
				view += fmt.Sprintf("\n")
				view += fmt.Sprintf("Please review the generated code snippet before using it in your project!\n")
			}
		}
	}

	if pm.generateRequest.IsStdout {
		return ""
	}

	return view
}

func (pm *Model) getExecutionCmd() tea.Cmd {
	pc := &prepareCommand{PromptModel: pm}
	gc := &generateCommand{PromptModel: pm}
	sc := &saveCommand{PromptModel: pm}

	return tea.Sequence(pc.Prepare, gc.GenerateCodeSnippet, sc.SaveCodeSnippet)
}

func (pm *Model) initNonInteractive() {
	if pm.generateRequest.Goal == "" {
		return
	}
}
