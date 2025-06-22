# データシステム設計

## 概要

ゴチャキャラバトルでは、ゲームデータをTOML形式で管理し、柔軟な設定変更とMOD対応を実現しています。ユニット定義、地形効果、ステージ設定などを外部ファイルで管理します。

## アーキテクチャ

### データ構造

```
DataManager
├── UnitsConfig    (ユニット定義)
├── TerrainsConfig (地形効果)
└── StagesConfig   (ステージ設定)

TOML Files
├── units.toml     (ユニット統計)
├── terrain.toml   (地形効果)
└── stages.toml    (ステージ定義)
```

## DataManager

### 基本構造

```go
type DataManager struct {
    Units    *UnitsConfig
    Terrains *TerrainsConfig
    Stages   *StagesConfig
}

func NewDataManager() *DataManager {
    return &DataManager{
        Units:    &UnitsConfig{UnitTypes: make(map[string]UnitTypeConfig)},
        Terrains: &TerrainsConfig{TerrainTypes: make(map[string]TerrainConfig)},
        Stages:   &StagesConfig{Stages: make(map[string]StageConfig)},
    }
}
```

### 読み込み処理

```go
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
```

## ユニットデータ

### 構造定義

```go
type UnitTypeConfig struct {
    Name       string  `toml:"name"`
    HP         int     `toml:"hp"`
    Attack     int     `toml:"attack"`
    Defense    int     `toml:"defense"`
    Speed      float64 `toml:"speed"`
    Range      float64 `toml:"range"`
    MagicPower int     `toml:"magic_power"`
}

type UnitsConfig struct {
    UnitTypes map[string]UnitTypeConfig `toml:"unit_types"`
}
```

### データファイル例

```toml
# assets/data/units.toml

[unit_types.infantry]
name = "歩兵"
hp = 100
attack = 15
defense = 10
speed = 2.0
range = 1.5
magic_power = 0

[unit_types.archer]
name = "弓兵"
hp = 70
attack = 12
defense = 5
speed = 1.8
range = 8.0
magic_power = 0

[unit_types.mage]
name = "魔術師"
hp = 50
attack = 8
defense = 3
speed = 1.5
range = 10.0
magic_power = 20

[unit_types.heavy_infantry]
name = "重装歩兵"
hp = 120
attack = 18
defense = 15
speed = 1.5
range = 1.5
magic_power = 0

[unit_types.cavalry]
name = "騎兵"
hp = 90
attack = 20
defense = 8
speed = 3.5
range = 2.0
magic_power = 0
```

### ユニット特性

| ユニット | HP | 攻撃 | 防御 | 速度 | 射程 | 魔力 | 特徴 |
|----------|----|----|----|----|----|----|------|
| 歩兵 | 100 | 15 | 10 | 2.0 | 1.5 | 0 | バランス型 |
| 弓兵 | 70 | 12 | 5 | 1.8 | 8.0 | 0 | 遠距離攻撃 |
| 魔術師 | 50 | 8 | 3 | 1.5 | 10.0 | 20 | 魔法攻撃 |
| 重装歩兵 | 120 | 18 | 15 | 1.5 | 1.5 | 0 | 高防御 |
| 騎兵 | 90 | 20 | 8 | 3.5 | 2.0 | 0 | 高機動 |

## 地形データ

### 構造定義

```go
type TerrainConfig struct {
    Name             string  `toml:"name"`
    MovementModifier float64 `toml:"movement_modifier"`
    DefenseModifier  float64 `toml:"defense_modifier"`
    ArcherBonus      float64 `toml:"archer_bonus"`
    MageBonus        float64 `toml:"mage_bonus"`
    InfantryBonus    float64 `toml:"infantry_bonus"`
}

type TerrainsConfig struct {
    TerrainTypes map[string]TerrainConfig `toml:"terrain_types"`
}
```

### データファイル例

```toml
# assets/data/terrain.toml

[terrain_types.forest]
name = "森"
movement_modifier = 0.7  # 移動速度70%
defense_modifier = 1.1   # 防御力110%
archer_bonus = 1.2       # 弓系攻撃力120%
mage_bonus = 1.0         # 魔術師系攻撃力100%
infantry_bonus = 0.9     # 歩兵系攻撃力90%

[terrain_types.mountain]
name = "山"
movement_modifier = 0.5  # 移動速度50%
defense_modifier = 1.3   # 防御力130%
archer_bonus = 1.1       # 弓系攻撃力110%
mage_bonus = 1.3         # 魔術師系攻撃力130%
infantry_bonus = 0.8     # 歩兵系攻撃力80%

[terrain_types.plain]
name = "平原"
movement_modifier = 1.2  # 移動速度120%
defense_modifier = 1.0   # 防御力100%
archer_bonus = 1.0       # 弓系攻撃力100%
mage_bonus = 1.0         # 魔術師系攻撃力100%
infantry_bonus = 1.1     # 歩兵系攻撃力110%
```

