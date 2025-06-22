# フォントシステム設計

## 概要

ゴチャキャラバトルでは、日本語表示を確実に行うためのフォント管理システムを実装しています。デフォルトで内蔵のMPlus1pフォントを使用し、必要に応じてカスタムフォントを読み込むことができます。

## アーキテクチャ

### コンポーネント構成

```
FontManager
├── デフォルトフォント管理
├── カスタムフォント読み込み
└── フォールバック機能

TextRenderer
├── 統一テキスト描画API
├── サイズ・色指定
└── 特殊効果（影付きなど）

Config System
├── TOML設定ファイル
├── フォント設定
└── 動的読み込み
```

## FontManager

### 基本機能

```go
type FontManager struct {
    defaultFont *text.GoTextFace
    fonts       map[string]*text.GoTextFace
}
```

#### 主要メソッド

- `LoadDefaultFont(size float64)`: MPlus1pフォント読み込み
- `LoadFontFromFile(path, size, name)`: カスタムフォント読み込み
- `GetDefaultFont()`: デフォルトフォント取得
- `CreateFontVariant(name, size)`: サイズ違いフォント生成

### フォント読み込み戦略

1. **デフォルトフォント**: 内蔵MPlus1pを使用
2. **カスタムフォント**: 指定パスから読み込み
3. **フォールバック**: 失敗時は自動的にデフォルトに切り替え

```go
// フォント読み込み例
func (fm *FontManager) LoadFontFromFile(fontPath string, size float64, name string) error {
    if fontPath == "" {
        return fm.LoadDefaultFont(size)
    }
    
    // ファイル存在チェック
    if _, err := os.Stat(fontPath); os.IsNotExist(err) {
        log.Printf("Font file not found: %s, using default font", fontPath)
        return fm.LoadDefaultFont(size)
    }
    
    // フォント読み込み処理...
}
```

## TextRenderer

### 描画機能

```go
type TextRenderer struct {
    fontManager *FontManager
}
```

#### 描画メソッド

- `DrawText()`: 基本テキスト描画
- `DrawTextWithSize()`: サイズ指定描画
- `DrawCenteredText()`: 中央揃え描画
- `DrawTextWithShadow()`: 影付きテキスト
- `MeasureText()`: テキストサイズ測定

### 使用例

```go
// 基本描画
textRenderer.DrawText(screen, "ゴチャキャラバトル", 100, 50, color.White)

// サイズ指定
textRenderer.DrawTextWithSize(screen, "タイトル", 200, 100, color.White, 32)

// 影付き
textRenderer.DrawTextWithShadow(screen, "選択項目", 150, 200, 
    color.RGBA{52, 152, 219, 255}, color.RGBA{0, 0, 0, 128})
```

## 設定システム

### Config構造

```go
type Config struct {
    Graphics GraphicsConfig `toml:"graphics"`
    Audio    AudioConfig    `toml:"audio"`
    Game     GameConfig     `toml:"game"`
}

type GraphicsConfig struct {
    FontPath     string  `toml:"font_path"`
    FontSize     int     `toml:"font_size"`
    UIScale      float64 `toml:"ui_scale"`
    ShowFPS      bool    `toml:"show_fps"`
    VSync        bool    `toml:"vsync"`
}
```

### 設定ファイル例

```toml
[graphics]
# フォントファイルのパス（空の場合はデフォルトのMPlus1pを使用）
font_path = ""
font_size = 16
ui_scale = 1.0
show_fps = false
vsync = true

[audio]
master_volume = 0.8
sfx_volume = 0.7
bgm_volume = 0.6
enabled = true

[game]
language = "ja"
auto_save = true
show_tutorial = true
```

## フォント対応

### デフォルトフォント

- **フォント**: MPlus1pRegular_ttf
- **ソース**: `github.com/hajimehoshi/ebiten/v2/examples/resources/fonts`
- **特徴**: 日本語完全対応、ゲーム用に最適化

### カスタムフォント

#### Windows推奨
```toml
font_path = "C:/Windows/Fonts/msgothic.ttc"  # MS ゴシック
```

#### Linux推奨
```toml
font_path = "/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc"
```

#### macOS推奨
```toml
font_path = "/System/Library/Fonts/ヒラギノ角ゴシック W3.ttc"
```

## エラーハンドリング

### フォント読み込みエラー

1. **ファイル不存在**: デフォルトフォントに自動切り替え
2. **パース失敗**: ログ出力後、デフォルトフォントを使用
3. **メモリ不足**: エラーログ出力、アプリケーション継続

### ログ出力例

```
Font loaded successfully: C:/Windows/Fonts/msgothic.ttc
Font file not found: invalid_path.ttf, using default font
Default font (MPlus1p) loaded successfully
```

## パフォーマンス考慮

### フォントキャッシュ

- 読み込んだフォントはメモリにキャッシュ
- サイズ違いフォントは動的生成
- 不要なフォントの自動解放

### 描画最適化

- テキスト測定結果のキャッシュ
- 描画オプションの再利用
- バッチ描画による高速化

## 実装詳細

### 初期化フロー

```go
// 1. 設定読み込み
config := config.LoadConfig("config.toml")

// 2. フォントマネージャー作成
fontManager := graphics.NewFontManager()

// 3. フォント読み込み
if config.Graphics.FontPath != "" {
    fontManager.LoadFontFromFile(config.Graphics.FontPath, fontSize, "default")
} else {
    fontManager.LoadDefaultFont(fontSize)
}

// 4. テキストレンダラー作成
textRenderer := graphics.NewTextRenderer(fontManager)
```

### シーン統合

各シーンでTextRendererを使用：

```go
// タイトルシーン
func NewTitleScene(sceneManager *SceneManager, textRenderer *graphics.TextRenderer) *TitleScene

// 戦闘シーン
func NewBattleSceneNew(sceneManager *SceneManager, dataManager *data.DataManager, textRenderer *graphics.TextRenderer) *BattleSceneNew
```

## トラブルシューティング

### よくある問題

1. **日本語が表示されない**
   - フォントパスの確認
   - ファイル権限の確認
   - ログでエラー内容を確認

2. **フォントが読み込まれない**
   - ファイル形式の確認（TTF/OTF対応）
   - パスの区切り文字確認（Windows: `\` または `/`）

3. **文字が崩れる**
   - フォントサイズの調整
   - UIスケールの調整

### デバッグ方法

```bash
# ログ出力でフォント読み込み状況を確認
./tinygocha 2>&1 | grep -i font

# 設定ファイルの構文チェック
# TOMLパーサーでエラー確認
```

## 今後の拡張

### 予定機能

1. **多言語対応**: 英語・中国語・韓国語
2. **フォント自動検出**: システムフォントの自動選択
3. **フォント品質設定**: アンチエイリアス・ヒンティング調整
4. **動的フォント切り替え**: ゲーム中の設定変更

### 技術的改善

1. **フォントストリーミング**: 大きなフォントファイルの部分読み込み
2. **GPU描画**: ハードウェアアクセラレーション対応
3. **フォント圧縮**: 内蔵フォントのサイズ最適化
