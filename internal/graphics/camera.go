package graphics

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// CameraManager manages the game camera position and zoom
type CameraManager struct {
	// Current position and zoom
	X, Y float64
	Zoom float64
	
	// Target position for smooth movement
	TargetX, TargetY float64
	TargetZoom       float64
	
	// Constraints
	MinX, MinY       float64
	MaxX, MaxY       float64
	MinZoom, MaxZoom float64
	
	// Viewport size
	ViewportWidth, ViewportHeight int
	
	// World size
	WorldWidth, WorldHeight float64
	
	// Settings
	ScrollSpeed float64
	ZoomSpeed   float64
	SmoothMove  bool
}

// NewCameraManager creates a new camera manager
func NewCameraManager(worldWidth, worldHeight float64, viewportWidth, viewportHeight int) *CameraManager {
	camera := &CameraManager{
		X:              worldWidth/2 - float64(viewportWidth)/2,  // Center initially
		Y:              worldHeight/2 - float64(viewportHeight)/2,
		Zoom:           1.0,
		TargetX:        worldWidth/2 - float64(viewportWidth)/2,
		TargetY:        worldHeight/2 - float64(viewportHeight)/2,
		TargetZoom:     1.0,
		MinX:           0,
		MinY:           0,
		MaxX:           worldWidth - float64(viewportWidth),
		MaxY:           worldHeight - float64(viewportHeight),
		MinZoom:        0.25,
		MaxZoom:        2.0,
		ViewportWidth:  viewportWidth,
		ViewportHeight: viewportHeight,
		WorldWidth:     worldWidth,
		WorldHeight:    worldHeight,
		ScrollSpeed:    800.0, // 100.0 -> 800.0 (8倍速)
		ZoomSpeed:      4.0,   // 2.0 -> 4.0 (2倍速)
		SmoothMove:     false, // true -> false (即座に移動)
	}
	
	camera.updateConstraints()
	return camera
}

// Update updates the camera position and zoom with smooth movement
func (c *CameraManager) Update(deltaTime float64) {
	if c.SmoothMove {
		// Smooth movement towards target
		moveSpeed := c.ScrollSpeed * deltaTime
		
		// Move X
		if math.Abs(c.TargetX-c.X) > 1.0 {
			if c.TargetX > c.X {
				c.X = math.Min(c.X+moveSpeed, c.TargetX)
			} else {
				c.X = math.Max(c.X-moveSpeed, c.TargetX)
			}
		} else {
			c.X = c.TargetX
		}
		
		// Move Y
		if math.Abs(c.TargetY-c.Y) > 1.0 {
			if c.TargetY > c.Y {
				c.Y = math.Min(c.Y+moveSpeed, c.TargetY)
			} else {
				c.Y = math.Max(c.Y-moveSpeed, c.TargetY)
			}
		} else {
			c.Y = c.TargetY
		}
		
		// Smooth zoom
		if math.Abs(c.TargetZoom-c.Zoom) > 0.01 {
			zoomSpeed := c.ZoomSpeed * deltaTime
			if c.TargetZoom > c.Zoom {
				c.Zoom = math.Min(c.Zoom+zoomSpeed, c.TargetZoom)
			} else {
				c.Zoom = math.Max(c.Zoom-zoomSpeed, c.TargetZoom)
			}
		} else {
			c.Zoom = c.TargetZoom
		}
	} else {
		// Immediate movement
		c.X = c.TargetX
		c.Y = c.TargetY
		c.Zoom = c.TargetZoom
	}
	
	// Apply constraints
	c.applyConstraints()
}

// SetPosition sets the camera position immediately
func (c *CameraManager) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
	c.TargetX = x
	c.TargetY = y
	c.applyConstraints()
}

// SetTargetPosition sets the target position for smooth movement
func (c *CameraManager) SetTargetPosition(x, y float64) {
	c.TargetX = x
	c.TargetY = y
	c.applyTargetConstraints()
}

// Move moves the camera by the specified offset
func (c *CameraManager) Move(dx, dy float64) {
	c.SetTargetPosition(c.TargetX+dx, c.TargetY+dy)
}

// SetZoom sets the zoom level immediately
func (c *CameraManager) SetZoom(zoom float64) {
	c.Zoom = math.Max(c.MinZoom, math.Min(c.MaxZoom, zoom))
	c.TargetZoom = c.Zoom
	c.updateConstraints()
}

