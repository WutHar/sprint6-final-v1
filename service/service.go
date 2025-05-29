package service

import (
	"strings"

	"github.com/WutHar/sprint6-final-v1/pkg/morse"
)

func DetectAndConvert(data string) (string, error) {
	isMorse := func(r rune) bool { return r == '.' || r == '-' || r == ' ' }

	if strings.ContainsFunc(data, isMorse) {
		// Для кода Морзе просто вызываем ToText
		return morse.ToText(data), nil
	}

	// Для обычного текста вызываем ToMorse
	return morse.ToMorse(data), nil
}
