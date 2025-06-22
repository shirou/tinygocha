package scenes

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/shirou/tinygocha/internal/data"
	"github.com/shirou/tinygocha/internal/game"
	"github.com/shirou/tinygocha/internal/graphics"
	"github.com/shirou/tinygocha/internal/input"
)

// BattleSceneUnified represents the unified battle screen with all features
type BattleSceneUnified struct {
	sceneManager     *SceneManager
	battleManager    *game.BattleManager
	dataManager      *data.DataManager
	textRenderer     *graphics.TextRenderer
	spriteGenerator  *graphics.SpriteGenerator
	
	// Camera and scrolling
	camera           *graphics.CameraManager
	scrollController *input.ScrollController
	minimap          *graphics.Minimap
	
	// Game state
	isPaused         bool
	selectedUnit     *game.Unit
	showDebugInfo    bool
	showHelp         bool
	
	// Timing
	lastUpdate       time.Time
	deltaTime        float64
	helpToggleTime   time.Time
}

// NewBattleSceneUnified creates a new unified battle scene
func NewBattleSceneUnified(sceneManager *SceneManager, dataManager *data.DataManager, textRenderer *graphics.TextRenderer) *BattleSceneUnified {
	// Create camera for 5000x5000 world with 1024x768 viewport
	camera := graphics.NewCameraManager(5000, 5000, 1024, 768)
	
	// Disable smooth movement for immediate response
	camera.SetSmoothMove(false)
	
	// Create scroll controller
	scrollController := input.NewScrollController(camera)
	
	fmt.Println("BattleSceneUnified: Camera and ScrollController initialized")
	
	return &BattleSceneUnified{
		sceneManager:     sceneManager,
		dataManager:      dataManager,
		textRenderer:     textRenderer,
		spriteGenerator:  graphics.NewSpriteGenerator(),
		camera:           camera,
		scrollController: scrollController,
		minimap:          graphics.NewMinimap(camera, 50, 620, 200, 150),
		isPaused:         false,
		showDebugInfo:    false,
		showHelp:         false,
		lastUpdate:       time.Now(),
	}
}

// OnEnter is called when entering the scene
func (bs *BattleSceneUnified) OnEnter(data interface{}) {
	bs.Initialize()
}

// OnExit is called when exiting the scene
func (bs *BattleSceneUnified) OnExit() {
	bs.battleManager = nil
}

