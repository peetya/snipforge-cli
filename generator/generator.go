package generator

import (
	"context"
	"github.com/golang-commonmark/markdown"
	"github.com/peetya/snipforge-cli/model"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

func GenerateCodeSnippet(req *model.GenerateRequest) (string, error) {
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
	})

	if err != nil {
		return "", err
	}

	return parseCodeFromMarkdown(resp.Choices[0].Message.Content)
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