### 地形効果一覧

| 地形 | 移動 | 防御 | 弓兵 | 魔術師 | 歩兵 | 戦術的意味 |
|------|------|------|------|--------|------|-----------|
| 森 | -30% | +10% | +20% | ±0% | -10% | 弓兵有利 |
| 山 | -50% | +30% | +10% | +30% | -20% | 魔術師有利、防御的 |
| 平原 | +20% | ±0% | ±0% | ±0% | +10% | 機動戦向き |
| 城塞 | -20% | +50% | +30% | +10% | +20% | 防御拠点 |
| 街 | ±0% | +20% | ±0% | +20% | ±0% | 魔術師支援 |

## ステージデータ

### 構造定義

```go
type DeploymentPoint struct {
    X float64 `toml:"x"`
    Y float64 `toml:"y"`
}

type StageConfig struct {
    Name              string            `toml:"name"`
    Terrain           string            `toml:"terrain"`
    DeploymentPointsA []DeploymentPoint `toml:"deployment_points_a"`
    DeploymentPointsB []DeploymentPoint `toml:"deployment_points_b"`
    TimeLimit         float64           `toml:"time_limit"`
    Width             int               `toml:"width"`
    Height            int               `toml:"height"`
}

type StagesConfig struct {
    Stages map[string]StageConfig `toml:"stages"`
}
```

### データファイル例

```toml
# assets/data/stages.toml

[stages.forest_battle]
name = "森の戦い"
terrain = "forest"
time_limit = 300.0  # 5分
width = 1024
height = 768

deployment_points_a = [
    { x = 100, y = 200 },
    { x = 150, y = 300 },
    { x = 100, y = 400 },
    { x = 200, y = 250 },
    { x = 200, y = 350 }
]

deployment_points_b = [
    { x = 900, y = 200 },
    { x = 850, y = 300 },
    { x = 900, y = 400 },
    { x = 800, y = 250 },
    { x = 800, y = 350 }
]

[stages.mountain_fortress]
name = "山岳要塞"
terrain = "mountain"
time_limit = 400.0  # 6分40秒
width = 1024
height = 768

deployment_points_a = [
    { x = 80, y = 150 },
    { x = 120, y = 250 },
    { x = 80, y = 350 },
    { x = 80, y = 450 },
    { x = 180, y = 300 }
]

deployment_points_b = [
    { x = 920, y = 150 },
    { x = 880, y = 250 },
    { x = 920, y = 350 },
    { x = 920, y = 450 },
    { x = 820, y = 300 }
]
```

### ステージ特性

| ステージ | 地形 | 制限時間 | 配置数 | 戦術的特徴 |
|----------|------|----------|--------|-----------|
| 森の戦い | 森 | 5分 | 5vs5 | 弓兵有利、視界制限 |
| 山岳要塞 | 山 | 6分40秒 | 5vs5 | 魔術師有利、防御的 |
| 平原決戦 | 平原 | 4分10秒 | 7vs7 | 機動戦、大規模 |

## データ活用

### 戦闘システム統合

```go
func (bm *BattleManager) applyTerrainModifiers(unit *Unit) {
    // 移動速度修正
    unit.Speed *= bm.TerrainData.MovementModifier
    
    // 防御力修正
    unit.Defense = int(float64(unit.Defense) * bm.TerrainData.DefenseModifier)
    
    // ユニット種別ボーナス
    switch unit.Type {
    case UnitTypeInfantry:
        unit.AttackPower = int(float64(unit.AttackPower) * bm.TerrainData.InfantryBonus)
    case UnitTypeArcher:
        unit.AttackPower = int(float64(unit.AttackPower) * bm.TerrainData.ArcherBonus)
    case UnitTypeMage:
        unit.AttackPower = int(float64(unit.AttackPower) * bm.TerrainData.MageBonus)
        unit.MagicPower = int(float64(unit.MagicPower) * bm.TerrainData.MageBonus)
    }
}
```

### プリセット軍勢生成

