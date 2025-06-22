package graphics

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteGenerator generates unit sprites programmatically
type SpriteGenerator struct {
	cache map[string]*ebiten.Image
}

// NewSpriteGenerator creates a new sprite generator
func NewSpriteGenerator() *SpriteGenerator {
	return &SpriteGenerator{
		cache: make(map[string]*ebiten.Image),
	}
}

// GenerateUnitSprite generates an animated sprite for a unit
func (sg *SpriteGenerator) GenerateUnitSprite(unitType string, baseColor color.RGBA, isLeader bool, animState *AnimationState) *ebiten.Image {
	size := 16
	if isLeader {
		size = 20
	}
	
	// Apply scale modifier from animation
	scale := animState.GetScaleModifier()
	actualSize := int(float64(size) * scale)
	
	// Create image
	img := ebiten.NewImage(actualSize*2, actualSize*2) // Extra space for effects
	
	// Get animation offsets
	offsetX, offsetY := animState.GetAnimationOffset()
	rotation := animState.GetRotationModifier()
	
	centerX := actualSize
	centerY := actualSize
	
	// Apply offsets
	centerX += int(offsetX)
	centerY += int(offsetY)
	
	// Draw unit shape based on type
	switch unitType {
	case "infantry":
		sg.drawAnimatedSquare(img, centerX, centerY, actualSize/2, baseColor, isLeader, animState, rotation)
	case "archer":
		sg.drawAnimatedTriangle(img, centerX, centerY, actualSize/2, baseColor, isLeader, animState, rotation)
	case "mage":
		sg.drawAnimatedDiamond(img, centerX, centerY, actualSize/2, baseColor, isLeader, animState, rotation)
	default:
		sg.drawAnimatedCircle(img, centerX, centerY, actualSize/2, baseColor, isLeader, animState, rotation)
	}
	
	return img
}

// drawAnimatedSquare draws an animated square (infantry)
func (sg *SpriteGenerator) drawAnimatedSquare(img *ebiten.Image, centerX, centerY, size int, baseColor color.RGBA, isLeader bool, animState *AnimationState, rotation float64) {
	// Animation-specific modifications
	var sizeModX, sizeModY int = size, size
	
	switch animState.Type {
	case AnimationWalk:
		// Slight stretching during walk
		if animState.Frame%2 == 0 {
			sizeModY = int(float64(size) * 0.9)
		}
	case AnimationAttack:
		// Stretch forward during attack
		if animState.Frame == 1 {
			sizeModX = int(float64(size) * 1.3)
		}
	}
	
	// Draw main body
	for dy := -sizeModY; dy <= sizeModY; dy++ {
		for dx := -sizeModX; dx <= sizeModX; dx++ {
			// Apply rotation if needed
			x, y := sg.rotatePoint(float64(dx), float64(dy), rotation)
			img.Set(centerX+int(x), centerY+int(y), baseColor)
		}
	}
	
	// Draw leader border
	if isLeader {
		borderColor := color.RGBA{255, 255, 255, 255}
		// Top and bottom borders
		for dx := -sizeModX; dx <= sizeModX; dx++ {
			x1, y1 := sg.rotatePoint(float64(dx), float64(-sizeModY), rotation)
			x2, y2 := sg.rotatePoint(float64(dx), float64(sizeModY), rotation)
			img.Set(centerX+int(x1), centerY+int(y1), borderColor)
			img.Set(centerX+int(x2), centerY+int(y2), borderColor)
		}
		// Left and right borders
		for dy := -sizeModY; dy <= sizeModY; dy++ {
			x1, y1 := sg.rotatePoint(float64(-sizeModX), float64(dy), rotation)
			x2, y2 := sg.rotatePoint(float64(sizeModX), float64(dy), rotation)
			img.Set(centerX+int(x1), centerY+int(y1), borderColor)
			img.Set(centerX+int(x2), centerY+int(y2), borderColor)
		}
	}
	
	// Add animation-specific effects
	sg.addAnimationEffects(img, centerX, centerY, size, animState)
}

// drawAnimatedTriangle draws an animated triangle (archer)
func (sg *SpriteGenerator) drawAnimatedTriangle(img *ebiten.Image, centerX, centerY, size int, baseColor color.RGBA, isLeader bool, animState *AnimationState, rotation float64) {
	// Animation-specific modifications
	heightMod := 1.0
	
	switch animState.Type {
	case AnimationAttack:
		// Point forward more during attack
		if animState.Frame == 1 {
			heightMod = 1.4
		}
	}
	
	actualSize := int(float64(size) * heightMod)
	
	// Draw triangle pointing up
	for dy := -actualSize; dy <= actualSize; dy++ {
		width := actualSize - int(math.Abs(float64(dy)))
		for dx := -width; dx <= width; dx++ {
			x, y := sg.rotatePoint(float64(dx), float64(dy), rotation)
			img.Set(centerX+int(x), centerY+int(y), baseColor)
		}
	}
	
	// Draw leader border
	if isLeader {
		borderColor := color.RGBA{255, 255, 255, 255}
		// Draw triangle outline
		for dy := -actualSize; dy <= actualSize; dy++ {
			width := actualSize - int(math.Abs(float64(dy)))
			if width >= 0 {
				x1, y1 := sg.rotatePoint(float64(-width), float64(dy), rotation)
				x2, y2 := sg.rotatePoint(float64(width), float64(dy), rotation)
				img.Set(centerX+int(x1), centerY+int(y1), borderColor)
				img.Set(centerX+int(x2), centerY+int(y2), borderColor)
			}
		}
	}
	
	sg.addAnimationEffects(img, centerX, centerY, size, animState)
}

