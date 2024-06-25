package registering

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/chromedp/cdproto/browser"
	cdpruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func NaverLogin(ctx context.Context) {
	err := chromedp.Run(ctx, chromedp.Navigate(`https://nid.naver.com/nidlogin.login`), chromedp.WaitVisible(`body`, chromedp.ByQuery))
	if err != nil {
		log.Fatal(err)
	}

	if err := clipboard.WriteAll(os.Getenv("NAVER_USERNAME")); err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	pasteText(ctx, `#id`)

	if err := clipboard.WriteAll(os.Getenv("NAVER_PASSWORD")); err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	pasteText(ctx, `#pw`)

	err = chromedp.Run(ctx, chromedp.Click(`#log\.login`, chromedp.ByQuery), chromedp.Sleep(2*time.Second))
	if err != nil {
		log.Fatal(err)
	}
}

func pasteText(ctx context.Context, selector string) {
	err := chromedp.Run(ctx,
		browser.GrantPermissions([]browser.PermissionType{"clipboardReadWrite"}),
		chromedp.Focus(selector),
		chromedp.Click(selector, chromedp.ByQuery),
		chromedp.Evaluate(generateExpression(selector), nil, func(ep *cdpruntime.EvaluateParams) *cdpruntime.EvaluateParams {
			return ep.WithAwaitPromise(true)
		}),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func generateExpression(selector string) string {
	expression := fmt.Sprintf(`
	async function setContent(){
		const content = await navigator.clipboard.readText();
		document.querySelector('%s').value = content;
	}
	setContent()
	`, selector)
	return expression
}
