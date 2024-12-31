package slider

import (
	"encoding/base64"
	"image"
	"strings"

	"github.com/gtkit/goerr"
	"github.com/gtkit/logger"
	"gocv.io/x/gocv"
	"golang.org/x/sync/errgroup"
)

func getDistance(img *Img, mode gocv.TemplateMatchMode) int {
	var (
		eg               errgroup.Group
		err              error
		block, bg, alpha gocv.Mat
	)

	// 验证码滑块预处理
	eg.Go(func() error {
		alpha, block, err = preProcess(img.BlockBase64, img.BlockWidth, img.BlockHeight)
		if err != nil {
			return goerr.Wrap(err, "Block preProcess err")
		}
		return nil
	})

	// 背景图预处理
	eg.Go(func() error {
		_, bg, err = preProcess(img.BgBase64, img.BgWidth, img.BgHeight)
		if err != nil {
			return goerr.Wrap(err, "Bg preProcess err")
		}
		return nil
	})

	// 等待所有协程完成
	if err = eg.Wait(); err != nil {
		logger.Error("preProcess err:", err)
		return 0
	}
	defer block.Close()
	defer bg.Close()

	return match(bg, block, alpha, mode).X
}

func decode(base64img string) []byte {
	i := strings.IndexByte(base64img, ',')
	if i == -1 {
		logger.Error(base64img)
		return nil
	}
	b, err := base64.StdEncoding.DecodeString(base64img[i+1:])
	if err != nil {
		logger.Error(err, base64img[i+1:])
		return nil
	}
	return b
}

func readBase64Image(b64Image string) (gocv.Mat, error) {
	origin, err := gocv.IMDecode(decode(b64Image), gocv.IMReadUnchanged)
	if err != nil {
		return gocv.Mat{}, err
	}
	return origin, nil
}

func resize(origin gocv.Mat, cols, rows int) gocv.Mat {
	resized := gocv.NewMatWithSize(cols, rows, origin.Type())
	gocv.Resize(origin, &resized, image.Pt(cols, rows), 0, 0, gocv.InterpolationNearestNeighbor)
	return resized
}

func gray(origin gocv.Mat) gocv.Mat {
	grayed := gocv.NewMat()
	gocv.CvtColor(origin, &grayed, gocv.ColorBGRToGray)
	return grayed
}

func match(bg, block, mask gocv.Mat, mode gocv.TemplateMatchMode) image.Point {
	result := gocv.NewMatWithSize(
		bg.Rows()-block.Rows()+1,
		bg.Cols()-block.Cols()+1,
		gocv.MatTypeCV32FC1)
	defer result.Close()

	// TmSqdiff 平方差匹配 0
	// TmSqdiffNormed 标准化平方差匹配 1
	// TmCcorr 相关匹配 2
	// TmCcorrNormed 标准化相关匹配 3
	// TmCcoeff 相关系数匹配 4
	// TmCcoeffNormed 标准化相关系数匹配 5

	gocv.MatchTemplate(bg, block, &result, mode, mask)
	gocv.Normalize(result, &result, 0, 1, gocv.NormMinMax)

	_, _, _, maxLoc := gocv.MinMaxLoc(result)

	return maxLoc
}

func preProcess(b64Image string, width, heigh int) (alpha, processed gocv.Mat, err error) {
	origin, err := readBase64Image(b64Image)
	if err != nil {
		return gocv.Mat{}, gocv.Mat{}, err
	}

	resized := resize(origin, width, heigh)
	grayed := gray(resized)

	logger.Debug(origin.Cols(), origin.Rows(), resized.Cols(), resized.Rows())

	if resized.Channels() == 4 {
		return gocv.Split(resized)[3], grayed, nil
	}

	return gocv.Mat{}, grayed, nil
}
