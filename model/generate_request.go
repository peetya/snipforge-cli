package model

type GenerateRequest struct {
	Goal        string
	Language    string
	Version     string
	Output      string
	OpenAIKey   string
	OpenAIModel string
}
