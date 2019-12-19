package action

import (
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
