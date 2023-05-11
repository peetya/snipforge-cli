package util

import (
	"os"
	"path/filepath"
)

func PrepareOutputFolderPath(output string) error {
	//dirPath := filepath.Dir(output)
	//if _, err := os.Stat(dirPath); os.IsNotExist(err) {
	//	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func SaveSnippet(snippet string, output string) error {
	dirPath := filepath.Dir(output)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	if err := os.WriteFile(output, []byte(snippet), 0644); err != nil {
		return err
	}

	return nil
}
