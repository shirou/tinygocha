package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/shirou/tinygocha/internal/graphics"
)

// ResultScene represents the battle result screen
type ResultScene struct {
	sceneManager *SceneManager
	textRenderer *graphics.TextRenderer
	winner       string
	selectedItem int
	menuItems    []string
}

// NewResultScene creates a new result scene
func NewResultScene(sceneManager *SceneManager, textRenderer *graphics.TextRenderer) *ResultScene {
	return &ResultScene{
		sceneManager: sceneManager,
		textRenderer: textRenderer,
		selectedItem: 0,
		menuItems:    []string{"再戦", "軍勢変更", "タイトル"},
	}
}

// Update updates the result scene
func (rs *ResultScene) Update() error {
	// Handle input
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		rs.selectedItem--
		if rs.selectedItem < 0 {
			rs.selectedItem = len(rs.menuItems) - 1
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		rs.selectedItem++
		if rs.selectedItem >= len(rs.menuItems) {
			rs.selectedItem = 0
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch rs.selectedItem {
		case 0: // 再戦
			rs.sceneManager.TransitionTo(SceneBattle, nil)
		case 1: // 軍勢変更
			rs.sceneManager.TransitionTo(SceneArmySetup, nil)
		case 2: // タイトル
			rs.sceneManager.TransitionTo(SceneTitle, nil)
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		rs.sceneManager.TransitionTo(SceneTitle, nil)
	}
	
	return nil
}

// Draw draws the result scene
func (rs *ResultScene) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{44, 62, 80, 255}) // #2C3E50
	
	// Draw winner announcement
	winnerText := fmt.Sprintf("%s 勝利！", rs.winner)
	if rs.winner == "引き分け" {
		winnerText = "引き分け！"
	}
	rs.textRenderer.DrawTextWithSize(screen, winnerText, 400, 150, color.RGBA{236, 240, 241, 255}, 32)
	
	// Draw battle statistics
	rs.drawStatistics(screen)
	
	// Draw menu items
	for i, item := range rs.menuItems {
		x := 350.0 + float64(i*100)
		y := 500.0
		
		// Highlight selected item
		if i == rs.selectedItem {
			rs.textRenderer.DrawTextWithShadow(screen, "> "+item+" <", x-20, y, 
				color.RGBA{52, 152, 219, 255}, color.RGBA{0, 0, 0, 128})
		} else {
			rs.textRenderer.DrawText(screen, item, x, y, color.RGBA{236, 240, 241, 255})
		}
	}
	
	// Draw controls hint
	controlsText := "↑↓: 選択  Enter: 決定  Esc: タイトル"
	rs.textRenderer.DrawText(screen, controlsText, 350, 600, color.RGBA{149, 165, 166, 255})
}

// drawStatistics draws battle statistics
func (rs *ResultScene) drawStatistics(screen *ebiten.Image) {
	// Statistics panel background
	panelX := 200
	panelY := 250
	panelWidth := 600
	panelHeight := 200
	
	// Draw panel background
	for dy := 0; dy < panelHeight; dy++ {
		for dx := 0; dx < panelWidth; dx++ {
			screen.Set(panelX+dx, panelY+dy, color.RGBA{52, 73, 94, 255}) // #34495E
		}
	}
	
	// Draw panel border
	borderColor := color.RGBA{236, 240, 241, 255} // #ECF0F1
	for dx := 0; dx < panelWidth; dx++ {
		screen.Set(panelX+dx, panelY, borderColor)
		screen.Set(panelX+dx, panelY+panelHeight-1, borderColor)
	}
	for dy := 0; dy < panelHeight; dy++ {
		screen.Set(panelX, panelY+dy, borderColor)
		screen.Set(panelX+panelWidth-1, panelY+dy, borderColor)
	}
	
	// Battle statistics (placeholder data)
	statsTitle := "戦闘統計"
	rs.textRenderer.DrawTextWithSize(screen, statsTitle, float64(panelX+20), float64(panelY+20), color.RGBA{236, 240, 241, 255}, 20)
	
	// Left column - General stats
	rs.textRenderer.DrawText(screen, "戦闘時間: 3:45", float64(panelX+20), float64(panelY+50), color.RGBA{236, 240, 241, 255})
	rs.textRenderer.DrawText(screen, "軍勢A生存: 8", float64(panelX+20), float64(panelY+70), color.RGBA{236, 240, 241, 255})
	rs.textRenderer.DrawText(screen, "軍勢B生存: 2", float64(panelX+20), float64(panelY+90), color.RGBA{236, 240, 241, 255})
	rs.textRenderer.DrawText(screen, "総ダメージ", float64(panelX+20), float64(panelY+110), color.RGBA{236, 240, 241, 255})
	rs.textRenderer.DrawText(screen, "A: 1200  B: 800", float64(panelX+20), float64(panelY+130), color.RGBA{236, 240, 241, 255})
	
	// Right column - MVP
	mvpTitle := "MVP"
	rs.textRenderer.DrawTextWithSize(screen, mvpTitle, float64(panelX+350), float64(panelY+50), color.RGBA{236, 240, 241, 255}, 18)
	rs.textRenderer.DrawText(screen, "弓兵リーダー", float64(panelX+350), float64(panelY+70), color.RGBA{236, 240, 241, 255})
	rs.textRenderer.DrawText(screen, "撃破数: 5", float64(panelX+350), float64(panelY+90), color.RGBA{236, 240, 241, 255})
	rs.textRenderer.DrawText(screen, "与ダメージ: 450", float64(panelX+350), float64(panelY+110), color.RGBA{236, 240, 241, 255})
}

// OnEnter is called when entering this scene
func (rs *ResultScene) OnEnter(data interface{}) {
	// Set winner from data
	if winner, ok := data.(string); ok {
		rs.winner = winner
	}
	rs.selectedItem = 0
}

// OnExit is called when exiting this scene
func (rs *ResultScene) OnExit() {
	// Nothing to clean up
}
