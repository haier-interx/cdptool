package action

import (
	"context"
	"github.com/chromedp/cdproto/performance"
	"github.com/chromedp/chromedp"
	"github.com/haier-interx/cdptool/models"
)

const (
	perfromance_js = `JSON.parse(JSON.stringify(performance.getEntriesByType('navigation')[0]))`
)

func Performance(ret *models.PerformanceTiming) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Evaluate(perfromance_js, ret),
	}
}

func Performance2(metrics []*performance.Metric) chromedp.Tasks {
	return []chromedp.Action{
		new(performance.EnableParams),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			metrics, err = performance.GetMetrics().Do(ctx)
			return err
		}),
	}
}
