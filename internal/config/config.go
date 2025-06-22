package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the global game configuration
type Config struct {
	Graphics GraphicsConfig `toml:"graphics"`
	Audio    AudioConfig    `toml:"audio"`
	Game     GameConfig     `toml:"game"`
}

// GraphicsConfig represents graphics settings
type GraphicsConfig struct {
	FontPath     string  `toml:"font_path"`
	FontSize     int     `toml:"font_size"`
	UIScale      float64 `toml:"ui_scale"`
	ShowFPS      bool    `toml:"show_fps"`
	VSync        bool    `toml:"vsync"`
}

// AudioConfig represents audio settings
type AudioConfig struct {
	MasterVolume float64 `toml:"master_volume"`
	SFXVolume    float64 `toml:"sfx_volume"`
	BGMVolume    float64 `toml:"bgm_volume"`
	Enabled      bool    `toml:"enabled"`
}

// GameConfig represents game settings
type GameConfig struct {
	Language     string `toml:"language"`
	AutoSave     bool   `toml:"auto_save"`
	ShowTutorial bool   `toml:"show_tutorial"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Graphics: GraphicsConfig{
			FontPath: "", // Empty means use default MPlus1p
			FontSize: 16,
			UIScale:  1.0,
			ShowFPS:  false,
			VSync:    true,
		},
		Audio: AudioConfig{
			MasterVolume: 0.8,
			SFXVolume:    0.7,
			BGMVolume:    0.6,
			Enabled:      true,
		},
		Game: GameConfig{
			Language:     "ja",
			AutoSave:     true,
			ShowTutorial: true,
		},
	}
}

// LoadConfig loads configuration from file
func LoadConfig(filename string) (*Config, error) {
	// Start with default config
	config := DefaultConfig()
	
	// Try to load from file
	data, err := os.ReadFile(filename)
	if err != nil {
		// If file doesn't exist, return default config
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}
	
	// Parse TOML
	if err := toml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	
	return config, nil
}

// SaveConfig saves configuration to file
func (c *Config) SaveConfig(filename string) error {
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}