```go
func (bm *BattleManager) createBalancedArmy(army *Army, deploymentPoints []Vector2D, dataManager *data.DataManager) {
    groupConfigs := []struct {
        leaderType string
        memberType string
        count      int
    }{
        {"infantry", "infantry", 4},
        {"archer", "archer", 3},
        {"mage", "infantry", 2},
    }
    
    for i, config := range groupConfigs {
        if i >= len(deploymentPoints) {
            break
        }
        
        group := bm.createGroup(army.ID, config.leaderType, config.memberType, 
                               config.count, deploymentPoints[i], dataManager)
        army.AddGroup(group)
    }
}
```

## エラーハンドリング

### ファイル読み込みエラー

```go
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
```

### データ検証

```go
func (dm *DataManager) ValidateData() error {
    // ユニットデータ検証
    for unitType, config := range dm.Units.UnitTypes {
        if config.HP <= 0 {
            return fmt.Errorf("invalid HP for unit %s: %d", unitType, config.HP)
        }
        if config.Speed <= 0 {
            return fmt.Errorf("invalid speed for unit %s: %f", unitType, config.Speed)
        }
    }
    
    // 地形データ検証
    for terrainType, config := range dm.Terrains.TerrainTypes {
        if config.MovementModifier <= 0 {
            return fmt.Errorf("invalid movement modifier for terrain %s: %f", 
                             terrainType, config.MovementModifier)
        }
    }
    
    return nil
}
```

## MOD対応

### カスタムデータディレクトリ

```go
func (dm *DataManager) LoadFromDirectory(dataDir string) error {
    unitsFile := filepath.Join(dataDir, "units.toml")
    terrainFile := filepath.Join(dataDir, "terrain.toml")
    stagesFile := filepath.Join(dataDir, "stages.toml")
    
    if err := dm.LoadUnits(unitsFile); err != nil {
        return err
    }
    // 他のファイルも同様に読み込み...
    
    return nil
}
```

### データ上書き機能

```go
func (dm *DataManager) MergeUnits(additionalUnits *UnitsConfig) {
    for unitType, config := range additionalUnits.UnitTypes {
        dm.Units.UnitTypes[unitType] = config
    }
}
```

## パフォーマンス最適化

### データキャッシュ

```go
type DataCache struct {
    unitConfigs    map[string]UnitTypeConfig
    terrainConfigs map[string]TerrainConfig
    stageConfigs   map[string]StageConfig
    mutex          sync.RWMutex
}

func (dc *DataCache) GetUnitConfig(unitType string) (UnitTypeConfig, bool) {
    dc.mutex.RLock()
    defer dc.mutex.RUnlock()
    
    config, exists := dc.unitConfigs[unitType]
    return config, exists
}
```

### 遅延読み込み

```go
func (dm *DataManager) GetUnitConfigLazy(unitType string) (UnitTypeConfig, error) {
    if config, exists := dm.Units.GetUnitConfig(unitType); exists {
        return config, nil
    }
    
    // 必要に応じて追加データを読み込み
    if err := dm.loadAdditionalUnits(); err != nil {
        return UnitTypeConfig{}, err
    }
    
    return dm.Units.GetUnitConfig(unitType)
}
```

## 今後の拡張

### 予定機能

1. **動的データ更新**: ゲーム中の設定変更
2. **バージョン管理**: データファイルの互換性チェック
3. **圧縮対応**: 大きなデータファイルの圧縮
4. **暗号化**: MOD保護のためのデータ暗号化

### データ拡張

```toml
# 将来の拡張例

[unit_types.dragon]
name = "ドラゴン"
hp = 500
attack = 50
defense = 30
speed = 4.0
range = 12.0
magic_power = 40
special_abilities = ["fire_breath", "flight"]
cost = 100
rarity = "legendary"

[terrain_types.lava]
name = "溶岩"
movement_modifier = 0.3
defense_modifier = 0.8
fire_resistance_required = true
damage_per_turn = 10
```

### スクリプト対応

```go
type ScriptableUnit struct {
    UnitTypeConfig
    ScriptPath string `toml:"script_path"`
    AIBehavior string `toml:"ai_behavior"`
}
```

## デバッグ支援

### データ検証ツール

```bash
# データファイル検証
go run tools/validate_data.go assets/data/

# 統計情報出力
go run tools/data_stats.go assets/data/units.toml
```

### ホットリロード

```go
func (dm *DataManager) WatchFiles() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()
    
    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    log.Println("Reloading data file:", event.Name)
                    dm.LoadAll()
                }
            }
        }
    }()
    
    watcher.Add("assets/data/")
}
```
