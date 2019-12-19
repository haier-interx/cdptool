package pipeline

import (
	"fmt"
	"github.com/haier-interx/cdptool/models"
)

type Result struct {
	parseErrorStep *ExecutingStep
	error          error
	ExecuteTrace   []*ExecutingStep

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

func (p *Result) ParseFailed(id string, index int, err error) {
	p.parseErrorStep = &ExecutingStep{id, index}
	p.error = err
}

func (p *Result) Error() error {
	if p.error == nil {
		return nil
	}

	if len(p.ExecuteTrace) == 0 {
		if p.parseErrorStep != nil {
			return fmt.Errorf("parse the %d step failed on task \"%s\"：%s", p.parseErrorStep.Index+1, p.parseErrorStep.Father, p.error)
		} else {
			return p.error
		}
	} else {
		return fmt.Errorf("execute the %d step failed on task \"%s\"：%w", p.LastExecutingStep().Index+1, p.LastExecutingStep().Father, p.error)
	}
}

func (p *Result) ErrorCN() error {
	if p.error == nil {
		return nil
	}

	if len(p.ExecuteTrace) == 0 {
		if p.parseErrorStep != nil {
			return fmt.Errorf("解析%s第%d步骤时失败：%s", p.parseErrorStep.Father, p.parseErrorStep.Index+1, ErrorCN(p.error))
		} else {
			return fmt.Errorf("%s", ErrorCN(p.error))
		}
	} else {
		last_idx := len(p.ExecuteTrace) - 1
		return fmt.Errorf("执行%s第%d步骤时失败：%s", p.ExecuteTrace[last_idx].Father, p.ExecuteTrace[last_idx].Index+1, ErrorCN(p.error))
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
