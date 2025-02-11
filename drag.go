package slider

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/gtkit/logger"
)

// DragSlider 拖动滑块.
// sel: 选择器, 如 `#slider`.
// xlap: 拖动的距离, 单位px.
func DragSlider(sel interface{}, xlap int) chromedp.QueryAction {
	logger.Info("开始拖动滑块")
	return chromedp.QueryAfter(sel, func(ctx context.Context, _ runtime.ExecutionContextID, node ...*cdp.Node) error {
		if len(node) == 0 {
			return fmt.Errorf("找不到相关 Node")
		}

		return MouseDragNode(node[0], xlap).Do(ctx)
	}, chromedp.ByQuery)
}

func MouseDragNode(n *cdp.Node, xlap int) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		boxes, err := dom.GetContentQuads().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			logger.Error("获取节点坐标失败", err)
			return err
		}
		if len(boxes) == 0 {
			logger.Error("节点坐标为空")
			return chromedp.ErrInvalidDimensions
		}

		box := boxes[0]
		c := len(box)
		if c%2 != 0 || c < 1 {
			logger.Error("节点坐标格式错误")
			return chromedp.ErrInvalidDimensions
		}

		var mx, my float64
		for i := 0; i < c; i += 2 {
			mx += box[i]
			my += box[i+1]
		}

		mx /= float64(c / 2)
		my /= float64(c / 2)

		p := &input.DispatchMouseEventParams{
			Type:       input.MousePressed, // 鼠标左键按下
			X:          mx,
			Y:          my,
			Button:     input.Left,
			ClickCount: 1,
		}

		// 鼠标左键按下
		if err = p.Do(ctx); err != nil {
			logger.Error("鼠标左键按下失败", err)
			return err
		}
		logger.Info("鼠标左键按下座标: ", mx, my)

		// 设置鼠标移动
		p.Type = input.MouseMoved

		t := rand.Intn(20) + 40
		totalX := 0
		// 生成随机的路径,模拟拖动
		for i := 0; i < t; i++ {
			// 随机等待
			rt := rand.Intn(20) + 40
			if err = chromedp.Run(ctx, chromedp.Sleep(time.Millisecond*time.Duration(rt))); err != nil {
				logger.Error("随机等待失败")
				continue
			}

			// 随机移动的距离
			x := rand.Intn(3) + 3

			// 限制移动的距离
			if totalX >= xlap {
				break
			}

			// 如果移动的距离大于剩余距离, 则移动到剩余距离
			if totalX+x >= xlap {
				x = xlap - totalX
			}
			totalX += x

			// 随机移动的方向
			y := rand.Intn(2)

			p.Y += float64(y)
			p.X += (float64(x) + 0.1)

			if err = p.Do(ctx); err != nil {
				logger.Error("拖动失败", err)
				return err
			}
		}
		logger.Info("拖动结束 totalX: ", totalX)
		// 鼠标松开
		p.Type = input.MouseReleased
		return p.Do(ctx)
	}
}
