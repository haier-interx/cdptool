package action

import (
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func DeviceScreen(width, height int64, mobile bool) chromedp.Action {
	return emulation.SetDeviceMetricsOverride(width, height, 1, mobile).
		WithScreenOrientation(&emulation.ScreenOrientation{
			Type:  emulation.OrientationTypePortraitPrimary,
			Angle: 0,
		})
}
