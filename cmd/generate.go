package cmd

import (
	"fmt"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/peetya/snipforge-cli/data"
	"github.com/peetya/snipforge-cli/generator"
	"github.com/peetya/snipforge-cli/model"
	"github.com/peetya/snipforge-cli/util"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var goal string
var language string
var version string
var output string
var openaiKey string
var openaiModel string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a code snippet based on a given description and programming language",
	Long: `SnipForge - A CLI tool for generating code snippets

Generate a code snippet based on a given description and programming language using SnipForge.
Provide a natural language description of your desired code functionality, and SnipForge will 
generate a corresponding code snippet in the specified programming language. 
Powered by ChatGPT, this tool helps you create custom code snippets quickly and efficiently.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if goal == "" {
			promptGoals()
		}

		if language == "" {
			promptLanguage()
		}

		if version == "" {
			promptVersion()
		}

		if output == "" {
			promptOutput()
		}

		rq := &model.GenerateRequest{
			Goal:        goal,
			Language:    language,
			Version:     version,
			Output:      output,
			OpenAIKey:   openaiKey,
			OpenAIModel: openaiModel,
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Suffix = " Generating code snippet...\n"
		s.Start()

		snippet, err := generator.GenerateCodeSnippet(rq)
		if err != nil {
			logrus.Fatal(err)
		}

		s.Stop()
		err = util.SaveSnippet(snippet, output)
		fmt.Printf("Snippet successfully generated and saved to %s\n", output)
	},
}

func init() {
	generateCmd.Flags().StringVarP(&goal, "goal", "g", "", "The functionality description for the code snippet")
	generateCmd.Flags().StringVarP(&language, "language", "l", "", "The programming language to generate code in (e.g. PHP 8.2, Golang 1.17, etc...)")
	generateCmd.Flags().StringVarP(&version, "version", "v", "", "The version of the programming language to generate code for (if applicable)")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "The output file path for the generated code snippet")
	generateCmd.Flags().StringVarP(&openaiKey, "key", "k", "", "The OpenAI API key")
	generateCmd.Flags().StringVarP(&openaiModel, "model", "m", openai.GPT3Dot5Turbo, "The OpenAI model to use")

	rootCmd.AddCommand(generateCmd)
}

func promptGoals() {
	var goals []string
	i := 1

	for {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Goal #%d", i),
		}

		res, err := prompt.Run()
		if err != nil {
			logrus.Fatal(err)
		}

		if res == "" {
			break
		}

		goals = append(goals, res)
		i++
	}

	goal = strings.Join(goals, "; ")
}

func promptLanguage() {
	prompt := promptui.Prompt{
		Label: "Language",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("language cannot be empty")
			}

			return nil
		},
	}

	res, err := prompt.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	language = res
}

func promptVersion() {
	prompt := promptui.Prompt{
		Label: "Version (optional)",
	}

	res, err := prompt.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	version = res
}

func promptOutput() {
	prompt := promptui.Prompt{
		Label:   "Output file path",
		Default: guessOutput(),
	}

	res, err := prompt.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	output = res
}

func guessOutput() string {
	var guessedLanguage data.Language
	var maxSimilarity float64

	for _, lang := range data.Languages {
		for _, name := range lang.Names {
			similarity := strutil.Similarity(language, name, metrics.NewLevenshtein())

			if similarity > maxSimilarity {
				maxSimilarity = similarity
				guessedLanguage = lang
			}
		}
	}

	return guessedLanguage.PreferredFileName
}
