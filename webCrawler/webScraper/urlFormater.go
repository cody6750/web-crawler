package webcrawler

import (
	"errors"
	"log"
	"strings"
)

func formatURL(url string, formatURLConfig FormatURLConfiguration) (string, error) {
	if IsEmpty(formatURLConfig) {
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
	} else {
		return "", errors.New("Does not match")
	}
	return url, nil
}
