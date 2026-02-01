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
-If the word/term doesn't exist leave the fields empty 
-Definition hint:%s`
	translationPrompt = `
Translate word/phrase %s from language %s to %s.
-If the word/phrase doesn't exist leave the fields empty
-Translation hint:%s`
)

type Client struct {
	client      *genai.Client
	logger      *zap.SugaredLogger
	geminiModel string
}

// New creates new gemini client
func New(ctx context.Context, logger *zap.SugaredLogger, geminiModel string) (*Client, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Client{client: client, logger: logger, geminiModel: geminiModel}, nil
}

// GenerateSentence sends a text-only request to Gemini
func (c *Client) GenerateSentence(ctx context.Context, word, wordLang, targetLang, translationHint string) (SentenceGenerationResponse, error) {
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
		genai.Text(formatSentenceGenPrompt(wordLang, targetLang, word, translationHint)),
		config,
	)
	if err != nil {
		return SentenceGenerationResponse{}, err
	}

	//Unmarshal response
	var resp SentenceGenerationResponse
	if err := json.Unmarshal([]byte(result.Text()), &resp); err != nil {
		return SentenceGenerationResponse{}, err
	}
	return resp, nil
}

func (c *Client) Translate(ctx context.Context, word, fromLang, toLang, translationHint string) (TranslationResponse, error) {
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
		genai.Text(formatTranslationPrompt(word, fromLang, toLang, translationHint)),
		config,
	)
	if err != nil {
		return TranslationResponse{}, err
	}

	//Unmarshal response
	var resp TranslationResponse
	if err := json.Unmarshal([]byte(result.Text()), &resp); err != nil {
		return TranslationResponse{}, err
	}
	return resp, nil
}

func (c *Client) GenerateDefinition(ctx context.Context, word, language, definitionHint string) (DefinitionResponse, error) {
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"definition": {
					Type:        genai.TypeString,
					Description: "Definition of the word/phrase",
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
		genai.Text(formatDefinitionPrompt(language, word, definitionHint)),
		config,
	)
	if err != nil {
		return DefinitionResponse{}, err
	}

	//Unmarshal response
	var resp DefinitionResponse
	if err := json.Unmarshal([]byte(result.Text()), &resp); err != nil {
		return DefinitionResponse{}, err
	}
	return resp, nil
}

func formatSentenceGenPrompt(fromLanguage, toLanguage, word, translationHint string) string {
	return fmt.Sprintf(generateSentencePrompt, fromLanguage, word, toLanguage, translationHint)
}

func formatTranslationPrompt(word, fromLang, toLang, translationHint string) string {
	return fmt.Sprintf(translationPrompt, word, fromLang, toLang, translationHint)
}

func formatDefinitionPrompt(language, word, definitionHint string) string {
	return fmt.Sprintf(generateDefinitionPrompt, language, word, definitionHint)
}
