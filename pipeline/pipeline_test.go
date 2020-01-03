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
network_enable: true
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

	fmt.Printf("--------------- javascript ----------------\n")
	fmt.Printf("javascript: %s\n", *ret.JavaScriptResult[0])

	fmt.Printf("--------------- performance ----------------\n")
	fmt.Printf("performance: %+v\n", ret.Performances[0].Metrics())

	for i, item := range ret.ScreenshotsFileName {
		fmt.Printf("--------------- screenshots ----------------\n")
		fmt.Printf("screenshots[%d]: %s\n", i, item)
	}

	fmt.Printf("--------------- stepTiming ----------------\n")
	fmt.Printf("init: %dms\n", ret.InitDuration.Milliseconds())
	total := ret.InitDuration
	for i, item := range ret.StepResult() {
		//fmt.Printf("step[%d] %s %dms: %+v\n", i, p.Steps[i].Type, item.Duration.Milliseconds(), item)
		fmt.Printf("step[%d] %s %dms\n", i, p.Steps[i].Type, item.Duration.Milliseconds())
		total += item.Duration
	}
	fmt.Printf("total: %dms\n", total.Milliseconds())

	fmt.Printf("--------------- network ----------------\n")
	for i, item := range ret.NetworkLogs.Items {
		fmt.Printf("[%d] docURL:%s\n", i, item.DocURL)
		fmt.Printf("[%d] method:%s name:%s\n", i, item.Method, item.Name)
		fmt.Printf("[%d] finish:%v failed:%v\n", i, item.IsFinished(), item.IsFailed())
		fmt.Printf("[%d] status:%d time:%d\n", i, item.Status, item.Time)
		fmt.Printf("[%d] detail:%+v\n", i, item.Metrics())
		fmt.Printf("---------------\n")
	}

}
