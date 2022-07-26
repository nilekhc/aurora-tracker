package utils

import (
	"fmt"
	"strconv"
)

// ConvertToInt converts a string to an int
func ConvertToInt(kpIndex string) int {
	intValue, err := strconv.ParseInt(kpIndex, 0, 32)
	if err != nil {
		panic(fmt.Errorf("Error converting kp index to int: %v", err))
	}

	return int(intValue)
}
