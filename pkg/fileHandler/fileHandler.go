package filehandler

import (
	"encoding/json"
	"fmt"
	"os"
)

func wrapError(err error) error {
	return fmt.Errorf("jsonHandler: %w", err)
}

func ReadFile[T any](filePath string) (*T, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, wrapError(err)
	}

	var jsonObject T

	err = json.Unmarshal(fileContent, &jsonObject)
	if err != nil {
		return nil, wrapError(err)
	}

	return &jsonObject, nil
}

func WriteFile(filePath string, jsonObject any) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return wrapError(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(jsonObject)
	if err != nil {
		return wrapError(err)
	}

	return nil
}
