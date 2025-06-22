package data

// TerrainConfig represents terrain configuration from TOML
type TerrainConfig struct {
	Name             string  `toml:"name"`
	MovementModifier float64 `toml:"movement_modifier"`
	DefenseModifier  float64 `toml:"defense_modifier"`
	ArcherBonus      float64 `toml:"archer_bonus"`
	MageBonus        float64 `toml:"mage_bonus"`
	InfantryBonus    float64 `toml:"infantry_bonus"`
}

// TerrainsConfig represents the entire terrain configuration
type TerrainsConfig struct {
	TerrainTypes map[string]TerrainConfig `toml:"terrain_types"`
}

// GetTerrainConfig returns the configuration for a specific terrain type
func (tc *TerrainsConfig) GetTerrainConfig(terrainType string) (TerrainConfig, bool) {
	config, exists := tc.TerrainTypes[terrainType]
	return config, exists
}
