package job

import (
	"context"
	"testing"
)

func TestInstance_Start(t *testing.T) {
	yml := `pipelines:
  - id: vue-element-admin
    timeout: 10s
    steps:
      - type: init
      - type: login
      - type: _screenshot_
        sel: '.app-main'
        screenshot: 
          quality: 90
      - type: _performance_

definitions:
  steps:
    - id: init
      steps:
        - type: _deviceScreen_
        - type: _language_
    - id: login
      steps: 
        - type: _navigate_
          url: 'https://panjiachen.gitee.io/vue-element-admin/#/login?redirect=%2Fdashboard'
        - type: _input_
          sel: '.el-input__inner'
          node_idx: 1
          input: 'admin123'
        - type: _click_
          sel: "button"
        - type: _wait_
          sel: ".app-main"`

	p, err := NewFromYaml([]byte(yml))
	if err != nil {
		t.Fatal(err)
	}

	rets := p.Start(context.Background())
	for idx, ret := range rets {
		if ret.Error() != nil {
			t.Fatalf("%s: %s", p.Pipelines[idx].Id, ret.ErrorCN())
		}
		t.Logf("%s: %+v", p.Pipelines[idx].Id, ret.Performances[0])
	}
}
