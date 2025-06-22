# AI戦闘システム設計

## 概要

ゴチャキャラバトルのAI戦闘システムは、各ユニットが自律的に戦闘行動を取るシステムです。距離ベースの行動決定により、戦術的で自然な戦闘を実現します。

## AI行動原理

### 基本行動サイクル

```
1. 敵探索・選択
2. 距離判定
3. 行動決定
4. 行動実行
```

### 距離ベース戦術

```
射程外 (> Range)     → 敵に接近 (Approach)
射程内 (≤ Range)     → 攻撃実行 (Attack)
近すぎる (< 理想距離) → 後退 (Retreat) ※遠距離ユニットのみ
```

## AIBehavior構造

### 基本構造

```go
type AIBehavior struct {
    TargetEnemy      *Unit    // 攻撃対象
    PreferredRange   float64  // 理想的な戦闘距離
    AggressionLevel  float64  // 攻撃性 (0.0-1.0)
    LastDecisionTime float64  // 最後の判断時刻
    DecisionCooldown float64  // 判断間隔（秒）
    
    CurrentAction    AIAction // 現在の行動
    ActionStartTime  float64  // 行動開始時刻
    ActionDuration   float64  // 行動持続時間
}
```

### AI行動種別

```go
type AIAction int

const (
    AIActionIdle     AIAction = iota // 待機
    AIActionApproach                 // 接近
    AIActionRetreat                  // 後退
    AIActionAttack                   // 攻撃
    AIActionHold                     // 位置保持
)
```

## ユニット種別別AI設定

### 設定パラメータ

| ユニット | 理想距離 | 攻撃性 | 戦術特徴 |
|----------|----------|--------|----------|
| 歩兵 | 1.5 | 0.7 | 積極的接近戦 |
| 弓兵 | 6.0 | 0.5 | 距離保持射撃 |
| 魔術師 | 8.0 | 0.4 | 最大射程活用 |
| 重装歩兵 | 1.5 | 0.8 | 高攻撃性接近 |
| 騎兵 | 2.0 | 0.9 | 最高攻撃性 |

### 理想距離の意味

- **近接系**: 攻撃射程ギリギリ
- **遠距離系**: 射程の60-80%地点（反撃回避）

## 敵選択アルゴリズム

### ターゲット優先度計算

```go
func calculateTargetScore(unit *Unit, enemy *Unit, distance float64) float64 {
    score := 100.0
    
    // 距離による減点（近い敵を優先）
    score -= distance * 2.0
    
    // 敵の体力による加点（体力が少ない敵を優先）
    healthPercent := enemy.GetHealthPercentage()
    score += (1.0 - healthPercent) * 30.0
    
    // リーダーボーナス
    if enemy.IsLeader {
        score += 50.0
    }
    
    // ユニット種別による優先度
    switch enemy.Type {
    case UnitTypeMage:      score += 20.0 // 魔術師を優先
    case UnitTypeArcher:    score += 15.0 // 弓兵を優先
    case UnitTypeInfantry:  score += 10.0
    }
    
    return score
}
```

### 優先度要素

1. **距離** - 近い敵ほど高優先度
2. **体力** - 体力が少ない敵ほど高優先度
3. **リーダー** - リーダーは最優先
4. **脅威度** - 魔術師 > 弓兵 > 歩兵

## 行動決定システム

### 距離判定ロジック

```go
func decideAction(unit *Unit, distance float64) {
    // 攻撃可能距離内かチェック
    if distance <= unit.Range && unit.CanAttack() {
        ai.CurrentAction = AIActionAttack
        return
    }
    
    // 理想的な距離と比較
    if distance > ai.PreferredRange * 1.2 {
        // 遠すぎる場合は接近
        ai.CurrentAction = AIActionApproach
    } else if distance < ai.PreferredRange * 0.8 && ai.isRangedUnit(unit) {
        // 近すぎる場合は後退（遠距離ユニットのみ）
        ai.CurrentAction = AIActionRetreat
    } else if distance <= unit.Range {
        // 射程内だが攻撃できない場合は位置保持
        ai.CurrentAction = AIActionHold
    } else {
        // その他の場合は接近
        ai.CurrentAction = AIActionApproach
    }
}
```

### 判定閾値

- **接近判定**: 理想距離の120%以上
- **後退判定**: 理想距離の80%以下（遠距離ユニットのみ）
- **攻撃判定**: 射程内 + 攻撃可能

## 移動制御

### 接近移動

```go
func moveTowardsTarget(unit *Unit, intensity float64) {
    direction := ai.TargetEnemy.Position.Sub(unit.Position).Normalize()
    
    // 理想的な距離まで移動
    targetDistance := ai.PreferredRange * 0.9 // 理想距離の90%地点
    moveDistance := unit.Position.Distance(ai.TargetEnemy.Position) - targetDistance
    
    if moveDistance > 0 {
        targetPos := unit.Position.Add(direction.Mul(moveDistance * intensity))
        unit.MoveTo(targetPos)
    }
}
```

### 後退移動

```go
func moveAwayFromTarget(unit *Unit, intensity float64) {
    direction := unit.Position.Sub(ai.TargetEnemy.Position).Normalize()
    
    // 理想的な距離まで後退
    currentDistance := unit.Position.Distance(ai.TargetEnemy.Position)
    targetDistance := ai.PreferredRange * 1.1 // 理想距離の110%地点
    moveDistance := targetDistance - currentDistance
    
    if moveDistance > 0 {
        targetPos := unit.Position.Add(direction.Mul(moveDistance * intensity))
        
        // 画面外に出ないようにクランプ
        targetPos.X = math.Max(50, math.Min(974, targetPos.X))
        targetPos.Y = math.Max(100, math.Min(700, targetPos.Y))
        
        unit.MoveTo(targetPos)
    }
}
```

