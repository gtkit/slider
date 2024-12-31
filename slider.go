package slider

import (
	"context"
	"time"

	"gocv.io/x/gocv"
)

type TemplateMatchMode = gocv.TemplateMatchMode

const (
	SleepTime = 3 * time.Second
)

const (
	// TmSqdiff maps to TM_SQDIFF 平方差匹配.
	TmSqdiff TemplateMatchMode = 0
	// TmSqdiffNormed maps to TM_SQDIFF_NORMED 标准化平方差匹配.
	TmSqdiffNormed TemplateMatchMode = 1
	// TmCcorr maps to TM_CCORR 相关匹配.
	TmCcorr TemplateMatchMode = 2
	// TmCcorrNormed maps to TM_CCORR_NORMED 标准化相关匹配.
	TmCcorrNormed TemplateMatchMode = 3
	// TmCcoeff maps to TM_CCOEFF 相关系数匹配.
	TmCcoeff TemplateMatchMode = 4
	// TmCcoeffNormed maps to TM_CCOEFF_NORMED 标准化相关系数匹配.
	TmCcoeffNormed TemplateMatchMode = 5
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
	ImgSize

	Mode TemplateMatchMode // 模板匹配模式
}

func NewSlider(s *Slide) Slider {
	if s.TryNum == 0 {
		s.TryNum = 10
	}
	if s.Mode == 0 {
		s.Mode = TmSqdiff
	}
	return s
}

type ImgSize struct {
	BgWidth  int `json:"bg_width" ` // 显示验证码背景图片的宽度
	BgHeight int `json:"bg_height"` // 显示验证码背景图片的高度

	BlockWidth  int `json:"block_width"`  // 显示验证码滑块的宽度
	BlockHeight int `json:"block_height"` // 显示验证码滑块的高度
}

type ImgBase64 struct {
	BgBase64    string `json:"bg_base64"`    // 验证码背景图片base64
	BlockBase64 string `json:"block_base64"` // 验证码滑块图片base64
}
type Img struct {
	ImgBase64
	ImgSize
}

type ImgURL struct {
	BgURL    string `json:"bg_url"`
	BlockURL string `json:"block_url"`
}

type ImgInterface interface {
	Save() *ImgBase64
	Set(bgurl, blockurl string) ImgInterface
}
