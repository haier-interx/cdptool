package action

import (
	"context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"math"
)

func FullScreenshot(quality int64, filename string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Printf("save fullScreenshot to %s ...", filename)

		// get layout metrics
		_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
		if err != nil {
			return err
		}

		width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))
		// force viewport emulation
		err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
			Do(ctx)
		if err != nil {
			return err
		}

		var buf []byte

		// capture screenshot
		buf, err = page.CaptureScreenshot().
			WithQuality(quality).
			WithClip(&page.Viewport{
				X: contentSize.X,
				Y: contentSize.Y,
				//Width:  float64(width),
				//Height: float64(height),
				Width:  contentSize.Width,
				Height: contentSize.Height,
				Scale:  1,
			}).Do(ctx)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
			log.Fatal(err)
		}

		return nil
	})
}