### 境界制御

- **X軸**: 50 ≤ x ≤ 974
- **Y軸**: 100 ≤ y ≤ 700

## 隊形システム統合

### リーダー追従

```go
func (g *Group) Update(deltaTime float64) {
    // リーダーを最初に更新
    g.Leader.Update(deltaTime)
    
    // リーダーの移動に応じて隊形目標を更新
    if g.Leader.Position.Distance(g.Leader.Target) > 5.0 {
        g.targetPosition = g.Leader.Target  // 移動中は目標位置
    } else {
        g.targetPosition = g.Leader.Position // 停止中は現在位置
    }
    
    // 隊形更新
    g.updateFormation()
    
    // メンバー更新
    for _, member := range g.Members {
        if member.IsAlive {
            member.Update(deltaTime)
        }
    }
}
```

### 隊形維持

- **リーダー**: AI行動に従って自由移動
- **メンバー**: リーダー中心の円形隊形を維持
- **追従距離**: 5.0ピクセル以上の移動で追従開始

## パフォーマンス最適化

### 判断間隔制御

```go
// 判断クールダウンチェック
ai.LastDecisionTime += deltaTime
if ai.LastDecisionTime < ai.DecisionCooldown {
    return // 0.5秒間隔で判断
}
```

### 計算負荷軽減

1. **判断頻度制限**: 0.5秒間隔
2. **射程外除外**: 射程の2倍以上の敵は対象外
3. **生存チェック**: 死亡ユニットは処理スキップ

## 戦術的特徴

### ユニット種別戦術

#### 歩兵系（Infantry, Heavy Infantry）
- **基本戦術**: 積極的接近戦
- **理想距離**: 攻撃射程ギリギリ
- **特徴**: 高い攻撃性、後退しない

#### 弓兵（Archer）
- **基本戦術**: 距離保持射撃
- **理想距離**: 射程の75%地点
- **特徴**: 敵が近づくと後退、射程を活かす

#### 魔術師（Mage）
- **基本戦術**: 最大射程活用
- **理想距離**: 射程の80%地点
- **特徴**: 最も慎重、最大射程から攻撃

#### 騎兵（Cavalry）
- **基本戦術**: 高速接近攻撃
- **理想距離**: やや長めの接近距離
- **特徴**: 最高攻撃性、機動力活用

### 集団戦術

#### リーダーシップ効果
- **リーダー戦死**: 部隊全体が撤退
- **リーダー優先**: 敵リーダーを優先攻撃
- **隊形維持**: リーダー中心の円形隊形

#### 地形活用
- **森**: 弓兵が攻撃力ボーナス
- **山**: 魔術師が攻撃力ボーナス
- **平原**: 全ユニットが機動力ボーナス

## デバッグ・監視

### AI状態表示

```go
func (ai *AIBehavior) GetActionName() string {
    switch ai.CurrentAction {
    case AIActionIdle:     return "待機"
    case AIActionApproach: return "接近"
    case AIActionRetreat:  return "後退"
    case AIActionAttack:   return "攻撃"
    case AIActionHold:     return "保持"
    default:               return "不明"
    }
}
```

### 情報表示項目

- **AI行動**: 現在の行動状態
- **目標距離**: ターゲットまでの距離
- **理想距離**: ユニットの理想戦闘距離
- **攻撃性**: ユニットの攻撃性レベル

## 今後の拡張

### 高度なAI行動

```go
type AdvancedAI struct {
    Personality  AIPersonality // 性格（攻撃的、守備的、支援的）
    Tactics      []Tactic      // 使用可能戦術
    Memory       AIMemory      // 戦闘記憶
    Cooperation  float64       // 協調性
}
```

### 戦術システム

```go
type Tactic struct {
    Name        string
    Conditions  []Condition
    Actions     []Action
    Cooldown    float64
    Priority    int
}
```

### 学習AI

```go
type LearningAI struct {
    Experience   map[string]float64 // 経験値
    Adaptability float64            // 適応性
    LearningRate float64            // 学習率
}
```

## 設定カスタマイズ

### AI設定ファイル

```toml
# assets/data/ai_config.toml

[ai_settings]
decision_cooldown = 0.5
max_target_distance = 2.0

[unit_ai.infantry]
preferred_range = 1.5
aggression_level = 0.7
retreat_enabled = false

[unit_ai.archer]
preferred_range = 6.0
aggression_level = 0.5
retreat_enabled = true
retreat_threshold = 0.8
```

### 難易度調整

```toml
[difficulty.easy]
ai_reaction_time = 1.0
ai_accuracy = 0.8

[difficulty.normal]
ai_reaction_time = 0.5
ai_accuracy = 1.0

[difficulty.hard]
ai_reaction_time = 0.3
ai_accuracy = 1.2
```

## トラブルシューティング

### よくある問題

1. **ユニットが動かない**
   - AI初期化の確認
   - 敵ユニットの存在確認
   - 判断クールダウンの確認

2. **不自然な動き**
   - 理想距離の調整
   - 攻撃性レベルの調整
   - 境界制御の確認

3. **パフォーマンス問題**
   - 判断間隔の延長
   - 対象範囲の制限
   - 不要な計算の削除

### デバッグ方法

```go
// AI状態のログ出力
log.Printf("Unit %d: Action=%s, Target=%v, Distance=%.2f", 
    unit.ID, ai.GetActionName(), ai.TargetEnemy, distance)
```

---

**実装状況**: ✅ 完成  
**最終更新**: 2024年6月21日  
**バージョン**: v1.0.0
