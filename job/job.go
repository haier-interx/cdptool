package job

import (
	"context"
	"github.com/haier-interx/cdptool/pipeline"
	"gopkg.in/yaml.v2"
)

type Instance struct {
	Pipelines   []*pipeline.Pipeline        `json:"pipelines" yaml:"pipelines"`
	Definitions *pipeline.CustomDefinitions `json:"definitions" yaml:"definitions"`
}

func NewFromYaml(body []byte) (*Instance, error) {
	i := new(Instance)
	err := yaml.Unmarshal(body, i)
	if err != nil {
		return nil, err
	}

	if err := i.Definitions.Init(); err != nil {
		return nil, err
	}

	return i, err
}

func (i *Instance) Start(ctx context.Context) []*pipeline.Result {
	cds := i.Definitions
	rets := make([]*pipeline.Result, len(i.Pipelines))
	for idx, p := range i.Pipelines {
		rets[idx] = p.Run(ctx, cds)
	}
	return rets
}

func (i *Instance) Check(ctx context.Context) error {
	for _, p := range i.Pipelines {
		_, ret_tmp := p.Parse(ctx, i.Definitions)
		if ret_tmp.Error() != nil {
			return ret_tmp.Error()
		}
	}
	return nil
}
