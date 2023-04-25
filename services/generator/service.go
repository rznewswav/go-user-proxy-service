package generator

import (
	"bytes"
	"fmt"
	"image/color"
	"os"

	"image/png"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

const mmPerPixel = 0.2645833333

type FontConfig string

func GetFirstNonOverflowingText(
	text string,
	width float64,
	maxHeight float64,
	opts ...interface{},
) (canvasText *canvas.Text) {
	var fillColor = canvas.Black
	var halign = canvas.Top
	var maxFontSize = 70
	var useFont = DefaultFontPath

	for _, opt := range opts {
		switch casted := opt.(type) {
		case color.RGBA:
			fillColor = casted
		case canvas.TextAlign:
			halign = casted
		case int:
			maxFontSize = casted
		case FontConfig:
			useFont = string(casted)
		}

	}

	if maxFontSize < 10 {
		panic(fmt.Errorf("cannot write in font size less than 10pt: %d", maxFontSize))
	}

	for sizePt := maxFontSize; sizePt >= 10; sizePt-- {
		face := LoadFont(useFont).Face(float64(sizePt), fillColor)
		canvasText = canvas.NewTextBox(face, text, width, 0, halign, canvas.Left, .0, .0)
		if canvasText.OutlineBounds().H < maxHeight && canvasText.OutlineBounds().W < width {
			return canvasText
		}
	}
	return canvasText
}

type ImageTemplate struct {
	ResourcePath string
	Config       []TextConfig
}

type M = map[string]string

type RunParamDebugLayout bool

const DebugLayout RunParamDebugLayout = true

func GenerateImage(
	template ImageTemplate,
	variableValues M,
	opts ...interface{},
) (buffer bytes.Buffer) {
	var debugLayout = false
	for _, opt := range opts {
		switch casted := opt.(type) {
		case RunParamDebugLayout:
			debugLayout = casted == DebugLayout
		}
	}

	fileHandler, fileReadError := os.Open(template.ResourcePath)
	if fileReadError != nil {
		panic(fileReadError)
	}

	textConfigs := template.Config

	pngImage, pngImageError := canvas.NewPNGImage(fileHandler)
	if pngImageError != nil {
		panic(pngImageError)
	}

	pngSize := pngImage.Image.Bounds().Max

	c := canvas.NewFromSize(canvas.Size{
		W: float64(pngSize.X) * mmPerPixel,
		H: float64(pngSize.Y) * mmPerPixel,
	})

	ctx := canvas.NewContext(c)

	ctx.DrawImage(0, 0, pngImage, canvas.DefaultResolution)

	for _, config := range textConfigs {
		var text = config.VariableName
		if value, variableInValueMap := variableValues[config.VariableName]; variableInValueMap {
			text = value
		}
		textBox := GetFirstNonOverflowingText(text,
			config.MaxWidth,
			config.MaxHeight,
			config.FontSize,
			config.Color,
			config.TextAlign,
			FontConfig(config.UseFont),
		)
		var x float64
		var y float64
		var maxWidth = config.MaxWidth
		if textBox.Bounds().W > maxWidth {
			maxWidth = textBox.Bounds().W
		}
		switch config.TextAlign {
		case canvas.Center:
			{
				x = c.W*config.RelativeX - maxWidth*0.5
				y = c.H*config.RelativeY + textBox.Bounds().H*0.5
			}
		case canvas.Right:
			{
				x = c.W*config.RelativeX - maxWidth
				y = c.H*config.RelativeY + textBox.Bounds().H*0.5
			}
		default:
			{
				x = c.W * config.RelativeX
				y = c.H*config.RelativeY + textBox.Bounds().H*0.5
			}

		}

		if debugLayout {
			ctx.SetStrokeWidth(2)
			ctx.SetFillColor(canvas.Lightblue)
			ctx.MoveTo(x, y)
			ctx.LineTo(x+maxWidth, y)
			ctx.LineTo(x+maxWidth, y-textBox.Bounds().H)
			ctx.LineTo(x, y-textBox.Bounds().H)
			ctx.LineTo(x, y)
			ctx.Close()
			ctx.FillStroke()
		}

		ctx.DrawText(
			x,
			y,
			textBox,
		)
	}

	img := rasterizer.Draw(c, canvas.DPMM(3.2), canvas.DefaultColorSpace)
	png.Encode(&buffer, img)
	return buffer
}
