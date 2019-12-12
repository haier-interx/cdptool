package pipeline

import (
	"cdptool/action"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

type Pipeline struct {
	Id      string        `json:"id"`
	Timeout time.Duration `json:"timeout"`
	Steps   []Step        `json:"steps"`
}

func (p *Pipeline) Run(ctx context.Context) (ret *Result) {
	ret = &Result{
		executingIdx: -1,
		performances: make([]*action.PerformanceResult, 0),
	}

	if p.Id == "" {
		ret.error = ERR_PIPELINE_ID_REQUIRED
		return
	}

	if p.Timeout == 0 {
		p.Timeout = 15 * time.Second
	}
	ctx_timeout, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	ctx_chromedp, cancel2 := chromedp.NewContext(context.Background())
	defer chromedp.Cancel(ctx_chromedp)
	defer cancel2()

	// actions
	actions := make([]chromedp.Action, 0)
	for i, step := range p.Steps {
		var data interface{}
		if step.Type == STEP_PERFORMANCE {
			pr := new(action.PerformanceResult)
			ret.performances = append(ret.performances, pr)
			data = pr
		}

		actions_tmp, err := step.ActionWithCtx(ctx_timeout, fmt.Sprintf("%s-%d", p.Id, i), data)
		if err != nil {
			ret.error = fmt.Errorf("Step%d: %w", i, err)
			return
		}
		actions = append(actions, SetExecutingIdxAction(i, &ret.executingIdx))
		actions = append(actions, actions_tmp...)
	}

	// run
	err := chromedp.Run(ctx_chromedp, actions...)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			ret.error = ERR_ELEMENT_NOTFOUND_OR_TIMEOUT
		default:
			ret.error = err
		}

		// 失败时做截屏操作
		ctx_tmp, cancel_tmp := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel_tmp()

		file_name := fmt.Sprintf("%s-%d-error.%d.png", p.Id, ret.ErrorStepIdx(), time.Now().UnixNano())
		ss_action := action.Wrap(ctx_tmp, action.FullScreenshot(90, file_name))
		err_tmp := chromedp.Run(ctx_chromedp, ss_action)
		if err_tmp != nil {
			log.Printf("screenshot action failed while execute error: %v", err_tmp)
		}

		return
	}

	return
}

func SetExecutingIdxAction(i int, executingIdx *int) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		*executingIdx = i
		return nil
	})
}

type Result struct {
	executingIdx int
	error        error
	performances []*action.PerformanceResult
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

func (p *Result) PerformanceResults() []*action.PerformanceResult {
	return p.performances
}
