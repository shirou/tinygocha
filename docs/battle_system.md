# 戦闘システム設計

## 概要

ゴチャキャラバトルの戦闘システムは、リアルタイムで進行する戦術的戦闘を実現します。ユニット、グループ、軍勢の階層構造により、複雑な戦術を可能にしています。

## アーキテクチャ

### 階層構造

```
Army (軍勢)
├── Group (部隊)
│   ├── Leader (リーダー)
│   └── Members[] (メンバー)
└── Group (部隊)
    ├── Leader (リーダー)
    └── Members[] (メンバー)

BattleManager
├── ArmyA (軍勢A)
├── ArmyB (軍勢B)
├── 戦闘ロジック
└── 勝利条件判定
```

## Unit（ユニット）

### 基本構造

```go
type Unit struct {
    ID           int
    Type         UnitType
    Name         string
    HP           int
    MaxHP        int
    AttackPower  int
    Defense      int
    Speed        float64
    Range        float64
    MagicPower   int
    Position     math.Vector2D
    Target       math.Vector2D
    IsLeader     bool
    IsAlive      bool
    IsRetreating bool
    GroupID      int
    ArmyID       int
    
    // 戦闘状態
    LastAttackTime float64
    AttackCooldown float64
    
    // アニメーション
    Animation *graphics.AnimationState
}
```

### ユニット種別

```go
type UnitType string

const (
    UnitTypeInfantry UnitType = "infantry"  // 歩兵
    UnitTypeArcher   UnitType = "archer"    // 弓兵
    UnitTypeMage     UnitType = "mage"      // 魔術師
)
```

### 戦闘メカニクス

#### 攻撃処理

```go
func (u *Unit) Attack(target *Unit) int {
    if !u.CanAttack() || !target.IsAlive {
        return 0
    }
    
    // 射程チェック
    distance := u.Position.Distance(target.Position)
    if distance > u.Range {
        return 0
    }
    
    // アニメーション開始
    u.Animation.SetAnimation(graphics.AnimationAttack)
    
    // ダメージ計算
    baseDamage := u.AttackPower
    if u.Type == UnitTypeMage {
        baseDamage += u.MagicPower
    }
    
    // 防御力適用
    damage := baseDamage - target.Defense
    if damage < 1 {
        damage = 1 // 最低ダメージ
    }
    
    // ダメージ適用
    target.TakeDamage(damage)
    
    // クールダウン設定
    u.LastAttackTime = u.AttackCooldown
    
    return damage
}
```

#### 移動処理

```go
func (u *Unit) Update(deltaTime float64) {
    if !u.IsAlive {
        // 戦死アニメーション
        if u.Animation.Type != graphics.AnimationDeath {
            u.Animation.SetAnimation(graphics.AnimationDeath)
        }
        u.Animation.Update(deltaTime)
        return
    }
    
    // アニメーション状態決定
    isMoving := u.Position.Distance(u.Target) > 1.0
    
    if u.LastAttackTime > u.AttackCooldown * 0.7 {
        u.Animation.SetAnimation(graphics.AnimationAttack)
    } else if isMoving {
        u.Animation.SetAnimation(graphics.AnimationWalk)
    } else {
        u.Animation.SetAnimation(graphics.AnimationIdle)
    }
    
    u.Animation.Update(deltaTime)
    
    // 移動処理
    if isMoving {
        direction := u.Target.Sub(u.Position).Normalize()
        movement := direction.Mul(u.Speed * deltaTime)
        u.Position = u.Position.Add(movement)
    }
}
```

## Group（部隊）

### 基本構造

```go
type Group struct {
    ID        int
    Leader    *Unit
    Members   []*Unit
    Formation Formation
    ArmyID    int
    
    targetPosition math.Vector2D
}

type Formation struct {
    Type    FormationType
    Radius  float64
    Spacing float64
}
```

### 隊形システム

#### 円形隊形

```go
func (g *Group) updateCircleFormation() {
    aliveMembers := g.getAliveMembers()
    if len(aliveMembers) == 0 {
        return
    }
    
    angleStep := 2 * math.Pi / float64(len(aliveMembers))
    
    for i, member := range aliveMembers {
        if member.IsRetreating {
            continue
        }
        
        angle := float64(i) * angleStep
        offsetX := math.Cos(angle) * g.Formation.Radius
        offsetY := math.Sin(angle) * g.Formation.Radius
        
        formationPos := g.targetPosition.Add(math.Vector2D{
            X: offsetX,
            Y: offsetY,
        })
        
        member.MoveTo(formationPos)
    }
}
```

### リーダーシップシステム

```go
func (g *Group) handleLeaderDeath() {
    // リーダー戦死時、全メンバーが撤退
    for _, member := range g.Members {
        if member.IsAlive && !member.IsRetreating {
            // 画面端への撤退
            exitPoint := math.Vector2D{X: -100, Y: member.Position.Y}
            if member.ArmyID == 1 {
                exitPoint.X = 1124 // 右端
            }
            member.StartRetreating(exitPoint)
        }
    }
}
```

## Army（軍勢）

### 基本構造

