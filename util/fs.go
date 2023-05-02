package util

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func PrepareOutputFolderPath(output string) error {
	dirPath := filepath.Dir(output)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		logrus.Debug("Creating necessary folders: ", dirPath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func SaveSnippet(snippet string, output string) error {
	logrus.Debug("Saving snippet to: ", output)
	if err := os.WriteFile(output, []byte(snippet), 0644); err != nil {
		return err
	}

	return nil
}
