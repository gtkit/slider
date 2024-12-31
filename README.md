## 图片滑块验证码破解
### 介绍
此功能借鉴于 [https://github.com/omigo/crack-slide-captcha](https://github.com/omigo/crack-slide-captcha) 项目，使用滑块验证码破解。
使用go语言编写，依赖 chromedp, gocv库。 通过 chromdp 操作 Chromium 浏览器,使用  OpenCV 匹配找出滑块位置，计算出滑动距离，然后模拟鼠标事件

示例网站: [https://ads.huawei.com/usermgtportal/home/index.html#/]("https://ads.huawei.com/usermgtportal/home/index.html#/")
### 使用方法
前期操作由 chromedp 控制浏览器，需要先安装 chromedp 库，并下载对应浏览器驱动。
```go
// 实例化滑块验证码破解
if err := slider.NewSlider(&slider.Slide{
		SliderImgSelector: "img.yidun_bg-img", // 用来滑块图片存在的节点选择器
		BgImgQuery:        "img.yidun_bg-img", // 滑块背景图片选择器
		BlockImgQuery:     "img.yidun_jigsaw", // 拼图图片选择器
		DragSelector:      "div.yidun_slider.yidun_slider--hover", // 拖动的滑块选择器
		ImgSize: slider.ImgSize{
			BgWidth:     260, // 滑块背景图片显示宽度
			BgHeight:    130, // 滑块背景图片显示高度
			BlockWidth:  49, // 拼图图片显示宽度
			BlockHeight: 130, // 拼图图片显示高度
		},
		Mode: slider.TmSqdiff, // 匹配模式，可选值：TmSqdiff TmSqdiffNormed TmCcorr TmCcorrNormed TmCcoeff TmCcoeffNormed
	}).Run(ctx); err != nil {
		logger.Error("slider run error: ", err)
	}
```
### 操作效果
![操作效果](./slider.gif)

### 参考
crack-slide-captcha: https://github.com/omigo/crack-slide-captcha

GoCV： https://gocv.io/computer-vision/

chromedp： https://github.com/chromedp/chromedp
