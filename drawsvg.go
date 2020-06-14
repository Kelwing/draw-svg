package dsvg

import (
	"fmt"
	"github.com/Kelwing/svgparser"
	"github.com/Kelwing/svgparser/utils"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"io"
	"math"
	"strconv"
)

func drawPath(ctx *gg.Context, path *utils.Path, xScale, yScale float64) {
	for _, subpath := range path.Subpaths {
		for _, command := range subpath.Commands {
			switch command.Symbol {
			case "M":
				ctx.MoveTo(command.Params[0]*xScale, command.Params[1]*yScale)
			case "m":
				p, _ := ctx.GetCurrentPoint()
				ctx.MoveTo(p.X+command.Params[0]*xScale, p.Y+command.Params[1]*yScale)
			case "L":
				ctx.LineTo(command.Params[0]*xScale, command.Params[1]*yScale)
			case "l":
				p, _ := ctx.GetCurrentPoint()
				ctx.LineTo(p.X+command.Params[0]*xScale, p.Y+command.Params[1]*yScale)
			case "H":
				p, _ := ctx.GetCurrentPoint()
				ctx.LineTo(command.Params[0]*xScale, p.Y)
			case "h":
				p, exists := ctx.GetCurrentPoint()
				if !exists {
					continue
				}

				endPoint := p.X + command.Params[0]*xScale
				ctx.LineTo(endPoint, p.Y)
			case "V":
				p, _ := ctx.GetCurrentPoint()
				ctx.LineTo(p.X, command.Params[0]*yScale)
			case "v":
				p, exists := ctx.GetCurrentPoint()
				if !exists {
					continue
				}

				endPoint := p.Y + command.Params[0]*yScale
				ctx.LineTo(p.X, endPoint)
			case "z", "Z":
				ctx.ClosePath()
			case "C":
				ctx.CubicTo(command.Params[0]*xScale, command.Params[1]*yScale, command.Params[2]*xScale, command.Params[3]*yScale, command.Params[4]*xScale, command.Params[5]*yScale)
			default:
				fmt.Println("Command: ", command.Symbol, " is unimplemented")
			}
		}
	}
	ctx.ClosePath()
}

func processElement(ctx *gg.Context, element *svgparser.Element, xScale, yScale float64) error {
	for _, child := range element.Children {
		switch child.Name {
		case "path":
			path, err := utils.PathParser(child.Attributes["d"])
			if err != nil {
				return err
			}
			drawPath(ctx, path, xScale, yScale)
			ctx.SetColor(color.Black)
			if fillColor, ok := child.Attributes["fill"]; ok {
				f, err := ParseColor(fillColor)
				if err != nil {
					return err
				}
				ctx.SetColor(f)
			}
			ctx.Fill()
			if strokeWidth, ok := child.Attributes["stroke-width"]; ok {
				lineWidth, err := strconv.Atoi(strokeWidth)
				if err != nil {
					return err
				}
				ctx.SetLineWidth(float64(lineWidth))
			}
			if strokeColor, ok := child.Attributes["stroke"]; ok {
				s, err := ParseColor(strokeColor)
				if err != nil {
					return err
				}
				ctx.SetColor(s)
				ctx.Stroke()
			}
		case "rect":
			width, err := strconv.Atoi(child.Attributes["width"])
			if err != nil {
				return err
			}
			height, err := strconv.Atoi(child.Attributes["height"])
			if err != nil {
				return err
			}

			var yOrigin, xOrigin int

			if yVal, ok := child.Attributes["y"]; ok {
				yOrigin, _ = strconv.Atoi(yVal)
				// a default of 0 is OK
			}

			if xVal, ok := child.Attributes["x"]; ok {
				xOrigin, _ = strconv.Atoi(xVal)
			}
			ctx.DrawRectangle(float64(xOrigin)*xScale, float64(yOrigin)*yScale, float64(width)*xScale, float64(height)*yScale)

			if fill, ok := child.Attributes["fill"]; ok {
				fillColor, err := ParseColor(fill)
				if err == nil {
					ctx.SetColor(fillColor)
					ctx.Fill()
				} else {
					return err
				}
			}
		case "g":
			return processElement(ctx, child, xScale, yScale)
		}
	}
	return nil
}

func DrawSVG(r io.Reader, target image.Rectangle, keepAspect bool) (*gg.Context, error) {
	element, _ := svgparser.Parse(r, false)
	var w, h int
	if viewbox, ok := element.Attributes["viewBox"]; ok {
		minx, miny, width, height := ParseViewBox(viewbox)
		w = width - minx
		h = height - miny
	}

	var xScale, yScale float64

	if keepAspect {
		xScale = math.Min(float64(target.Dx())/float64(w), float64(target.Dy())/float64(h))
		yScale = xScale
	} else {
		xScale = float64(target.Dx()) / float64(w)
		yScale = float64(target.Dy()) / float64(h)
	}

	ctx := gg.NewContext(int(float64(w)*xScale), int(float64(h)*yScale))

	err := processElement(ctx, element, xScale, yScale)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
