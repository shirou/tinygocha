package graphics

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Minimap represents the minimap display
type Minimap struct {
	camera *CameraManager
	
	// Position and size
	X, Y          int
	Width, Height int
	Scale         float64 // World to minimap scale
	
	// Display settings
	Visible       bool
	ShowUnits     bool
	ShowTerrain   bool
	ShowViewport  bool
	
	// Images
	backgroundImage *ebiten.Image
	minimapImage    *ebiten.Image
	
	// Update control
	needUpdate    bool
	updateCounter int
	updateFreq    int // Update every N frames
	
	// Interaction
	isDragging    bool
	dragStartX    int
	dragStartY    int
	
	// Colors
	backgroundColor   color.Color
	viewportColor     color.Color
	friendlyUnitColor color.Color
	enemyUnitColor    color.Color
	terrainColors     map[string]color.Color
}

// NewMinimap creates a new minimap
func NewMinimap(camera *CameraManager, x, y, width, height int) *Minimap {
	worldWidth := camera.WorldWidth
	worldHeight := camera.WorldHeight
	
	// Calculate scale to fit world in minimap
	scaleX := float64(width) / worldWidth
	scaleY := float64(height) / worldHeight
	scale := math.Min(scaleX, scaleY)
	
	minimap := &Minimap{
		camera:            camera,
		X:                 x,
		Y:                 y,
		Width:             width,
		Height:            height,
		Scale:             scale,
		Visible:           true,
		ShowUnits:         true,
		ShowTerrain:       true,
		ShowViewport:      true,
		needUpdate:        true,
		updateFreq:        2, // Update every 2 frames (30 FPS when main is 60 FPS)
		backgroundColor:   color.RGBA{40, 40, 40, 200},
		viewportColor:     color.RGBA{255, 255, 255, 255},
		friendlyUnitColor: color.RGBA{0, 255, 0, 255},
		enemyUnitColor:    color.RGBA{255, 0, 0, 255},
		terrainColors: map[string]color.Color{
			"plain":    color.RGBA{100, 150, 100, 255},
			"forest":   color.RGBA{50, 100, 50, 255},
			"mountain": color.RGBA{120, 120, 120, 255},
			"water":    color.RGBA{50, 50, 150, 255},
			"fortress": color.RGBA{150, 150, 100, 255},
			"town":     color.RGBA{120, 100, 80, 255},
		},
	}
	
	// Create images
	minimap.backgroundImage = ebiten.NewImage(width, height)
	minimap.minimapImage = ebiten.NewImage(width, height)
	
	// Fill background
	minimap.backgroundImage.Fill(minimap.backgroundColor)
	
	return minimap
}

// Update updates the minimap
func (m *Minimap) Update() {
	if !m.Visible {
		return
	}
	
	// Handle input
	m.handleInput()
	
	// Update minimap image periodically
	m.updateCounter++
	if m.updateCounter >= m.updateFreq || m.needUpdate {
		m.updateMinimapImage()
		m.updateCounter = 0
		m.needUpdate = false
	}
}

// Draw draws the minimap
func (m *Minimap) Draw(screen *ebiten.Image) {
	if !m.Visible {
		return
	}
	
	// Draw background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(m.backgroundImage, op)
	
	// Draw minimap content
	screen.DrawImage(m.minimapImage, op)
	
	// Draw viewport rectangle
	if m.ShowViewport {
		m.drawViewport(screen)
	}
	
	// Draw border
	m.drawBorder(screen)
}

// updateMinimapImage updates the minimap image content
func (m *Minimap) updateMinimapImage() {
	m.minimapImage.Clear()
	
	// Draw terrain (simplified)
	if m.ShowTerrain {
		m.drawTerrain()
	}
	
	// Draw units would go here when unit system is integrated
	if m.ShowUnits {
		// TODO: Draw units when unit system is available
	}
}

// drawTerrain draws simplified terrain on minimap
func (m *Minimap) drawTerrain() {
	// For now, draw a simple terrain pattern
	// This would be replaced with actual terrain data
	
	// Draw some sample terrain areas
	terrainAreas := []struct {
		x, y, w, h int
		terrainType string
	}{
		{int(1000 * m.Scale), int(1000 * m.Scale), int(1000 * m.Scale), int(1000 * m.Scale), "forest"},
		{int(3000 * m.Scale), int(1500 * m.Scale), int(800 * m.Scale), int(800 * m.Scale), "mountain"},
		{int(2000 * m.Scale), int(3000 * m.Scale), int(1500 * m.Scale), int(500 * m.Scale), "water"},
	}
	
	for _, area := range terrainAreas {
		if color, exists := m.terrainColors[area.terrainType]; exists {
			// Create a small image for the terrain area
			terrainImg := ebiten.NewImage(area.w, area.h)
			terrainImg.Fill(color)
			
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(area.x), float64(area.y))
			m.minimapImage.DrawImage(terrainImg, op)
		}
	}
}

