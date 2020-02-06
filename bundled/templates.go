package bundled

import (
	"errors"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

type FedRAMPBaseline int

const (
	FedRAMPUnknown = iota
	FedRAMPLow
	FedRAMPModerate
	FedRAMPHigh
)

func TemplateDOCX(baseline FedRAMPBaseline) (pkging.File, error) {
	switch baseline {
	case FedRAMPLow:
		return pkger.Open("/bundled/templates/FedRAMP-SSP-Low-Baseline-Template.docx")
	case FedRAMPModerate:
		return pkger.Open("/bundled/templates/FedRAMP-SSP-Moderate-Baseline-Template.docx")
	case FedRAMPHigh:
		return pkger.Open("/bundled/templates/FedRAMP-SSP-High-Baseline-Template.docx")
	}
	return nil, errors.New("Not supported")
}
