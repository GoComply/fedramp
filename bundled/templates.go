package bundled

import (
	"errors"
	"github.com/gocomply/fedramp/pkg/fedramp/common"
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

func TemplateOSCAL() (pkging.File, error) {
	return pkger.Open("/bundled/templates/FedRAMP-SSP-OSCAL-Template.xml")
}

func CatalogOSCAL(baseline common.BaselineLevel) (pkging.File, error) {
	switch baseline {
	case common.LevelLow:
		return pkger.Open("/bundled/catalogs/FedRAMP_LOW-baseline-resolved-profile_catalog.xml")
	case common.LevelModerate:
		return pkger.Open("/bundled/catalogs/FedRAMP_MODERATE-baseline-resolved-profile_catalog.xml")
	case common.LevelHigh:
		return pkger.Open("/bundled/catalogs/FedRAMP_HIGH-baseline-resolved-profile_catalog.xml")
	}
	return nil, errors.New("Not supported")
}
