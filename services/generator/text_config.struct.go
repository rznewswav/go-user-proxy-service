package generator

import (
	"image/color"

	"github.com/tdewolff/canvas"
)

type TextConfig struct {
	VariableName string
	RelativeX    float64
	RelativeY    float64
	FontSize     int
	MaxWidth     float64
	MaxHeight    float64
	TextAlign    canvas.TextAlign
	Color        color.RGBA
	UseFont      string
}
