package util

import (
	"os"
)

func SaveSnippet(snippet string, output string) error {
	if err := os.WriteFile(output, []byte(snippet), 0644); err != nil {
		return err
	}

	return nil
}
