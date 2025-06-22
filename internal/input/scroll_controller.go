package input

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/shirou/tinygocha/internal/graphics"
)

// ScrollController handles camera scrolling input
type ScrollController struct {
	camera *graphics.CameraManager
	
	// Settings
	EdgeScrolling bool    // Enable edge scrolling
	KeyScrolling  bool    // Enable keyboard scrolling
	DragScrolling bool    // Enable middle mouse drag scrolling
	
	// Edge scrolling settings
	EdgeWidth    int     // Edge width in pixels
	EdgeSpeed    float64 // Base edge scroll speed
	EdgeAccel    float64 // Edge scroll acceleration multiplier
	
	// Keyboard scrolling settings
	KeySpeed     float64 // Keyboard scroll speed
	
	// Drag scrolling state
	isDragging   bool
	dragStartX   int
	dragStartY   int
	dragLastX    int
	dragLastY    int
	
	// Zoom settings
	ZoomStep     float64 // Zoom step per wheel tick
	
	// Key states for smooth scrolling
	keyStates    map[ebiten.Key]float64 // Key press duration
}

// NewScrollController creates a new scroll controller
func NewScrollController(camera *graphics.CameraManager) *ScrollController {
	fmt.Println("ScrollController created successfully")
	return &ScrollController{
		camera:        camera,
		EdgeScrolling: true,
		KeyScrolling:  true,
		DragScrolling: true,
		EdgeWidth:     50,
		EdgeSpeed:     400.0,  // 100.0 -> 400.0 (4倍速)
		EdgeAccel:     3.0,    // 2.0 -> 3.0 (加速度アップ)
		KeySpeed:      500.0,  // 150.0 -> 500.0 (3.3倍速)
		ZoomStep:      0.25,
		keyStates:     make(map[ebiten.Key]float64),
	}
}

// Update processes input and updates camera accordingly
func (sc *ScrollController) Update(deltaTime float64) {
	// Debug: Check if Update is being called
	if deltaTime > 0 {
		// Only log occasionally to avoid spam
		if int(deltaTime*1000)%1000 < 50 { // Log roughly once per second
			fmt.Printf("ScrollController.Update called with deltaTime=%.3f\n", deltaTime)
		}
	}
	
	// Handle edge scrolling
	if sc.EdgeScrolling {
		sc.handleEdgeScrolling(deltaTime)
	}
	
	// Handle keyboard scrolling
	if sc.KeyScrolling {
		sc.handleKeyboardScrolling(deltaTime)
	}
	
	// Handle drag scrolling
	if sc.DragScrolling {
		sc.handleDragScrolling()
	}
	
	// Handle zoom
	sc.handleZoom()
}

// handleEdgeScrolling processes mouse edge scrolling
func (sc *ScrollController) handleEdgeScrolling(deltaTime float64) {
	mouseX, mouseY := ebiten.CursorPosition()
	screenWidth, screenHeight := ebiten.WindowSize()
	
	var scrollX, scrollY float64
	
	// Left edge
	if mouseX < sc.EdgeWidth {
		intensity := float64(sc.EdgeWidth-mouseX) / float64(sc.EdgeWidth)
		scrollX = -sc.EdgeSpeed * (1 + sc.EdgeAccel*intensity) * deltaTime
	}
	// Right edge
	if mouseX > screenWidth-sc.EdgeWidth {
		intensity := float64(mouseX-(screenWidth-sc.EdgeWidth)) / float64(sc.EdgeWidth)
		scrollX = sc.EdgeSpeed * (1 + sc.EdgeAccel*intensity) * deltaTime
	}
	
	// Top edge
	if mouseY < sc.EdgeWidth {
		intensity := float64(sc.EdgeWidth-mouseY) / float64(sc.EdgeWidth)
		scrollY = -sc.EdgeSpeed * (1 + sc.EdgeAccel*intensity) * deltaTime
	}
	// Bottom edge
	if mouseY > screenHeight-sc.EdgeWidth {
		intensity := float64(mouseY-(screenHeight-sc.EdgeWidth)) / float64(sc.EdgeWidth)
		scrollY = sc.EdgeSpeed * (1 + sc.EdgeAccel*intensity) * deltaTime
	}
	
	if scrollX != 0 || scrollY != 0 {
		sc.camera.Move(scrollX, scrollY)
	}
}

// handleKeyboardScrolling processes keyboard scrolling
func (sc *ScrollController) handleKeyboardScrolling(deltaTime float64) {
	keys := []ebiten.Key{
		ebiten.KeyW, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD,
		ebiten.KeyArrowUp, ebiten.KeyArrowLeft, ebiten.KeyArrowDown, ebiten.KeyArrowRight,
	}
	
	// Check if any movement keys are pressed
	anyKeyPressed := false
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			anyKeyPressed = true
			break
		}
	}
	
	if anyKeyPressed {
		fmt.Println("Movement keys detected!")
	}
	
	// Update key states
	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			sc.keyStates[key] += deltaTime
		} else {
			sc.keyStates[key] = 0
		}
	}
	
	var scrollX, scrollY float64
	
	// Calculate scroll based on pressed keys
	// Up movement
	if sc.keyStates[ebiten.KeyW] > 0 || sc.keyStates[ebiten.KeyArrowUp] > 0 {
		scrollY = -sc.KeySpeed * deltaTime
		fmt.Printf("Moving up: scrollY=%.2f\n", scrollY)
	}
	// Down movement
	if sc.keyStates[ebiten.KeyS] > 0 || sc.keyStates[ebiten.KeyArrowDown] > 0 {
		scrollY = sc.KeySpeed * deltaTime
		fmt.Printf("Moving down: scrollY=%.2f\n", scrollY)
	}
	// Left movement
	if sc.keyStates[ebiten.KeyA] > 0 || sc.keyStates[ebiten.KeyArrowLeft] > 0 {
		scrollX = -sc.KeySpeed * deltaTime
		fmt.Printf("Moving left: scrollX=%.2f\n", scrollX)
	}
	// Right movement
	if sc.keyStates[ebiten.KeyD] > 0 || sc.keyStates[ebiten.KeyArrowRight] > 0 {
		scrollX = sc.KeySpeed * deltaTime
		fmt.Printf("Moving right: scrollX=%.2f\n", scrollX)
	}
	
	// Apply zoom-adjusted scrolling
	zoomFactor := 1.0 / sc.camera.GetZoom()
	if scrollX != 0 || scrollY != 0 {
		fmt.Printf("Applying camera movement: (%.2f, %.2f) with zoom factor %.2f\n", scrollX, scrollY, zoomFactor)
		sc.camera.Move(scrollX*zoomFactor, scrollY*zoomFactor)
	}
}

