package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"github.com/disintegration/imaging"
)

type Config struct {
	BgImg             ImageConfig     `json:"bgImg"`
	Output            OutputConfig    `json:"output"`
	CompositeItemList []CompositeItem `json:"compositeItemList"`
}

type ImageConfig struct {
	FilePath string `json:"filePath"`
}

type CompositeItem struct {
	CommonParam   CommonParam   `json:"commonParam"`
	SpecificParam SpecificParam `json:"specificParam"`
}

type CommonParam struct {
	Type  string   `json:"type"`
	Depth int      `json:"depth"`
	Scale float64  `json:"scale"`
	Pos   Position `json:"pos"`
	Align string   `json:"align,omitempty"`
}

type SpecificParam struct {
	FilePath string `json:"filePath,omitempty"`
	Text     string `json:"text,omitempty"`
	Font     string `json:"font,omitempty"`
	Color    string `json:"color,omitempty"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type OutputConfig struct {
	Size struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"size"`
}

func loadImage(filePath string) (image.Image, error) {
	img, err := imaging.Open(filePath)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func processConfig(configPath string) error {
	// 設定ファイルを開く
	fmt.Println("conf: ", configPath)
	configFile, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("- 設定ファイルを開く際のエラー: %v", err)
	}
	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return fmt.Errorf("- 設定ファイルをデコードする際のエラー: %v", err)
	}

	// 設定ファイルのパスを基に出力ファイルのパスを決定
	relPath, err := filepath.Rel("conf", configPath)
	if err != nil {
		return fmt.Errorf("- 設定ファイルの相対パスを取得する際のエラー: %v", err)
	}
	outputFilePath := filepath.Join("dst", strings.TrimSuffix(relPath, filepath.Ext(relPath))+".png")

	// 出力ディレクトリが存在しない場合は作成
	outputDir := filepath.Dir(outputFilePath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("- 出力ディレクトリを作成する際のエラー: %v", err)
		}
	}

	// 背景画像の読み込み
	bgImage, err := loadImage(filepath.Join("src", config.BgImg.FilePath))
	if err != nil {
		return fmt.Errorf("- 背景画像を読み込む際のエラー: %v", err)
	}

	// 出力画像の作成
	dstImage := imaging.Clone(bgImage)

	// 合成アイテムの処理
	for _, item := range config.CompositeItemList {
		itemImage, err := loadImage(filepath.Join("src", item.SpecificParam.FilePath))
		if err != nil {
			return fmt.Errorf("- アイテム画像を読み込む際のエラー: %v", err)
		}
		// 画像のリサイズ
		if item.CommonParam.Scale != 1.0 {
			itemImage = imaging.Resize(itemImage, int(float64(itemImage.Bounds().Dx())*item.CommonParam.Scale), 0, imaging.Lanczos)
		}

		bgWidth := dstImage.Bounds().Dx()
		bgHeight := dstImage.Bounds().Dy()
		posX := int(float64(bgWidth) * (item.CommonParam.Pos.X / 100.0))
		posY := int(float64(bgHeight) * (item.CommonParam.Pos.Y / 100.0))
		posX -= itemImage.Bounds().Dx() / 2
		posY -= itemImage.Bounds().Dy() / 2
		pos := image.Pt(posX, posY)

		dstImage = imaging.Overlay(dstImage, itemImage, pos, 1.0)
	}

	// リサイズ処理
	if config.Output.Size.X > 0 && config.Output.Size.Y > 0 {
		dstImage = imaging.Resize(dstImage, config.Output.Size.X, config.Output.Size.Y, imaging.Lanczos)
	}

	// 画像を保存
	if err := imaging.Save(dstImage, outputFilePath); err != nil {
		return fmt.Errorf("- 出力画像を保存する際のエラー: %v", err)
	}

	fmt.Println("- 出力画像を保存しました:")
	return nil
}

func main() {
	// フラグの追加
	configPath := flag.String("conf", "", "設定JSONファイルのパス")
	configDir := flag.String("confDir", "", "設定JSONファイルのディレクトリ")
	flag.Parse()

	// 両方が指定された場合はエラー
	if *configPath != "" && *configDir != "" {
		fmt.Println("confオプションとconfDirオプションは同時に指定できません")
		return
	}

	// confのみが指定された場合
	if *configPath != "" {
		if err := processConfig(*configPath); err != nil {
			fmt.Println("エラー:", err)
		}
		return
	}

	// confDirのみが指定された場合
	if *configDir != "" {
		err := filepath.Walk(*configDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".json") {
				if err := processConfig(path); err != nil {
					fmt.Println("- 処理中にエラー:\n", err)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("ディレクトリ処理中のエラー:", err)
		}
		return
	}

	// どちらも指定されていない場合
	fmt.Println("confオプションまたはconfDirオプションのいずれかを指定してください")
}

