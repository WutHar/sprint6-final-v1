package service

import (
	"github.com/WutHar/sprint6-final/pkg/morse"
)

func DetectAndConvert(data string) (string, error) {
	isMorse := true
	for _, char := range data {
		if char != '.' && char != '-' && char != ' ' {
			isMorse = false
			break
		}
	}

	if isMorse {
		text := morse.ToText(data)
		return text, nil
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
