package action

import (
	"github.com/chromedp/chromedp"
	"github.com/haier-interx/cdptool/models"
)

const (
	network_js = `JSON.parse(JSON.stringify(performance.getEntriesByType('resource')))`
	//network_js = `performance.getEntriesByType('resource')`
)

func Network(ret *[]*models.PerformanceTiming) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Evaluate(network_js, ret),
	}
}
