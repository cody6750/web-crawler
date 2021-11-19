package webcrawler

import (
	"log"
	"strings"
)

func formatURL(url string, formatURLConfig FormatURLConfiguration) (string, error) {
	if isEmptyFormatURLConfiguration(formatURLConfig) {
		log.Print("formatURLConfig is empty")
		return url, nil
	}
	if strings.HasPrefix(url, formatURLConfig.PrefixExist) && strings.HasSuffix(url, formatURLConfig.SuffixExist) && strings.Contains(url, formatURLConfig.ReplaceOldString) {
		if formatURLConfig.ReplaceOldString != "" && formatURLConfig.ReplaceNewString != "" {
			url = strings.ReplaceAll(url, formatURLConfig.ReplaceOldString, formatURLConfig.ReplaceNewString)
		}
		if strings.HasPrefix(url, formatURLConfig.PrefixToRemove) {
			url = strings.TrimPrefix(url, formatURLConfig.PrefixToRemove)
		}
		if strings.HasSuffix(url, formatURLConfig.SuffixToRemove) {
			url = strings.TrimSuffix(url, formatURLConfig.SuffixToRemove)
		}
		if formatURLConfig.PrefixToAdd != "" {
			url = formatURLConfig.PrefixToAdd + url
		}
		if formatURLConfig.SuffixToAdd != "" {
			url = url + formatURLConfig.SuffixToAdd
		}
	}
	return url, nil
}

func isEmptyFormatURLConfiguration(formatURLConfiguration FormatURLConfiguration) bool {
	if formatURLConfiguration == (FormatURLConfiguration{}) {
		return true
	}
	return false
}
