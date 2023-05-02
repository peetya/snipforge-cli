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

var (
	goal            string
	language        string
	languageVersion string
	output          string
	stdout          bool
	openaiKey       string
	openaiModel     string
	isQuiet         bool
	isDryRun        bool
	isInteractive   bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a code snippet based on a given description and programming or tooling language",
	Long: `Generate a code snippet based on a given description and programming or tooling language

  Find more information at: https://github.com/peetya/snipforge-cli`,
	RunE: func(cmd *cobra.Command, args []string) error {
		isInteractive = !areMandatoryFlagsProvided()

		if stdout {
			isQuiet = true
		}

		if isInteractive && !isQuiet {
			logrus.Info("Welcome to SnipForge's interactive mode! In this mode, you will be prompted to provide the necessary information for generating a code snippet.\n")
		}

		if goal == "" {
			if isQuiet {
				return fmt.Errorf("goal is required")
			}
			promptGoals()
		}

		if language == "" {
			if isQuiet {
				return fmt.Errorf("language is required")
			}
			promptLanguage()

			if languageVersion == "" && !isQuiet {
				promptVersion()
			}
		}

		detectedLanguage := detectLanguage()

		if output == "" {
			if isQuiet {
				output = guessOutput(detectedLanguage)
			} else {
				promptOutput(detectedLanguage)
			}
		}

		if openaiKey == "" {
			if isQuiet {
				return fmt.Errorf("openai-key is required")
			}
			promptOpenAIKey()
		}

		rq := &model.GenerateRequest{
			Goal:            goal,
			Language:        language,
			LanguageVersion: languageVersion,
			Output:          output,
			OpenAIKey:       openaiKey,
			OpenAIModel:     openaiModel,
		}

		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Generating code snippet...\n"
		s.Start()

		if !stdout {
			if err := util.PrepareOutputFolderPath(output); err != nil {
				return err
			}
		}

		if isDryRun {
			logrus.WithFields(logrus.Fields{
				"goal":            goal,
				"language":        language,
				"languageVersion": languageVersion,
				"output":          output,
			}).Warningf("Dry run enabled, skipping generating code snippet\n")
			return nil
		}

		snippet, err := generator.GenerateCodeSnippet(rq, detectedLanguage)
		if err != nil {
			return err
		}

		s.Stop()

		if stdout {
			fmt.Println(snippet)
			return nil
		}

		if err = util.SaveSnippet(snippet, output); err != nil {
			return err
		}

		logrus.Infof("Snippet successfully generated and saved to %s\n", output)
		logrus.Warn("Please review the generated code snippet before using it in your project!")

		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&goal, "goal", "g", "", "the functionality description for the code snippet")
	generateCmd.Flags().StringVarP(&language, "language", "l", "", "the programming or tooling language to generate code in (e.g. PHP, Golang, etc...)")
	generateCmd.Flags().StringVarP(&languageVersion, "language-version", "v", "", "the version of the programming or tooling language to generate code for (if applicable)")

	generateCmd.Flags().StringVarP(&output, "output", "o", "", "the output file path for the generated code snippet")
	generateCmd.Flags().BoolVar(&stdout, "stdout", false, "print the generated code snippet to stdout instead of saving to a file")

	generateCmd.Flags().StringVarP(&openaiKey, "openai-key", "k", "", "the OpenAI API key")
	generateCmd.Flags().StringVarP(&openaiModel, "openai-model", "m", openai.GPT3Dot5Turbo, "the OpenAI model to use")

	generateCmd.Flags().BoolVarP(&isQuiet, "quiet", "q", false, "suppress all output except for the generated code snippet")
	generateCmd.Flags().BoolVarP(&isDryRun, "dry-run", "d", false, "do not generate a code snippet, only print the generated description")

	rootCmd.AddCommand(generateCmd)
}

func promptGoals() {
	var goals []string
	i := 1

	logrus.Info("First, please enter your goals one by one. These goals will help SnipForge understand the functionality you want in your code snippet. After entering a goal, press Enter to input the next one. When you're done, simply press Enter on an empty line to proceed to the next step.")
	logrus.Info("Enter your goals:")

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
				return fmt.Errorf("Language cannot be empty. Please provide a programming or tooling language, e.g. PHP, Golang, Docker, etc...")
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
		Label: "LanguageVersion (optional)",
	}

	res, err := prompt.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	languageVersion = res
}

func promptOutput(detectedLanguage *data.Language) {
	prompt := promptui.Prompt{
		Label:     "Output file path",
		Default:   guessOutput(detectedLanguage),
		AllowEdit: true,
	}

	res, err := prompt.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	output = res
}

func guessOutput(detectedLanguage *data.Language) string {
	if detectedLanguage != nil {
		return detectedLanguage.PreferredFileName
	}

	return "snippet.txt"
}

func promptOpenAIKey() {
	prompt := promptui.Prompt{
		Label: "OpenAI API Key",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("OpenAI API Key cannot be empty. Please provide a valid OpenAI API Key. More info: https://platform.openai.com/account/api-keys")
			}

			return nil
		},
	}

	res, err := prompt.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	openaiKey = res
}

func areMandatoryFlagsProvided() bool {
	return goal != "" && language != "" && output != "" && openaiKey != ""
}

func detectLanguage() *data.Language {
	var detectedLanguage data.Language
	var maxSimilarity float64

	similarityScoreThreshold := 0.5

	for _, lang := range data.Languages {
		for _, name := range lang.Names {
			similarity := strutil.Similarity(strings.ToLower(language), strings.ToLower(name), metrics.NewLevenshtein())

			if similarity > maxSimilarity {
				maxSimilarity = similarity
				detectedLanguage = lang
			}
		}
	}

	logrus.Debugf("Detected language is %s with similarity score of %f", detectedLanguage.Names[0], maxSimilarity)

	if maxSimilarity < similarityScoreThreshold {
		logrus.Debugf("Similarity score is lower than %f, skipping language detection", similarityScoreThreshold)
		return nil
	}

	return &detectedLanguage
}
