# 衝突判定システム設計書

## 概要

TinyGocha Battleにユニット間の衝突判定と重なり防止機能を実装。各ユニットに「大きさ」ステータスを追加し、円形の衝突判定により重なりを防止する。

## 実装内容

### 1. ユニットサイズシステム

#### サイズパラメータの追加
```toml
# assets/data/units.toml
[unit_types.infantry]
size = 1.0  # 基本サイズ

[unit_types.heavy_infantry]
size = 1.2  # やや大きめ（重装備）

[unit_types.cavalry]
size = 1.5  # 大きめ（馬に乗っている）
```

#### ユニット種別別サイズ設定
| ユニット種別 | サイズ値 | 衝突半径 | 設計意図 |
|-------------|---------|----------|----------|
| 歩兵 | 1.0 | 10px | 基本サイズ |
| 弓兵 | 1.0 | 10px | 基本サイズ |
| 魔術師 | 1.0 | 10px | 基本サイズ |
| 重装歩兵 | 1.2 | 12px | 重装備で大きめ |
| 騎兵 | 1.5 | 15px | 馬に乗って大きめ |

### 2. 衝突判定システム

#### 衝突半径の計算
```go
func (u *Unit) GetCollisionRadius() float64 {
    baseRadius := 10.0  // 基本半径（ピクセル）
    return baseRadius * u.Size
}
```

#### 衝突判定
```go
func (u *Unit) IsCollidingWith(other *Unit) bool {
    if !u.IsAlive || !other.IsAlive {
        return false
    }
    
    distance := u.Position.Distance(other.Position)
    combinedRadius := u.GetCollisionRadius() + other.GetCollisionRadius()
    
    return distance < combinedRadius
}
```

### 3. 重なり解消システム

#### 衝突解消処理
```go
func (u *Unit) ResolveCollision(other *Unit) {
    distance := u.Position.Distance(other.Position)
    combinedRadius := u.GetCollisionRadius() + other.GetCollisionRadius()
    
    if distance < combinedRadius && distance > 0 {
        // 重なりを解消するために押し出す
        overlap := combinedRadius - distance
        direction := other.Position.Sub(u.Position).Normalize()
        
        // 両方のユニットを半分ずつ押し出す
        pushDistance := overlap * 0.5
        u.Position = u.Position.Sub(direction.Mul(pushDistance))
        other.Position = other.Position.Add(direction.Mul(pushDistance))
    }
}
```

#### 特徴
- **相互押し出し**: 両方のユニットを半分ずつ移動
- **生存ユニットのみ**: 死亡ユニットは衝突判定対象外
- **距離ゼロ対策**: 完全重複時の除算エラー防止

### 4. 戦闘システム統合

#### 衝突処理の統合
```go
func (bm *BattleManager) Update(deltaTime float64) {
    // ... 既存の処理 ...
    
    // Handle unit collisions
    bm.handleCollisions()
    
    // ... 残りの処理 ...
}

func (bm *BattleManager) handleCollisions() {
    allUnits := append(bm.ArmyA.GetAliveUnits(), bm.ArmyB.GetAliveUnits()...)
    
    // Check collisions between all pairs of units
    for i := 0; i < len(allUnits); i++ {
        for j := i + 1; j < len(allUnits); j++ {
            unit1 := allUnits[i]
            unit2 := allUnits[j]
            
            if unit1.IsCollidingWith(unit2) {
                unit1.ResolveCollision(unit2)
            }
        }
    }
}
```

#### 処理順序
1. 軍勢更新（移動・AI）
2. **衝突判定・解消** ← 新規追加
3. 戦闘処理
4. 勝利条件判定

### 5. 移動システムの改善

#### 移動判定の改善
```go
// 従来: 固定値での移動判定
isMoving := u.Position.Distance(u.Target) > 2.0

// 改善: 衝突半径を考慮した移動判定
isMoving := u.Position.Distance(u.Target) > u.GetCollisionRadius()
```

#### 効果
- **自然な停止**: ユニットサイズに応じた適切な停止距離
- **重なり防止**: 目標地点での重なりを防止
- **サイズ差対応**: 大きなユニットは早めに停止

## 技術仕様

### データ構造の変更

#### UnitTypeConfig構造体
```go
type UnitTypeConfig struct {
    Name       string
    HP         int
    Attack     int
    Defense    int
    Speed      float64
    Range      float64
    MagicPower int
    Size       float64  // 新規追加
}
```

#### Unit構造体
```go
type Unit struct {
    // ... 既存フィールド ...
    Size float64  // 新規追加
    // ... 残りのフィールド ...
}
```

### パフォーマンス考慮

#### 計算量
- **衝突判定**: O(n²) - 全ユニット間の総当たり
- **現在のユニット数**: 約24ユニット
- **計算回数**: 約276回/フレーム（60FPS）

