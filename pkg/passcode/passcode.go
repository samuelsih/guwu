package passcode

import (
	"fmt"
	"time"
)

func Generate(length int) string {
	if length <= 0 {
		return ""
	}

	currentTime := time.Now().Nanosecond()
	passcode := fmt.Sprint(currentTime)

	if len(passcode) <= length {
		passcode += passcode[:1]
	}

	return passcode[:length]
}