package utilities

import (
	"errors"
	"strings"
)

func GetJsonPartOfResponse(response string) (string, error) {
	var jsonPart string
	start := strings.Index(string(response), "{")
	end := strings.LastIndex(string(response), "}")
	if start >= 0 && end >= 0 {
		jsonPart = response[start : end+1]
	} else {
		if start < 0 && end < 0 {
			return "", errors.New("json part not found in response")
		} else if start < 0 {
			return "", errors.New(`"{" not found in response`)
		} else if end < 0 {
			return "", errors.New(`"}" not found in response`)
		}
	}

	return jsonPart, nil
}
