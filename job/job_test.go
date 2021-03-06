package job

import (
	"context"
	"testing"
)

func TestInstance_Start(t *testing.T) {
	yml := `pipelines:
  - id: vue-element-admin
    timeout: 10s
    network_enable: true
    steps:
      - type: init
      - type: login
      - type: _screenshots_
        sel: '.app-main'
        screenshots: 
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
			for i, item := range ret.NetworkLogs.Items {
				t.Logf("network[%d] docURL:%s\n", i, item.DocURL)
				t.Logf("network[%d] method:%s name:%s\n", i, item.Method, item.Name)
				t.Logf("network[%d] finish:%v failed:%v\n", i, item.IsFinished(), item.IsFailed())
				t.Logf("network[%d] status:%d time:%d\n", i, item.Status, item.Time)
				t.Logf("network[%d] detail:%+v\n", i, item.Metrics())
			}

			t.Fatalf("%s: %s", p.Pipelines[idx].Id, ret.Error())
		}
		t.Logf("%s: %+v", p.Pipelines[idx].Id, ret.Performances[0])
	}
}
