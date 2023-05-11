package prompt

import "github.com/charmbracelet/bubbles/list"

type OpenaiModelItem string

func (i OpenaiModelItem) FilterValue() string { return "" }

func ConvertItemToString(item list.Item) string {
	return string(item.(OpenaiModelItem))
}