// Initialize initializes the battle scene
func (bs *BattleSceneUnified) Initialize() {
	if bs.battleManager == nil {
		fmt.Println("=== Battle Scene Initialize ===")
		
		// Get stage and preset from scene manager's game data
		stageName := bs.sceneManager.gameData.CurrentStage
		presetName := bs.sceneManager.gameData.CurrentPreset
		
		if stageName == "" {
			stageName = "森の戦い" // Default
		}
		if presetName == "" {
			presetName = "バランス型" // Default
		}
		
		fmt.Printf("Selected Stage: %s\n", stageName)
		fmt.Printf("Selected Preset: %s\n", presetName)
		
		// Map stage names to config names
		stageConfigMap := map[string]string{
			"森の戦い": "forest_battle",
			"山岳要塞": "mountain_fortress", 
			"平原決戦": "plain_battle",
		}
		
		terrainConfigMap := map[string]string{
			"森の戦い": "forest",
			"山岳要塞": "mountain",
			"平原決戦": "plain",
		}
		
		stageConfigName := stageConfigMap[stageName]
		terrainConfigName := terrainConfigMap[stageName]
		
		if stageConfigName == "" {
			fmt.Printf("Warning: Unknown stage name '%s', using default\n", stageName)
			stageConfigName = "forest_battle" // Default
		}
		if terrainConfigName == "" {
			fmt.Printf("Warning: Unknown terrain name for stage '%s', using default\n", stageName)
			terrainConfigName = "forest" // Default
		}
		
		fmt.Printf("Looking for stage config: %s\n", stageConfigName)
		fmt.Printf("Looking for terrain config: %s\n", terrainConfigName)
		
		// Debug: List all available stages
		fmt.Println("Available stages in data manager:")
		// This would require adding a method to list all stages, but for now let's try the configs directly
		
		// Set up stage
		stageConfig, err := bs.dataManager.GetStageConfig(stageConfigName)
		if err != nil {
			fmt.Printf("Error loading stage config '%s': %v\n", stageConfigName, err)
			fmt.Println("Falling back to forest_battle")
			stageConfig, err = bs.dataManager.GetStageConfig("forest_battle")
			if err != nil {
				fmt.Printf("Error loading fallback stage config: %v\n", err)
				return
			}
		}
		fmt.Printf("Stage loaded: %s\n", stageConfig.Name)
		
		terrainConfig, err := bs.dataManager.GetTerrainConfig(terrainConfigName)
		if err != nil {
			fmt.Printf("Error loading terrain config '%s': %v\n", terrainConfigName, err)
			fmt.Println("Falling back to forest terrain")
			terrainConfig, err = bs.dataManager.GetTerrainConfig("forest")
			if err != nil {
				fmt.Printf("Error loading fallback terrain config: %v\n", err)
				return
			}
		}
		fmt.Printf("Terrain loaded: %s\n", terrainConfig.Name)
		
		// Create battle manager with stage and terrain
		bs.battleManager = game.NewBattleManager(stageConfig, terrainConfig)
		if bs.battleManager == nil {
			fmt.Println("Error: Failed to create battle manager")
			return
		}
		fmt.Println("Battle manager created successfully")
		
		// Create armies with selected preset
		fmt.Printf("Creating armies with preset: %s\n", presetName)
		err1 := bs.battleManager.CreatePresetArmy(0, presetName, bs.dataManager)
		if err1 != nil {
			fmt.Printf("Error creating army A: %v\n", err1)
		}
		
		err2 := bs.battleManager.CreatePresetArmy(1, presetName, bs.dataManager)
		if err2 != nil {
			fmt.Printf("Error creating army B: %v\n", err2)
		}
		
		if err1 != nil || err2 != nil {
			fmt.Printf("Army creation had errors, but continuing...\n")
		}
		
		// Verify armies were created
		armyAUnits := bs.battleManager.ArmyA.GetAllUnits()
		armyBUnits := bs.battleManager.ArmyB.GetAllUnits()
		fmt.Printf("Army A has %d units, Army B has %d units\n", len(armyAUnits), len(armyBUnits))
		
		if len(armyAUnits) == 0 || len(armyBUnits) == 0 {
			fmt.Println("Warning: One or both armies have no units!")
		}
		
		// Start battle
		bs.battleManager.StartBattle()
		fmt.Println("Battle started!")
		
		// Center camera on battlefield
		bs.camera.SetPosition(2500, 2500) // Center of 5000x5000 world
	}
}

// Update updates the battle scene
func (bs *BattleSceneUnified) Update() error {
	// Calculate delta time
	now := time.Now()
	if !bs.lastUpdate.IsZero() {
		bs.deltaTime = now.Sub(bs.lastUpdate).Seconds()
	}
	bs.lastUpdate = now
	
	// Update camera first
	if bs.camera != nil {
		bs.camera.Update(bs.deltaTime)
	}
	
	// Update scroll controller (after camera update)
	if bs.scrollController != nil {
		bs.scrollController.Update(bs.deltaTime)
	}
	
	// Handle input
	bs.handleInput()
	
	// Update battle if not paused
	if !bs.isPaused && bs.battleManager != nil {
		bs.battleManager.Update(bs.deltaTime)
		
		// Check if battle ended
		if !bs.battleManager.IsActive {
			winner := bs.battleManager.GetWinnerName()
			bs.sceneManager.TransitionTo(SceneResult, winner)
			return nil
		}
	}
	
	return nil
}

