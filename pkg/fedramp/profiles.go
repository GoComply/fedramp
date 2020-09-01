package fedramp

import (
	"fmt"
	"github.com/gocomply/fedramp/bundled"
	"github.com/gocomply/fedramp/pkg/fedramp/common"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/pkg/oscal_source"
	"github.com/gocomply/oscalkit/types/oscal/catalog"
	"github.com/gocomply/oscalkit/types/oscal/profile"
)

type Baseline struct {
	Level   common.BaselineLevel
	profile *profile.Profile
	catalog *catalog.Catalog
}

func NewBaseline(baselineLevel common.BaselineLevel) (*Baseline, error) {
	var result Baseline
	result.Level = baselineLevel
	file, err := bundled.ProfileOSCAL(baselineLevel)
	if err != nil {
		return nil, fmt.Errorf("could not initiate FedRAMP: could not open internal files: %v", err)
	}
	defer file.Close()
	source, err := oscal_source.OpenFromReader(file.Name(), file)
	if err != nil {
		return nil, err
	}
	defer source.Close()
	oscal := source.OSCAL()
	if oscal.DocumentType() != constants.ProfileDocument {
		return nil, fmt.Errorf("Could not initiate FedRAMP. Expected profile element in %s", file.Name())
	}
	result.profile = oscal.Profile

	file, err = bundled.CatalogOSCAL(baselineLevel)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	source, err = oscal_source.OpenFromReader(file.Name(), file)
	if err != nil {
		return nil, err
	}
	defer source.Close()
	oscal = source.OSCAL()
	if oscal.DocumentType() != constants.CatalogDocument {
		return nil, fmt.Errorf("Could not initiate FedRAMP. Expected catalog element in %s", file.Name())
	}
	result.catalog = oscal.Catalog

	return &result, nil
}

func AvailableBaselines() ([]Baseline, error) {
	var result []Baseline
	var level common.BaselineLevel
	for level = common.LevelLow; level <= common.LevelHigh; level++ {
		baseline, err := NewBaseline(level)
		if err != nil {
			return nil, err
		}
		result = append(result, *baseline)
	}
	return result, nil
}

func (b *Baseline) ProfileURL() string {
	return common.ProfileUrls[b.Level]
}

func (b *Baseline) Controls() []catalog.Control {
	return b.catalog.Controls
}

func (b *Baseline) ControlGroups() []catalog.Group {
	return b.catalog.Groups
}

func (b *Baseline) FindSetParam(id string) *profile.SetParameter {
	if b.profile.Modify == nil {
		return nil
	}
	for _, setParam := range b.profile.Modify.ParameterSettings {
		if setParam.ParamId == id {
			return &setParam
		}
	}
	return nil

}
