package webcrawler

import (
	"strings"
)

func formatURL(url string, formatURLConfig FormatURLConfiguration) string {
	if strings.HasPrefix(url, formatURLConfig.PrefixExist) && strings.HasSuffix(url, formatURLConfig.SuffixExist) && strings.Contains(url, formatURLConfig.ReplaceOldString) {
		if formatURLConfig.ReplaceOldString != "" && formatURLConfig.ReplaceNewString != "" {
			url = strings.ReplaceAll(url, formatURLConfig.ReplaceOldString, formatURLConfig.ReplaceNewString)
		}
		url = strings.TrimPrefix(url, formatURLConfig.PrefixToRemove)
		url = strings.TrimSuffix(url, formatURLConfig.SuffixToRemove)
		if formatURLConfig.PrefixToAdd != "" {
			url = formatURLConfig.PrefixToAdd + url
		}
		if formatURLConfig.SuffixToAdd != "" {
			url = url + formatURLConfig.SuffixToAdd
		}
	} else {
		return ""
	}
	return url
}
