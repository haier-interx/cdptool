package action

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

func Debug(content interface{}) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf("%v", content)
			return nil
		}),
	}
}
