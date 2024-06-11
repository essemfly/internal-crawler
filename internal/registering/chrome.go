package registering

import (
	"context"

	"github.com/chromedp/chromedp"
)

func OpenChrome() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocatorCtx)

	// 여기서는 cancel 함수를 직접 호출하지 않고, 반환하여 호출자가 제어하도록 합니다.
	return ctx, func() {
		cancel()
		allocatorCancel() // allocator context도 취소합니다.
	}
}
