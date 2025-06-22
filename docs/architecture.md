# アーキテクチャ設計

## 概要

ゴチャキャラバトルは、Go言語とEbitengineを使用したリアルタイム戦術バトルゲームです。モジュラー設計により、拡張性と保守性を重視したアーキテクチャを採用しています。

## 全体アーキテクチャ

### システム構成

```
┌─────────────────────────────────────────────────────────────┐
│                        Main Game                            │
├─────────────────────────────────────────────────────────────┤
│  Scene Management  │  Config System  │  Font Management    │
├─────────────────────────────────────────────────────────────┤
│     Game Logic     │   Data System   │  Graphics System    │
├─────────────────────────────────────────────────────────────┤
│                      Ebitengine v2                          │
├─────────────────────────────────────────────────────────────┤
│                     Go Runtime                              │
└─────────────────────────────────────────────────────────────┘
```

### ディレクトリ構造

```
tinygocha/
├── main.go                    # エントリーポイント
├── config.toml               # 設定ファイル
├── config_sample.toml        # 設定サンプル
├── Makefile                  # ビルド設定
├── README.md                 # プロジェクト説明
├── internal/                 # 内部パッケージ
│   ├── config/              # 設定管理
│   │   └── config.go
│   ├── data/                # データ管理
│   │   ├── loader.go
│   │   ├── units.go
│   │   ├── terrain.go
│   │   └── stages.go
│   ├── game/                # ゲームロジック
│   │   ├── unit.go
│   │   ├── group.go
│   │   ├── army.go
│   │   ├── battle.go
│   │   └── config.go
│   ├── graphics/            # 描画システム
│   │   ├── animation.go
│   │   ├── sprite_generator.go
│   │   ├── font_manager.go
│   │   └── text_renderer.go
│   ├── math/                # 数学ライブラリ
│   │   └── vector.go
│   └── scenes/              # シーン管理
│       ├── scene.go
│       ├── title.go
│       ├── army_setup.go
│       ├── battle_new.go
│       ├── battle_draw_methods.go
│       └── result.go
├── assets/                  # ゲームアセット
│   └── data/               # データファイル
│       ├── units.toml
│       ├── terrain.toml
│       └── stages.toml
├── build/                   # ビルド成果物
├── docs/                    # ドキュメント
└── tools/                   # 開発ツール
```

## レイヤー設計

### 1. プレゼンテーション層

#### Scene Management
```go
type SceneManager struct {
    currentScene SceneType
    scenes       map[SceneType]Scene
    gameData     *GameData
    transition   *SceneTransition
}
```

**責務:**
- シーン間の遷移管理
- ゲーム状態の保持
- 入力処理の委譲

#### Graphics System
```go
type TextRenderer struct {
    fontManager *FontManager
}

type SpriteGenerator struct {
    cache map[string]*ebiten.Image
}
```

**責務:**
- テキスト描画
- スプライト生成
- アニメーション管理

### 2. アプリケーション層

#### Game Logic
```go
type BattleManager struct {
    ArmyA        *Army
    ArmyB        *Army
    Stage        data.StageConfig
    TerrainData  data.TerrainConfig
    // ...
}
```

**責務:**
- 戦闘ロジック
- ゲームルール適用
- 勝利条件判定

#### Configuration
```go
type Config struct {
    Graphics GraphicsConfig
    Audio    AudioConfig
    Game     GameConfig
}
```

**責務:**
- 設定管理
- 環境固有設定
- ユーザー設定

### 3. ドメイン層

#### Core Entities
```go
type Unit struct {
    // 基本属性
    ID, Type, Name, HP, Attack, Defense
    // 位置・移動
    Position, Target, Speed
    // 状態
    IsAlive, IsLeader, IsRetreating
    // アニメーション
    Animation *graphics.AnimationState
}

type Group struct {
    Leader    *Unit
    Members   []*Unit
    Formation Formation
}

type Army struct {
    Groups []*Group
    Side   int
}
```

**責務:**
- ビジネスロジック
- ドメインルール
- エンティティ関係

### 4. インフラストラクチャ層

#### Data Access
```go
type DataManager struct {
    Units    *UnitsConfig
    Terrains *TerrainsConfig
    Stages   *StagesConfig
}
```

**責務:**
- データ永続化
- 外部ファイル読み込み
- データ変換

## 設計パターン

### 1. Strategy Pattern