// handleDragScrolling processes middle mouse button drag scrolling
func (sc *ScrollController) handleDragScrolling() {
	// Check for middle mouse button
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		sc.isDragging = true
		sc.dragStartX, sc.dragStartY = ebiten.CursorPosition()
		sc.dragLastX, sc.dragLastY = sc.dragStartX, sc.dragStartY
	}
	
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		sc.isDragging = false
	}
	
	if sc.isDragging {
		mouseX, mouseY := ebiten.CursorPosition()
		
		// Calculate movement delta
		deltaX := float64(sc.dragLastX - mouseX)
		deltaY := float64(sc.dragLastY - mouseY)
		
		// Apply zoom factor and sensitivity multiplier for faster drag scrolling
		zoomFactor := 1.0 / sc.camera.GetZoom()
		sensitivity := 2.0 // 2倍の感度
		
		if deltaX != 0 || deltaY != 0 {
			sc.camera.Move(deltaX*zoomFactor*sensitivity, deltaY*zoomFactor*sensitivity)
		}
		
		sc.dragLastX, sc.dragLastY = mouseX, mouseY
	}
}

// handleZoom processes mouse wheel zoom
func (sc *ScrollController) handleZoom() {
	_, wheelY := ebiten.Wheel()
	
	if wheelY != 0 {
		fmt.Printf("Mouse wheel detected: wheelY=%.2f\n", wheelY)
		mouseX, mouseY := ebiten.CursorPosition()
		zoomDelta := wheelY * sc.ZoomStep
		fmt.Printf("Applying zoom: delta=%.2f at (%d, %d)\n", zoomDelta, mouseX, mouseY)
		sc.camera.ZoomAt(mouseX, mouseY, zoomDelta)
	}
	
	// Handle keyboard zoom
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) || inpututil.IsKeyJustPressed(ebiten.KeyKPAdd) {
		fmt.Println("Zoom in key pressed")
		// Zoom in at screen center
		screenWidth, screenHeight := ebiten.WindowSize()
		sc.camera.ZoomAt(screenWidth/2, screenHeight/2, sc.ZoomStep)
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) || inpututil.IsKeyJustPressed(ebiten.KeyKPSubtract) {
		fmt.Println("Zoom out key pressed")
		// Zoom out at screen center
		screenWidth, screenHeight := ebiten.WindowSize()
		sc.camera.ZoomAt(screenWidth/2, screenHeight/2, -sc.ZoomStep)
	}
}

// SetEdgeScrolling enables or disables edge scrolling
func (sc *ScrollController) SetEdgeScrolling(enabled bool) {
	sc.EdgeScrolling = enabled
}

// SetKeyScrolling enables or disables keyboard scrolling
func (sc *ScrollController) SetKeyScrolling(enabled bool) {
	sc.KeyScrolling = enabled
}

// SetDragScrolling enables or disables drag scrolling
func (sc *ScrollController) SetDragScrolling(enabled bool) {
	sc.DragScrolling = enabled
}

// SetEdgeWidth sets the edge scrolling width
func (sc *ScrollController) SetEdgeWidth(width int) {
	sc.EdgeWidth = width
}

// SetScrollSpeed sets the scrolling speed
func (sc *ScrollController) SetScrollSpeed(edgeSpeed, keySpeed float64) {
	sc.EdgeSpeed = edgeSpeed
	sc.KeySpeed = keySpeed
}

// SetZoomStep sets the zoom step per wheel tick
func (sc *ScrollController) SetZoomStep(step float64) {
	sc.ZoomStep = step
}

// IsScrolling returns true if any scrolling is currently active
func (sc *ScrollController) IsScrolling() bool {
	// Check if any scroll keys are pressed
	scrollKeys := []ebiten.Key{
		ebiten.KeyW, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD,
		ebiten.KeyArrowUp, ebiten.KeyArrowLeft, ebiten.KeyArrowDown, ebiten.KeyArrowRight,
	}
	
	for _, key := range scrollKeys {
		if ebiten.IsKeyPressed(key) {
			return true
		}
	}
	
	// Check if dragging
	if sc.isDragging {
		return true
	}
	
	// Check edge scrolling
	if sc.EdgeScrolling {
		mouseX, mouseY := ebiten.CursorPosition()
		screenWidth, screenHeight := ebiten.WindowSize()
		
		if mouseX < sc.EdgeWidth || mouseX > screenWidth-sc.EdgeWidth ||
			mouseY < sc.EdgeWidth || mouseY > screenHeight-sc.EdgeWidth {
			return true
		}
	}
	
	return false
}