// drawAnimatedDiamond draws an animated diamond (mage)
func (sg *SpriteGenerator) drawAnimatedDiamond(img *ebiten.Image, centerX, centerY, size int, baseColor color.RGBA, isLeader bool, animState *AnimationState, rotation float64) {
	// Animation-specific modifications
	pulseMod := 1.0
	
	switch animState.Type {
	case AnimationIdle:
		// Gentle pulsing for mages
		pulseMod = 1.0 + math.Sin(float64(animState.Frame)*math.Pi/2)*0.1
	case AnimationAttack:
		// Bright flash during attack
		if animState.Frame == 1 {
			pulseMod = 1.3
			// Make color brighter
			baseColor.R = uint8(math.Min(255, float64(baseColor.R)*1.2))
			baseColor.G = uint8(math.Min(255, float64(baseColor.G)*1.2))
			baseColor.B = uint8(math.Min(255, float64(baseColor.B)*1.2))
		}
	}
	
	actualSize := int(float64(size) * pulseMod)
	
	// Draw diamond
	for dy := -actualSize; dy <= actualSize; dy++ {
		width := actualSize - int(math.Abs(float64(dy)))
		for dx := -width; dx <= width; dx++ {
			x, y := sg.rotatePoint(float64(dx), float64(dy), rotation)
			img.Set(centerX+int(x), centerY+int(y), baseColor)
		}
	}
	
	// Draw leader border
	if isLeader {
		borderColor := color.RGBA{255, 255, 255, 255}
		for dy := -actualSize; dy <= actualSize; dy++ {
			width := actualSize - int(math.Abs(float64(dy)))
			if width >= 0 {
				x1, y1 := sg.rotatePoint(float64(-width), float64(dy), rotation)
				x2, y2 := sg.rotatePoint(float64(width), float64(dy), rotation)
				img.Set(centerX+int(x1), centerY+int(y1), borderColor)
				img.Set(centerX+int(x2), centerY+int(y2), borderColor)
			}
		}
	}
	
	sg.addAnimationEffects(img, centerX, centerY, size, animState)
}

// drawAnimatedCircle draws an animated circle
func (sg *SpriteGenerator) drawAnimatedCircle(img *ebiten.Image, centerX, centerY, size int, baseColor color.RGBA, isLeader bool, animState *AnimationState, rotation float64) {
	// Animation-specific modifications
	radiusMod := 1.0
	
	switch animState.Type {
	case AnimationWalk:
		// Slight oval shape during walk
		radiusMod = 1.0 + math.Sin(float64(animState.Frame)*math.Pi/2)*0.1
	}
	
	radius := int(float64(size) * radiusMod)
	
	// Draw circle
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				x, y := sg.rotatePoint(float64(dx), float64(dy), rotation)
				img.Set(centerX+int(x), centerY+int(y), baseColor)
			}
		}
	}
	
	// Draw leader border
	if isLeader {
		borderColor := color.RGBA{255, 255, 255, 255}
		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				dist := dx*dx + dy*dy
				if dist <= radius*radius && dist > (radius-2)*(radius-2) {
					x, y := sg.rotatePoint(float64(dx), float64(dy), rotation)
					img.Set(centerX+int(x), centerY+int(y), borderColor)
				}
			}
		}
	}
	
	sg.addAnimationEffects(img, centerX, centerY, size, animState)
}

// addAnimationEffects adds special effects based on animation state
func (sg *SpriteGenerator) addAnimationEffects(img *ebiten.Image, centerX, centerY, size int, animState *AnimationState) {
	switch animState.Type {
	case AnimationAttack:
		if animState.Frame == 1 {
			// Add attack flash effect
			flashColor := color.RGBA{255, 255, 0, 128} // Yellow flash
			for i := 0; i < 3; i++ {
				for angle := 0.0; angle < 2*math.Pi; angle += math.Pi / 4 {
					x := centerX + int(math.Cos(angle)*float64(size+i+2))
					y := centerY + int(math.Sin(angle)*float64(size+i+2))
					img.Set(x, y, flashColor)
				}
			}
		}
	case AnimationDeath:
		// Add fading effect
		alpha := uint8(255 * (1.0 - float64(animState.Frame)/float64(animState.TotalFrames)))
		fadeColor := color.RGBA{100, 100, 100, alpha}
		
		// Overlay fade effect
		for dy := -size-2; dy <= size+2; dy++ {
			for dx := -size-2; dx <= size+2; dx++ {
				img.Set(centerX+dx, centerY+dy, fadeColor)
			}
		}
	}
}

// rotatePoint rotates a point around the origin
func (sg *SpriteGenerator) rotatePoint(x, y, angle float64) (float64, float64) {
	if angle == 0 {
		return x, y
	}
	
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	
	newX := x*cos - y*sin
	newY := x*sin + y*cos
	
	return newX, newY
}
