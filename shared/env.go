package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// getEnvBool converts string environment variables to booleans.
func GetEnvBool(envVar string) (bool, error) {
	s := os.Getenv(envVar)
	if s == "" {
		return false, fmt.Errorf("")
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

// getEnvBool converts string environment variables to integers.
func GetEnvInt(envVar string) (int, error) {
	s := os.Getenv(envVar)
	if s == "" {
		return 0, fmt.Errorf("")
	}
	strconv.Atoi(s)
	v, err := strconv.Atoi(s)
	if err != nil {
		return v, err
	}
	return v, nil
}

//GetEnvTime converts string to time duration
func GetEnvTime(input string) (time.Duration, error) {
	duration, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return time.Duration(float64(duration) * float64(time.Second)), nil
}
