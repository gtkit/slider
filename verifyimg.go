package slider

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gtkit/goerr"
	"github.com/gtkit/logger"
)

// 判断是否有验证图片.
func (s *Slide) hasVerifyImg(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func(c context.Context) {
		if err := chromedp.Run(ctx,
			// 等待元素可见（即页面已加载）账号框
			chromedp.WaitVisible(s.SliderImgSelector),
		); err != nil {
			if c.Err() == nil {
				errChan <- err
			}
			return
		}
		errChan <- nil
	}(ctx)
	select {
	case e := <-errChan:
		return e
	case <-time.After(SleepTime):
		return goerr.Err("没有验证图片")
	case <-ctx.Done():
		return goerr.Err("hasVerifyImg context done")
	}
}

// 保存验证码图片.
func (s *Slide) saveVerifyImg(ctx context.Context, img ImgInterface) (*ImgBase64, error) {
	errChan := make(chan error, 1)
	saveImg := make(chan *ImgBase64, 1)

	go func(c context.Context) {
		var bgimg, blockimg string
		// 等待图片验证码加载
		if err := chromedp.Run(ctx,
			chromedp.Evaluate("document.querySelector('"+s.BgImgQuery+"').src", &bgimg),       // 验证码背景
			chromedp.Evaluate("document.querySelector('"+s.BlockImgQuery+"').src", &blockimg), // 验证码滑块
		); err != nil {
			if c.Err() == nil {
				errChan <- err
			}
			return
		}

		// 保存验证码图片
		if bgimg != "" && blockimg != "" {
			imgBase := img.Set(bgimg, blockimg).Save()
			// logger.Info("验证码图片保存成功:", imgBase)
			saveImg <- imgBase
		}
	}(ctx)

	select {
	case e := <-errChan:
		return nil, e
	case <-time.After(SleepTime):
		return nil, goerr.Err("等待图片验证码超时")
	case <-ctx.Done():
		return nil, goerr.Err("等待图片验证码超时 context done")
	case img := <-saveImg:
		return img, nil
	}
}

// 处理图片验证.
func (s *Slide) handleVerifyImg(ctx context.Context, imgbase64 *ImgBase64) error {
	errChan := make(chan error, 1)
	go func() {
		img := s.sliderimg(imgbase64)
		if err := chromedp.Run(ctx,
			// 拖动滑块进行验证
			DragSlider(s.DragSelector, getDistance(img, s.Mode)),
		); err != nil {
			logger.Error("图片验证失败:", err)
			errChan <- err
			return
		}
		errChan <- nil
	}()
	select {
	case e := <-errChan:
		return e
	case <-time.After(SleepTime):
		return goerr.Err("等待图片验证码超时")
	case <-ctx.Done():
		return goerr.Err("等待图片验证码超时 context done")
	}
}
