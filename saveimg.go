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

type ImgURL struct {
	BgURL    string `json:"bg_url"`
	BlockURL string `json:"block_url"`
}

// 下载验证码图片.
func (i *ImgURL) Save(bgurl, blockurl string) *ImgBase64 {
	sib := &ImgBase64{}
	// 分别获取图片
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	// 下载验证码背景图片
	eg.Go(func() error {
		data, err := downloadImg(ctx, bgurl)
		if err != nil {
			logger.Error("downloadImg BgURL error: %w", err)
			return err
		}

		str := base64.StdEncoding.EncodeToString(data)
		// fmt.Println("验证码背景图片base64: ", str)
		sib.BgBase64 = Base64Prefix + str
		return nil
	})

	// 下载验证码滑块图片
	eg.Go(func() error {
		data, err := downloadImg(ctx, blockurl)
		if err != nil {
			logger.Error("downloadImg BlockURL error: %w", err)
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

// 下载图片.
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
		if resp.StatusCode != http.StatusOK {
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
