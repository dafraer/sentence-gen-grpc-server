package tts

import (
	"context"
	"errors"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"go.uber.org/zap"
)

const (
	Chirp3HD = "Chirp3-HD"
	Standard = "Standard"
	Male     = "MALE"
	Female   = "FEMALE"
)

var (
	ErrNoSuchVoice = errors.New("no such voice")
)

type Client struct {
	tts    *texttospeech.Client
	logger *zap.SugaredLogger
}

// New creates new tts client
func New(ctx context.Context, logger *zap.SugaredLogger) (*Client, error) {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		tts:    client,
		logger: logger,
	}, nil
}

// Close closes tts client
func (c *Client) Close() error {
	if err := c.tts.Close(); err != nil {
		return err
	}
	return nil
}

// Generate generates mp3 audio based on the text and language provided
func (c *Client) Generate(ctx context.Context, text, languageCode, gender, model string) ([]byte, error) {
	//Select a voice
	voices, err := c.tts.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{
		LanguageCode: languageCode,
	})
	if err != nil {
		return nil, err
	}

	name := ""
	for _, v := range voices.Voices {
		if strings.Contains(v.Name, model) && v.SsmlGender.String() == gender {
			name = v.Name
			break
		}
	}
	if name == "" {
		return nil, ErrNoSuchVoice
	}

	// Perform the text-to-speech request on the text input with the selected voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},

		// Build the voice request, select the language code (e.g. "en-US") and the SSML voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			Name:         name,
			LanguageCode: languageCode,
		},

		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	//Generate speech
	resp, err := c.tts.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp.AudioContent, nil
}
