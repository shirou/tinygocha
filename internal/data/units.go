package data

// UnitTypeConfig represents unit configuration from TOML
type UnitTypeConfig struct {
	Name       string  `toml:"name"`
	HP         int     `toml:"hp"`
	Attack     int     `toml:"attack"`
	Defense    int     `toml:"defense"`
	Speed      float64 `toml:"speed"`
	Range      float64 `toml:"range"`
	SightRange float64 `toml:"sight_range"` // 知覚範囲
	MagicPower int     `toml:"magic_power"`
	Size       float64 `toml:"size"`  // ユニットの大きさ（衝突判定用）
}

// UnitsConfig represents the entire units configuration
type UnitsConfig struct {
	UnitTypes map[string]UnitTypeConfig `toml:"unit_types"`
}

// GetUnitConfig returns the configuration for a specific unit type
func (uc *UnitsConfig) GetUnitConfig(unitType string) (UnitTypeConfig, bool) {
	config, exists := uc.UnitTypes[unitType]
	return config, exists
}
