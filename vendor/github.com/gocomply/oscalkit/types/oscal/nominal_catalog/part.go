package nominal_catalog

import (
	"regexp"
	"strings"
)

type ParamResolver interface {
	FindParam(id string) *Param
}

func (p *Part) ResolveInserts(resolver ParamResolver) string {
	re := regexp.MustCompile(`<insert param-id="([a-z0-9._-]*)"/>`)
	resolved := p.Prose.Raw

	for {
		match := re.FindStringSubmatch(resolved)
		if len(match) == 0 {
			break
		}
		resolved = strings.ReplaceAll(resolved, match[0], resolver.FindParam(match[1]).TextRepresentation())
	}

	return resolved
}
