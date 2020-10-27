package nominal_catalog

import (
	"fmt"
	"strings"
)

func (param *Param) TextRepresentation() string {
	if param.Label != nil {
		return fmt.Sprintf("[Assignment: %s]", param.Label.Raw)

	}
	if param.Select != nil {
		howMany := ""
		if len(param.Select.HowMany) != 0 {
			howMany = fmt.Sprintf(" (%s)", param.Select.HowMany)
		}

		choicesList := make([]string, len(param.Select.Alternatives))
		for i, v := range param.Select.Alternatives {
			choicesList[i] = strings.TrimSpace(string(v.Raw))
		}
		choices := strings.Join(choicesList, ", ")

		return fmt.Sprintf("[Selection%s: %s]", howMany, choices)
	}

	return "[TODO: Not implemented parameter]"

}
