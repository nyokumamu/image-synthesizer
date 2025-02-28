# 画像合成ツール

このツールは、設定ファイル（JSON）に従い、素材画像を合成して出力画像を生成するGo言語のプログラムです。

## 概要

- `conf`ディレクトリ内のJSON設定ファイルに基づいて、`src`ディレクトリに配置された素材画像を合成し、最終的な画像を`dst`ディレクトリに出力します。
- 背景画像や複数のアイテム画像を指定した位置・スケールで合成できます。

## 事前準備

- Go言語の実行環境が必要です。
  - Goのインストール：[https://go.dev/doc/install](https://go.dev/doc/install)
- 素材画像を`src`ディレクトリに配置しておく必要があります。
- 設定ファイル（JSON）は`conf`ディレクトリ内に配置してください。

## ディレクトリ構成

```
.
├── main.go
├── conf          # JSON設定ファイル配置ディレクトリ
├── src           # 素材画像配置ディレクトリ
└── dst           # 出力画像生成ディレクトリ（自動生成）
```

## 設定ファイルの構造（例）

```json
{
  "bgImg": {
    "filePath": "background.png"
  },
  "output": {
    "size": { "x": 800, "y": 600 }
  },
  "compositeItemList": [
    {
      "commonParam": {
        "type": "image",
        "depth": 1,
        "scale": 0.5,
        "pos": { "x": 50, "y": 50 }
      },
      "specificParam": {
        "filePath": "icon.png"
      }
    }
  ]
}
```

| 項目 | 説明 | 備考 |
|------|------|------|
| bgImg.filePath | 背景画像のパス（srcからの相対パス）| 必須 |
| output.size | 出力画像サイズ（px単位） | 任意（未指定の場合は背景画像サイズ）|
| compositeItemList | 合成アイテムのリスト | 任意（複数指定可能）|

### compositeItemList の設定項目

| 項目 | 説明 | 備考 |
|------|------|------|
| commonParam.type | アイテムのタイプ（現状は `image` のみ対応）| 必須 |
| commonParam.depth | レイヤーの深さ（合成順序）| 数値が小さいほど下層 |
| commonParam.scale | アイテム画像の拡大縮小率 | 1.0で原寸 |
| commonParam.pos | 配置位置（%） | 画像中央を基準として配置 |
| specificParam.filePath | アイテム画像のパス（srcからの相対パス）| 必須 |

## 実行方法

### 単一の設定ファイルを指定する場合

```bash
go run main.go --conf=conf/example.json
```

### 設定ファイルが入ったディレクトリを指定して一括処理する場合

```bash
go run main.go --confDir=conf
```

### コマンドオプション

| オプション名 | 説明 |
|--------------|------|
| `--conf`     | 単一の設定ファイル（JSON）を指定 |
| `--confDir`  | 設定ファイルが含まれるディレクトリを指定（*.json をすべて処理）|

※ `--conf` と `--confDir` は同時に指定できません。

## 出力結果

設定ファイルの名前と同じ名前の画像ファイルが、`dst`ディレクトリ内に生成されます。  
例：  
- 設定ファイル: `conf/example.json` → 出力画像: `dst/example.png`

## ライブラリ依存

このプログラムは、以下の外部ライブラリに依存しています。  
- [imaging](https://github.com/disintegration/imaging)

インストール方法（Goモジュールを利用している場合）:

```bash
go get github.com/disintegration/imaging
```

## ライセンス

本プログラムのライセンスについては、プロジェクト管理者にお問い合わせください。

