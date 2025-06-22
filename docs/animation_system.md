# アニメーションシステム設計

## 概要

ゴチャキャラバトルでは、プログラム生成によるアニメーション付きスプライトシステムを実装しています。画像ファイルを使用せず、リアルタイムで動的にユニットのアニメーションを生成します。

## アーキテクチャ

### コンポーネント構成

```
AnimationState
├── アニメーション種別管理
├── フレーム進行制御
└── 状態遷移管理

SpriteGenerator
├── プログラム描画
├── 形状生成
└── 特殊効果

Unit Integration
├── アニメーション状態保持
├── 自動状態遷移
└── 描画連携
```

## AnimationState

### アニメーション種別

```go
type AnimationType int

const (
    AnimationIdle   AnimationType = iota  // 待機
    AnimationWalk                         // 移動
    AnimationAttack                       // 攻撃
    AnimationDeath                        // 戦死
)
```

### 状態管理

```go
type AnimationState struct {
    Type          AnimationType
    Frame         int     // 現在フレーム
    FrameTime     float64 // フレーム経過時間
    FrameDuration float64 // フレーム持続時間
    TotalFrames   int     // 総フレーム数
    Loop          bool    // ループ再生
    Finished      bool    // 完了フラグ
}
```

### アニメーション特性

| アニメーション | フレーム数 | 持続時間 | ループ | 効果 |
|---------------|-----------|----------|--------|------|
| Idle | 4 | 0.5s | ○ | ゆっくりとした上下動 |
| Walk | 4 | 0.15s | ○ | バウンス効果 |
| Attack | 3 | 0.1s | × | 前方突進、フラッシュ |
| Death | 5 | 0.2s | × | 回転・縮小・フェード |

## SpriteGenerator

### 基本構造

```go
type SpriteGenerator struct {
    cache map[string]*ebiten.Image
}

func (sg *SpriteGenerator) GenerateUnitSprite(
    unitType string, 
    baseColor color.RGBA, 
    isLeader bool, 
    animState *AnimationState
) *ebiten.Image
```

### ユニット形状

#### 歩兵（Infantry）- 四角形
```go
func (sg *SpriteGenerator) drawAnimatedSquare(img *ebiten.Image, centerX, centerY, size int, baseColor color.RGBA, isLeader bool, animState *AnimationState, rotation float64) {
    // アニメーション固有の変形
    var sizeModX, sizeModY int = size, size
    
    switch animState.Type {
    case AnimationWalk:
        // 歩行時の軽い縦圧縮
        if animState.Frame%2 == 0 {
            sizeModY = int(float64(size) * 0.9)
        }
    case AnimationAttack:
        // 攻撃時の前方伸長
        if animState.Frame == 1 {
            sizeModX = int(float64(size) * 1.3)
        }
    }
    
    // 形状描画...
}
```

#### 弓兵（Archer）- 三角形
```go
func (sg *SpriteGenerator) drawAnimatedTriangle(img *ebiten.Image, centerX, centerY, size int, baseColor color.RGBA, isLeader bool, animState *AnimationState, rotation float64) {
    // 攻撃時の縦方向伸長
    heightMod := 1.0
    if animState.Type == AnimationAttack && animState.Frame == 1 {
        heightMod = 1.4
    }
    
    // 三角形描画...
}
```

#### 魔術師（Mage）- ダイヤモンド
```go
func (sg *SpriteGenerator) drawAnimatedDiamond(img *ebiten.Image, centerX, centerY, size int, baseColor color.RGBA, isLeader bool, animState *AnimationState, rotation float64) {
    // 脈動効果
    pulseMod := 1.0
    switch animState.Type {
    case AnimationIdle:
        // 待機時の脈動
        pulseMod = 1.0 + math.Sin(float64(animState.Frame)*math.Pi/2)*0.1
    case AnimationAttack:
        // 攻撃時の光る効果
        if animState.Frame == 1 {
            pulseMod = 1.3
            // 色を明るく
            baseColor.R = uint8(math.Min(255, float64(baseColor.R)*1.2))
        }
    }
    
    // ダイヤモンド描画...
}
```

