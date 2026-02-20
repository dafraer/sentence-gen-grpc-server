package gemini

type SentenceGenerationRequest struct {
	Word                string
	WordLanguage        string
	TranslationLanguage string
	TranslationHint     string
}

type SentenceGenerationResponse struct {
	OriginalSentence   string `json:"original_sentence"`
	TranslatedSentence string `json:"translated_sentence"`
}

type TranslationRequest struct {
	Word            string
	FromLanguage    string
	ToLanguage      string
	TranslationHint string
}

type TranslationResponse struct {
	Translation string `json:"translation"`
}

type DefinitionRequest struct {
	Word           string
	Language       string
	DefinitionHint string
}

type DefinitionResponse struct {
	Definition string `json:"definition"`
}

type Tokens struct {
	InputTokens  int64
	OutputTokens int64
}
