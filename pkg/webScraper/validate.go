package webcrawler

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

//ComparisonConfiguration ...
type FilterConfiguration struct {
	IsLessThan            interface{} `json:"IsLessThan"`
	IsGreaterThan         interface{} `json:"IsGreaterThan"`
	IsEqualTo             interface{} `json:"IsEqualTo"`
	IsNotEqualTo          interface{} `json:"IsNotEqualTo"`
	Contains              string      `json:"Contains"`
	ConvertStringToNumber string      `json:"ConvertStringToNumber"`
}

func Validate(v interface{}, c *FilterConfiguration) bool {

	switch v := v.(type) {
	case int:
		log.Print("int")
		return ValidateInt(v, c)
	case string:
		if c.ConvertStringToNumber == "true" {
			number := ConvertStringToNunber(v)
			return Validate(number, c)
		}
		return ValidateString(v, c)
	case float64:
		log.Print("float")
		return ValidateFloat64(v, c)
	default:
		return false
	}
}

func ValidateString(s string, c *FilterConfiguration) bool {
	if !strings.Contains(s, c.Contains) {
		return false
	}

	switch t := c.IsEqualTo.(type) {
	case string:
		if t != s && t != "" {
			return false
		}
	}
	switch t := c.IsNotEqualTo.(type) {
	case string:
		if t == s && t != "" {
			return false
		}
	}
	return true
}

func ValidateInt(i int, c *FilterConfiguration) bool {
	switch t := c.IsEqualTo.(type) {
	case int:
		if t != i && t != 0 {
			return false
		}
	}
	switch t := c.IsNotEqualTo.(type) {
	case int:
		if t == i && t != 0 {
			return false
		}
	}

	switch t := c.IsLessThan.(type) {
	case int:
		if t < i && t != 0 {
			return false
		}
	}

	switch t := c.IsGreaterThan.(type) {
	case int:
		if t > i && t != 0 {
			return false
		}
	}

	return true
}

func ValidateFloat64(f float64, c *FilterConfiguration) bool {
	switch t := c.IsEqualTo.(type) {
	case float64:
		if t != f && t != 0 {
			return false
		}
	}
	switch t := c.IsNotEqualTo.(type) {
	case float64:
		if t == f && t != 0 {
			return false
		}
	}

	switch t := c.IsLessThan.(type) {
	case float64:
		if t < f && t != 0.0 {
			return false
		}
	}
	switch t := c.IsGreaterThan.(type) {
	case float64:
		if t > f && t != 0.0 {
			return false
		}
	}
	return true
}

func ConvertStringToNunber(s string) interface{} {
	s = strings.ReplaceAll(s, ",", "")
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	match := re.FindAllString(s, -1)
	if match != nil {
		var numb interface{}
		if strings.Contains(match[0], ".") {
			numb, err = strconv.ParseFloat(match[0], 64)

		} else {
			numb, err = strconv.Atoi(match[0])
		}
		if err != nil {
			log.Print(err)
			return err
		}
		return numb

	}
	return nil
}
