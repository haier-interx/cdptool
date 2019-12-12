package action

import (
	"context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func InitDeviceMetrics() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		return emulation.SetDeviceMetricsOverride(1600, 900, 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
			Do(ctx)
	})
}
