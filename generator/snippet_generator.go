package generator

import (
	"context"
	"github.com/golang-commonmark/markdown"
	"github.com/peetya/snipforge-cli/data"
	"github.com/peetya/snipforge-cli/model"
	"github.com/sashabaranov/go-openai"
	"strings"
)

type TokenUsage struct {
	TotalTokens      int
	PromptTokens     int
	CompletionTokens int
}

func GenerateCodeSnippet(req *model.GenerateRequest, detectedLanguage *data.Language) (string, TokenUsage, error) {
	client := openai.NewClient(req.OpenAIKey)

	ccr := openai.ChatCompletionRequest{
		Model: req.OpenAIModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: getSystemPrompt()},
			{Role: openai.ChatMessageRoleUser, Content: getUserPrompt(req)},
		},
		N: 1,
	}

	if req.OpenAIMaxTokens > 0 {
		ccr.MaxTokens = req.OpenAIMaxTokens
	}

	if req.OpenAITemperature > 0 {
		ccr.Temperature = req.OpenAITemperature
	}

	resp, err := client.CreateChatCompletion(context.Background(), ccr)
	if err != nil {
		return "", TokenUsage{}, err
	}

	content := resp.Choices[0].Message.Content
	parsedContent, err := parseCodeFromMarkdown(content)
	if err != nil {
		return "", TokenUsage{}, err
	}

	if detectedLanguage != nil && detectedLanguage.Format != nil {
		parsedContent = detectedLanguage.Format(parsedContent)
	}

	return parsedContent, TokenUsage{
		TotalTokens:      resp.Usage.TotalTokens,
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
	}, nil
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
