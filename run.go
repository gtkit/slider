package slider

import (
	"context"
	"time"

	"github.com/gtkit/goerr"
	"github.com/gtkit/logger"
)

func (s *Slide) Run(ctx context.Context) error {
	var (
		err     error
		imgbase *ImgBase64
	)
	errchan := make(chan error, 1)

	// 滑动验证
	go func() {
		for i := 0; i < s.TryNum; i++ {
			time.Sleep(SleepTime)
			// 判断是否有图片验证码
			if err = s.hasVerifyImg(ctx); err != nil {
				logger.Error("not has verify img, error:", err.Error())
				errchan <- err
				return
			}
			// 保存图片验证码
			imgbase, err = s.saveVerifyImg(ctx)
			if err != nil {
				if i < s.TryNum-1 {
					continue
				}
				errchan <- err
				return
			}
			// 处理图片验证
			if err = s.handleVerifyImg(ctx, imgbase); err != nil {
				if i < s.TryNum-1 {
					// 刷新验证图片
					// if err := refresh(ctx); err != nil {
					// 	logger.Error("refresh failed, error:", err.Error(), ", try num:", i+1)
					// }
					logger.Error("handleVerifyImg failed, error:", err.Error(), "; try num:", i+1)
					continue
				}
				errchan <- goerr.WithMsg(ErrSliderVerify, err.Error())
				return
			}
			errchan <- nil

			// if err = s.hasVerifyImg(ctx); err == nil {
			// 	// 保存图片验证码
			// 	imgbase, err = s.saveVerifyImg(ctx)
			// 	if err != nil {
			// 		logger.Error("save verify img failed, error:", err.Error(), ", try num:", i+1)
			// 		errchan <- goerr.WithMsg(ErrSliderSave, err.Error())
			// 		return
			// 	}
			// 	// 处理图片验证
			// 	if err = s.handleVerifyImg(ctx, imgbase); err != nil {
			// 		logger.Error("verify failed, error:", err.Error(), ", try num:", i+1)
			// 		if i < s.TryNum-1 {
			// 			// 刷新验证图片
			// 			// if err := refresh(ctx); err != nil {
			// 			// 	logger.Error("refresh failed, error:", err.Error(), ", try num:", i+1)
			// 			// }
			//
			// 			// logger.Error("verify failed, error:", err.Error(), ", try num:", i+1)
			// 			continue
			// 		}
			// 		errchan <- goerr.WithMsg(ErrSliderVerify, err.Error())
			// 		return
			// 	}
			// 	errchan <- nil
			// 	return
			// }
			// errchan <- err
		}
	}()

	// 尝试失败次数过多, 需点击 "失败过多，点此重试"
	go func() {
		if err := s.tryFailed.TryFail(ctx); err != nil {
			logger.Error("尝试失败次数过多, 错误 :", err)
			return
		}
	}()

	select {
	case err = <-errchan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Second):
		return goerr.Wrap(ctx.Err(), "timeout")
	}
}

func (s *Slide) sliderimg(imgbase *ImgBase64) *Img {
	return &Img{
		ImgBase64: ImgBase64{
			BgBase64:    imgbase.BgBase64,
			BlockBase64: imgbase.BlockBase64,
		},
		ImgSize: ImgSize{
			BgWidth:     s.BgWidth,
			BgHeight:    s.BgHeight,
			BlockWidth:  s.BlockWidth,
			BlockHeight: s.BlockHeight,
		},
	}
}
