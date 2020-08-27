/*
 Copyright (C) 2018 OpenControl Contributors. See LICENSE.md for license.
*/

package standards

import (
	"github.com/opencontrol/compliance-masonry/pkg/lib/common"
	v1_0_0 "github.com/opencontrol/compliance-masonry/pkg/lib/standards/versions/1_0_0"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Load will read the file at the given path and attempt to return a standard object.
func Load(path string) (common.Standard, error) {
	var standard v1_0_0.Standard
	standardData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, common.ErrReadFile
	}
	err = yaml.Unmarshal(standardData, &standard)
	if err != nil {
		return nil, common.ErrStandardSchema
	}
	return standard, nil
}
