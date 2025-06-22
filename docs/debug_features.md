# デバッグ機能仕様書

## 概要

TinyGocha Battleに実装されたデバッグ機能の詳細仕様。開発・テスト・バランス調整に使用。

## 機能一覧

### 1. ターゲット可視化システム

#### 機能概要
選択ユニットがターゲットしている敵ユニットを視覚的に識別可能にする。

#### 実装詳細
```go
// battle_draw_methods.go
func (bs *BattleSceneNew) drawUnit(screen *ebiten.Image, unit *game.Unit, baseColor color.RGBA) {
    // ターゲット判定
    isTargeted := false
    if bs.selectedUnit != nil && bs.selectedUnit.AI != nil && bs.selectedUnit.AI.TargetEnemy == unit {
        isTargeted = true
        baseColor = color.RGBA{255, 255, 0, 255} // 黄色
    }
}
```

#### 視覚効果
- **通常ユニット**: 軍勢色（青/赤）
- **ターゲットユニット**: 明るい黄色（#FFFF00）
- **選択ユニット**: 白い枠線

#### 操作方法
1. 左クリックでユニット選択
2. 選択ユニットのターゲットが黄色で表示
3. 別ユニット選択で表示更新

### 2. ターゲットライン描画

#### 機能概要
選択ユニットから目標ユニットへの方向を白い線で表示。

#### 実装詳細
```go
// battle_draw_methods.go
func (bs *BattleSceneNew) drawTargetLine(screen *ebiten.Image, unit *Unit, target *Unit) {
    // Bresenham's line algorithm
    // 半透明白線で描画
    lineColor := color.RGBA{255, 255, 255, 128}
}
```

#### 視覚効果
- **線の色**: 半透明白（#FFFFFF80）
- **線の太さ**: 3ピクセル
- **描画条件**: 選択ユニットにターゲットが存在する場合

#### 用途
- AI行動方向の確認
- ターゲット変更の検証
- 戦術分析

### 3. 詳細AI情報パネル

#### 機能概要
選択ユニットのAI状態を詳細表示するパネル。

#### 表示項目
```go
// 表示される情報
- AI行動: 現在の行動状態
- 目標: ターゲットIDと距離
- 理想距離: AI設定値
- 攻撃性: AI設定値
- 判断CD: クールダウン状況
```

#### パネル仕様
- **位置**: 画面左下（50, 480）
- **サイズ**: 400×200ピクセル
- **背景**: 半透明黒（#00000080）
- **文字色**: 白・黄・赤（状況に応じて）

#### 表示内容詳細

##### AI行動状態
- **接近**: 敵に向かって移動中
- **後退**: 敵から距離を取る中
- **攻撃**: 射程内で攻撃中
- **保持**: 理想位置で待機
- **待機**: ターゲットなし

##### 目標情報
- **目標あり**: `目標: ID15 距離45.2`
- **目標なし**: `目標: なし`（赤文字）

##### AI設定値
- **理想距離**: ユニット種別による設定値
- **攻撃性**: 0.0-1.0の数値
- **判断CD**: `現在値/最大値`形式

### 4. コンソールデバッグログ

#### 機能概要
ターミナル/コマンドプロンプトにAI行動の詳細ログを出力。

#### ログレベル

##### 戦闘初期化ログ
```
=== Battle Scene OnEnter ===
Stage loaded: 森の戦い
Terrain loaded: 森
Creating armies...
Creating preset army 0 (バランス型)
Army 0 created with 12 units:
  Unit ID=1, Type=infantry, HP=100/100, Alive=true, Army=0
```

##### AI更新ログ
```
AI Update - Army A: 12 units, Army B: 12 units
Unit 1 (Army 0) selecting target from 12 enemies:
  Valid enemies: 12/12
    Enemy ID=13: Distance=800.0, Score=50.0
Unit 1 selected target: ID=13 (score: 50.0)
Unit 1: Target=13, Distance=800.0, Action=接近
```

##### ユニット作成ログ
```
Creating group: Leader=infantry (HP=100), Members=infantry (HP=100), Count=4
Created Unit ID=1, Type=infantry, HP=100/100, Alive=true, Army=0
```

