package service

import (
	"errors"

	"golang.org/x/text/language"
)

const (
	maxWordLength = 100
	maxHintLength = 200
)

var (
	ErrInvalidRequest  = errors.New("validation error")
	ErrHintTooLong     = errors.New("hint too long")
	ErrEmptyWord       = errors.New("empty word")
	ErrWordTooLong     = errors.New("word too long")
	ErrInvalidResponse = errors.New("invalid response")
)

func (req *GenerateSentenceRequest) validate() error {
	if err := validateWord(req.Word); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validateLanguageCode(req.WordLanguage); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validateLanguageCode(req.TranslationLanguage); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validateHint(req.TranslationHint); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}
	return nil
}

func (resp *GenerateSentenceResponse) validate() error {
	if resp.OriginalSentence == "" || resp.TranslatedSentence == "" {
		return ErrInvalidResponse
	}
	return nil
}

func (req *GenerateDefinitionRequest) validate() error {
	if err := validateWord(req.Word); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validateLanguageCode(req.Language); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validateHint(req.DefinitionHint); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}
	return nil
}

func (resp *GenerateDefinitionResponse) validate() error {
	if resp.Definition == "" {
		return ErrInvalidResponse
	}
	return nil
}

func (req *TranslateRequest) validate() error {
	if err := validateWord(req.Word); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validateLanguageCode(req.FromLanguage); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validateLanguageCode(req.ToLanguage); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}

	if err := validateHint(req.TranslationHint); err != nil {
		return errors.Join(err, ErrInvalidRequest)
	}
	return nil
}

func (resp *TranslateResponse) validate() error {
	if resp.Translation == "" {
		return ErrInvalidResponse
	}
	return nil
}

func validateWord(word string) error {
	switch {
	case word == "":
		return ErrEmptyWord
	case len([]rune(word)) > maxWordLength:
		return ErrWordTooLong
	}
	return nil
}

func validateLanguageCode(languageCode string) error {
	_, err := language.Parse(languageCode)
	return err
}

func validateHint(hint string) error {
	if len([]rune(hint)) > maxHintLength {
		return ErrHintTooLong
	}
	return nil
}
