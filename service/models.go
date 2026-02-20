package service

const (
	Female = iota
	Male
)

type Gender int

type GenerateSentenceRequest struct {
	Word                string
	WordLanguage        string
	TranslationLanguage string
	TranslationHint     string
	IncludeAudio        bool
	VoiceGender         Gender
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
	VoiceGender    Gender
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
	VoiceGender     Gender
}

type TranslateResponse struct {
	Translation string
	Audio       []byte
}

type AddDailySpendingParams struct {
	GeminiInputTokens  int64
	GeminiOutputTokens int64
	Characters         int64
	TTSModel           string
}