#### 最適化の余地
- **空間分割**: 近接ユニットのみ判定
- **フレーム分散**: 複数フレームに分散処理
- **早期終了**: 距離による事前フィルタリング

### 設定可能項目

#### 基本設定
```go
const (
    BaseCollisionRadius = 10.0  // 基本衝突半径
    CollisionPushRatio  = 0.5   // 押し出し比率
)
```

#### カスタマイズ可能
- **基本半径**: ユニットサイズの基準値
- **押し出し比率**: 衝突解消の強度
- **サイズ倍率**: ユニット種別ごとの大きさ

## 動作確認

### テスト項目
- [ ] ユニット作成時のサイズ設定確認
- [ ] 衝突判定の動作確認
- [ ] 重なり解消の動作確認
- [ ] 移動停止距離の確認
- [ ] パフォーマンスの確認

### デバッグ情報
```
Created Unit ID=1, Type=infantry, HP=100/100, Alive=true, Army=0, Size=1.0
Created Unit ID=2, Type=cavalry, HP=90/90, Alive=true, Army=0, Size=1.5
```

### 期待される動作
1. **重なり防止**: ユニット同士が重ならない
2. **自然な移動**: サイズに応じた適切な停止
3. **戦術性向上**: 大きなユニットの配置戦略
4. **視覚的改善**: より自然な戦闘シーン

## 修正履歴

### 問題: 攻撃範囲内に近寄れない問題

#### 原因
- 衝突半径(10px) > 攻撃範囲(歩兵1.5px)
- ユニット同士が重ならないため、攻撃できない距離で停止

#### 解決策

##### 1. 衝突半径の縮小
```go
// 修正前
baseRadius := 10.0

// 修正後  
baseRadius := 3.0  // 基本半径を縮小
```

##### 2. 攻撃判定の改善
```go
// 修正前
if distance > u.Range {
    return 0
}

// 修正後
effectiveRange := u.Range + u.GetCollisionRadius() + target.GetCollisionRadius()
if distance > effectiveRange {
    return 0
}
```

##### 3. AI行動の改善
```go
// 修正前
if distance <= unit.Range && unit.CanAttack() {
    ai.CurrentAction = AIActionAttack
}

// 修正後
effectiveDistance := distance - unit.GetCollisionRadius() - ai.TargetEnemy.GetCollisionRadius()
if effectiveDistance <= unit.Range && unit.CanAttack() {
    ai.CurrentAction = AIActionAttack
}
```

#### 修正後の仕様

##### 新しい衝突半径
| ユニット種別 | サイズ値 | 衝突半径 | 攻撃範囲 | 実効攻撃範囲 |
|-------------|---------|----------|----------|-------------|
| 歩兵 | 1.0 | 3px | 1.5px | 7.5px |
| 弓兵 | 1.0 | 3px | 8.0px | 14.0px |
| 魔術師 | 1.0 | 3px | 10.0px | 16.0px |
| 重装歩兵 | 1.2 | 3.6px | 1.5px | 8.1px |
| 騎兵 | 1.5 | 4.5px | 2.0px | 11.0px |

##### 実効攻撃範囲の計算
```
実効攻撃範囲 = 基本攻撃範囲 + 攻撃者衝突半径 + 目標衝突半径
```

#### 効果
- ✅ ユニット同士の重なりを防止
- ✅ 攻撃範囲内での戦闘が可能
- ✅ サイズ差による戦術的多様性
- ✅ 自然な戦闘距離の維持

---

### 短期改善
- [ ] 衝突エフェクトの追加
- [ ] サイズ可視化（デバッグ用）
- [ ] パフォーマンス監視

### 中期改善
- [ ] 空間分割による最適化
- [ ] 隊形維持との統合
- [ ] 地形との衝突判定

### 長期改善
- [ ] 物理エンジンの導入
- [ ] 複雑な形状の衝突判定
- [ ] 動的サイズ変更

## 設計思想

### 1. シンプルさ重視
- **円形判定**: 計算が簡単で高速
- **基本パラメータ**: 設定が容易
- **直感的**: 理解しやすい動作

### 2. ゲームプレイ向上
- **戦術性**: ユニットサイズによる配置戦略
- **視覚性**: より自然な戦闘シーン
- **快適性**: 重なりによる混乱を防止

### 3. 拡張性確保
- **設定可能**: パラメータによる調整
- **モジュール化**: 独立した衝突システム
- **最適化余地**: 将来の性能改善に対応

## 実装完了

- ✅ ユニットサイズパラメータの追加
- ✅ 衝突判定システムの実装
- ✅ 重なり解消システムの実装
- ✅ 戦闘システムへの統合
- ✅ 移動システムの改善
- ✅ データ構造の更新
- ✅ ビルド・動作確認

---

**実装日**: 2024年6月21日  
**バージョン**: v0.1.1  
**機能**: ユニット衝突判定システム
