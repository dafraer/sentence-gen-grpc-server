package gemini

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/genai"
)

const (
	generateSentencePrompt = `
Generate a simple sentence in %s using the word %s.  
-The sentence should make it easy to understand the word from context.  
-If the word doesn't exist or if it is from another language, leave the fields empty 
-Otherwise, return the sentence and it's translation to %s language.
-Translation hint:%s`
	generateDefinitionPrompt = `
Generate a simple definition in %s for the word/term %s.  
-If the word/term doesn't exist in the language leave the fields empty 
-Definition hint:%s`
	translationPrompt = `
Translate word/phrase %s from language %s to %s.
-If the word/phrase doesn't exist in the language leave the fields empty
-Translation hint:%s`
)

type Client struct {
	client      *genai.Client
	logger      *zap.SugaredLogger
	geminiModel string
}

// New creates new gemini client
func New(ctx context.Context, logger *zap.SugaredLogger, geminiModel string) (*Client, error) {
	logger.Infow("initializing gemini client", "model", geminiModel)
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		logger.Errorw("failed to initialize gemini client", "error", err)
		return nil, err
	}
	logger.Infow("gemini client initialized", "model", geminiModel)
	return &Client{client: client, logger: logger, geminiModel: geminiModel}, nil
}

// GenerateSentence sends a text-only request to Gemini
func (c *Client) GenerateSentence(ctx context.Context, req *SentenceGenerationRequest) (*SentenceGenerationResponse, *Tokens, error) {
	c.logger.Debugw("gemini generate sentence request started", "word", req.Word, "word_language", req.WordLanguage, "translation_language", req.TranslationLanguage)

	//Create a config for structured output
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"original_sentence": {
					Type:        genai.TypeString,
					Description: "Sentence in the language of the provided word.",
				},
				"translated_sentence": {
					Type:        genai.TypeString,
					Description: "Translated sentence.",
				},
			},
			Required:         []string{"original_sentence", "translated_sentence"},
			PropertyOrdering: []string{"original_sentence", "translated_sentence"},
		},
	}

	//Generate response
	result, err := c.client.Models.GenerateContent(
		ctx,
		c.geminiModel,
		genai.Text(formatSentenceGenPrompt(req.WordLanguage, req.TranslationLanguage, req.Word, req.TranslationHint)),
		config,
	)
	if err != nil {
		c.logger.Errorw("gemini generate sentence request failed", "error", err)
		return nil, nil, err
	}

	//Unmarshal response
	resp := &SentenceGenerationResponse{}
	if err := json.Unmarshal([]byte(result.Text()), resp); err != nil {
		c.logger.Errorw("failed to unmarshal gemini sentence response", "error", err)
		return nil, nil, err
	}

	//Calculate tokens spent
	tokens := &Tokens{
		OutputTokens: int64(result.UsageMetadata.CandidatesTokenCount),
		InputTokens:  int64(result.UsageMetadata.PromptTokenCount),
	}
	c.logger.Debugw("gemini generate sentence request completed", "input_tokens", tokens.InputTokens, "output_tokens", tokens.OutputTokens)
	return resp, tokens, nil
}

func (c *Client) Translate(ctx context.Context, req *TranslationRequest) (*TranslationResponse, *Tokens, error) {
	c.logger.Debugw("gemini translate request started", "word", req.Word, "from_language", req.FromLanguage, "to_language", req.ToLanguage)

	//Create a config for structured output
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"translation": {
					Type:        genai.TypeString,
					Description: "Translated word/phrase",
				},
			},
			Required:         []string{"translation"},
			PropertyOrdering: []string{"translation"},
		},
	}

	//Generate response
	result, err := c.client.Models.GenerateContent(
		ctx,
		c.geminiModel,
		genai.Text(formatTranslationPrompt(req.Word, req.FromLanguage, req.ToLanguage, req.TranslationHint)),
		config,
	)
	if err != nil {
		c.logger.Errorw("gemini translate request failed", "error", err)
		return nil, nil, err
	}

	//Unmarshal response
	resp := &TranslationResponse{}
	if err := json.Unmarshal([]byte(result.Text()), resp); err != nil {
		c.logger.Errorw("failed to unmarshal gemini translation response", "error", err)
		return nil, nil, err
	}

	//Calculate tokens spent
	tokens := &Tokens{
		OutputTokens: int64(result.UsageMetadata.CandidatesTokenCount),
		InputTokens:  int64(result.UsageMetadata.PromptTokenCount),
	}
	c.logger.Debugw("gemini translate request completed", "input_tokens", tokens.InputTokens, "output_tokens", tokens.OutputTokens)

	return resp, tokens, nil
}

func (c *Client) GenerateDefinition(ctx context.Context, req *DefinitionRequest) (*DefinitionResponse, *Tokens, error) {
	c.logger.Debugw("gemini generate definition request started", "word", req.Word, "language", req.Language)

	//Create a config for structured output
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"definition": {
					Type:        genai.TypeString,
					Description: "Definition of the word/phrase without the word/phrase itself",
				},
			},
			Required:         []string{"definition"},
			PropertyOrdering: []string{"definition"},
		},
	}

	//Generate response
	result, err := c.client.Models.GenerateContent(
		ctx,
		c.geminiModel,
		genai.Text(formatDefinitionPrompt(req.Language, req.Word, req.DefinitionHint)),
		config,
	)
	if err != nil {
		c.logger.Errorw("gemini generate definition request failed", "error", err)
		return nil, nil, err
	}

	//Unmarshal response
	resp := &DefinitionResponse{}
	if err := json.Unmarshal([]byte(result.Text()), resp); err != nil {
		c.logger.Errorw("failed to unmarshal gemini definition response", "error", err)
		return nil, nil, err
	}
	//Calculate tokens spent
	tokens := &Tokens{
		OutputTokens: int64(result.UsageMetadata.CandidatesTokenCount),
		InputTokens:  int64(result.UsageMetadata.PromptTokenCount),
	}
	c.logger.Debugw("gemini generate definition request completed", "input_tokens", tokens.InputTokens, "output_tokens", tokens.OutputTokens)
	return resp, tokens, nil
}

func formatSentenceGenPrompt(wordLanguage, translationLanguage, word, translationHint string) string {
	return fmt.Sprintf(generateSentencePrompt, wordLanguage, word, translationLanguage, translationHint)
}

func formatTranslationPrompt(word, fromLanguage, toLanguage, translationHint string) string {
	return fmt.Sprintf(translationPrompt, word, fromLanguage, toLanguage, translationHint)
}

func formatDefinitionPrompt(language, word, definitionHint string) string {
	return fmt.Sprintf(generateDefinitionPrompt, language, word, definitionHint)
}
