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
	Steps   []*Step       `json:"steps"`
}

func (p *Pipeline) Run(ctx context.Context, cds *CustomDefinitions) (ret *Result) {
	ret = NewResult()
	if p.Id == "" {
		ret.error = ERR_PIPELINE_ID_REQUIRED
		return
	}

	if p.Timeout == 0 {
		p.Timeout = 15 * time.Second
	}
	ctx_timeout, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	// actions
	actions := make([]chromedp.Action, 0)
	for i, step := range p.Steps {
		step.SetId(p.GenerateStepId(step, i))

		actions_tmp, err := step.ActionWithCtx(ctx_timeout, ret, cds)
		if err != nil {
			ret.executingIdx = i
			ret.error = fmt.Errorf("%s: %w", step.Id(), err)
			return
		}

		actions = append(actions, SetExecutingIdxAction(i, &ret.executingIdx))
		actions = append(actions, actions_tmp...)
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

func (p *Pipeline) GenerateStepId(s *Step, idx int) string {
	return fmt.Sprintf("%s_%d_%s", p.Id, idx, s.Type)
}

func SetExecutingIdxAction(i int, executingIdx *int) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		*executingIdx = i
		return nil
	})
}