// drawViewport draws the current viewport rectangle
func (m *Minimap) drawViewport(screen *ebiten.Image) {
	// Calculate viewport position and size in minimap coordinates
	camX, camY := m.camera.GetPosition()
	zoom := m.camera.GetZoom()
	
	viewWidth := float64(m.camera.ViewportWidth) / zoom
	viewHeight := float64(m.camera.ViewportHeight) / zoom
	
	// Convert to minimap coordinates
	minimapX := int(camX * m.Scale) + m.X
	minimapY := int(camY * m.Scale) + m.Y
	minimapW := int(viewWidth * m.Scale)
	minimapH := int(viewHeight * m.Scale)
	
	// Ensure viewport rectangle stays within minimap bounds
	if minimapX < m.X {
		minimapW -= m.X - minimapX
		minimapX = m.X
	}
	if minimapY < m.Y {
		minimapH -= m.Y - minimapY
		minimapY = m.Y
	}
	if minimapX + minimapW > m.X + m.Width {
		minimapW = m.X + m.Width - minimapX
	}
	if minimapY + minimapH > m.Y + m.Height {
		minimapH = m.Y + m.Height - minimapY
	}
	
	// Draw viewport rectangle outline
	if minimapW > 0 && minimapH > 0 {
		ebitenutil.DrawRect(screen, float64(minimapX), float64(minimapY), float64(minimapW), 2, m.viewportColor)
		ebitenutil.DrawRect(screen, float64(minimapX), float64(minimapY+minimapH-2), float64(minimapW), 2, m.viewportColor)
		ebitenutil.DrawRect(screen, float64(minimapX), float64(minimapY), 2, float64(minimapH), m.viewportColor)
		ebitenutil.DrawRect(screen, float64(minimapX+minimapW-2), float64(minimapY), 2, float64(minimapH), m.viewportColor)
	}
}

// drawBorder draws the minimap border
func (m *Minimap) drawBorder(screen *ebiten.Image) {
	borderColor := color.RGBA{200, 200, 200, 255}
	
	// Draw border
	ebitenutil.DrawRect(screen, float64(m.X-1), float64(m.Y-1), float64(m.Width+2), 1, borderColor)
	ebitenutil.DrawRect(screen, float64(m.X-1), float64(m.Y+m.Height), float64(m.Width+2), 1, borderColor)
	ebitenutil.DrawRect(screen, float64(m.X-1), float64(m.Y-1), 1, float64(m.Height+2), borderColor)
	ebitenutil.DrawRect(screen, float64(m.X+m.Width), float64(m.Y-1), 1, float64(m.Height+2), borderColor)
}

// handleInput handles minimap input
func (m *Minimap) handleInput() {
	mouseX, mouseY := ebiten.CursorPosition()
	
	// Check if mouse is over minimap
	if mouseX >= m.X && mouseX < m.X+m.Width && mouseY >= m.Y && mouseY < m.Y+m.Height {
		// Handle left click - move camera to clicked position
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			m.handleMinimapClick(mouseX, mouseY)
		}
		
		// Handle drag - start dragging viewport
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			m.isDragging = true
			m.dragStartX = mouseX
			m.dragStartY = mouseY
		}
	}
	
	// Handle dragging
	if m.isDragging {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			m.handleMinimapDrag(mouseX, mouseY)
		} else {
			m.isDragging = false
		}
	}
	
	// Handle right click - toggle minimap visibility
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if mouseX >= m.X && mouseX < m.X+m.Width && mouseY >= m.Y && mouseY < m.Y+m.Height {
			m.Visible = !m.Visible
		}
	}
}

// handleMinimapClick handles clicking on the minimap
func (m *Minimap) handleMinimapClick(mouseX, mouseY int) {
	// Convert minimap coordinates to world coordinates
	relativeX := mouseX - m.X
	relativeY := mouseY - m.Y
	
	worldX := float64(relativeX) / m.Scale
	worldY := float64(relativeY) / m.Scale
	
	// Center camera on clicked position
	viewWidth := float64(m.camera.ViewportWidth) / m.camera.GetZoom()
	viewHeight := float64(m.camera.ViewportHeight) / m.camera.GetZoom()
	
	targetX := worldX - viewWidth/2
	targetY := worldY - viewHeight/2
	
	m.camera.SetTargetPosition(targetX, targetY)
}

// handleMinimapDrag handles dragging on the minimap
func (m *Minimap) handleMinimapDrag(mouseX, mouseY int) {
	// Calculate drag delta
	deltaX := mouseX - m.dragStartX
	deltaY := mouseY - m.dragStartY
	
	// Convert to world coordinates
	worldDeltaX := float64(deltaX) / m.Scale
	worldDeltaY := float64(deltaY) / m.Scale
	
	// Move camera
	camX, camY := m.camera.GetPosition()
	m.camera.SetTargetPosition(camX+worldDeltaX, camY+worldDeltaY)
	
	// Update drag start position
	m.dragStartX = mouseX
	m.dragStartY = mouseY
}

// WorldToMinimap converts world coordinates to minimap coordinates
func (m *Minimap) WorldToMinimap(worldX, worldY float64) (int, int) {
	minimapX := int(worldX*m.Scale) + m.X
	minimapY := int(worldY*m.Scale) + m.Y
	return minimapX, minimapY
}

// MinimapToWorld converts minimap coordinates to world coordinates
func (m *Minimap) MinimapToWorld(minimapX, minimapY int) (float64, float64) {
	worldX := float64(minimapX-m.X) / m.Scale
	worldY := float64(minimapY-m.Y) / m.Scale
	return worldX, worldY
}

// SetVisible sets the minimap visibility
func (m *Minimap) SetVisible(visible bool) {
	m.Visible = visible
}

// IsVisible returns whether the minimap is visible
func (m *Minimap) IsVisible() bool {
	return m.Visible
}

// SetShowUnits sets whether to show units on minimap
func (m *Minimap) SetShowUnits(show bool) {
	m.ShowUnits = show
	m.needUpdate = true
}

// SetShowTerrain sets whether to show terrain on minimap
func (m *Minimap) SetShowTerrain(show bool) {
	m.ShowTerrain = show
	m.needUpdate = true
}

// SetPosition sets the minimap position
func (m *Minimap) SetPosition(x, y int) {
	m.X = x
	m.Y = y
}

// GetBounds returns the minimap bounds
func (m *Minimap) GetBounds() (x, y, width, height int) {
	return m.X, m.Y, m.Width, m.Height
}
