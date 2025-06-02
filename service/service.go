package service

import (
	"strings"

	"github.com/WutHar/sprint6-final-v1/pkg/morse"
)

func DetectAndConvert(data string) (string, error) {

	isMorse := func(r rune) bool { return r == '.' || r == '-' || r == ' ' || r == '/' }

	if strings.ContainsFunc(data, isMorse) {

		return morse.ToText(data), nil

	}

	return morse.ToMorse(data), nil
}
