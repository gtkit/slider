package slider

import (
	"context"
	"time"

	"github.com/gtkit/goerr"
	"github.com/gtkit/logger"
)

func (s *Slide) Run(ctx context.Context) error {
	if err := s.verifyParams(); err != nil {
		return err
	}
	errchan := make(chan error, 1)
	go func() {
		for i := 0; i < 10; i++ {
			if err := s.hasVerifyImg(ctx); err == nil {
				logger.Info("chrome has verifyImg")
				imgbase, err := s.saveVerifyImg(ctx, &ImgUrl{})
				if err != nil {
					errchan <- goerr.Wrap(err, "saveVerifyImg error")
					return
				}

				// 处理图片验证
				if err := s.handleVerifyImg(ctx, imgbase); err != nil {
					errchan <- goerr.Wrap(err, "handleVerifyImg error")
					return
				}
			} else {
				logger.Info("没有图片验证码: ", err)
				errchan <- nil
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case err := <-errchan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		return goerr.Wrap(ctx.Err(), "timeout")
	}
}
func (s *Slide) verifyParams() error {
	if s.DragSelector == "" {
		return goerr.Err("dragSelector is empty")
	}
	if s.SliderImgSelector == "" {
		return goerr.Err("SliderImgSelector is empty")
	}
	if s.BgImgQuery == "" {
		return goerr.Err("BgImgQuery is empty")
	}
	if s.BlockImgQuery == "" {
		return goerr.Err("BlockImgQuery is empty")
	}
	if s.BgWidth <= 0 {
		return goerr.Err("BgWidth is 0")
	}
	if s.BgHeight <= 0 {
		return goerr.Err("BgHeight is 0")
	}
	if s.BlockWidth <= 0 {
		return goerr.Err("BlockWidth is 0")
	}
	if s.BlockHeight <= 0 {
		return goerr.Err("BlockHeight is 0")
	}
	return nil
}

func (s *Slide) sliderimg(imgbase *SliderImgBase64) *SliderImg {
	return &SliderImg{
		SliderImgBase64: SliderImgBase64{
			BgBase64:    imgbase.BgBase64,
			BlockBase64: imgbase.BlockBase64,
		},
		SliderImgSize: SliderImgSize{
			BgWidth:     s.BgWidth,
			BgHeight:    s.BgHeight,
			BlockWidth:  s.BlockWidth,
			BlockHeight: s.BlockHeight,
		},
	}

}
