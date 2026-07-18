package utils

import "encoding/json"

func ToJson(v interface{}) (string, error) {
	json, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(json), nil
}