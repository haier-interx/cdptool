package pipeline

import (
	"cdptool/models"
	"fmt"
)

type Result struct {
	executingIdx int
	error        error

	JavaScriptResult []*[]byte
	Performances     []*models.PerformanceTiming
}

func NewResult() *Result {
	return &Result{
		executingIdx:     -1,
		JavaScriptResult: make([]*[]byte, 0),
		Performances:     make([]*models.PerformanceTiming, 0),
	}
}

func (p *Result) Error() error {
	return p.error
}

func (p *Result) ErrorCN() string {
	if p.error == nil {
		return ""
	}

	return fmt.Sprintf("执行步骤%d时失败：%s", p.executingIdx+1, ErrorCN(p.error))
}

func (p *Result) ErrorStepIdx() int {
	if p.error == nil {
		return -1
	}
	return p.executingIdx
}

func (p *Result) ExecutingIdx() int {
	return p.executingIdx
}
