package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/shirou/tinygocha/internal/config"
	"github.com/shirou/tinygocha/internal/data"
	"github.com/shirou/tinygocha/internal/graphics"
	"github.com/shirou/tinygocha/internal/scenes"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

// Game represents the main game structure
type Game struct {
	sceneManager   *scenes.SceneManager
	dataManager    *data.DataManager
	config         *config.Config
	fontManager    *graphics.FontManager
	textRenderer   *graphics.TextRenderer
}

// NewGame creates a new game instance
func NewGame() *Game {
	// Load configuration
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Printf("Warning: Failed to load config: %v, using defaults", err)
		cfg = config.DefaultConfig()
	}
	
	// Create font manager and load fonts
	fontManager := graphics.NewFontManager()
	fontSize := float64(cfg.Graphics.FontSize)
	
	if cfg.Graphics.FontPath != "" {
		// Load custom font
		if err := fontManager.LoadFontFromFile(cfg.Graphics.FontPath, fontSize, "default"); err != nil {
			log.Printf("Warning: Failed to load custom font, using default: %v", err)
		}
	} else {
		// Load default MPlus1p font
		if err := fontManager.LoadDefaultFont(fontSize); err != nil {
			log.Printf("Error: Failed to load default font: %v", err)
		}
	}
	
	// Create text renderer
	textRenderer := graphics.NewTextRenderer(fontManager)
	
	// Create data manager and load all data
	dataManager := data.NewDataManager()
	if err := dataManager.LoadAll(); err != nil {
		log.Printf("Warning: Failed to load data files: %v", err)
		// Continue with default/empty data
	}
	
	sceneManager := scenes.NewSceneManager()
	
	// Register all scenes with text renderer
	sceneManager.RegisterScene(scenes.SceneTitle, scenes.NewTitleScene(sceneManager, textRenderer))
	sceneManager.RegisterScene(scenes.SceneArmySetup, scenes.NewArmySetupScene(sceneManager, textRenderer))
	sceneManager.RegisterScene(scenes.SceneBattle, scenes.NewBattleSceneUnified(sceneManager, dataManager, textRenderer))
	sceneManager.RegisterScene(scenes.SceneResult, scenes.NewResultScene(sceneManager, textRenderer))
	
	return &Game{
		sceneManager: sceneManager,
		dataManager:  dataManager,
		config:       cfg,
		fontManager:  fontManager,
		textRenderer: textRenderer,
	}
}

// Update updates the game logic
func (g *Game) Update() error {
	return g.sceneManager.Update()
}

// Draw draws the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
	
	// Draw FPS if enabled
	if g.config.Graphics.ShowFPS {
		fpsText := "FPS: " + fmt.Sprintf("%.1f", ebiten.ActualFPS())
		g.textRenderer.DrawText(screen, fpsText, 10, 10, color.RGBA{255, 255, 255, 255})
	}
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Set window properties
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ゴチャキャラバトル - Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	// Create and run the game
	game := NewGame()
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
