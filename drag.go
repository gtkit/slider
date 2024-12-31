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
)

func DragSlider(sel interface{}, xlap int) chromedp.QueryAction {
	return chromedp.QueryAfter(sel, func(ctx context.Context, _ runtime.ExecutionContextID, node ...*cdp.Node) error {
		if len(node) == 0 {
			return fmt.Errorf("找不到相关 Node")
		}

		return MouseDragNode(node[0], xlap).Do(ctx)
	})
}

func MouseDragNode(n *cdp.Node, xlap int) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		boxes, err := dom.GetContentQuads().WithNodeID(n.NodeID).Do(ctx)
		if err != nil {
			return err
		}
		if len(boxes) == 0 {
			return chromedp.ErrInvalidDimensions
		}

		box := boxes[0]
		c := len(box)
		if c%2 != 0 || c < 1 {
			return chromedp.ErrInvalidDimensions
		}

		var x, y float64
		for i := 0; i < c; i += 2 {
			x += box[i]
			y += box[i+1]
		}
		x /= float64(c / 2)
		y /= float64(c / 2)

		p := &input.DispatchMouseEventParams{
			Type:       input.MousePressed,
			X:          x,
			Y:          y,
			Button:     input.Left,
			ClickCount: 1,
		}

		// 鼠标左键按下
		if err = p.Do(ctx); err != nil {
			return err
		}

		// 拖动
		p.Type = input.MouseMoved

		t := rand.Intn(20) + 40
		totalX := 0
		// 生成随机的路径,模拟拖动
		for i := 0; i < t; i++ {
			rt := rand.Intn(20) + 20
			if err = chromedp.Run(ctx, chromedp.Sleep(time.Millisecond*time.Duration(rt))); err != nil {
				continue
			}
			x := rand.Intn(2) + 4
			if totalX >= xlap {
				break
			}
			if totalX+x >= xlap {
				x = xlap - totalX
			}

			totalX += x
			y := rand.Intn(2)

			p.Y += float64(y)
			p.X += float64(x)

			if err = p.Do(ctx); err != nil {
				return err
			}
		}
		// 鼠标松开
		p.Type = input.MouseReleased
		return p.Do(ctx)
	}
}
