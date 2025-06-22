package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TextRenderer handles text rendering with proper fonts
type TextRenderer struct {
	fontManager *FontManager
}

// NewTextRenderer creates a new text renderer
func NewTextRenderer(fontManager *FontManager) *TextRenderer {
	return &TextRenderer{
		fontManager: fontManager,
	}
}

// DrawText draws text at the specified position
func (tr *TextRenderer) DrawText(screen *ebiten.Image, str string, x, y float64, clr color.Color) {
	font := tr.fontManager.GetDefaultFont()
	if font == nil {
		return
	}
	
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	
	text.Draw(screen, str, font, op)
}

// DrawTextWithFont draws text with a specific font
func (tr *TextRenderer) DrawTextWithFont(screen *ebiten.Image, str string, x, y float64, clr color.Color, fontName string) {
	font := tr.fontManager.GetFont(fontName)
	if font == nil {
		font = tr.fontManager.GetDefaultFont()
	}
	
	if font == nil {
		return
	}
	
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	
	text.Draw(screen, str, font, op)
}

// DrawTextWithSize draws text with a specific size
func (tr *TextRenderer) DrawTextWithSize(screen *ebiten.Image, str string, x, y float64, clr color.Color, size float64) {
	font := tr.fontManager.CreateFontVariant("default", size)
	if font == nil {
		font = tr.fontManager.GetDefaultFont()
	}
	
	if font == nil {
		return
	}
	
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	
	text.Draw(screen, str, font, op)
}

// MeasureText measures the size of text
func (tr *TextRenderer) MeasureText(str string) (float64, float64) {
	font := tr.fontManager.GetDefaultFont()
	if font == nil {
		return 0, 0
	}
	
	width, height := text.Measure(str, font, 0)
	return width, height
}

// MeasureTextWithFont measures text with a specific font
func (tr *TextRenderer) MeasureTextWithFont(str string, fontName string) (float64, float64) {
	font := tr.fontManager.GetFont(fontName)
	if font == nil {
		font = tr.fontManager.GetDefaultFont()
	}
	
	if font == nil {
		return 0, 0
	}
	
	width, height := text.Measure(str, font, 0)
	return width, height
}

// DrawCenteredText draws text centered at the specified position
func (tr *TextRenderer) DrawCenteredText(screen *ebiten.Image, str string, centerX, centerY float64, clr color.Color) {
	width, height := tr.MeasureText(str)
	x := centerX - width/2
	y := centerY - height/2
	tr.DrawText(screen, str, x, y, clr)
}

// DrawTextWithShadow draws text with a shadow effect
func (tr *TextRenderer) DrawTextWithShadow(screen *ebiten.Image, str string, x, y float64, textColor, shadowColor color.Color) {
	// Draw shadow (offset by 1 pixel)
	tr.DrawText(screen, str, x+1, y+1, shadowColor)
	// Draw main text
	tr.DrawText(screen, str, x, y, textColor)
}
