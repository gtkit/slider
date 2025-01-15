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
			chromedp.WaitVisible(s.Selector),
		); err != nil {
			if c.Err() == nil {
				errChan <- ErrSliderNotExists
			}
			return
		}
		errChan <- nil
	}(ctx)
	select {
	case e := <-errChan:
		return e
	case <-time.After(SleepTime):
		return ErrSliderNotExists
	case <-ctx.Done():
		return ErrSliderCtxDone
	}
}

// 保存验证码图片.
func (s *Slide) saveVerifyImg(ctx context.Context) (*ImgBase64, error) {
	errChan := make(chan error, 1)
	saveImg := make(chan *ImgBase64, 1)

	go func(c context.Context) {
		var bgimg, blockimg string
		// 等待图片验证码加载
		if err := chromedp.Run(ctx,
			chromedp.Evaluate("document.querySelector('"+s.BgImgSelector+"').src", &bgimg),       // 验证码背景
			chromedp.Evaluate("document.querySelector('"+s.BlockImgSelector+"').src", &blockimg), // 验证码滑块
		); err != nil {
			if c.Err() == nil {
				errChan <- err
			}
			return
		}

		// 保存验证码图片
		if bgimg != "" && blockimg != "" {
			imgBase := s.imgSave.Save(bgimg, blockimg)
			// logger.Info("验证码图片保存成功:", imgBase)
			saveImg <- imgBase
		}
	}(ctx)

	select {
	case e := <-errChan:
		return nil, e
	case <-time.After(SleepTime):
		return nil, goerr.Err("slider timeout")
	case <-ctx.Done():
		return nil, goerr.Err("save slider context done")
	case img := <-saveImg:
		return img, nil
	}
}

// yidun_tips__icon 叉号
// yidun_tips__text yidun-fallback__tip 失败过多，点此重试
// yidun--error
// yidun--loading

// 处理图片验证.
func (s *Slide) handleVerifyImg(ctx context.Context, imgbase64 *ImgBase64) error {
	errChan := make(chan error, 1)

	// 验证失败
	go func() {
		if err := chromedp.Run(ctx,
			chromedp.WaitVisible(s.ErrorSelector),
		); err != nil {
			errChan <- err
			return
		}
		logger.Info("图片验证失败, " + s.ErrorSelector)
		errChan <- ErrSliderVerify
	}()

	// 处理失败次数过多, 由调用者实现自己的处理逻辑
	go func() {
		if err := s.tryFailed.TryFail(ctx); err != nil {
			errChan <- err
		}
	}()

	// 拖动滑块进行验证
	go func() {
		img := s.sliderimg(imgbase64)
		if err := chromedp.Run(ctx,
			DragSlider(s.DragSelector, getDistance(img, s.Mode)),
		); err != nil {
			logger.Error("图片验证失败:", err)
			errChan <- err
			return
		}
		logger.Info("图片验证重新加载")
		errChan <- nil
	}()

	select {
	case e := <-errChan:
		return e
	case <-time.After(SleepTime * 100):
		return goerr.Err("等待图片验证码超时")
	case <-ctx.Done():
		return goerr.Err("等待图片验证码超时 context done")
	}
}

type tryFailed struct{}
type TryFailer interface {
	TryFail(ctx context.Context) error
}

func (t tryFailed) TryFail(ctx context.Context) error {
	return nil
}
