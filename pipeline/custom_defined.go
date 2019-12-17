package pipeline

import (
	"fmt"
)

type CustomDefinitions struct {
	StepGroups    []*StepGroup `json:"steps" yaml:"steps"`
	stepGroupById map[string]*StepGroup
}

type StepGroup struct {
	Id    string  `json:"id" yaml:"id"`
	Steps []*Step `json:"steps" yaml:"steps"`
}

func (cds *CustomDefinitions) Init() error {
	cds.stepGroupById = make(map[string]*StepGroup)

	if cds.StepGroups != nil {
		for _, sg := range cds.StepGroups {
			_, found := cds.stepGroupById[sg.Id]
			if found {
				return fmt.Errorf("%s%w", sg.Id, ERR_STEPDEFINED_REPEAT)
			}
			cds.stepGroupById[sg.Id] = sg

			// 初始化step id
			for i, s := range sg.Steps {
				s.SetId(sg.GenerateStepId(i))
			}
		}
	}

	return nil
}

func (cds *CustomDefinitions) Steps(step_type string) (*StepGroup, bool) {
	sg, found := cds.stepGroupById[step_type]
	return sg, found
}

func (sg *StepGroup) GenerateStepId(idx int) string {
	return fmt.Sprintf("%s_%d_%s", sg.Id, idx, sg.Steps[idx].Type)
}
