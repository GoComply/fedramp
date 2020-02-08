package bundled

import (
	"errors"
	"github.com/GoComply/fedramp/pkg/fedramp/common"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
)

func TemplateDOCX(baseline common.BaselineLevel) (pkging.File, error) {
	switch baseline {
	case common.LevelLow:
		return pkger.Open("/bundled/templates/FedRAMP-SSP-Low-Baseline-Template.docx")
	case common.LevelModerate:
		return pkger.Open("/bundled/templates/FedRAMP-SSP-Moderate-Baseline-Template.docx")
	case common.LevelHigh:
		return pkger.Open("/bundled/templates/FedRAMP-SSP-High-Baseline-Template.docx")
	}
	return nil, errors.New("Not supported")
}
