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
	logger.Infow("initializing tts client")
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		logger.Errorw("failed to initialize tts client", "error", err)
		return nil, err
	}
	logger.Infow("tts client initialized")
	return &Client{
		tts:    client,
		logger: logger,
	}, nil
}

// Close closes tts client
func (c *Client) Close() error {
	c.logger.Infow("closing tts client")
	if err := c.tts.Close(); err != nil {
		c.logger.Errorw("failed to close tts client", "error", err)
		return err
	}
	c.logger.Debugw("tts client closed")
	return nil
}

// Generate generates mp3 audio based on the text and language provided
func (c *Client) Generate(ctx context.Context, text, languageCode, gender, model string) ([]byte, error) {
	c.logger.Debugw("tts generation started", "language_code", languageCode, "gender", gender, "model", model, "text_len", len([]rune(text)))

	//Select a voice
	voices, err := c.tts.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{
		LanguageCode: languageCode,
	})
	if err != nil {
		c.logger.Errorw("failed to list tts voices", "error", err)
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
		c.logger.Debugw("no matching tts voice found", "language_code", languageCode, "gender", gender, "model", model)
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
		c.logger.Errorw("tts synthesize speech failed", "error", err)
		return nil, err
	}
	c.logger.Debugw("tts generation completed", "audio_size_bytes", len(resp.AudioContent))
	return resp.AudioContent, nil
}
