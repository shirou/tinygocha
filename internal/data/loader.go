package data

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// DataManager manages all game data
type DataManager struct {
	Units    *UnitsConfig
	Terrains *TerrainsConfig
	Stages   *StagesConfig
}

// NewDataManager creates a new data manager
func NewDataManager() *DataManager {
	return &DataManager{
		Units:    &UnitsConfig{UnitTypes: make(map[string]UnitTypeConfig)},
		Terrains: &TerrainsConfig{TerrainTypes: make(map[string]TerrainConfig)},
		Stages:   &StagesConfig{Stages: make(map[string]StageConfig)},
	}
}

// LoadAll loads all data files
func (dm *DataManager) LoadAll() error {
	if err := dm.LoadUnits("assets/data/units.toml"); err != nil {
		return fmt.Errorf("failed to load units: %w", err)
	}
	
	if err := dm.LoadTerrains("assets/data/terrain.toml"); err != nil {
		return fmt.Errorf("failed to load terrains: %w", err)
	}
	
	if err := dm.LoadStages("assets/data/stages.toml"); err != nil {
		return fmt.Errorf("failed to load stages: %w", err)
	}
	
	return nil
}

// LoadUnits loads unit configurations from TOML file
func (dm *DataManager) LoadUnits(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	
	var config UnitsConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse TOML in %s: %w", filename, err)
	}
	
	dm.Units = &config
	return nil
}

// LoadTerrains loads terrain configurations from TOML file
func (dm *DataManager) LoadTerrains(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	
	var config TerrainsConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse TOML in %s: %w", filename, err)
	}
	
	dm.Terrains = &config
	return nil
}

// LoadStages loads stage configurations from TOML file
func (dm *DataManager) LoadStages(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	
	var config StagesConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse TOML in %s: %w", filename, err)
	}
	
	dm.Stages = &config
	return nil
}

// GetUnitConfig returns unit configuration by type
func (dm *DataManager) GetUnitConfig(unitType string) (UnitTypeConfig, error) {
	config, exists := dm.Units.GetUnitConfig(unitType)
	if !exists {
		return UnitTypeConfig{}, fmt.Errorf("unit type %s not found", unitType)
	}
	return config, nil
}

// GetTerrainConfig returns terrain configuration by type
func (dm *DataManager) GetTerrainConfig(terrainType string) (TerrainConfig, error) {
	config, exists := dm.Terrains.GetTerrainConfig(terrainType)
	if !exists {
		return TerrainConfig{}, fmt.Errorf("terrain type %s not found", terrainType)
	}
	return config, nil
}

// GetStageConfig returns stage configuration by name
func (dm *DataManager) GetStageConfig(stageName string) (StageConfig, error) {
	config, exists := dm.Stages.GetStageConfig(stageName)
	if !exists {
		return StageConfig{}, fmt.Errorf("stage %s not found", stageName)
	}
	return config, nil
}
