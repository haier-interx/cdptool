package pipeline

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/haier-interx/cdptool/action"
	"log"
	"time"
)

type Pipeline struct {
	Id                    string        `json:"id" yaml:"id"`
	Timeout               time.Duration `json:"timeout" yaml:"timeout"`
	Steps                 []*Step       `json:"steps" yaml:"steps"`
	ScreenshotsWhenFailed bool          `json:"screenshots_when_failed" yaml:"screenshots_when_failed"`
}

func (p *Pipeline) Run(ctx context.Context, cds *CustomDefinitions) (ret *Result) {
	ctx_timeout, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	var actions chromedp.Tasks
	actions, ret = p.Parse(ctx_timeout, cds)
	if ret.Error() != nil {
		return
	}

	// run
	ctx_chromedp, cancel2 := chromedp.NewContext(context.Background())
	defer cancel2()
	defer chromedp.Cancel(ctx_chromedp)
	err := chromedp.Run(ctx_chromedp, actions...)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			ret.error = ERR_ELEMENT_NOTFOUND_OR_TIMEOUT
		default:
			ret.error = err
		}

		// 失败时做截屏操作
		if p.ScreenshotsWhenFailed {
			ctx_tmp, cancel_tmp := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel_tmp()
			file_name := fmt.Sprintf("%s-error.%d.png", ret.LastExecutingStep().Id(), time.Now().UnixNano())
			ss_action := action.Wrap(ctx_tmp, action.FullScreenshot(90, file_name))
			err_tmp := chromedp.Run(ctx_chromedp, ss_action)
			if err_tmp != nil {
				log.Printf("screenshot action failed while execute error: %v", err_tmp)
			}
		}

		return
	}

	return
}

func (p *Pipeline) Parse(ctx context.Context, cds *CustomDefinitions) (actions chromedp.Tasks, ret *Result) {
	ret = NewResult()
	if p.Id == "" {
		ret.Failed(ERR_PIPELINE_ID_REQUIRED)
		return
	}

	if p.Timeout == 0 {
		p.Timeout = 15 * time.Second
	}

	// actions
	actions = make([]chromedp.Action, 0)
	for i, step := range p.Steps {
		step.SetId(p.GenerateStepId(step, i))
		actions_tmp := step.ActionWithCtx(ctx, ret, cds)
		if ret.Error() != nil {
			ret.ParseFailed(p.Id, i, ret.Error())
			return
		}

		// save execute record
		e_tmp := &ExecutingStep{p.Id, i}
		actions = append(actions,
			chromedp.ActionFunc(func(ctx context.Context) error {
				ret.PutExecuting(e_tmp)
				return nil
			}),
		)

		// real action
		actions = append(actions, actions_tmp...)
	}

	return
}

func (p *Pipeline) GenerateStepId(s *Step, idx int) string {
	return fmt.Sprintf("%s_%d_%s", p.Id, idx, s.Type)
}
