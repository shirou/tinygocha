package data

import (
	gamemath "github.com/shirou/tinygocha/internal/math"
)

// DeploymentPoint represents a deployment position
type DeploymentPoint struct {
	X float64 `toml:"x"`
	Y float64 `toml:"y"`
}

// ToVector2D converts DeploymentPoint to Vector2D
func (dp DeploymentPoint) ToVector2D() gamemath.Vector2D {
	return gamemath.Vector2D{X: dp.X, Y: dp.Y}
}

// StageConfig represents stage configuration from TOML
type StageConfig struct {
	Name              string            `toml:"name"`
	Terrain           string            `toml:"terrain"`
	DeploymentPointsA []DeploymentPoint `toml:"deployment_points_a"`
	DeploymentPointsB []DeploymentPoint `toml:"deployment_points_b"`
	TimeLimit         float64           `toml:"time_limit"`
	Width             int               `toml:"width"`
	Height            int               `toml:"height"`
}

// StagesConfig represents the entire stages configuration
type StagesConfig struct {
	Stages map[string]StageConfig `toml:"stages"`
}

// GetStageConfig returns the configuration for a specific stage
func (sc *StagesConfig) GetStageConfig(stageName string) (StageConfig, bool) {
	config, exists := sc.Stages[stageName]
	return config, exists
}

// GetDeploymentPointsA returns deployment points for Army A as Vector2D slice
func (sc StageConfig) GetDeploymentPointsA() []gamemath.Vector2D {
	points := make([]gamemath.Vector2D, len(sc.DeploymentPointsA))
	for i, dp := range sc.DeploymentPointsA {
		points[i] = dp.ToVector2D()
	}
	return points
}

// GetDeploymentPointsB returns deployment points for Army B as Vector2D slice
func (sc StageConfig) GetDeploymentPointsB() []gamemath.Vector2D {
	points := make([]gamemath.Vector2D, len(sc.DeploymentPointsB))
	for i, dp := range sc.DeploymentPointsB {
		points[i] = dp.ToVector2D()
	}
	return points
}
