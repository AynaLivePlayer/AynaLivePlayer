package config

import (
	"encoding/json"
	"os"
)

func LoadJson(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func SaveJson(path string, dst any) error {
	data, err := json.MarshalIndent(dst, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0666)
}
