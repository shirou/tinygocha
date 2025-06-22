# ゴチャキャラバトル (TinyGocha Battle)

**このプロジェクトは技術検証です**

リアルタイム戦術バトルゲーム - Ebitengine製デモ版

## 特徴

- **リアルタイム戦闘**: 60FPSでの同時行動戦闘
- **アニメーション付きユニット**: プログラム生成による滑らかなアニメーション
- **戦術的要素**: 地形効果、ユニット種別、隊形システム
- **日本語対応**: カスタムフォント対応
- **クロスプラットフォーム**: Windows/Linux/macOS対応

## システム要件

- **OS**: Windows 10+, Linux, macOS
- **メモリ**: 512MB以上
- **ディスク**: 50MB以上

## インストール・実行

### Windows
```bash
# ダウンロード後
tinygocha.exe
```

### Linux/macOS
```bash
# 実行権限を付与
chmod +x tinygocha-linux
./tinygocha-linux
```

## 設定

### フォント設定
日本語表示のため、以下の方法でフォントを設定できます：

1. **デフォルト（推奨）**: 設定不要、内蔵のMPlus1pフォントを使用
2. **カスタムフォント**: `config.toml`でフォントパスを指定

#### config.toml例
```toml
[graphics]
# カスタムフォント（例）
font_path = "C:/Windows/Fonts/msgothic.ttc"  # Windows
# font_path = "/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc"  # Linux
font_size = 16
```

#### 推奨フォント
- **Windows**: MS ゴシック (`msgothic.ttc`)
- **Linux**: Noto Sans CJK (`NotoSansCJK-Regular.ttc`)
- **macOS**: ヒラギノ角ゴシック

### 設定ファイル作成
```bash
# サンプルをコピー
cp config_sample.toml config.toml
```

## 操作方法

### メニュー操作
- **↑↓**: 選択
- **←→**: ステージ変更（設定画面）
- **Enter/Space**: 決定
- **Escape**: 戻る

### 戦闘画面
- **左クリック**: ユニット選択
- **P/Esc**: 一時停止
- **R**: 設定画面に戻る

## ゲームシステム

### ユニット種別
- **歩兵** (□): バランス型、近接戦闘
- **弓兵** (△): 遠距離攻撃、射程が長い
- **魔術師** (◇): 魔法攻撃、高威力・長射程
- **重装歩兵**: 高防御力、移動が遅い
- **騎兵**: 高機動力、突撃攻撃

### 地形効果
- **森**: 移動速度↓、弓兵攻撃力↑
- **山**: 移動速度↓↓、防御力↑、魔術師攻撃力↑
- **平原**: 移動速度↑、全攻撃力↑
- **城塞**: 防御力↑↑、弓兵攻撃力↑
- **街**: 防御力↑、魔術師攻撃力↑

### 戦術要素
- **隊形システム**: リーダー中心の円形隊形
- **リーダーシップ**: リーダー戦死で部隊逃走
- **射程管理**: ユニット選択で射程表示
- **地形活用**: 地形効果を活かした配置

## 開発・ビルド

### 必要環境
- Go 1.24+
- Git

### ビルド
```bash
# 依存関係インストール
make deps

# Windows版ビルド（デフォルト）
make build

# 全プラットフォーム版ビルド
make build-all

# 開発版ビルド
make dev

# 実行（開発用）
make run
```

### プロジェクト構造
```
tinygocha/
├── main.go                    # エントリーポイント
├── config.toml               # 設定ファイル
├── internal/
│   ├── config/              # 設定管理
│   ├── data/                # データローダー
│   ├── game/                # ゲームロジック
│   ├── graphics/            # 描画・アニメーション
│   ├── input/               # 入力処理
│   ├── math/                # 数学ユーティリティ
│   └── scenes/              # シーン管理
├── assets/
│   ├── data/                # ゲームデータ（TOML）
│   ├── images/              # 画像リソース
│   └── sounds/              # サウンドリソース
└── docs/                    # 設計ドキュメント
```

## 技術仕様

- **エンジン**: Ebitengine v2.8.8
- **言語**: Go 1.24
- **データ形式**: TOML
- **フォント**: MPlus1p (内蔵) + カスタム対応
- **解像度**: 1024x768 (リサイズ対応)
- **フレームレート**: 60 FPS

## ライセンス

MIT License

## 開発者

- **エンジン**: [Ebitengine](https://ebitengine.org/)
- **フォント**: [M+ FONTS](https://mplus-fonts.osdn.jp/)
- **設定**: [go-toml](https://github.com/pelletier/go-toml)

---

**注意**: これはデモ版です。完全版では追加機能（配置画面、AI改善、サウンド等）が実装予定です。
