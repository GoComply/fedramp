package template

import (
	"fmt"
	"regexp"
)

func parseControlId(str string) (name string, err error) {
	// matches controls and control enhancements, e.g. `AC-2`, `AC-2 (1)`, etc.
	regex := regexp.MustCompile(`[A-Z]{2}-\d+( +\(\d+\))?`)
	name = regex.FindString(str)
	if name == "" {
		err = fmt.Errorf("control name not found in '%s'", str)
	}
	return
}