```go
type Army struct {
    ID     int
    Name   string
    Groups []*Group
    Side   int // 0: A軍, 1: B軍
}
```

### 軍勢管理

```go
func (a *Army) GetTotalHealth() float64 {
    units := a.GetAllUnits()
    if len(units) == 0 {
        return 0
    }
    
    totalHealth := 0.0
    for _, unit := range units {
        totalHealth += unit.GetHealthPercentage()
    }
    
    return totalHealth / float64(len(units))
}

func (a *Army) IsDefeated() bool {
    for _, group := range a.Groups {
        if !group.IsDefeated() {
            return false
        }
    }
    return true
}
```

## BattleManager（戦闘管理）

### 基本構造

```go
type BattleManager struct {
    ArmyA        *Army
    ArmyB        *Army
    Stage        data.StageConfig
    TerrainData  data.TerrainConfig
    BattleTime   float64
    TimeLimit    float64
    IsActive     bool
    Winner       int // -1: 未決定, 0: A軍勝利, 1: B軍勝利, 2: 引き分け
    
    nextUnitID   int
}
```

### 戦闘進行

```go
func (bm *BattleManager) Update(deltaTime float64) {
    if !bm.IsActive {
        return
    }
    
    // 戦闘時間更新
    bm.BattleTime += deltaTime
    
    // 軍勢更新
    bm.ArmyA.Update(deltaTime)
    bm.ArmyB.Update(deltaTime)
    
    // 戦闘処理
    bm.processCombat()
    
    // 勝利条件チェック
    bm.checkWinConditions()
}
```

### 戦闘処理

```go
func (bm *BattleManager) processCombat() {
    unitsA := bm.ArmyA.GetAliveUnits()
    unitsB := bm.ArmyB.GetAliveUnits()
    
    // A軍の攻撃
    for _, unitA := range unitsA {
        if !unitA.CanAttack() {
            continue
        }
        
        // 最も近い敵を探索
        var target *Unit
        minDistance := float64(unitA.Range + 1)
        
        for _, unitB := range unitsB {
            distance := unitA.Position.Distance(unitB.Position)
            if distance <= unitA.Range && distance < minDistance {
                target = unitB
                minDistance = distance
            }
        }
        
        // 攻撃実行
        if target != nil {
            unitA.Attack(target)
        }
    }
    
    // B軍の攻撃（同様の処理）
    // ...
}
```

### 勝利条件

```go
func (bm *BattleManager) checkWinConditions() {
    // 時間切れチェック
    if bm.BattleTime >= bm.TimeLimit {
        bm.IsActive = false
        healthA := bm.ArmyA.GetTotalHealth()
        healthB := bm.ArmyB.GetTotalHealth()
        
        if healthA > healthB {
            bm.Winner = 0 // A軍勝利
        } else if healthB > healthA {
            bm.Winner = 1 // B軍勝利
        } else {
            bm.Winner = 2 // 引き分け
        }
        return
    }
    
    // 全滅チェック
    if bm.ArmyA.IsDefeated() && bm.ArmyB.IsDefeated() {
        bm.IsActive = false
        bm.Winner = 2 // 引き分け
    } else if bm.ArmyA.IsDefeated() {
        bm.IsActive = false
        bm.Winner = 1 // B軍勝利
    } else if bm.ArmyB.IsDefeated() {
        bm.IsActive = false
        bm.Winner = 0 // A軍勝利
    }
}
```

## 地形効果システム

### 地形修正適用

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

## プリセット軍勢

### バランス型

```go
func (bm *BattleManager) createBalancedArmy(army *Army, deploymentPoints []Vector2D, dataManager *data.DataManager) {
    groupConfigs := []struct {
        leaderType string
        memberType string
        count      int
    }{
        {"infantry", "infantry", 4},  // 歩兵部隊
        {"archer", "archer", 3},      // 弓兵部隊
        {"mage", "infantry", 2},      // 魔術師+歩兵部隊
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

### 攻撃重視型

```go
func (bm *BattleManager) createOffensiveArmy(army *Army, deploymentPoints []Vector2D, dataManager *data.DataManager) {
    groupConfigs := []struct {
        leaderType string
        memberType string
        count      int
    }{
        {"cavalry", "cavalry", 2},    // 騎兵部隊
        {"archer", "archer", 4},      // 弓兵部隊
        {"infantry", "infantry", 3},  // 歩兵部隊
    }
    // 実装は同様...
}
```

### 防御重視型

```go
func (bm *BattleManager) createDefensiveArmy(army *Army, deploymentPoints []Vector2D, dataManager *data.DataManager) {
    groupConfigs := []struct {
        leaderType string
        memberType string
        count      int
    }{
        {"heavy_infantry", "heavy_infantry", 3}, // 重装歩兵部隊
        {"infantry", "archer", 4},               // 混成部隊
        {"mage", "mage", 2},                     // 魔術師部隊
    }
    // 実装は同様...
}
```

## AI行動

### 基本AI

```go
func (bm *BattleManager) updateAI() {
    // 簡易AI: 最も近い敵を攻撃
    for _, army := range []*Army{bm.ArmyA, bm.ArmyB} {
        for _, group := range army.Groups {
            if group.Leader != nil && group.Leader.IsAlive {
                bm.moveTowardsEnemy(group.Leader)
            }
            
            for _, member := range group.Members {
                if member.IsAlive {
                    bm.moveTowardsEnemy(member)
                }
            }
        }
    }
}