// handleInput handles user input
func (bs *BattleSceneUnified) handleInput() {
	// Handle return to setup (works even if battleManager is nil)
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		bs.sceneManager.TransitionTo(SceneArmySetup, nil)
		return
	}
	
	// Handle force reinitialize (F5 key)
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		fmt.Println("Force reinitializing battle scene...")
		bs.battleManager = nil
		bs.Initialize()
		return
	}
	
	// Direct camera control test (temporary)
	if bs.camera != nil {
		moveSpeed := 200.0 * bs.deltaTime
		
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			fmt.Println("Direct camera move: UP")
			bs.camera.Move(0, -moveSpeed)
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			fmt.Println("Direct camera move: DOWN")
			bs.camera.Move(0, moveSpeed)
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			fmt.Println("Direct camera move: LEFT")
			bs.camera.Move(-moveSpeed, 0)
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			fmt.Println("Direct camera move: RIGHT")
			bs.camera.Move(moveSpeed, 0)
		}
		
		// Direct zoom test
		_, wheelY := ebiten.Wheel()
		if wheelY != 0 {
			fmt.Printf("Direct zoom: wheelY=%.2f\n", wheelY)
			mouseX, mouseY := ebiten.CursorPosition()
			bs.camera.ZoomAt(mouseX, mouseY, wheelY*0.25)
		}
	}
	
	// Other input handling only if battleManager exists
	if bs.battleManager == nil {
		return
	}
	
	// Handle pause (but not Escape if it's used for camera)
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		bs.isPaused = !bs.isPaused
	}
	
	// Handle pause with Escape only if not used for camera movement
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		bs.isPaused = !bs.isPaused
	}
	
	// Handle debug info toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		bs.showDebugInfo = !bs.showDebugInfo
	}
	
	// Handle help toggle
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		now := time.Now()
		if now.Sub(bs.helpToggleTime) > 200*time.Millisecond {
			bs.showHelp = !bs.showHelp
			bs.helpToggleTime = now
		}
	}
	
	// Handle unit selection (only left mouse button, middle button is for camera drag)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		bs.handleUnitSelection()
	}
}

// handleUnitSelection handles unit selection with mouse
func (bs *BattleSceneUnified) handleUnitSelection() {
	if bs.battleManager == nil {
		return
	}
	
	// Get mouse position
	mouseX, mouseY := ebiten.CursorPosition()
	
	// Convert screen coordinates to world coordinates
	worldX, worldY := bs.camera.ScreenToWorld(mouseX, mouseY)
	
	// Find unit at position
	bs.selectedUnit = nil
	
	// Check Army A units
	for _, unit := range bs.battleManager.ArmyA.GetAllUnits() {
		if unit.IsAlive && bs.isUnitAtPosition(unit, worldX, worldY) {
			bs.selectedUnit = unit
			return
		}
	}
	
	// Check Army B units
	for _, unit := range bs.battleManager.ArmyB.GetAllUnits() {
		if unit.IsAlive && bs.isUnitAtPosition(unit, worldX, worldY) {
			bs.selectedUnit = unit
			return
		}
	}
}

// isUnitAtPosition checks if a unit is at the given world position
func (bs *BattleSceneUnified) isUnitAtPosition(unit *game.Unit, worldX, worldY float64) bool {
	size := 16.0 // Default unit size
	
	return math.Abs(unit.Position.X-worldX) < size && 
		   math.Abs(unit.Position.Y-worldY) < size
}

// Draw draws the battle scene
func (bs *BattleSceneUnified) Draw(screen *ebiten.Image) {
	if bs.battleManager == nil {
		// Show loading message with more details
		screen.Fill(color.RGBA{44, 62, 80, 255})
		bs.textRenderer.DrawCenteredText(screen, "戦闘準備中...", 512, 300, color.RGBA{236, 240, 241, 255})
		
		// Show selected stage and preset
		if bs.sceneManager.gameData.CurrentStage != "" {
			stageText := fmt.Sprintf("ステージ: %s", bs.sceneManager.gameData.CurrentStage)
			bs.textRenderer.DrawCenteredText(screen, stageText, 512, 350, color.RGBA{149, 165, 166, 255})
		}
		
		if bs.sceneManager.gameData.CurrentPreset != "" {
			presetText := fmt.Sprintf("編成: %s", bs.sceneManager.gameData.CurrentPreset)
			bs.textRenderer.DrawCenteredText(screen, presetText, 512, 380, color.RGBA{149, 165, 166, 255})
		}
		
		// Show hint to return
		bs.textRenderer.DrawCenteredText(screen, "Rキーで設定に戻る  F5キーで再初期化", 512, 450, color.RGBA{149, 165, 166, 255})
		return
	}
	
	// Clear screen
	screen.Fill(color.RGBA{20, 40, 20, 255}) // Dark green background
	
	// Get camera transform
	transform := bs.camera.GetTransform()
	
	// Draw battlefield
	bs.drawBattlefield(screen, transform)
	
	// Draw units
	bs.drawUnits(screen, transform)
	
	// Draw selected unit range
	if bs.selectedUnit != nil && bs.selectedUnit.IsAlive {
		bs.drawUnitRange(screen, transform)
	}
	
	// Draw UI (not affected by camera transform)
	bs.drawStatusBar(screen)
	bs.drawUI(screen)
	
	// Draw overlays
	if bs.showDebugInfo {
		bs.drawDebugInfo(screen)
	}
	
	if bs.showHelp {
		bs.drawHelp(screen)
	}
	
	if bs.isPaused {
		bs.drawPauseOverlay(screen)
	}
}

