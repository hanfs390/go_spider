package headless

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"time"
)

func GetHTMLByChromebp(url string, flag bool) string{
	//增加选项，允许chrome窗口显示出来
	options := []chromedp.ExecAllocatorOption{
		//chromedp.ExecPath("C:\\Users\\hanfushun\\AppData\\Local\\Google\\Chrome\\Application\\chrome.exe"),
		chromedp.ExecPath("/usr/bin/google-chrome"),
		chromedp.Flag("headless", flag),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	//创建chrome窗口
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	var res string
	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.Sleep(time.Second * 8 ),
			chromedp.OuterHTML("html", &res),
			chromedp.Stop(),
		},
	); err != nil {
		fmt.Println("chromebp: ", err)
	}

	chromedp.Stop()
	return res
}