### アニメーション効果

#### オフセット計算
```go
func (as *AnimationState) GetAnimationOffset() (float64, float64) {
    switch as.Type {
    case AnimationIdle:
        // ゆっくりとした上下動
        bob := math.Sin(float64(as.Frame) * math.Pi / 2) * 1.0
        return 0, bob
        
    case AnimationWalk:
        // 歩行時のバウンス
        bounce := math.Abs(math.Sin(float64(as.Frame) * math.Pi / 2)) * 2.0
        return 0, -bounce
        
    case AnimationAttack:
        // 攻撃時の前方突進
        thrust := 0.0
        if as.Frame == 1 {
            thrust = 3.0
        }
        return thrust, 0
        
    case AnimationDeath:
        // 戦死時の落下
        fall := float64(as.Frame) * 2.0
        return 0, fall
    }
    
    return 0, 0
}
```

#### スケール変更
```go
func (as *AnimationState) GetScaleModifier() float64 {
    switch as.Type {
    case AnimationAttack:
        if as.Frame == 1 {
            return 1.2 // 攻撃時に拡大
        }
    case AnimationDeath:
        // 戦死時に縮小
        return 1.0 - (float64(as.Frame) / float64(as.TotalFrames) * 0.3)
    }
    
    return 1.0
}
```

#### 回転効果
```go
func (as *AnimationState) GetRotationModifier() float64 {
    switch as.Type {
    case AnimationDeath:
        // 戦死時の回転
        return float64(as.Frame) * math.Pi / 8
    }
    
    return 0.0
}
```

## 特殊効果

### 攻撃フラッシュ
```go
func (sg *SpriteGenerator) addAnimationEffects(img *ebiten.Image, centerX, centerY, size int, animState *AnimationState) {
    switch animState.Type {
    case AnimationAttack:
        if animState.Frame == 1 {
            // 黄色いフラッシュ効果
            flashColor := color.RGBA{255, 255, 0, 128}
            for i := 0; i < 3; i++ {
                for angle := 0.0; angle < 2*math.Pi; angle += math.Pi / 4 {
                    x := centerX + int(math.Cos(angle)*float64(size+i+2))
                    y := centerY + int(math.Sin(angle)*float64(size+i+2))
                    img.Set(x, y, flashColor)
                }
            }
        }
    }
}
```

### 戦死フェード
```go
case AnimationDeath:
    // フェード効果
    alpha := uint8(255 * (1.0 - float64(animState.Frame)/float64(animState.TotalFrames)))
    fadeColor := color.RGBA{100, 100, 100, alpha}
    
    // オーバーレイ適用
    for dy := -size-2; dy <= size+2; dy++ {
        for dx := -size-2; dx <= size+2; dx++ {
            img.Set(centerX+dx, centerY+dy, fadeColor)
        }
    }
```

## ユニット統合

### Unit構造体拡張

```go
type Unit struct {
    // 既存フィールド...
    Animation *graphics.AnimationState
}

func NewUnit(...) *Unit {
    return &Unit{
        // 既存初期化...
        Animation: graphics.NewAnimationState(graphics.AnimationIdle),
    }
}
```

