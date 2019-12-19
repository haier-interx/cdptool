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
      quality: 90
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
}
