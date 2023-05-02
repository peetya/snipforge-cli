package generator

import (
	"context"
	"github.com/golang-commonmark/markdown"
	"github.com/peetya/snipforge-cli/data"
	"github.com/peetya/snipforge-cli/model"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func GenerateCodeSnippet(req *model.GenerateRequest, detectedLanguage *data.Language) (string, error) {
	client := openai.NewClient(req.OpenAIKey)

	systemPrompt, err := getSystemPrompt()
	if err != nil {
		return "", err
	}

	userPrompt, err := getUserPrompt(req)
	if err != nil {
		return "", err
	}

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: req.OpenAIModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
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
	}).Debug("Received GPT response")
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

func getSystemPrompt() (string, error) {
	c, err := os.ReadFile("./prompts/system_prompt.txt")
	if err != nil {
		return "", err
	}
	return string(c), nil

}

func getUserPrompt(req *model.GenerateRequest) (string, error) {
	c, err := os.ReadFile("./prompts/user_prompt.txt")
	if err != nil {
		return "", err
	}
	prompt := string(c)
	prompt = strings.Replace(prompt, "{{LANG}}", req.Language, -1)
	prompt = strings.Replace(prompt, "{{GOAL}}", req.Goal, -1)
	return prompt, nil
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