#### アニメーション戦略
```go
type AnimationState struct {
    Type AnimationType
    // ...
}

func (as *AnimationState) GetAnimationOffset() (float64, float64) {
    switch as.Type {
    case AnimationIdle:
        return as.getIdleOffset()
    case AnimationWalk:
        return as.getWalkOffset()
    case AnimationAttack:
        return as.getAttackOffset()
    }
}
```

### 2. Factory Pattern

#### ユニット生成
```go
func (bm *BattleManager) createUnit(unitType UnitType, config UnitTypeConfig, isLeader bool, armyID int) *Unit {
    unit := NewUnit(bm.nextUnitID, unitType, config, isLeader, 0, armyID)
    bm.nextUnitID++
    
    // 地形効果適用
    bm.applyTerrainModifiers(unit)
    
    return unit
}
```

### 3. Observer Pattern

#### 戦闘イベント
```go
type BattleEvent struct {
    Type   EventType
    Source *Unit
    Target *Unit
    Data   interface{}
}

type EventListener interface {
    OnBattleEvent(event BattleEvent)
}
```

### 4. State Pattern

#### ユニット状態
```go
type UnitState interface {
    Update(unit *Unit, deltaTime float64)
    OnEnter(unit *Unit)
    OnExit(unit *Unit)
}

type IdleState struct{}
type MovingState struct{}
type AttackingState struct{}
type DeadState struct{}
```

### 5. Command Pattern

#### ユーザー入力
```go
type Command interface {
    Execute()
    Undo()
}

type MoveCommand struct {
    unit   *Unit
    target Vector2D
}

type AttackCommand struct {
    attacker *Unit
    target   *Unit
}
```

## データフロー

### 1. 初期化フロー

```
main() 
  ↓
NewGame()
  ↓
Config.LoadConfig() → FontManager.LoadFont() → DataManager.LoadAll()
  ↓
SceneManager.RegisterScenes()
  ↓
ebiten.RunGame()
```

### 2. ゲームループ

```
Game.Update()
  ↓
SceneManager.Update()
  ↓
CurrentScene.Update()
  ↓
BattleManager.Update() (戦闘シーンの場合)
  ↓
Army.Update() → Group.Update() → Unit.Update()
```

### 3. 描画フロー

```
Game.Draw()
  ↓
SceneManager.Draw()
  ↓
CurrentScene.Draw()
  ↓
SpriteGenerator.GenerateUnitSprite() → TextRenderer.DrawText()
```

## 依存関係管理

### パッケージ依存関係

```
main
  ↓
scenes ← config ← graphics ← data
  ↓      ↓         ↓
game ← math    font_manager
```

### 依存性注入

```go
// コンストラクタ注入
func NewBattleSceneNew(
    sceneManager *SceneManager, 
    dataManager *data.DataManager, 
    textRenderer *graphics.TextRenderer
) *BattleSceneNew

// インターフェース分離
type TextDrawer interface {
    DrawText(screen *ebiten.Image, text string, x, y float64, color color.Color)
}
```

## エラーハンドリング

### エラー戦略

1. **Graceful Degradation**: 機能低下での継続
2. **Fail Fast**: 早期エラー検出
3. **Logging**: 詳細なエラーログ

```go
func (dm *DataManager) LoadAll() error {
    if err := dm.LoadUnits("assets/data/units.toml"); err != nil {
        log.Printf("Warning: Failed to load units: %v", err)
        // デフォルトデータで継続
        dm.loadDefaultUnits()
    }
    
    // 他のデータも同様...
    return nil
}
```

### エラー分類

```go
type GameError struct {
    Type    ErrorType
    Message string
    Cause   error
}

const (
    ErrorTypeConfig ErrorType = iota
    ErrorTypeData
    ErrorTypeGraphics
    ErrorTypeAudio
)
```

## パフォーマンス設計

### メモリ管理

```go
// オブジェクトプール
type UnitPool struct {
    pool sync.Pool
}

func (up *UnitPool) Get() *Unit {
    if unit := up.pool.Get(); unit != nil {
        return unit.(*Unit)
    }
    return &Unit{}
}

func (up *UnitPool) Put(unit *Unit) {
    unit.Reset()
    up.pool.Put(unit)
}
```

### 計算最適化

```go
// 空間分割
type SpatialGrid struct {
    cellSize int
    cells    map[string][]*Unit
}

// フレームレート制御
const (
    TargetFPS = 60
    FrameTime = time.Second / TargetFPS
)
```

## 拡張性設計

### プラグインアーキテクチャ

