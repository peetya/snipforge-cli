package model

type GenerateRequest struct {
	Goal            string
	Language        string
	LanguageVersion string
	Output          string
	OpenAIKey       string
	OpenAIModel     string
}
