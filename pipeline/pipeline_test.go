package pipeline

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestPipeline_Execute(t *testing.T) {
	yml := `id: vue-element-admin
timeout: 5s
steps:
  - type: deviceScreen
  - type: language
  - type: navigate
    url: 'https://panjiachen.gitee.io/vue-element-admin/#/login?redirect=%2Fdashboard'
  #- type: dump
  - type: input
    sel: '.el-input__inner'
    node_idx: 1
    input: 'admin123'
  - type: click
    sel: "button"
  - type: screenshot
    sel: ".app-main"
    screenshot: 
      quality: 90
  - type: performance`

	p := new(Pipeline)
	err := yaml.Unmarshal([]byte(yml), p)
	if err != nil {
		t.Fatal(err)
	}

	ret := p.Run(context.Background())
	if ret.ErrorStepIdx() != -1 {
		t.Fatal(ret.ErrorCN())
	}

	fmt.Printf("performance: %+v\n", ret.PerformanceResults()[0])

}
