package server_routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"service/services/generator"
)

var ErrStatError = errors.New("stat error is not nil")
var ErrPathIsDir = errors.New("path is a directory")
var ErrFileRead = errors.New("file read error is not nil")

func getTemplate(name string) (template generator.ImageTemplate, err error) {
	fullpath := fmt.Sprintf("resources/templates/%s", name)
	fileInfo, statError := os.Stat(fullpath)

	if statError != nil {
		err = fmt.Errorf("%w: %w", ErrStatError, statError)
		return template, err
	}

	if fileInfo.IsDir() {
		return template, ErrPathIsDir
	}

	content, readError := os.ReadFile(fullpath)
	if readError != nil {
		err = fmt.Errorf("%w: %w", ErrStatError, readError)
		return template, err
	}

	unmarshalError := json.Unmarshal(content, &template)
	if unmarshalError != nil {
		return template, unmarshalError
	}

	return template, nil
}
