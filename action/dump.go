package action

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

func Dump() chromedp.Tasks {
	var html string
	return chromedp.Tasks{
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf("%s", html)
			return nil
		}),
	}
}
