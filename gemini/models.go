package gemini

type SentenceGenerationResponse struct {
	OriginalSentence   string `json:"original_sentence"`
	TranslatedSentence string `json:"translated_sentence"`
}

type TranslationResponse struct {
	Translation string `json:"translation"`
}

type DefinitionResponse struct {
	Definition string `json:"definition"`
}
