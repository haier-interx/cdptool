package action

import (
	"context"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

const (
	DEVICE_CHROME_OSX = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"
)

func UserAgent(ua_language ...string) chromedp.Action {
	var ua, l string
	if len(ua_language) < 1 || ua_language[0] == "" {
		ua = DEVICE_CHROME_OSX
	} else {
		ua = ua_language[1]
	}
	if len(ua_language) < 2 || ua_language[1] == "" {
		l = "zh-CN,zh"
	} else {
		l = ua_language[1]
	}

	return chromedp.ActionFunc(func(ctx context.Context) error {
		return emulation.SetUserAgentOverride(ua).
			WithAcceptLanguage(l). //中文
			Do(ctx)
	})
}
