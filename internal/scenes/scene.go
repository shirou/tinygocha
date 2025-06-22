package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// SceneType represents different types of scenes
type SceneType int

const (
	SceneTitle SceneType = iota
	SceneArmySetup
	SceneDeployment
	SceneBattle
	SceneResult
	ScenePause
)

// Scene interface that all scenes must implement
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	OnEnter(data interface{})
	OnExit()
}

// GameData holds data that needs to be passed between scenes
type GameData struct {
	// Will be expanded as we implement more features
	CurrentStage  string
	CurrentPreset string
	// ArmyA        *ArmyConfig
	// ArmyB        *ArmyConfig
	// BattleResult *BattleResult
}

// SceneTransition handles smooth transitions between scenes
type SceneTransition struct {
	IsTransitioning bool
	FromScene      SceneType
	ToScene        SceneType
	Progress       float64
	Duration       float64
}

// SceneManager manages all scenes and transitions
type SceneManager struct {
	currentScene SceneType
	scenes       map[SceneType]Scene
	gameData     *GameData
	transition   *SceneTransition
}

// NewSceneManager creates a new scene manager
func NewSceneManager() *SceneManager {
	return &SceneManager{
		currentScene: SceneTitle,
		scenes:       make(map[SceneType]Scene),
		gameData:     &GameData{},
		transition: &SceneTransition{
			IsTransitioning: false,
			Duration:        0.5, // 0.5 seconds transition
		},
	}
}

// RegisterScene registers a scene with the manager
func (sm *SceneManager) RegisterScene(sceneType SceneType, scene Scene) {
	sm.scenes[sceneType] = scene
}

// TransitionTo starts a transition to a new scene
func (sm *SceneManager) TransitionTo(sceneType SceneType, data interface{}) {
	if sm.currentScene == sceneType {
		return
	}

	sm.transition.IsTransitioning = true
	sm.transition.FromScene = sm.currentScene
	sm.transition.ToScene = sceneType
	sm.transition.Progress = 0.0

	// Pass data to the new scene
	if data != nil {
		// Update game data based on the passed data
		if battleData, ok := data.(map[string]interface{}); ok {
			if stage, exists := battleData["stage"]; exists {
				if stageStr, ok := stage.(string); ok {
					sm.gameData.CurrentStage = stageStr
				}
			}
			if preset, exists := battleData["preset"]; exists {
				if presetStr, ok := preset.(string); ok {
					sm.gameData.CurrentPreset = presetStr
				}
			}
		}
	}
}

// Update updates the current scene and handles transitions
func (sm *SceneManager) Update() error {
	if sm.transition.IsTransitioning {
		sm.transition.Progress += 1.0 / 60.0 / sm.transition.Duration // Assuming 60 FPS
		
		if sm.transition.Progress >= 1.0 {
			// Transition complete
			if currentScene := sm.scenes[sm.currentScene]; currentScene != nil {
				currentScene.OnExit()
			}
			
			sm.currentScene = sm.transition.ToScene
			
			if newScene := sm.scenes[sm.currentScene]; newScene != nil {
				newScene.OnEnter(sm.gameData)
			}
			
			sm.transition.IsTransitioning = false
		}
		return nil
	}

	// Update current scene
	if scene := sm.scenes[sm.currentScene]; scene != nil {
		return scene.Update()
	}
	
	return nil
}

// Draw draws the current scene with transition effects
func (sm *SceneManager) Draw(screen *ebiten.Image) {
	if sm.transition.IsTransitioning {
		// During transition, we could implement fade effects here
		// For now, just draw the current scene
		if scene := sm.scenes[sm.currentScene]; scene != nil {
			scene.Draw(screen)
		}
		
		// Apply fade effect based on transition progress
		// This will be implemented later with proper graphics
		return
	}

	// Draw current scene
	if scene := sm.scenes[sm.currentScene]; scene != nil {
		scene.Draw(screen)
	}
}

// GetCurrentScene returns the current scene type
func (sm *SceneManager) GetCurrentScene() SceneType {
	return sm.currentScene
}

// GetGameData returns the shared game data
func (sm *SceneManager) GetGameData() *GameData {
	return sm.gameData
}
