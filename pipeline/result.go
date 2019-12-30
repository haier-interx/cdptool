package pipeline

import (
	"fmt"
	"github.com/haier-interx/cdptool/models"
	"time"
)

type Result struct {
	PipelineId     string
	parseErrorStep *StepResult
	error          error
	ExecuteTrace   []*StepResult

	JavaScriptResult          []*[]byte
	Performances              []*models.PerformanceTiming
	NetworkPerformances       []*NetworkPerformance
	ScreenshotsFileName       []string
	FailedScreenshotsFileName string
}

type NetworkPerformance = []*models.PerformanceTiming

type StepResult struct {
	Father string
	Index  int
	Type   string

	StarTime time.Time
	EndTime  time.Time
	Duration time.Duration
}

func NewStepResult(father string, idx int, stepType string) *StepResult {
	return &StepResult{
		Father: father,
		Index:  idx,
		Type:   stepType,
	}
}

func (es *StepResult) Id() string {
	return fmt.Sprintf("%s-%d", es.Father, es.Index)
}

func NewResult(pipelineId string) *Result {
	return &Result{
		PipelineId:          pipelineId,
		ExecuteTrace:        make([]*StepResult, 0),
		JavaScriptResult:    make([]*[]byte, 0),
		Performances:        make([]*models.PerformanceTiming, 0),
		NetworkPerformances: make([]*NetworkPerformance, 0),
		ScreenshotsFileName: make([]string, 0),
	}
}

func (p *Result) SetStepStarted(e *StepResult) {
	e.StarTime = time.Now()
	p.ExecuteTrace = append(p.ExecuteTrace, e)
}

func (p *Result) SetStepOver(e *StepResult) {
	e.EndTime = time.Now()
	e.Duration = e.EndTime.Sub(e.StarTime)
}

func (p *Result) Failed(err error) {
	p.error = err
}

func (p *Result) ParseFailed(id string, index int, err error) {
	p.parseErrorStep = NewStepResult(id, index, "")
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

func (p *Result) LastExecutingStep() *StepResult {
	if len(p.ExecuteTrace) == 0 {
		return NewStepResult("", 0, "")
	} else {
		last_idx := len(p.ExecuteTrace) - 1
		return p.ExecuteTrace[last_idx]
	}
}

func (p *Result) StepResult() []*StepResult {
	rets := make([]*StepResult, 0)
	for _, ret_tmp := range p.ExecuteTrace {
		if ret_tmp.Father == p.PipelineId {
			rets = append(rets, ret_tmp)
		}
	}

	return rets
}
