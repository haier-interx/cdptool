package pipeline

import (
	"fmt"
	"github.com/haier-interx/cdptool/models"
)

type Result struct {
	error        error
	ExecuteTrace []*ExecutingStep

	JavaScriptResult []*[]byte
	Performances     []*models.PerformanceTiming
}

type ExecutingStep struct {
	Father string
	Index  int
}

func (es *ExecutingStep) Id() string {
	return fmt.Sprintf("%s-%d", es.Father, es.Index)
}

func NewResult() *Result {
	return &Result{
		ExecuteTrace:     make([]*ExecutingStep, 0),
		JavaScriptResult: make([]*[]byte, 0),
		Performances:     make([]*models.PerformanceTiming, 0),
	}
}

func (p *Result) PutExecuting(e *ExecutingStep) {
	p.ExecuteTrace = append(p.ExecuteTrace, e)
}

func (p *Result) Failed(err error) {
	p.error = err
}

func (p *Result) Error() error {
	return p.error
}

func (p *Result) ErrorCN() string {
	if p.error == nil {
		return ""
	}

	if len(p.ExecuteTrace) == 0 {
		return fmt.Sprintf("执行失败：%s", ErrorCN(p.error))
	} else {
		last_idx := len(p.ExecuteTrace) - 1
		return fmt.Sprintf("执行%s第%d步骤时失败：%s", p.ExecuteTrace[last_idx].Father, p.ExecuteTrace[last_idx].Index, ErrorCN(p.error))
	}
}

func (p *Result) LastExecutingStep() *ExecutingStep {
	if len(p.ExecuteTrace) == 0 {
		return &ExecutingStep{"", 0}
	} else {
		last_idx := len(p.ExecuteTrace) - 1
		return p.ExecuteTrace[last_idx]
	}
}
