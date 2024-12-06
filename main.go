package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type Config struct {
	BgImg             ImageConfig     `json:"bgImg"`
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

func loadImage(filePath string) (image.Image, error) {
	img, err := imaging.Open(filePath)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func parseHexColor(s string) (color.RGBA, error) {
	c := color.RGBA{A: 0xff}
	if s[0] != '#' {
		return c, fmt.Errorf("invalid color format")
	}
	switch len(s) {
	case 7:
		_, err := fmt.Sscanf(s[1:], "%02x%02x%02x", &c.R, &c.G, &c.B)
		return c, err
	case 4:
		_, err := fmt.Sscanf(s[1:], "%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
		return c, err
	default:
		return c, fmt.Errorf("invalid color format")
	}
}

func main() {
	// 設定ファイルパスのフラグを追加
	configPath := flag.String("conf", "", "設定JSONファイルのパス")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("設定ファイルのパスを指定してください")
		return
	}

	// 設定ファイルを開く
	configFile, err := os.Open(*configPath)
	if err != nil {
		fmt.Println("設定ファイルを開く際のエラー:", err)
		return
	}
	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		fmt.Println("設定ファイルをデコードする際のエラー:", err)
		return
	}

	// 設定ファイルのパスを基に出力ファイルのパスを決定
	relPath, err := filepath.Rel("conf", *configPath)
	if err != nil {
		fmt.Println("設定ファイルの相対パスを取得する際のエラー:", err)
		return
	}
	outputFilePath := filepath.Join("dst", strings.TrimSuffix(relPath, filepath.Ext(relPath))+".png")

	// 出力ディレクトリが存在しない場合は作成
	outputDir := filepath.Dir(outputFilePath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			fmt.Println("出力ディレクトリを作成する際のエラー:", err)
			return
		}
	}

	// 背景画像の読み込み (srcディレクトリを前置)
	bgImage, err := loadImage(filepath.Join("src", config.BgImg.FilePath))
	if err != nil {
		fmt.Println("背景画像を読み込む際のエラー:", err)
		return
	}

	// 出力画像の作成
	dstImage := imaging.Clone(bgImage)

	// 合成アイテムの処理 (srcディレクトリを前置)
	for _, item := range config.CompositeItemList {
		itemImage, err := loadImage(filepath.Join("src", item.SpecificParam.FilePath))
		if err != nil {
			fmt.Println("アイテム画像を読み込む際のエラー:", err)
			return
		}
		// スケールパラメータに基づいて画像をリサイズ
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

		dstImage = imaging.Overlay(dstImage, itemImage, pos, 1.0) // 透過度を1.0に設定して完全に表示
	}

	// 画像を保存
	if err := imaging.Save(dstImage, outputFilePath); err != nil {
		fmt.Println("出力画像を保存する際のエラー:", err)
		return
	}

	fmt.Println("出力画像を保存しました:", outputFilePath)
}