#### 出力条件
- **リーダーユニットのみ**: ログ量を制限
- **重要イベント**: ターゲット変更、行動変更
- **エラー状況**: 設定読み込み失敗など

### 5. ユニット数表示

#### 機能概要
画面上部に各軍勢のユニット数をリアルタイム表示。

#### 実装詳細
```go
// battle_new.go
armyACount := len(bs.battleManager.ArmyA.GetAllUnits())
armyBCount := len(bs.battleManager.ArmyB.GetAllUnits())
debugText := fmt.Sprintf("ユニット数 A:%d B:%d", armyACount, armyBCount)
```

#### 表示仕様
- **位置**: 画面上部中央（200, 40）
- **色**: 黄色（#FFFF00）
- **フォーマット**: `ユニット数 A:12 B:12`

#### 用途
- 軍勢バランスの確認
- ユニット生成の検証
- 戦況把握

### 6. 射程表示システム

#### 機能概要
選択ユニットの攻撃射程を円形で表示。

#### 実装詳細
```go
// battle_draw_methods.go
func (bs *BattleSceneNew) drawRangeIndicator(screen *ebiten.Image, x, y int, range float64) {
    // 半透明円で射程表示
    rangeColor := color.RGBA{255, 255, 255, 64}
}
```

#### 視覚効果
- **色**: 半透明白（#FFFFFF40）
- **形状**: 円形（中心がユニット位置）
- **半径**: ユニットの射程値

## デバッグ操作方法

### 基本操作
1. **ユニット選択**: 左クリック
2. **情報確認**: 選択後に自動表示
3. **ターゲット確認**: 黄色ユニットと白線
4. **ログ確認**: ターミナル/コマンドプロンプト

### 推奨デバッグ手順
1. **戦闘開始**: コンソールで初期化ログ確認
2. **ユニット選択**: リーダーユニットを選択
3. **AI状態確認**: 情報パネルで行動状態確認
4. **ターゲット確認**: 黄色表示と白線確認
5. **移動確認**: ユニットの実際の移動観察

## 開発者向け情報

### デバッグ機能の有効化
```go
// 本番環境では無効化可能
const DEBUG_MODE = true

if DEBUG_MODE {
    // デバッグ情報表示
}
```

### パフォーマンス考慮
- **ログ出力**: リーダーユニットのみに制限
- **描画処理**: 選択時のみ実行
- **更新頻度**: 必要時のみ更新

### カスタマイズ可能項目
- **ログレベル**: 詳細度調整
- **表示色**: 色覚対応
- **パネル位置**: UI配置調整
- **更新頻度**: パフォーマンス調整

## トラブルシューティング

### よくある問題

#### 1. ターゲットが表示されない
- **原因**: AIがターゲットを選択していない
- **確認**: コンソールログで`No valid target found`
- **対策**: 敵軍の存在とAI設定を確認

#### 2. 移動が遅い/止まる
- **原因**: 移動速度設定またはAI判断頻度
- **確認**: 情報パネルの判断CD値
- **対策**: units.tomlの速度値を調整

#### 3. ログが出力されない
- **原因**: リーダーユニット以外を選択
- **確認**: 選択ユニットがリーダーか確認
- **対策**: 各グループのリーダーを選択

### デバッグ情報の読み方

#### AI行動の判断
```
Unit 1: Target=13, Distance=800.0, Action=接近
```
- **正常**: ターゲットありで適切な行動
- **異常**: ターゲットなしまたは不適切な行動

#### スコア計算の確認
```
Enemy ID=13: Distance=800.0, Score=50.0
```
- **高スコア**: 優先ターゲット
- **低スコア**: 低優先度
- **負スコア**: 除外対象

## 今後の拡張予定

### 追加予定機能
- [ ] AI思考過程の可視化
- [ ] 戦闘統計の表示
- [ ] リプレイ機能
- [ ] パフォーマンス監視

### 改善予定項目
- [ ] UI配置の最適化
- [ ] 色覚対応の改善
- [ ] ログフィルタリング機能
- [ ] 設定ファイルでの制御
