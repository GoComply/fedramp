package fedramp

type ImplementationStatus uint8

// implemented, partial, planned, alternate, not-applicable
const (
	StatusNoStatus ImplementationStatus = iota
	StatusImplemented
	StatusPartiallyImplemented
	StatusPlanned
	StatusAlternativeImplementation
	StatusNotApplicable
)

var fromOSCAL = map[string]ImplementationStatus{
	"unknown":        StatusNoStatus,
	"implemented":    StatusImplemented,
	"partial":        StatusPartiallyImplemented,
	"planned":        StatusPlanned,
	"alternate":      StatusAlternativeImplementation,
	"not-applicable": StatusNotApplicable,
}

var humanString = map[ImplementationStatus]string{
	StatusNoStatus:                  "Unknown",
	StatusImplemented:               "Implemented",
	StatusPartiallyImplemented:      "Partially implemented",
	StatusPlanned:                   "Planned",
	StatusAlternativeImplementation: "Alternative Implementation",
	StatusNotApplicable:             "Not applicable",
}

func StatusFromOSCAL(status string) ImplementationStatus {
	s, found := fromOSCAL[status]
	if !found {
		return StatusNoStatus
	}
	return s
}

func (is ImplementationStatus) HumanString() string {
	return humanString[is]
}
