package service

type GenerateSentenceRequest struct {
	Word            string
	FromLanguage    string
	ToLanguage      string
	TranslationHint string
}

type GenerateSentenceResponse struct {
	OriginalSentence   string
	TranslatedSentence string
	Audio              []byte
}

type GenerateDefinitionRequest struct {
	Word           string
	Language       string
	DefinitionHint string
}

type GenerateDefinitionResponse struct {
	Definition string
	Audio      []byte
}

type TranslateRequest struct {
	Word            string
	FromLanguage    string
	ToLanguage      string
	TranslationHint string
}

type TranslateResponse struct {
	Translation string
}
