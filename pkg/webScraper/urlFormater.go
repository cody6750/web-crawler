package webcrawler

import (
	"strings"
)

//FormatURLConfig ...
type FormatURLConfig struct {
	SuffixExist      string `json:"SuffixExist"`
	SuffixToAdd      string `json:"SuffixToAdd"`
	SuffixToRemove   string `json:"SuffixToRemove"`
	PrefixToAdd      string `json:"PrefixToAdd"`
	PrefixExist      string `json:"PrefixExist"`
	PrefixToRemove   string `json:"PrefixToRemove"`
	ReplaceOldString string `json:"ReplaceOldString"`
	ReplaceNewString string `json:"ReplaceNewString"`
}

func formatURL(url string, config FormatURLConfig) string {
	if strings.HasPrefix(url, config.PrefixExist) && strings.HasSuffix(url, config.SuffixExist) && strings.Contains(url, config.ReplaceOldString) {
		if config.ReplaceOldString != "" && config.ReplaceNewString != "" {
			url = strings.ReplaceAll(url, config.ReplaceOldString, config.ReplaceNewString)
		}

		url = strings.TrimPrefix(url, config.PrefixToRemove)
		url = strings.TrimSuffix(url, config.SuffixToRemove)

		if config.PrefixToAdd != "" {
			url = config.PrefixToAdd + url
		}

		if config.SuffixToAdd != "" {
			url = url + config.SuffixToAdd
		}
	} else {
		return ""
	}
	return url
}
