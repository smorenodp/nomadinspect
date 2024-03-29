package screens

import (
	"bytes"
	"encoding/json"
)

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
