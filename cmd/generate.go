package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peetya/snipforge-cli/model"
	"github.com/peetya/snipforge-cli/prompt"
	"github.com/spf13/cobra"
)

var (
	goal              string
	language          string
	languageVersion   string
	output            string
	openaiKey         string
	openaiModel       string
	openaiMaxTokens   int
	openaiTemperature float32
	isQuiet           bool
	isDryRun          bool
	isStdout          bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a code snippet based on the required description and programming or tooling language",
	Long: `Generate a code snippet based on the required description and programming or tooling language

  Find more information at: https://github.com/peetya/snipforge-cli`,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := &model.GenerateRequest{
			Goal:              goal,
			Language:          language,
			LanguageVersion:   languageVersion,
			Output:            output,
			OpenAIKey:         openaiKey,
			OpenAIModel:       openaiModel,
			OpenAIMaxTokens:   openaiMaxTokens,
			OpenAITemperature: openaiTemperature,

			IsQuiet:  isQuiet,
			IsDryRun: isDryRun,
			IsStdout: isStdout,
		}

		var po []tea.ProgramOption

		if req.IsStdout {
			req.IsQuiet = true
			po = append(po, tea.WithoutRenderer())
		}

		p := tea.NewProgram(prompt.InitializeModel(req), po...)
		if _, err := p.Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&goal, "goal", "g", "", "the functionality description for the code snippet")
	generateCmd.Flags().StringVarP(&language, "language", "l", "", "the programming or tooling language to generate code in (e.g. PHP, Golang, etc...)")
	generateCmd.Flags().StringVarP(&languageVersion, "language-version", "v", "", "the version of the programming or tooling language to generate code for (if applicable)")

	generateCmd.Flags().StringVarP(&output, "output", "o", "", "the output file path for the generated code snippet")
	generateCmd.Flags().BoolVar(&isStdout, "stdout", false, "print the generated code snippet to isStdout instead of saving to a file")

	generateCmd.Flags().StringVarP(&openaiKey, "openai-key", "k", "", "the OpenAI API key")
	generateCmd.Flags().StringVarP(&openaiModel, "openai-model", "m", "", "the OpenAI model to use")
	generateCmd.Flags().IntVar(&openaiMaxTokens, "openai-max-tokens", 0, "the maximum number of tokens to generate")
	generateCmd.Flags().Float32Var(&openaiTemperature, "openai-temperature", 0.0, "the sampling temperature for the OpenAI model (between 0.0 and 2.0)")

	generateCmd.Flags().BoolVarP(&isQuiet, "quiet", "q", false, "suppress all output except for the generated code snippet")
	generateCmd.Flags().BoolVarP(&isDryRun, "dry-run", "d", false, "do not generate a code snippet, only print the generated description")

	rootCmd.AddCommand(generateCmd)
}
