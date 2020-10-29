package fedramp

import (
	"fmt"
	"github.com/gocomply/fedramp/bundled"
	"github.com/gocomply/fedramp/pkg/fedramp/common"
	"github.com/gocomply/oscalkit/pkg/oscal/constants"
	"github.com/gocomply/oscalkit/pkg/oscal_source"
	"github.com/gocomply/oscalkit/types/oscal/catalog"
)

type Baseline struct {
	Level   common.BaselineLevel
	catalog *catalog.Catalog
}

func NewBaseline(baselineLevel common.BaselineLevel) (*Baseline, error) {
	var result Baseline
	result.Level = baselineLevel
	return &result, result.loadCatalog()
}

func (baseline *Baseline) loadCatalog() error {
	file, err := bundled.CatalogOSCAL(baseline.Level)
	if err != nil {
		return err
	}
	defer file.Close()
	source, err := oscal_source.OpenFromReader(file.Name(), file)
	if err != nil {
		return err
	}
	defer source.Close()
	oscal := source.OSCAL()
	if oscal.DocumentType() != constants.CatalogDocument {
		return fmt.Errorf("Could not initiate FedRAMP. Expected catalog element in %s", file.Name())
	}
	baseline.catalog = oscal.Catalog

	return nil
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

func (b *Baseline) FindParam(controlId, id string) (*catalog.Param, error) {
	ctrl := b.catalog.FindControlById(controlId)
	if ctrl == nil {
		fmt.Errorf("could not find control '%s' in FedRAMP %s baseline", controlId, b.Level.Name())
	}
	return ctrl.FindParamById(id), nil
}
