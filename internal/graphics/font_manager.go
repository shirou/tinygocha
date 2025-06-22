package graphics

import (
	"bytes"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// FontManager manages fonts for the game
type FontManager struct {
	defaultFont *text.GoTextFace
	fonts       map[string]*text.GoTextFace
}

// NewFontManager creates a new font manager
func NewFontManager() *FontManager {
	return &FontManager{
		fonts: make(map[string]*text.GoTextFace),
	}
}

// LoadDefaultFont loads the default MPlus1p font
func (fm *FontManager) LoadDefaultFont(size float64) error {
	// Load MPlus1p font from ebiten examples
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return err
	}
	
	fm.defaultFont = &text.GoTextFace{
		Source: source,
		Size:   size,
	}
	
	log.Printf("Default font (MPlus1p) loaded successfully")
	return nil
}

// LoadFontFromFile loads a font from a file
func (fm *FontManager) LoadFontFromFile(fontPath string, size float64, name string) error {
	if fontPath == "" {
		// Use default font
		return fm.LoadDefaultFont(size)
	}
	
	// Check if file exists
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		log.Printf("Font file not found: %s, using default font", fontPath)
		return fm.LoadDefaultFont(size)
	}
	
	// Read font file
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		log.Printf("Failed to read font file: %s, using default font", fontPath)
		return fm.LoadDefaultFont(size)
	}
	
	// Create font source
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Printf("Failed to parse font file: %s, using default font", fontPath)
		return fm.LoadDefaultFont(size)
	}
	
	goTextFace := &text.GoTextFace{
		Source: source,
		Size:   size,
	}
	
	if name == "default" {
		fm.defaultFont = goTextFace
	} else {
		fm.fonts[name] = goTextFace
	}
	
	log.Printf("Font loaded successfully: %s", fontPath)
	return nil
}

// GetDefaultFont returns the default font
func (fm *FontManager) GetDefaultFont() *text.GoTextFace {
	if fm.defaultFont == nil {
		// Fallback: load default font if not loaded
		if err := fm.LoadDefaultFont(16); err != nil {
			log.Printf("Failed to load fallback font: %v", err)
		}
	}
	return fm.defaultFont
}

// GetFont returns a named font or default if not found
func (fm *FontManager) GetFont(name string) *text.GoTextFace {
	if font, exists := fm.fonts[name]; exists {
		return font
	}
	return fm.GetDefaultFont()
}

// CreateFontVariant creates a font variant with different size
func (fm *FontManager) CreateFontVariant(baseFontName string, size float64) *text.GoTextFace {
	baseFont := fm.GetFont(baseFontName)
	if baseFont == nil {
		return nil
	}
	
	// Create new face with different size
	newFace := &text.GoTextFace{
		Source: baseFont.Source,
		Size:   size,
	}
	
	return newFace
}
