package generator

import (
	"context"
	"github.com/golang-commonmark/markdown"
	"github.com/peetya/snipforge-cli/data"
	"github.com/peetya/snipforge-cli/model"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"strings"
)

func GenerateCodeSnippet(req *model.GenerateRequest, detectedLanguage *data.Language) (string, error) {
	client := openai.NewClient(req.OpenAIKey)

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: req.OpenAIModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: getSystemPrompt()},
			{Role: openai.ChatMessageRoleUser, Content: getUserPrompt(req)},
		},
		N: 1,
	})
	if err != nil {
		return "", err
	}

	content := resp.Choices[0].Message.Content

	logrus.WithFields(logrus.Fields{
		"finishReason":     resp.Choices[0].FinishReason,
		"promptTokens":     resp.Usage.PromptTokens,
		"completionTokens": resp.Usage.CompletionTokens,
		"totalTokens":      resp.Usage.TotalTokens,
	}).Debug("Received OpenAI API response")
	logrus.WithField("content", content).Trace("Received OpenAI API response content")

	parsedContent, err := parseCodeFromMarkdown(content)
	if err != nil {
		return "", err
	}

	if detectedLanguage != nil && detectedLanguage.Format != nil {
		logrus.Debugf("Apply formatting for detected language: %s", detectedLanguage.Names[0])
		parsedContent = detectedLanguage.Format(parsedContent)
	}

	return parsedContent, nil
}

func getSystemPrompt() string {
	return `You are a Code Snippet Generator called SnipForge.
You are given a goal and a programming language.
You generate code snippets that achieve the goal in the given programming language.`
}

func getUserPrompt(req *model.GenerateRequest) string {
	prompt := `Please provide a code snippet in {{LANG}} in markdown format that achieves the following goals: {{GOAL}}.
Return only the code itself, without any additional text or explanation or note.`
	prompt = strings.Replace(prompt, "{{LANG}}", req.Language, -1)
	prompt = strings.Replace(prompt, "{{GOAL}}", req.Goal, -1)
	return prompt
}

func parseCodeFromMarkdown(mdContent string) (string, error) {
	md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
	tokens := md.Parse([]byte(mdContent))
	content := ""

	for _, token := range tokens {
		if _, ok := token.(*markdown.CodeBlock); ok {
			content = token.(*markdown.CodeBlock).Content
			break
		}

		if _, ok := token.(*markdown.Fence); ok {
			content = token.(*markdown.Fence).Content
			break
		}
	}

	return content, nil
}
