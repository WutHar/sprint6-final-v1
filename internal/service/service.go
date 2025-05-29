package service

import (
	"log"
	"strings"
	"unicode"

	"github.com/WutHar/sprint6-final-v1/pkg/morse"
)

func DetectAndConvert(data string) (string, error) {
	log.Printf("Входящие данные в DetectAndConvert: %q", data)

	isMorse := true
	for _, char := range data {
		if !unicode.Is(unicode.Dash, char) && char != '.' && char != ' ' {
			isMorse = false
			break
		}
	}

	if isMorse && strings.ContainsAny(data, ".-") { // Дополнительная проверка на наличие хотя бы одного символа Морзе
		text := morse.ToText(data)
		log.Printf("Определено как код Морзе, результат ToText: %q", text)
		return text, nil
	} else {
		morseCode := morse.ToMorse(data)
		log.Printf("Определено как обычный текст, результат ToMorse: %q", morseCode)
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