func (bm *BattleManager) moveTowardsEnemy(unit *Unit) {
    var enemyArmy *Army
    if unit.ArmyID == 0 {
        enemyArmy = bm.ArmyB
    } else {
        enemyArmy = bm.ArmyA
    }
    
    // 最も近い敵を探索
    var closestEnemy *Unit
    minDistance := math.MaxFloat64
    
    for _, enemy := range enemyArmy.GetAliveUnits() {
        distance := unit.Position.Distance(enemy.Position)
        if distance < minDistance {
            closestEnemy = enemy
            minDistance = distance
        }
    }
    
    // 敵に向かって移動
    if closestEnemy != nil {
        unit.MoveTo(closestEnemy.Position)
    }
}
```

## ダメージ計算

### 基本ダメージ

```
最終ダメージ = (基本攻撃力 + 魔力ボーナス + 地形ボーナス) - 防御力
最低ダメージ = 1
```

### 地形ボーナス

```
弓兵の森での攻撃 = 基本攻撃力 × 1.2
魔術師の山での攻撃 = (基本攻撃力 + 魔力) × 1.3
```

### クリティカル（将来実装）

```go
func calculateCritical(attacker *Unit, target *Unit) float64 {
    critChance := 0.05 // 5%基本確率
    
    // リーダーボーナス
    if attacker.IsLeader {
        critChance += 0.05
    }
    
    // 種別ボーナス
    switch attacker.Type {
    case UnitTypeArcher:
        critChance += 0.03 // 弓兵は精密攻撃
    }
    
    if rand.Float64() < critChance {
        return 2.0 // 2倍ダメージ
    }
    
    return 1.0
}
```

## パフォーマンス最適化

### 空間分割

```go
type SpatialGrid struct {
    cellSize int
    cells    map[string][]*Unit
}

func (sg *SpatialGrid) GetNearbyUnits(position Vector2D, radius float64) []*Unit {
    // セル座標計算
    cellX := int(position.X) / sg.cellSize
    cellY := int(position.Y) / sg.cellSize
    
    var nearbyUnits []*Unit
    
    // 周辺セルをチェック
    for dx := -1; dx <= 1; dx++ {
        for dy := -1; dy <= 1; dy++ {
            cellKey := fmt.Sprintf("%d,%d", cellX+dx, cellY+dy)
            if units, exists := sg.cells[cellKey]; exists {
                for _, unit := range units {
                    if position.Distance(unit.Position) <= radius {
                        nearbyUnits = append(nearbyUnits, unit)
                    }
                }
            }
        }
    }
    
    return nearbyUnits
}
```

### 戦闘最適化

```go
func (bm *BattleManager) processCombatOptimized() {
    // 空間分割を使用した高速戦闘処理
    grid := NewSpatialGrid(50) // 50ピクセルセル
    
    // 全ユニットをグリッドに登録
    allUnits := append(bm.ArmyA.GetAliveUnits(), bm.ArmyB.GetAliveUnits()...)
    for _, unit := range allUnits {
        grid.AddUnit(unit)
    }
    
    // 各ユニットの戦闘処理
    for _, unit := range allUnits {
        if !unit.CanAttack() {
            continue
        }
        
        // 近隣ユニットのみチェック
        nearbyUnits := grid.GetNearbyUnits(unit.Position, unit.Range)
        
        var target *Unit
        minDistance := unit.Range + 1
        
        for _, nearby := range nearbyUnits {
            if nearby.ArmyID == unit.ArmyID {
                continue // 味方は除外
            }
            
            distance := unit.Position.Distance(nearby.Position)
            if distance <= unit.Range && distance < minDistance {
                target = nearby
                minDistance = distance
            }
        }
        
        if target != nil {
            unit.Attack(target)
        }
    }
}
```

## 今後の拡張

### 予定機能

1. **特殊能力**: ユニット固有のスキル
2. **状態異常**: 毒、麻痺、混乱など
3. **範囲攻撃**: 複数ユニットへの同時攻撃
4. **魔法システム**: 詠唱時間、マナ消費

### 高度なAI

```go
type AIBehavior interface {
    Update(unit *Unit, battlefield *Battlefield) Action
}

type AggressiveAI struct{}
type DefensiveAI struct{}
type SupportAI struct{}

func (ai *AggressiveAI) Update(unit *Unit, battlefield *Battlefield) Action {
    // 積極的に敵を攻撃
    return AttackNearestEnemy
}
```

### 戦術システム

```go
type Tactic struct {
    Name        string
    Conditions  []Condition
    Actions     []Action
    Cooldown    float64
}

type Formation struct {
    Type     FormationType
    Bonuses  map[string]float64
    Penalty  map[string]float64
}
```
