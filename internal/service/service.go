package service

import (
	"strings"
	"unicode"

	"github.com/WutHar/sprint6-final-v1/pkg/morse"
)

func DetectAndConvert(data string) (string, error) {

	textFromMorse := morse.ToText(data)
	textFromMorse = strings.TrimSpace(textFromMorse)

	isLikelyMorse := false
	for _, r := range textFromMorse {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			isLikelyMorse = true
			break
		}
	}

	if isLikelyMorse {
		return textFromMorse, nil
	} else {

		morseCode := morse.ToMorse(data)
		return morseCode, nil
	}
}

type Converter struct {
	morseConverter morse.Converter
}

func NewConverter() *Converter {
	return &Converter{
		morseConverter: morse.DefaultConverter,
	}
}

func (c *Converter) TextToMorse(text string) string {
	return c.morseConverter.ToMorse(text)
}

func (c *Converter) MorseToText(morseCode string) string {
	return c.morseConverter.ToText(morseCode)
}
