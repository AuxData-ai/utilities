package utilities

import (
	"net/url"
	"strings"
)

func ReplaceParametersInCommand(command string, parameters map[string]string) string {

	if len(parameters) > 0 {

		for key, value := range parameters {

			key := "${" + key + "}"
			command = strings.ReplaceAll(command, key, value)
		}

	}

	return command
}

func ReplaceParametersInUrl(command string, parameters map[string]string) string {

	if len(parameters) > 0 {

		for key, value := range parameters {

			key := "${" + key + "}"
			command = strings.ReplaceAll(command, key, CleanAndEncodeURL(value))
		}

	}

	return command
}

func CleanAndEncodeURL(str string) string {
	// Schritt 1: Lösche ungültige Zeichen
	//str = replaceInvalidCharacters(str)

	// Schritt 2: Entferne Leerzeichen
	//str = strings.ReplaceAll(str, " ", "%20")

	// URL Encoding
	urlStr := url.QueryEscape(str)

	return urlStr
}

func replaceInvalidCharacters(str string) string {
	// Liste der ungültigen Zeichen
	invalidCharacters := []string{";", "?", "#", "[", "]", "{", "}", "(", ")", "*", "+", ",", "/", ":", "@", "!", "$", "&", "'", "=", "|"}

	// Ersetzen ungültiger Zeichen durch das leerzeichen
	for _, char := range invalidCharacters {
		str = strings.ReplaceAll(str, char, "")
	}

	return str
}
