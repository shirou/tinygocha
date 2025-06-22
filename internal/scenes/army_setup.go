package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/shirou/tinygocha/internal/graphics"
)

// ArmySetupScene represents the army setup screen
type ArmySetupScene struct {
	sceneManager     *SceneManager
	textRenderer     *graphics.TextRenderer
	selectedItem     int
	presetArmies     []string
	selectedPreset   int
	selectedStage    int
	stages           []string
}

// NewArmySetupScene creates a new army setup scene
func NewArmySetupScene(sceneManager *SceneManager, textRenderer *graphics.TextRenderer) *ArmySetupScene {
	return &ArmySetupScene{
		sceneManager:   sceneManager,
		textRenderer:   textRenderer,
		selectedItem:   0,
		presetArmies:   []string{"バランス型", "攻撃重視", "防御重視"},
		selectedPreset: 0,
		selectedStage:  0,
		stages:         []string{"森の戦い", "山岳要塞", "平原決戦"},
	}
}

// Update updates the army setup scene
func (as *ArmySetupScene) Update() error {
	// Handle input
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		as.selectedItem--
		if as.selectedItem < 0 {
			as.selectedItem = 5 // Total number of selectable items - 1
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		as.selectedItem++
		if as.selectedItem > 5 {
			as.selectedItem = 0
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		switch as.selectedItem {
		case 0: // Stage selection
			as.selectedStage--
			if as.selectedStage < 0 {
				as.selectedStage = len(as.stages) - 1
			}
		case 1, 2, 3: // Preset army selection
			as.selectedPreset--
			if as.selectedPreset < 0 {
				as.selectedPreset = len(as.presetArmies) - 1
			}
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		switch as.selectedItem {
		case 0: // Stage selection
			as.selectedStage++
			if as.selectedStage >= len(as.stages) {
				as.selectedStage = 0
			}
		case 1, 2, 3: // Preset army selection
			as.selectedPreset++
			if as.selectedPreset >= len(as.presetArmies) {
				as.selectedPreset = 0
			}
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch as.selectedItem {
		case 4: // 戦闘開始
			// Set selected stage and preset in game data
			as.sceneManager.gameData.CurrentStage = as.stages[as.selectedStage]
			// Pass both stage and preset information to battle scene
			battleData := map[string]interface{}{
				"stage":  as.stages[as.selectedStage],
				"preset": as.presetArmies[as.selectedPreset],
			}
			as.sceneManager.TransitionTo(SceneBattle, battleData)
		case 5: // 戻る
			as.sceneManager.TransitionTo(SceneTitle, nil)
		}
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		as.sceneManager.TransitionTo(SceneTitle, nil)
	}
	
	return nil
}

// Draw draws the army setup scene
func (as *ArmySetupScene) Draw(screen *ebiten.Image) {
	// Clear screen with dark background
	screen.Fill(color.RGBA{44, 62, 80, 255}) // #2C3E50
	
	// Draw title
	titleText := "軍勢設定"
	as.textRenderer.DrawTextWithSize(screen, titleText, 450, 50, color.RGBA{236, 240, 241, 255}, 24)
	
	// Draw stage selection
	stageText := "ステージ選択:"
	as.textRenderer.DrawText(screen, stageText, 100, 120, color.RGBA{236, 240, 241, 255})
	
	stageSelectionText := "< " + as.stages[as.selectedStage] + " >"
	if as.selectedItem == 0 {
		as.textRenderer.DrawTextWithShadow(screen, "> "+stageSelectionText, 80, 150, 
			color.RGBA{52, 152, 219, 255}, color.RGBA{0, 0, 0, 128})
	} else {
		as.textRenderer.DrawText(screen, stageSelectionText, 100, 150, color.RGBA{236, 240, 241, 255})
	}
	
	// Draw stage effects
	effectsText := "地形効果:"
	as.textRenderer.DrawText(screen, effectsText, 100, 180, color.RGBA{149, 165, 166, 255})
	
	switch as.selectedStage {
	case 0: // 森の戦い
		as.textRenderer.DrawText(screen, "・移動速度-30%", 100, 200, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・弓兵攻撃+20%", 100, 220, color.RGBA{149, 165, 166, 255})
	case 1: // 山岳要塞
		as.textRenderer.DrawText(screen, "・移動速度-50%", 100, 200, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・防御力+30%", 100, 220, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・魔術師攻撃+30%", 100, 240, color.RGBA{149, 165, 166, 255})
	case 2: // 平原決戦
		as.textRenderer.DrawText(screen, "・移動速度+20%", 100, 200, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・全ユニット攻撃+10%", 100, 220, color.RGBA{149, 165, 166, 255})
	}
	
	// Draw preset armies
	presetText := "プリセット軍勢:"
	as.textRenderer.DrawText(screen, presetText, 100, 300, color.RGBA{236, 240, 241, 255})
	
	// Show current selected preset
	currentPresetText := "< " + as.presetArmies[as.selectedPreset] + " >"
	if as.selectedItem >= 1 && as.selectedItem <= 3 {
		as.textRenderer.DrawTextWithShadow(screen, "> "+currentPresetText, 80, 330, 
			color.RGBA{52, 152, 219, 255}, color.RGBA{0, 0, 0, 128})
	} else {
		as.textRenderer.DrawText(screen, currentPresetText, 100, 330, color.RGBA{236, 240, 241, 255})
	}
	
	// Show preset details
	as.drawPresetDetails(screen, as.selectedPreset)
	
	// Draw buttons
	buttons := []string{"戦闘開始", "戻る"}
	for i, button := range buttons {
		x := 400.0 + float64(i*150)
		y := 500.0
		if as.selectedItem == i+4 {
			as.textRenderer.DrawTextWithShadow(screen, "> "+button+" <", x-20, y, 
				color.RGBA{52, 152, 219, 255}, color.RGBA{0, 0, 0, 128})
		} else {
			as.textRenderer.DrawText(screen, button, x, y, color.RGBA{236, 240, 241, 255})
		}
	}
	
	// Draw controls hint
	controlsText := "↑↓: 選択  ←→: ステージ・編成変更  Enter: 決定  Esc: 戻る"
	as.textRenderer.DrawText(screen, controlsText, 200, 600, color.RGBA{149, 165, 166, 255})
}

// OnEnter is called when entering this scene
func (as *ArmySetupScene) OnEnter(data interface{}) {
	// Reset selection
	as.selectedItem = 0
	as.selectedStage = 0
	as.selectedPreset = 0
}

// OnExit is called when exiting this scene
func (as *ArmySetupScene) OnExit() {
	// Nothing to clean up
}

// drawPresetDetails draws details about the selected preset
func (as *ArmySetupScene) drawPresetDetails(screen *ebiten.Image, presetIndex int) {
	detailsText := "編成詳細:"
	as.textRenderer.DrawText(screen, detailsText, 100, 360, color.RGBA{149, 165, 166, 255})
	
	switch presetIndex {
	case 0: // バランス型
		as.textRenderer.DrawText(screen, "・歩兵: 3部隊", 100, 380, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・弓兵: 2部隊", 100, 400, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・魔術師: 1部隊", 100, 420, color.RGBA{149, 165, 166, 255})
	case 1: // 攻撃重視
		as.textRenderer.DrawText(screen, "・歩兵: 2部隊", 100, 380, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・弓兵: 3部隊", 100, 400, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・魔術師: 2部隊", 100, 420, color.RGBA{149, 165, 166, 255})
	case 2: // 防御重視
		as.textRenderer.DrawText(screen, "・歩兵: 4部隊", 100, 380, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・弓兵: 1部隊", 100, 400, color.RGBA{149, 165, 166, 255})
		as.textRenderer.DrawText(screen, "・魔術師: 1部隊", 100, 420, color.RGBA{149, 165, 166, 255})
	}
}
