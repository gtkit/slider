package slider

import (
	"context"
	"time"
)

type Slider interface {
	Run(ctx context.Context) error
}

type Slide struct {
	SliderImgSelector string // 验证码滑块图片选择器 示例: "img.yidun_bg-img"
	BgImgQuery        string // 验证码背景图片选择器查询 示例: "document.querySelector('img.yidun_bg-img').src"
	BlockImgQuery     string // 验证码滑块选择器查询 示例: "document.querySelector('img.yidun_jigsaw').src"
	DragSelector      string // 拖动选择器 示例: "div.yidun_slider.yidun_slider--hover"
	TryNum            int    // 尝试次数
	SliderImgSize
}

func NewSlider(s *Slide) Slider {
	if s.TryNum == 0 {
		s.TryNum = 10
	}
	return s
}

const (
	SleepTime = 3 * time.Second
)

type SliderImgSize struct {
	BgWidth  int `json:"bg_width" ` // 显示验证码背景图片的宽度
	BgHeight int `json:"bg_height"` // 显示验证码背景图片的高度

	BlockWidth  int `json:"block_width"`  // 显示验证码滑块的宽度
	BlockHeight int `json:"block_height"` // 显示验证码滑块的高度
}

type SliderImgBase64 struct {
	BgBase64    string `json:"bg_base64"`    // 验证码背景图片base64
	BlockBase64 string `json:"block_base64"` // 验证码滑块图片base64
}
type SliderImg struct {
	SliderImgBase64
	SliderImgSize
}

type ImgUrl struct {
	BgUrl    string `json:"bg_url"`
	BlockUrl string `json:"block_url"`
}

type ImgInterface interface {
	Save() *SliderImgBase64
	Set(bgurl, blockurl string) ImgInterface
}
