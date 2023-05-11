package model

type GenerateRequest struct {
	Goal              string
	Language          string
	LanguageVersion   string
	Output            string
	OpenAIKey         string
	OpenAIModel       string
	OpenAIMaxTokens   int
	OpenAITemperature float32

	IsQuiet  bool
	IsDryRun bool
	IsStdout bool
}

func (r *GenerateRequest) IsMandatoryParamMissing() bool {
	return r.Goal == "" || r.Language == "" || r.Output == "" || r.OpenAIKey == "" || r.OpenAIModel == ""
}
