# データ構造設計

## TOMLファイル構成

### ユニット定義ファイル (units.toml)

```toml
# 近接系ユニット
[unit_types.infantry]
name = "歩兵"
hp = 100
attack = 15
defense = 10
speed = 2.0
range = 1.0
magic_power = 0

# 弓系ユニット
[unit_types.archer]
name = "弓兵"
hp = 70
attack = 12
defense = 5
speed = 1.8
range = 8.0
magic_power = 0

# 魔術師系ユニット
[unit_types.mage]
name = "魔術師"
hp = 50
attack = 8
defense = 3
speed = 1.5
range = 10.0
magic_power = 20
```

### 地形効果定義ファイル (terrain.toml)

```toml
[terrain_types.forest]
name = "森"
movement_modifier = 0.7  # 移動速度70%
defense_modifier = 1.1   # 防御力110%
archer_bonus = 1.2       # 弓系攻撃力120%

[terrain_types.mountain]
name = "山"
movement_modifier = 0.5  # 移動速度50%
defense_modifier = 1.3   # 防御力130%
mage_bonus = 1.3         # 魔術師系攻撃力130%

[terrain_types.plain]
name = "平原"
movement_modifier = 1.2  # 移動速度120%
defense_modifier = 1.0   # 防御力100%
```

### ステージ定義ファイル (stages.toml)

```toml
[stages.forest_battle]
name = "森の戦い"
terrain = "forest"
deployment_points_a = [
    { x = 100, y = 200 },
    { x = 150, y = 250 },
    { x = 200, y = 300 }
]
deployment_points_b = [
    { x = 700, y = 200 },
    { x = 650, y = 250 },
    { x = 600, y = 300 }
]
time_limit = 300  # 秒
```

## Go言語での基本クラス構造

### Unit（個別ユニット）
```go
type Unit struct {
    ID           int
    Type         string
    Name         string
    HP           int
    MaxHP        int
    Attack       int
    Defense      int
    Speed        float64
    Range        float64
    MagicPower   int
    Position     Vector2D
    IsLeader     bool
    IsAlive      bool
    IsRetreating bool
    GroupID      int
}
```

### Group（グループ管理）
```go
type Group struct {
    ID        int
    Leader    *Unit
    Members   []*Unit
    Formation FormationType
    ArmyID    int
}
```

### Formation（隊形管理）
```go
type FormationType int

const (
    CircleFormation FormationType = iota
    // 将来追加予定: LineFormation, WedgeFormation, etc.
)

type Formation struct {
    Type     FormationType
    Radius   float64  // 円形隊形の半径
    Spacing  float64  // ユニット間の間隔
}
```

### Army（軍勢管理）
```go
type Army struct {
    ID     int
    Name   string
    Groups []*Group
    Side   int  // 0: A軍, 1: B軍
}
```

### Game（ゲーム全体管理）
```go
type Game struct {
    ArmyA        *Army
    ArmyB        *Army
    Stage        *Stage
    GameTime     float64
    TimeLimit    float64
    GameState    GameState
    Winner       int  // -1: 未決定, 0: A軍勝利, 1: B軍勝利, 2: 引き分け
}

type GameState int

const (
    GameStateDeployment GameState = iota
    GameStateBattle
    GameStateResult
)
```

### Stage（ステージ管理）
```go
type Stage struct {
    Name              string
    TerrainType       string
    DeploymentPointsA []Vector2D
    DeploymentPointsB []Vector2D
    TimeLimit         float64
    Width             int
    Height            int
}
```

### Vector2D（座標管理）
```go
type Vector2D struct {
    X float64
    Y float64
}

func (v Vector2D) Distance(other Vector2D) float64 {
    dx := v.X - other.X
    dy := v.Y - other.Y
    return math.Sqrt(dx*dx + dy*dy)
}
```

## データローダー

### UnitTypeLoader
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

### TerrainLoader
```go
type TerrainConfig struct {
    Name             string  `toml:"name"`
    MovementModifier float64 `toml:"movement_modifier"`
    DefenseModifier  float64 `toml:"defense_modifier"`
    ArcherBonus      float64 `toml:"archer_bonus"`
    MageBonus        float64 `toml:"mage_bonus"`
}

type TerrainsConfig struct {
    TerrainTypes map[string]TerrainConfig `toml:"terrain_types"`
}
```
