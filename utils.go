package dsvg

import (
	"fmt"
	"golang.org/x/image/colornames"
	"image/color"
	"strconv"
	"strings"
)

const (
	PixelsPerInch = 96
)

func ParseUnits(input string) (int, error) {
	p, err := strconv.Atoi(input)
	if err == nil {
		return p, nil
	}

	if strings.HasSuffix(input, "in") {
		input = strings.TrimRight(input, "in")
		inches, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return 0, err
		}

		return int(inches * PixelsPerInch), nil
	}

	return 0, nil
}

func ParseViewBox(viewbox string) (minx, miny, width, height int) {
	split := strings.Split(viewbox, " ")

	minx, _ = strconv.Atoi(split[0])
	miny, _ = strconv.Atoi(split[1])
	width, _ = strconv.Atoi(split[2])
	height, _ = strconv.Atoi(split[3])

	return
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

func ParseColor(color string) (color.Color, error) {
	if strings.HasPrefix(color, "#") {
		return ParseHexColor(color)
	}
	return colornames.Map[color], nil
}
