# 画像合成ツール

このツールは、設定ファイルに基づいて画像とテキストを合成することができます。設定ファイルには、合成する画像やテキスト、それらの位置、拡大縮小などの属性が指定されています。

## ディレクトリ構成

```
project-root/
├── conf/
│   └── config.json
├── dst/
├── src/
│   ├── background.png
│   ├── image1.png
│   └── image2.png
└── main.go
```

- `conf/`: 設定ファイルを含むディレクトリ。
- `dst/`: 出力画像が保存されるディレクトリ。
- `src/`: ソース画像を含むディレクトリ。
- `main.go`: メインの Go プログラムファイル。

## 設定ファイル

設定ファイルは以下の構造を持つ JSON ファイルです：

```json
{
    "bgImg": {
        "filePath": "background.png"
    },
    "compositeItemList": [
        {
            "commonParam": {
                "type": "image",
                "depth": 1,
                "scale": 1.0,
                "pos": {
                    "x": 40,
                    "y": 50
                }
            },
            "specificParam": {
                "filePath": "image1.png"
            }
        },
        {
            "commonParam": {
                "type": "image",
                "depth": 2,
                "scale": 1.0,
                "pos": {
                    "x": 60,
                    "y": 50
                }
            },
            "specificParam": {
                "filePath": "image1.png"
            }
        },
        {
            "commonParam": {
                "type": "text",
                "depth": 3,
                "scale": 1.0,
                "pos": {
                    "x": 50,
                    "y": 75
                },
                "align": "center"
            },
            "specificParam": {
                "text": "Hello, World!",
                "font": "Arial Unicode",
                "color": "#FF0000"
            }
        }
    ],
    "outputImg": {
        "fileName": "output.png"
    }
}
```

- `bgImg`: 背景画像の設定。
  - `filePath`: 背景画像のファイルパス。
- `compositeItemList`: 背景画像に合成するアイテムのリスト。
  - `commonParam`: 各アイテムの共通パラメータ。
    - `type`: アイテムのタイプ（`image` または `text`）。
    - `depth`: アイテムの深さ（z-index）。
    - `scale`: アイテムのスケールファクター。
    - `pos`: アイテムの位置。
      - `x`: x座標（パーセンテージ）。
      - `y`: y座標（パーセンテージ）。
    - `align`: テキストの整列（`left`、`center`、`right`）。
  - `specificParam`: 各アイテムの特定パラメータ。
    - `filePath`: 画像のファイルパス（`image` タイプの場合）。
    - `text`: テキスト内容（`text` タイプの場合）。
    - `font`: フォント名（`text` タイプの場合）。
    - `color`: テキストの色（`text` タイプの場合）。
- `outputImg`: 出力画像の設定。
  - `fileName`: 出力画像のファイル名。

## プログラムの実行

1. 背景画像と合成したい画像を `src/` ディレクトリに配置します。
2. 設定ファイルを `conf/` ディレクトリに作成します。
3. 次のコマンドを使用してプログラムを実行します：

```
go run main.go --conf your_config.json
```

`your_config.json` を設定ファイルの名前に置き換えてください。

## 例

`conf/` ディレクトリに `example_config.json` という例の設定ファイルが提供されています。

## フォント

指定されたフォントがシステム上に存在することを確認してください。プログラムは次のディレクトリでフォントを探します：

- Windows: `C:\\Windows\\Fonts`
- macOS: `/Library/Fonts`

フォントが見つからない場合、プログラムはエラーを返します。

## ライセンス

このプロジェクトは MIT ライセンスの下でライセンスされています。詳細は LICENSE ファイルを参照してください。

