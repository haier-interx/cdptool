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
  - type: _deviceScreen_
  - type: _language_
  - type: _navigate_
    url: 'https://panjiachen.gitee.io/vue-element-admin/#/login?redirect=%2Fdashboard'
  #- type: _dump_
  - type: _input_
    sel: '.el-input__inner'
    node_idx: 1
    input: 'admin123'
  - type: _click_
    sel: "button"
  - type: _javascript_
    javascript: "window.location.href"
  - type: _screenshots_
    sel: ".app-main"
    screenshots: 
      quality: 20
  - type: _network_
  - type: _performance_`

	p := new(Pipeline)
	err := yaml.Unmarshal([]byte(yml), p)
	if err != nil {
		t.Fatal(err)
	}

	ret := p.Run(context.Background(), nil)
	t.Logf("error step index was %s", ret.LastExecutingStep().Id())
	if ret.error != nil {
		t.Fatal(ret.ErrorCN())
	}

	fmt.Printf("javascript: %s\n", *ret.JavaScriptResult[0])
	fmt.Printf("performance: %+v\n", ret.Performances[0].Metrics())
	for i, item := range *ret.NetworkPerformances[0] {
		fmt.Printf("network[%d]: %s %s %+v\n", i, item.Name, item.InitiatorType, item.Metrics())
	}

	for i, item := range ret.StepResult() {
		fmt.Printf("step[%d] %s: %+v\n", i, p.Steps[i].Type, item)
	}

	for i, item := range ret.ScreenshotsFileName {
		fmt.Printf("screenshots[%d]: %s\n", i, item)
	}
}