### 自動状態遷移

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
    
    // 状態に応じたアニメーション切り替え
    isMoving := u.Position.Distance(u.Target) > 1.0
    
    if u.LastAttackTime > u.AttackCooldown * 0.7 {
        // 攻撃アニメーション
        if u.Animation.Type != graphics.AnimationAttack {
            u.Animation.SetAnimation(graphics.AnimationAttack)
        }
    } else if isMoving {
        // 移動アニメーション
        if u.Animation.Type != graphics.AnimationWalk {
            u.Animation.SetAnimation(graphics.AnimationWalk)
        }
    } else {
        // 待機アニメーション
        if u.Animation.Type != graphics.AnimationIdle {
            u.Animation.SetAnimation(graphics.AnimationIdle)
        }
    }
    
    u.Animation.Update(deltaTime)
}
```

## 描画統合

### 戦闘シーンでの使用

```go
func (bs *BattleSceneNew) drawUnit(screen *ebiten.Image, unit *game.Unit, baseColor color.RGBA) {
    // 体力に応じた色調整
    healthPercent := unit.GetHealthPercentage()
    if healthPercent < 0.5 {
        factor := 0.5 + healthPercent
        baseColor.R = uint8(float64(baseColor.R) * factor)
        baseColor.G = uint8(float64(baseColor.G) * factor)
        baseColor.B = uint8(float64(baseColor.B) * factor)
    }
    
    // アニメーション付きスプライト生成
    sprite := bs.spriteGenerator.GenerateUnitSprite(
        string(unit.Type), 
        baseColor, 
        unit.IsLeader, 
        unit.Animation
    )
    
    // 描画オプション設定
    op := &ebiten.DrawImageOptions{}
    bounds := sprite.Bounds()
    op.GeoM.Translate(-float64(bounds.Dx())/2, -float64(bounds.Dy())/2)
    op.GeoM.Translate(float64(x), float64(y))
    
    screen.DrawImage(sprite, op)
}
```

## パフォーマンス最適化

### スプライトキャッシュ

```go
type SpriteGenerator struct {
    cache map[string]*ebiten.Image
}

func (sg *SpriteGenerator) getCacheKey(unitType string, frame int, isLeader bool) string {
    return fmt.Sprintf("%s_%d_%t", unitType, frame, isLeader)
}
```

### 描画最適化

1. **オフスクリーン描画**: 複雑な効果は事前生成
2. **バッチ処理**: 同種ユニットの一括処理
3. **LOD**: 距離に応じた詳細度調整

## 拡張可能性

### 新アニメーション追加

```go
const (
    AnimationIdle AnimationType = iota
    AnimationWalk
    AnimationAttack
    AnimationDeath
    AnimationCast    // 魔法詠唱
    AnimationBlock   // 防御
    AnimationRetreat // 撤退
)
```

### カスタム効果

```go
type EffectType int

const (
    EffectFlash EffectType = iota
    EffectParticle
    EffectTrail
    EffectAura
)

func (sg *SpriteGenerator) addCustomEffect(img *ebiten.Image, effectType EffectType, params EffectParams)
```

## デバッグ機能

### アニメーション状態表示

```go
func (bs *BattleSceneNew) drawSelectedUnitInfo(screen *ebiten.Image) {
    // アニメーション状態表示
    animText := fmt.Sprintf("状態: %s", bs.getAnimationStateName(unit.Animation.Type))
    bs.textRenderer.DrawText(screen, animText, x, y, color.White)
}

func (bs *BattleSceneNew) getAnimationStateName(animType graphics.AnimationType) string {
    switch animType {
    case graphics.AnimationIdle:   return "待機"
    case graphics.AnimationWalk:   return "移動"
    case graphics.AnimationAttack: return "攻撃"
    case graphics.AnimationDeath:  return "戦死"
    default:                       return "不明"
    }
}
```

## 今後の改善

### 予定機能

1. **パーティクルシステム**: 攻撃時の火花、魔法エフェクト
2. **トレイル効果**: 移動軌跡の表示
3. **オーラ効果**: リーダーや強化ユニットの視覚効果
4. **ダメージ表示**: 数値の飛び出しアニメーション

### 技術的改善

1. **GPU描画**: シェーダーによる高速化
2. **物理演算**: より自然な動きの実現
3. **補間**: フレーム間の滑らかな補間
4. **圧縮**: アニメーションデータの最適化
