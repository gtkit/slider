package slider

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/gtkit/goerr"
	"github.com/gtkit/logger"
	"golang.org/x/sync/errgroup"
)

const (
	Base64Prefix = "data:image/jpeg;base64,"
)

// 下载验证码图片
func (i *ImgUrl) Save() *SliderImgBase64 {
	sib := &SliderImgBase64{}
	// 分别获取图片
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	// 下载验证码背景图片
	eg.Go(func() error {
		data, err := downloadImg(ctx, i.BgUrl)
		if err != nil {
			logger.Error("downloadImg BgUrl error: %w", err)
			return err
		}

		str := base64.StdEncoding.EncodeToString(data)
		// fmt.Println("验证码背景图片base64: ", str)
		sib.BgBase64 = Base64Prefix + str
		return nil
	})

	// 下载验证码滑块图片
	eg.Go(func() error {
		data, err := downloadImg(ctx, i.BlockUrl)
		if err != nil {
			logger.Error("downloadImg BlockUrl error: %w", err)
			return err
		}

		str := base64.StdEncoding.EncodeToString(data)
		// fmt.Println("验证码滑块图片base64: ", str)
		sib.BlockBase64 = Base64Prefix + str
		return nil

	})

	if err := eg.Wait(); err != nil {
		logger.Error("download img error: %w", err)
		return nil
	}

	return sib
}

// 设置url和path
func (i *ImgUrl) Set(bgurl, blockurl string) ImgInterface {
	i.BgUrl = bgurl
	i.BlockUrl = blockurl
	return i
}

// 下载图片
func downloadImg(ctx context.Context, url string) ([]byte, error) {
	imgchan := make(chan []byte, 1)
	errchan := make(chan error, 1)

	go func() {
		// 判断文件是否存在
		resp, err := http.Head(url)
		if err != nil {
			errchan <- err
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			errchan <- goerr.Err("图片访问失败: %s" + resp.Status)
			return
		}

		// 获取图片
		res, err := http.Get(url)
		if err != nil {
			errchan <- err
			return
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			errchan <- err
			return
		}
		imgchan <- data
	}()
	select {
	case data := <-imgchan:
		return data, nil
	case err := <-errchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