// drawBattlefield draws the battlefield background
func (bs *BattleSceneUnified) drawBattlefield(screen *ebiten.Image, transform ebiten.GeoM) {
	// Draw terrain-based background
	var bgColor color.RGBA
	
	switch bs.battleManager.TerrainData.Name {
	case "森":
		bgColor = color.RGBA{34, 139, 34, 255} // Forest green
	case "山":
		bgColor = color.RGBA{139, 69, 19, 255} // Saddle brown
	case "平原":
		bgColor = color.RGBA{124, 252, 0, 255} // Lawn green
	case "城塞":
		bgColor = color.RGBA{105, 105, 105, 255} // Dim gray
	case "街":
		bgColor = color.RGBA{160, 82, 45, 255} // Saddle brown
	default:
		bgColor = color.RGBA{34, 139, 34, 255} // Default green
	}
	
	// Create a large background image
	bg := ebiten.NewImage(5000, 5000)
	bg.Fill(bgColor)
	
	// Draw with camera transform
	op := &ebiten.DrawImageOptions{}
	op.GeoM = transform
	screen.DrawImage(bg, op)
	
	// Draw grid pattern for reference
	bs.drawGrid(screen, transform)
}

// drawGrid draws a reference grid
func (bs *BattleSceneUnified) drawGrid(screen *ebiten.Image, transform ebiten.GeoM) {
	gridSize := 100
	gridColor := color.RGBA{255, 255, 255, 32} // Very transparent white
	
	// Draw vertical lines
	for x := 0; x < 5000; x += gridSize {
		line := ebiten.NewImage(1, 5000)
		line.Fill(gridColor)
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), 0)
		op.GeoM.Concat(transform)
		screen.DrawImage(line, op)
	}
	
	// Draw horizontal lines
	for y := 0; y < 5000; y += gridSize {
		line := ebiten.NewImage(5000, 1)
		line.Fill(gridColor)
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(y))
		op.GeoM.Concat(transform)
		screen.DrawImage(line, op)
	}
}

// drawUnits draws all units
func (bs *BattleSceneUnified) drawUnits(screen *ebiten.Image, transform ebiten.GeoM) {
	// Draw Army A units (red)
	for _, unit := range bs.battleManager.ArmyA.GetAllUnits() {
		if unit.IsAlive {
			bs.drawUnit(screen, unit, transform, color.RGBA{231, 76, 60, 255})
		}
	}
	
	// Draw Army B units (blue)
	for _, unit := range bs.battleManager.ArmyB.GetAllUnits() {
		if unit.IsAlive {
			bs.drawUnit(screen, unit, transform, color.RGBA{41, 128, 185, 255})
		}
	}
}

// drawUnit draws a single unit
func (bs *BattleSceneUnified) drawUnit(screen *ebiten.Image, unit *game.Unit, transform ebiten.GeoM, baseColor color.RGBA) {
	// Determine unit color
	unitColor := baseColor
	
	// Highlight selected unit
	if bs.selectedUnit == unit {
		unitColor = color.RGBA{255, 255, 0, 255} // Yellow
	} else {
		// Adjust color based on health
		healthPercent := unit.GetHealthPercentage()
		if healthPercent < 0.5 {
			factor := 0.5 + healthPercent
			unitColor.R = uint8(float64(unitColor.R) * factor)
			unitColor.G = uint8(float64(unitColor.G) * factor)
			unitColor.B = uint8(float64(unitColor.B) * factor)
		}
	}
	
	// Generate unit sprite
	sprite := bs.spriteGenerator.GenerateUnitSprite(string(unit.Type), unitColor, unit.IsLeader, unit.Animation)
	
	// Draw unit
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(unit.Position.X-8, unit.Position.Y-8) // Center the sprite
	op.GeoM.Concat(transform)
	screen.DrawImage(sprite, op)
	
	// Draw health bar
	bs.drawHealthBar(screen, unit, transform)
}

