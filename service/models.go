package service

type GenerateSentenceRequest struct {
	Word                string
	WordLanguage        string
	TranslationLanguage string
	TranslationHint     string
	IncludeAudio        bool
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
	IncludeAudio   bool
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
	IncludeAudio    bool
}

type TranslateResponse struct {
	Translation string
	Audio       []byte
}
