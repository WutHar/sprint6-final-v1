package service

import (
	"strings"

	"github.com/WutHar/sprint6-final/pkg/morse"
)

// DetectAndConvert определяет тип входящей строки (текст или код Морзе) и конвертирует её.
func DetectAndConvert(data string) (string, error) {

	isMorse := true
	for _, char := range data {
		if !strings.Contains(morse.Alphabet, strings.ToUpper(string(char))) && char != ' ' && char != '.' && char != '-' {
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
