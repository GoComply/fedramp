package oc2oscal

import (
	"fmt"
	"github.com/GoComply/fedramp/pkg/utils"
	"github.com/opencontrol/compliance-masonry/pkg/lib/common"
)

type Component struct {
	component common.Component
	satisfies map[string]common.Satisfies
}

func NewComponent(component common.Component) (*Component, error) {
	result := Component{
		component: component,
		satisfies: map[string]common.Satisfies{},
	}

	for _, sat := range component.GetAllSatisfies() {
		id := utils.ControlKeyToOSCAL(sat.GetControlKey())

		if _, ok := result.satisfies[id]; ok {
			return nil, fmt.Errorf("Duplicate key %s found in component %s", sat.GetControlKey(), component.GetKey())
		}
		result.satisfies[id] = sat
	}
	return &result, nil
}

func (c *Component) GetSatisfy(id string) common.Satisfies {
	return c.satisfies[id]
}

func (c *Component) GetKey() string {
	return c.component.GetKey()
}

func (c *Component) GetName() string {
	return c.component.GetName()
}