// drawHealthBar draws a unit's health bar
func (bs *BattleSceneUnified) drawHealthBar(screen *ebiten.Image, unit *game.Unit, transform ebiten.GeoM) {
	size := 16.0
	barWidth := int(size)
	barHeight := 3
	
	// Create health bar background
	bgBar := ebiten.NewImage(barWidth, barHeight)
	bgBar.Fill(color.RGBA{100, 100, 100, 255})
	
	// Create health bar fill
	healthPercent := unit.GetHealthPercentage()
	fillWidth := int(float64(barWidth) * healthPercent)
	if fillWidth > 0 {
		fillBar := ebiten.NewImage(fillWidth, barHeight)
		
		// Color based on health
		var fillColor color.RGBA
		if healthPercent > 0.6 {
			fillColor = color.RGBA{0, 255, 0, 255} // Green
		} else if healthPercent > 0.3 {
			fillColor = color.RGBA{255, 255, 0, 255} // Yellow
		} else {
			fillColor = color.RGBA{255, 0, 0, 255} // Red
		}
		fillBar.Fill(fillColor)
		
		// Draw fill bar
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(unit.Position.X-size/2, unit.Position.Y-size/2-8)
		op.GeoM.Concat(transform)
		screen.DrawImage(fillBar, op)
	}
	
	// Draw background bar
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(unit.Position.X-size/2, unit.Position.Y-size/2-8)
	op.GeoM.Concat(transform)
	screen.DrawImage(bgBar, op)
}

// drawUnitRange draws the selected unit's attack range
func (bs *BattleSceneUnified) drawUnitRange(screen *ebiten.Image, transform ebiten.GeoM) {
	if bs.selectedUnit == nil {
		return
	}
	
	attackRange := bs.selectedUnit.Range
	radius := int(attackRange)
	
	// Create range circle
	rangeImg := ebiten.NewImage(radius*2, radius*2)
	rangeColor := color.RGBA{255, 255, 255, 64} // Semi-transparent white
	
	// Draw circle outline
	for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
		x := int(float64(radius) + float64(radius-2)*math.Cos(angle))
		y := int(float64(radius) + float64(radius-2)*math.Sin(angle))
		if x >= 0 && x < radius*2 && y >= 0 && y < radius*2 {
			rangeImg.Set(x, y, rangeColor)
		}
	}
	
	// Draw range indicator
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bs.selectedUnit.Position.X-float64(radius), bs.selectedUnit.Position.Y-float64(radius))
	op.GeoM.Concat(transform)
	screen.DrawImage(rangeImg, op)
}

// drawStatusBar draws the top status bar
func (bs *BattleSceneUnified) drawStatusBar(screen *ebiten.Image) {
	// Background for status bar
	statusBarHeight := 60
	statusBar := ebiten.NewImage(1024, statusBarHeight)
	statusBar.Fill(color.RGBA{52, 73, 94, 255}) // #34495E
	screen.DrawImage(statusBar, nil)
	
	// Time display
	remainingTime := bs.battleManager.TimeLimit - bs.battleManager.BattleTime
	minutes := int(remainingTime) / 60
	seconds := int(remainingTime) % 60
	timeText := fmt.Sprintf("時間: %02d:%02d", minutes, seconds)
	bs.textRenderer.DrawText(screen, timeText, 20, 20, color.RGBA{236, 240, 241, 255})
	
	// Stage name
	stageText := bs.battleManager.Stage.Name + " (" + bs.battleManager.TerrainData.Name + ")"
	bs.textRenderer.DrawText(screen, stageText, 200, 20, color.RGBA{236, 240, 241, 255})
	
	// Army A info
	armyAText := "軍勢A"
	bs.textRenderer.DrawText(screen, armyAText, 500, 20, color.RGBA{236, 240, 241, 255})
	bs.drawArmyHealthBar(screen, 580, 25, bs.battleManager.ArmyA.GetTotalHealth(), color.RGBA{231, 76, 60, 255})
	
	// Army B info
	armyBText := "軍勢B"
	bs.textRenderer.DrawText(screen, armyBText, 750, 20, color.RGBA{236, 240, 241, 255})
	bs.drawArmyHealthBar(screen, 830, 25, bs.battleManager.ArmyB.GetTotalHealth(), color.RGBA{41, 128, 185, 255})
	
	// Unit counts
	armyACount := len(bs.battleManager.ArmyA.GetAllUnits())
	armyBCount := len(bs.battleManager.ArmyB.GetAllUnits())
	countText := fmt.Sprintf("ユニット数 A:%d B:%d", armyACount, armyBCount)
	bs.textRenderer.DrawText(screen, countText, 200, 40, color.RGBA{255, 255, 0, 255})
}

