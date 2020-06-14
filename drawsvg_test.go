package dsvg

import (
	"image"
	"os"
	"testing"
)

func TestDrawSVG(t *testing.T) {
	f, err := os.Open("test.svg")
	if err != nil {
		t.Fatal("Unable to open test SVG: ", err)
	}
	ctx, err := DrawSVG(f, image.Rect(0, 0, 400, 400), true)
	if err != nil {
		t.Fatal("Failed to draw SVG: ", err)
	}

	err = ctx.SavePNG("test.png")
	if err != nil {
		t.Fatal("Failed to write PNG: ", err)
	}
}
