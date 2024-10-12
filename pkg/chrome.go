package pkg

import (
	"context"

	"github.com/chromedp/chromedp"
)

func OpenChrome() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),            // GUI 모드로 실행
		chromedp.Flag("start-maximized", false),    // 브라우저 창 최대화
		chromedp.Flag("enable-automation", false),  // 자동화 배너 비활성화
		chromedp.Flag("disable-extensions", false), // 확장 프로그램 활성화 (필요한 경우)
	)
	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocatorCtx)

	// 여기서는 cancel 함수를 직접 호출하지 않고, 반환하여 호출자가 제어하도록 합니다.
	return ctx, func() {
		cancel()
		allocatorCancel() // allocator context도 취소합니다.
	}
}
