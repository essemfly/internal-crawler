package registering

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

func KakaoLogin(ctx context.Context, cafeName, username, password string) error {
	var loginURL = fmt.Sprintf("https://m.cafe.daum.net/%s", cafeName)

	err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
		chromedp.Click(`#daumMinidaum > a`),
		chromedp.WaitVisible(`#loginId--1`),
		chromedp.SendKeys(`#loginId--1`, username),
		chromedp.SendKeys(`#password--2`, password),
		chromedp.Click(`#mainContent > div > div > form > div.confirm_btn > button.btn_g.highlight.submit`),
		chromedp.Sleep(15*time.Second),
	)
	if err != nil {
		return err
	}

	return nil
}
