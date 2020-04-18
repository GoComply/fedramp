package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func ControlKeyToOSCAL(controlKey string) string {
	lower := strings.ToLower(controlKey)
	re := regexp.MustCompile("^([a-z][a-z])-([0-9]+)(\\s+\\(([0-9]+)\\))?$")
	match := re.FindStringSubmatch(lower)
	result := fmt.Sprintf("%s-%s", match[1], match[2])
	if match[4] != "" {
		result = fmt.Sprintf("%s.%s", result, match[4])
	}
	return result
}
