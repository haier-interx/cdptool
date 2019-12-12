package action

import (
	"context"
	"github.com/chromedp/chromedp"
)

type ActionCtx struct {
	myCtx  context.Context
	Action chromedp.Action
}

func Wrap(ctx context.Context, ac chromedp.Action) chromedp.Action {
	return &ActionCtx{ctx, ac}
}

func (ac *ActionCtx) Do(ctx context.Context) error {
	done := make(chan struct{})
	var execErr error
	go func() {
		defer close(done)
		execErr = ac.Action.Do(ctx)
	}()
	select {
	case <-done:
		return execErr
	case <-ac.myCtx.Done():
		return ac.myCtx.Err()
	}
}