```go
type Plugin interface {
    Name() string
    Version() string
    Initialize(game *Game) error
    Update(deltaTime float64) error
    Shutdown() error
}

type PluginManager struct {
    plugins []Plugin
}
```

### MOD対応

```go
type ModLoader struct {
    modPaths []string
    loadedMods map[string]*Mod
}

type Mod struct {
    Metadata ModMetadata
    Data     ModData
    Scripts  []Script
}
```

### 設定拡張

```toml
[mods]
enabled = true
mod_directory = "mods/"
allowed_mods = ["balance_mod", "graphics_mod"]

[developer]
debug_mode = false
show_collision = false
log_level = "info"
```

## セキュリティ設計

### データ検証

```go
func (dm *DataManager) ValidateData() error {
    for unitType, config := range dm.Units.UnitTypes {
        if err := validateUnitConfig(unitType, config); err != nil {
            return fmt.Errorf("invalid unit config %s: %w", unitType, err)
        }
    }
    return nil
}

func validateUnitConfig(unitType string, config UnitTypeConfig) error {
    if config.HP <= 0 {
        return errors.New("HP must be positive")
    }
    if config.Speed < 0 {
        return errors.New("speed cannot be negative")
    }
    return nil
}
```

### 入力サニタイゼーション

```go
func sanitizeFileName(filename string) string {
    // パストラバーサル攻撃防止
    filename = filepath.Base(filename)
    // 危険な文字除去
    filename = strings.ReplaceAll(filename, "..", "")
    return filename
}
```

## テスト設計

### テスト戦略

1. **Unit Tests**: 個別コンポーネント
2. **Integration Tests**: システム間連携
3. **Performance Tests**: パフォーマンス検証

```go
func TestUnitAttack(t *testing.T) {
    attacker := NewTestUnit(UnitTypeInfantry, 100, 15, 10)
    target := NewTestUnit(UnitTypeArcher, 70, 12, 5)
    
    damage := attacker.Attack(target)
    
    expectedDamage := 15 - 5 // 攻撃力 - 防御力
    assert.Equal(t, expectedDamage, damage)
    assert.Equal(t, 70-expectedDamage, target.HP)
}
```

### モック設計

```go
type MockDataManager struct {
    units    map[string]UnitTypeConfig
    terrains map[string]TerrainConfig
}

func (m *MockDataManager) GetUnitConfig(unitType string) (UnitTypeConfig, error) {
    if config, exists := m.units[unitType]; exists {
        return config, nil
    }
    return UnitTypeConfig{}, errors.New("unit not found")
}
```

## 監視・ログ設計

### ログレベル

```go
type LogLevel int

const (
    LogLevelDebug LogLevel = iota
    LogLevelInfo
    LogLevelWarn
    LogLevelError
    LogLevelFatal
)
```

### メトリクス収集

```go
type GameMetrics struct {
    FPS           float64
    MemoryUsage   uint64
    ActiveUnits   int
    BattleTime    float64
    UserActions   int
}

func (gm *GameMetrics) Update() {
    gm.FPS = ebiten.ActualFPS()
    gm.MemoryUsage = getMemoryUsage()
    // その他のメトリクス更新...
}
```

## 今後の拡張計画

### Phase 1: 基本機能完成
- [x] 基本戦闘システム
- [x] アニメーションシステム
- [x] フォントシステム
- [x] データ管理システム

### Phase 2: 機能拡張
- [ ] サウンドシステム
- [ ] 配置画面
- [ ] AI改善
- [ ] パーティクル効果

### Phase 3: 高度な機能
- [ ] ネットワーク対戦
- [ ] リプレイシステム
- [ ] MOD対応
- [ ] ランキングシステム

### Phase 4: 最適化・品質向上
- [ ] パフォーマンス最適化
- [ ] メモリ使用量削減
- [ ] バグ修正
- [ ] ユーザビリティ向上

## 技術的負債管理

### 既知の技術的負債

1. **循環インポート**: graphics ↔ game パッケージ
2. **ハードコーディング**: 画面サイズ、色定数
3. **エラーハンドリング**: 一部で不十分
4. **テストカバレッジ**: 低い状態

### 改善計画

```go
// 1. インターフェース分離
type UnitRenderer interface {
    RenderUnit(unit UnitData, position Vector2D) *ebiten.Image
}

// 2. 設定外部化
type DisplayConfig struct {
    Width  int `toml:"width"`
    Height int `toml:"height"`
    Colors ColorScheme `toml:"colors"`
}

// 3. エラー型定義
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}
```
