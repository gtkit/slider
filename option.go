package slider

import (
	"time"

	"github.com/gtkit/goerr"
)

type Opt func(*Slide) error

func WithTryNum(n int) Opt {
	return func(s *Slide) error {
		if n < 0 {
			return goerr.Err("slider.WithTryNum: negative number")
		}
		s.TryNum = n
		return nil
	}
}

func WithMode(m TemplateMatchMode) Opt {
	return func(s *Slide) error {
		s.Mode = m
		return nil
	}
}

func WithSleepTime(t time.Duration) Opt {
	return func(s *Slide) error {
		if t < 0 {
			return goerr.Err("slider.WithSleepTime: negative duration")
		}
		s.SleepTime = t
		return nil
	}
}

func WithImgSaver(sa ImgSaver) Opt {
	return func(s *Slide) error {
		if sa == nil {
			return goerr.Err("nil ImgSaver")
		}
		s.imgSave = sa
		return nil
	}
}

func WithTryFailed(try TryFailer) Opt {
	return func(s *Slide) error {
		if try == nil {
			return goerr.Err("nil TryFailer")
		}
		s.tryFailed = try
		return nil
	}
}