// drawArmyHealthBar draws an army's total health bar
func (bs *BattleSceneUnified) drawArmyHealthBar(screen *ebiten.Image, x, y int, health float64, barColor color.Color) {
	barWidth := 120
	barHeight := 15
	
	// Background
	bgBar := ebiten.NewImage(barWidth, barHeight)
	bgBar.Fill(color.RGBA{100, 100, 100, 255})
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(bgBar, op)
	
	// Health fill
	filledWidth := int(float64(barWidth) * health)
	if filledWidth > 0 {
		fillBar := ebiten.NewImage(filledWidth, barHeight)
		fillBar.Fill(barColor)
		
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(fillBar, op)
	}
	
	// Border
	border := ebiten.NewImage(barWidth, 1)
	border.Fill(color.RGBA{255, 255, 255, 255})
	
	// Top and bottom borders
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(border, op1)
	
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(float64(x), float64(y+barHeight-1))
	screen.DrawImage(border, op2)
	
	// Side borders
	sideBorder := ebiten.NewImage(1, barHeight)
	sideBorder.Fill(color.RGBA{255, 255, 255, 255})
	
	op3 := &ebiten.DrawImageOptions{}
	op3.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(sideBorder, op3)
	
	op4 := &ebiten.DrawImageOptions{}
	op4.GeoM.Translate(float64(x+barWidth-1), float64(y))
	screen.DrawImage(sideBorder, op4)
}

// drawUI draws the user interface
func (bs *BattleSceneUnified) drawUI(screen *ebiten.Image) {
	// Draw minimap
	if bs.minimap != nil {
		bs.minimap.Draw(screen)
	}
	
	// Draw selected unit info
	if bs.selectedUnit != nil && bs.selectedUnit.IsAlive {
		bs.drawSelectedUnitInfo(screen)
	}
	
	// Draw controls
	controlsText := "P/Esc: 一時停止  R: 設定に戻る  F1: デバッグ  F2: ヘルプ"
	bs.textRenderer.DrawText(screen, controlsText, 300, 740, color.RGBA{255, 255, 255, 255})
}

// drawSelectedUnitInfo draws information about the selected unit
func (bs *BattleSceneUnified) drawSelectedUnitInfo(screen *ebiten.Image) {
	unit := bs.selectedUnit
	if unit == nil || !unit.IsAlive {
		return
	}
	
	// Background
	infoX := 300
	infoY := 620
	infoWidth := 300
	infoHeight := 100
	
	infoBg := ebiten.NewImage(infoWidth, infoHeight)
	infoBg.Fill(color.RGBA{52, 73, 94, 200}) // Semi-transparent
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(infoX), float64(infoY))
	screen.DrawImage(infoBg, op)
	
	// Unit info
	y := infoY + 10
	bs.textRenderer.DrawText(screen, "選択ユニット:", float64(infoX+10), float64(y), color.RGBA{236, 240, 241, 255})
	y += 20
	
	unitTypeText := fmt.Sprintf("種別: %s", unit.Type)
	if unit.IsLeader {
		unitTypeText += " (リーダー)"
	}
	bs.textRenderer.DrawText(screen, unitTypeText, float64(infoX+10), float64(y), color.RGBA{236, 240, 241, 255})
	y += 15
	
	healthText := fmt.Sprintf("HP: %d/%d", unit.HP, unit.MaxHP)
	bs.textRenderer.DrawText(screen, healthText, float64(infoX+10), float64(y), color.RGBA{236, 240, 241, 255})
	y += 15
	
	attackText := fmt.Sprintf("攻撃力: %d  射程: %.0f", unit.AttackPower, unit.Range)
	bs.textRenderer.DrawText(screen, attackText, float64(infoX+10), float64(y), color.RGBA{236, 240, 241, 255})
}

