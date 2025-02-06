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
	logger.Info("Check if slider image exists...")
	errChan := make(chan error, 1)
	length := 0
	go func(c context.Context) {
		if err := chromedp.Run(ctx,
			// 等待元素可见（即页面已加载）账号框
			chromedp.WaitVisible(s.Selector),
			chromedp.Evaluate("document.querySelectorAll('"+s.Selector+"').length", &length),
		); err != nil {
			if c.Err() == nil {
				errChan <- ErrSliderNotExists
			}
			return
		}
		if length > 0 {
			errChan <- nil
			return
		}
		errChan <- ErrSliderNotExists
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

var tempBgImg string

// 保存验证码图片.
func (s *Slide) saveVerifyImg(ctx context.Context) (*ImgBase64, error) {
	logger.Info("Save image to storage...")
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
		logger.Info("----temp bg img before:", tempBgImg)
		if bgimg == tempBgImg {
			logger.Info("------验证码图片需要刷新------")
			refresh(ctx)
			errChan <- ErrSliderRefresh
			return
		}
		tempBgImg = bgimg
		logger.Info("---temp bg img after:", tempBgImg)

		logger.Info("验证码 图片 加载成功:", bgimg)
		logger.Info("验证码 滑块 加载成功:", blockimg)

		// 保存验证码图片
		if bgimg != "" && blockimg != "" {
			imgBase := s.imgSave.Save(bgimg, blockimg)
			logger.Info("验证码背景和滑块保存成功")
			saveImg <- imgBase
			return
		}
		errChan <- ErrSliderSave
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
// yidun--error 验证失败
// yidun--success 验证成功
// yidun--loading

// 处理图片验证.
func (s *Slide) handleVerifyImg(ctx context.Context, imgbase64 *ImgBase64) error {
	errChan := make(chan error, 1)
	errlength := 0

	logger.Info("-----------开始处理图片验证-----------")
	// 滑动验证失败监测
	go func() {
		logger.Info("-----------开始验证失败监测------------", s.ErrorSelector)
		if err := chromedp.Run(ctx,
			chromedp.WaitVisible(s.ErrorSelector, chromedp.ByQuery),
			chromedp.Evaluate("document.querySelectorAll('"+s.ErrorSelector+"').length", &errlength),
		); err != nil {
			logger.Error("监测 `验证失败` 错误 :", err)
			errChan <- err
			return
		}
		// 有失败提示元素
		if errlength > 0 {
			logger.Info("监测 `验证失败` 元素数量, "+s.ErrorSelector+" 元素数量:", errlength)
			errChan <- ErrSliderVerify
			return
		}
		logger.Info("监测 `验证失败` 无失败提示, " + s.ErrorSelector)
		errChan <- nil
	}()

	time.Sleep(time.Second * 2)

	// 拖动滑块进行验证
	go func() {
		logger.Info("-----------开始滑动验证-----------")
		img := s.sliderimg(imgbase64)
		distance := getDistance(img, s.Mode)
		logger.Info("滑块距离:", distance)
		if distance == 0 {
			errChan <- ErrSliderDistance
			return
		}

		if err := chromedp.Run(ctx,
			DragSlider(s.DragSelector, distance), /* 拖动滑块*/
			chromedp.Sleep(time.Second),
		); err != nil {
			logger.Error("滑块拖动失败:", err)
			errChan <- err
			return
		}
		logger.Info("Drag Slider has done.")
		// time.Sleep(SleepTime)
		// errChan <- nil
	}()

	select {
	case e := <-errChan:
		return e
	case <-time.After(SleepTime):
		return goerr.Err("等待图片验证码超时")
	case <-ctx.Done():
		return goerr.Err("图片验证码 context done")
	}
}

func refresh(ctx context.Context) error {
	errchan := make(chan error, 1)
	go func() {
		if err := chromedp.Run(ctx,
			chromedp.WaitVisible(`button.yidun_refresh`),
			chromedp.Click(`button.yidun_refresh`),
			chromedp.Sleep(time.Second),
		); err != nil {
			errchan <- err
			return
		}
		errchan <- nil
	}()
	select {
	case e := <-errchan:
		return e
	case <-time.After(SleepTime):
		return goerr.Err("刷新页面超时")
	case <-ctx.Done():
		return goerr.Err("刷新页面 context done")
	}
}

type tryFailed struct{}
type TryFailer interface {
	TryFail(ctx context.Context) error
}

func (t tryFailed) TryFail(ctx context.Context) error {
	return nil
}
