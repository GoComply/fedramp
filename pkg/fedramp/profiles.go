package fedramp

import (
	"github.com/GoComply/fedramp/bundled"
	"github.com/GoComply/fedramp/pkg/fedramp/common"
	"github.com/docker/oscalkit/types/oscal/profile"
)

type Baseline struct {
	level   common.BaselineLevel
	profile profile.Profile
}

func New(baselineLevel common.BaselineLevel) (*Baseline, error) {
	var result Baseline
	result.level = baselineLevel
	file, err := bundled.ProfileOSCAL(baselineLevel)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return &result, nil
}

func (b *Baseline) ProfileURL() string {
	return common.ProfileUrls[b.level]
}
