# AI行動システム開発ログ

## 概要

TinyGocha BattleのAI行動システムの開発過程と問題解決の記録。

## 開発フェーズ

### フェーズ1: AI行動が動かない問題の発見

#### 問題
- 全ユニットの目標が「なし」状態
- AI行動が全く実行されない
- ユニットが静止したまま

#### 調査アプローチ
1. **デバッグ機能の実装**
2. **ターゲット可視化システム**
3. **詳細ログ出力**

### フェーズ2: デバッグ機能の実装

#### 実装したデバッグ機能

##### 1. ターゲット可視化
```go
// 選択ユニットのターゲットを黄色で表示
if bs.selectedUnit != nil && bs.selectedUnit.AI != nil && bs.selectedUnit.AI.TargetEnemy == unit {
    isTargeted = true
    baseColor = color.RGBA{255, 255, 0, 255} // 黄色
}
```

##### 2. ターゲットライン描画
```go
// 選択ユニットから目標への白い線
func (bs *BattleSceneNew) drawTargetLine(screen *ebiten.Image, unit *Unit, target *Unit) {
    // Bresenham's line algorithm implementation
}
```

##### 3. 詳細AI情報パネル
```go
// AI行動状態の詳細表示
aiText := fmt.Sprintf("AI行動: %s", unit.AI.GetActionName())
targetText := fmt.Sprintf("目標: ID%d 距離%.1f", unit.AI.TargetEnemy.ID, targetDistance)
preferredText := fmt.Sprintf("理想距離: %.1f", unit.AI.PreferredRange)
aggressionText := fmt.Sprintf("攻撃性: %.1f", unit.AI.AggressionLevel)
cooldownText := fmt.Sprintf("判断CD: %.2f/%.2f", unit.AI.LastDecisionTime, unit.AI.DecisionCooldown)
```

##### 4. コンソールデバッグログ
```go
// リーダーユニットのAI行動をコンソール出力
if unit.IsLeader {
    fmt.Printf("AI Update: Unit %d, Enemies: %d\n", unit.ID, len(enemies))
    fmt.Printf("Unit %d: Target=%d, Distance=%.2f, Action=%s\n", 
        unit.ID, ai.TargetEnemy.ID, distance, ai.GetActionName())
}
```

### フェーズ3: 根本原因の特定

#### 問題の発見
コンソール出力で以下が判明：
```
Unit 1 selecting target from 12 enemies:
  Enemy[0]: ID=13, Alive=true, Retreating=false, Pos=(900.0,200.0)
  ...
Unit 1: No valid target found!
```

#### 原因分析
`calculateTargetScore`関数の距離制限が問題：
```go
// 問題のあるコード
if distance > unit.Range * 2.0 {
    return -1.0 // 射程の2倍以上離れている敵は対象外
}
```

- **歩兵の射程**: 1.5
- **射程の2倍**: 3.0
- **実際の敵との距離**: 約800（画面端から端）
- **結果**: 全ての敵が除外される

### フェーズ4: AI行動システムの修正

#### 修正内容

##### 1. ターゲットスコア計算の改善
```go
// 修正後のコード
func (ai *AIBehavior) calculateTargetScore(unit *Unit, enemy *Unit, distance float64) float64 {
    // 距離制限を撤廃
    score := 100.0
    
    // 距離による減点を緩和
    score -= distance * 0.1  // 2.0 → 0.1
    
    // 射程内ボーナス追加
    if distance <= unit.Range {
        score += 100.0
    }
    
    return score
}
```

##### 2. AI判断頻度の向上
```go
// 判断クールダウンの短縮
DecisionCooldown: 0.5 → 0.1  // 5倍高速化
```

##### 3. 移動処理の改善
```go
// より積極的な移動
moveDistance := stdmath.Min(currentDistance - targetDistance, 50.0) // 最大50ピクセル移動
targetPos := unit.Position.Add(direction.Mul(20.0 * intensity)) // 移動距離増加
```

### フェーズ5: 移動速度の大幅改善

#### 問題
- AI行動は動作するが移動が非常に遅い
- ユニット間の速度差が不明確

#### 解決策

##### 1. 基本移動速度の大幅向上
```toml
# units.toml の修正
[unit_types.infantry]
speed = 2.0 → 80.0   # 40倍高速化

[unit_types.cavalry]  
speed = 3.5 → 120.0  # 34倍高速化
```

##### 2. ユニット種別による速度差の明確化
- **騎兵**: 120.0（最高速度）
- **歩兵**: 80.0（標準速度）
- **弓兵**: 70.0（やや遅め）
- **魔術師**: 60.0（遅め）
- **重装歩兵**: 50.0（最も遅い）

##### 3. 地形効果の緩和
```toml
# terrain.toml の修正
[terrain_types.forest]
movement_modifier = 0.7 → 0.9  # 森での減速を緩和

[terrain_types.mountain]
movement_modifier = 0.5 → 0.8  # 山での減速を緩和
```

## 実装結果

### デバッグ機能
- ✅ ターゲット可視化（黄色表示）
- ✅ ターゲットライン描画
- ✅ 詳細AI情報パネル
- ✅ コンソールデバッグログ
- ✅ ユニット数表示

### AI行動システム
- ✅ ターゲット選択の修正
- ✅ 距離制限の撤廃
- ✅ スコア計算の改善
- ✅ 判断頻度の向上

### 移動システム
- ✅ 移動速度の大幅向上（40倍高速化）
- ✅ ユニット種別による速度差
- ✅ 地形効果の最適化
- ✅ AI移動処理の改善

## 技術的学習

### 1. デバッグ駆動開発
問題を可視化することで根本原因を特定できた：
- **仮説**: AIロジックの問題
- **実際**: 距離制限による除外

### 2. 段階的改善
1. **動作確認** → デバッグ機能実装
2. **問題特定** → ログ出力分析
3. **根本修正** → アルゴリズム改善
4. **体験向上** → 速度調整

### 3. バランス調整
ゲームプレイの快適性を重視：
- **リアリズム** < **プレイアビリティ**
- **戦術性** + **レスポンシブ性**

## 今後の改善点

### 短期
- [ ] 攻撃アニメーションの改善
- [ ] 射程表示の最適化
- [ ] AI行動パターンの多様化

### 中期
- [ ] 隊形システムの実装
- [ ] 地形を活用したAI戦術
- [ ] ユニット配置画面の追加

### 長期
- [ ] マルチプレイヤー対応
- [ ] カスタムAI設定
- [ ] リプレイ機能

## 開発時間
- **デバッグ機能実装**: 2時間
- **問題特定・修正**: 1時間
- **移動速度改善**: 1時間
- **ドキュメント化**: 1時間
- **合計**: 5時間

## 使用技術
- **言語**: Go 1.24
- **エンジン**: Ebitengine v2.8.8
- **デバッグ**: fmt.Printf, 可視化
- **データ**: TOML設定ファイル