// drawDebugInfo draws debug information
func (bs *BattleSceneUnified) drawDebugInfo(screen *ebiten.Image) {
	camX, camY := bs.camera.GetPosition()
	zoom := bs.camera.GetZoom()
	
	debugText := fmt.Sprintf("Camera: (%.0f, %.0f) Zoom: %.2f", camX, camY, zoom)
	bs.textRenderer.DrawText(screen, debugText, 10, 80, color.RGBA{255, 255, 0, 255})
	
	// Show mouse position for debugging
	mouseX, mouseY := ebiten.CursorPosition()
	worldX, worldY := bs.camera.ScreenToWorld(mouseX, mouseY)
	mouseText := fmt.Sprintf("Mouse: Screen(%d, %d) World(%.0f, %.0f)", mouseX, mouseY, worldX, worldY)
	bs.textRenderer.DrawText(screen, mouseText, 10, 100, color.RGBA{255, 255, 0, 255})
	
	if bs.selectedUnit != nil {
		unitDebug := fmt.Sprintf("Selected: %s at (%.0f, %.0f)", 
			bs.selectedUnit.Type, bs.selectedUnit.Position.X, bs.selectedUnit.Position.Y)
		bs.textRenderer.DrawText(screen, unitDebug, 10, 120, color.RGBA{255, 255, 0, 255})
	}
	
	fpsText := fmt.Sprintf("FPS: %.1f", 1.0/bs.deltaTime)
	bs.textRenderer.DrawText(screen, fpsText, 10, 140, color.RGBA{255, 255, 0, 255})
	
	// Show scroll controller status
	if bs.scrollController != nil {
		scrollText := fmt.Sprintf("Scroll: Edge=%t Key=%t Drag=%t", 
			bs.scrollController.EdgeScrolling, bs.scrollController.KeyScrolling, bs.scrollController.DragScrolling)
		bs.textRenderer.DrawText(screen, scrollText, 10, 160, color.RGBA{255, 255, 0, 255})
	}
}

// drawHelp draws help information
func (bs *BattleSceneUnified) drawHelp(screen *ebiten.Image) {
	// Semi-transparent background
	helpBg := ebiten.NewImage(400, 300)
	helpBg.Fill(color.RGBA{0, 0, 0, 200})
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(312, 234) // Center on screen
	screen.DrawImage(helpBg, op)
	
	// Help text
	helpLines := []string{
		"=== 操作方法 ===",
		"",
		"マウス: ユニット選択",
		"WASD/矢印キー: カメラ移動",
		"マウスホイール: ズーム",
		"中ボタンドラッグ: カメラドラッグ",
		"画面端: エッジスクロール",
		"+/-キー: ズームイン/アウト",
		"P: 一時停止",
		"R: 設定画面に戻る",
		"F1: デバッグ情報表示",
		"F2: このヘルプ表示",
		"F5: 戦闘再初期化",
		"",
		"=== ユニット記号 ===",
		"□: 歩兵  △: 弓兵  ◇: 魔術師",
		"",
		"F2でヘルプを閉じる",
	}
	
	y := 250
	for _, line := range helpLines {
		bs.textRenderer.DrawText(screen, line, 330, float64(y), color.RGBA{255, 255, 255, 255})
		y += 18
	}
}

// drawPauseOverlay draws the pause overlay
func (bs *BattleSceneUnified) drawPauseOverlay(screen *ebiten.Image) {
	// Semi-transparent overlay
	overlay := ebiten.NewImage(1024, 768)
	overlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(overlay, nil)
	
	// Pause text
	bs.textRenderer.DrawCenteredText(screen, "一時停止", 512, 350, color.RGBA{255, 255, 255, 255})
	bs.textRenderer.DrawCenteredText(screen, "P/Escで再開", 512, 400, color.RGBA{255, 255, 255, 255})
}