// SetTargetZoom sets the target zoom for smooth zooming
func (c *CameraManager) SetTargetZoom(zoom float64) {
	c.TargetZoom = math.Max(c.MinZoom, math.Min(c.MaxZoom, zoom))
	c.updateConstraints()
}

// ZoomAt zooms at a specific screen point
func (c *CameraManager) ZoomAt(screenX, screenY int, zoomDelta float64) {
	// Convert screen point to world coordinates before zoom
	worldX, worldY := c.ScreenToWorld(screenX, screenY)
	
	// Apply zoom
	newZoom := c.TargetZoom + zoomDelta
	c.SetTargetZoom(newZoom)
	
	// Convert world point back to screen coordinates after zoom
	newScreenX, newScreenY := c.WorldToScreen(worldX, worldY)
	
	// Adjust camera position to keep the point under the cursor
	c.Move(float64(newScreenX-screenX)/c.TargetZoom, float64(newScreenY-screenY)/c.TargetZoom)
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *CameraManager) ScreenToWorld(screenX, screenY int) (float64, float64) {
	worldX := c.X + float64(screenX)/c.Zoom
	worldY := c.Y + float64(screenY)/c.Zoom
	return worldX, worldY
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *CameraManager) WorldToScreen(worldX, worldY float64) (int, int) {
	screenX := int((worldX - c.X) * c.Zoom)
	screenY := int((worldY - c.Y) * c.Zoom)
	return screenX, screenY
}

// IsVisible checks if a world rectangle is visible on screen
func (c *CameraManager) IsVisible(worldX, worldY, width, height float64) bool {
	// Add margin for smooth scrolling
	margin := 100.0
	
	left := c.X - margin
	right := c.X + float64(c.ViewportWidth)/c.Zoom + margin
	top := c.Y - margin
	bottom := c.Y + float64(c.ViewportHeight)/c.Zoom + margin
	
	return worldX+width >= left && worldX <= right && worldY+height >= top && worldY <= bottom
}

// GetViewBounds returns the current view bounds in world coordinates
func (c *CameraManager) GetViewBounds() (left, top, right, bottom float64) {
	left = c.X
	top = c.Y
	right = c.X + float64(c.ViewportWidth)/c.Zoom
	bottom = c.Y + float64(c.ViewportHeight)/c.Zoom
	return
}

// GetTransform returns the transformation matrix for rendering
func (c *CameraManager) GetTransform() ebiten.GeoM {
	var transform ebiten.GeoM
	
	// Apply zoom
	transform.Scale(c.Zoom, c.Zoom)
	
	// Apply camera translation
	transform.Translate(-c.X*c.Zoom, -c.Y*c.Zoom)
	
	return transform
}

// updateConstraints updates the camera movement constraints based on zoom
func (c *CameraManager) updateConstraints() {
	viewWidth := float64(c.ViewportWidth) / c.Zoom
	viewHeight := float64(c.ViewportHeight) / c.Zoom
	
	c.MaxX = c.WorldWidth - viewWidth
	c.MaxY = c.WorldHeight - viewHeight
	
	// Ensure min constraints don't exceed max
	if c.MaxX < c.MinX {
		c.MaxX = c.MinX
	}
	if c.MaxY < c.MinY {
		c.MaxY = c.MinY
	}
}

// applyConstraints applies position and zoom constraints
func (c *CameraManager) applyConstraints() {
	c.X = math.Max(c.MinX, math.Min(c.MaxX, c.X))
	c.Y = math.Max(c.MinY, math.Min(c.MaxY, c.Y))
	c.Zoom = math.Max(c.MinZoom, math.Min(c.MaxZoom, c.Zoom))
}

// applyTargetConstraints applies constraints to target position
func (c *CameraManager) applyTargetConstraints() {
	c.TargetX = math.Max(c.MinX, math.Min(c.MaxX, c.TargetX))
	c.TargetY = math.Max(c.MinY, math.Min(c.MaxY, c.TargetY))
}

// GetPosition returns the current camera position
func (c *CameraManager) GetPosition() (float64, float64) {
	return c.X, c.Y
}

// GetZoom returns the current zoom level
func (c *CameraManager) GetZoom() float64 {
	return c.Zoom
}

// SetScrollSpeed sets the camera scroll speed
func (c *CameraManager) SetScrollSpeed(speed float64) {
	c.ScrollSpeed = speed
}

// SetSmoothMove enables or disables smooth camera movement
func (c *CameraManager) SetSmoothMove(smooth bool) {
	c.SmoothMove = smooth
}
