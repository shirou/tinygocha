package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/shirou/tinygocha/internal/graphics"
)

// TitleScene represents the title screen
type TitleScene struct {
	sceneManager *SceneManager
	textRenderer *graphics.TextRenderer
	selectedItem int
	menuItems    []string
}

// NewTitleScene creates a new title scene
func NewTitleScene(sceneManager *SceneManager, textRenderer *graphics.TextRenderer) *TitleScene {
	return &TitleScene{
		sceneManager: sceneManager,
		textRenderer: textRenderer,
		selectedItem: 0,
		menuItems:    []string{"戦闘開始", "終了"},
	}
}

// Update updates the title scene
func (ts *TitleScene) Update() error {
	// Handle input
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		ts.selectedItem--
		if ts.selectedItem < 0 {
			ts.selectedItem = len(ts.menuItems) - 1
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		ts.selectedItem++
		if ts.selectedItem >= len(ts.menuItems) {
			ts.selectedItem = 0
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch ts.selectedItem {
		case 0: // 戦闘開始
			ts.sceneManager.TransitionTo(SceneArmySetup, nil)
		case 1: // 終了
			return ebiten.Termination
		}
	}
	
	return nil
}

// Draw draws the title scene
func (ts *TitleScene) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{44, 62, 80, 255}) // #2C3E50
	
	// Draw title
	titleText := "ゴチャキャラバトル"
	ts.textRenderer.DrawTextWithSize(screen, titleText, 320, 200, color.RGBA{236, 240, 241, 255}, 32)
	
	// Draw version
	versionText := "Version 0.1.0 (Demo)"
	ts.textRenderer.DrawText(screen, versionText, 400, 250, color.RGBA{149, 165, 166, 255})
	
	// Draw menu items
	for i, item := range ts.menuItems {
		x := 450.0
		y := 350.0 + float64(i*50)
		
		// Highlight selected item
		if i == ts.selectedItem {
			// Draw selection indicator with shadow
			selectedText := "> " + item + " <"
			ts.textRenderer.DrawTextWithShadow(screen, selectedText, x-20, y, 
				color.RGBA{52, 152, 219, 255}, color.RGBA{0, 0, 0, 128})
		} else {
			ts.textRenderer.DrawText(screen, item, x, y, color.RGBA{236, 240, 241, 255})
		}
	}
	
	// Draw controls hint
	controlsText := "↑↓: 選択  Enter/Space: 決定"
	ts.textRenderer.DrawText(screen, controlsText, 350, 500, color.RGBA{149, 165, 166, 255})
}

// OnEnter is called when entering this scene
func (ts *TitleScene) OnEnter(data interface{}) {
	// Reset selection
	ts.selectedItem = 0
}

// OnExit is called when exiting this scene
func (ts *TitleScene) OnExit() {
	// Nothing to clean up
}
