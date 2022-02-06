package tools

import (
	"os"
	"strconv"
	"time"
)

func ConvertStringToTime(input string) (time.Duration, error) {
	duration, err := strconv.Atoi(os.Getenv(input))
	if err != nil {
		return 0, err
	}
	return time.Duration(float64(duration) * float64(time.Second)), nil
}
