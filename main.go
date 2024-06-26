package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Config struct {
	BgImg            ImageConfig     `json:"bgImg"`
	CompositeItemList []CompositeItem `json:"compositeItemList"`
	OutputImg        OutputImage     `json:"outputImg"`
}

type ImageConfig struct {
	FilePath string `json:"filePath"`
}

type CompositeItem struct {
	CommonParam   CommonParam   `json:"commonParam"`
	SpecificParam SpecificParam `json:"specificParam"`
}

type CommonParam struct {
	Type   string   `json:"type"`
	Depth  int      `json:"depth"`
	Scale  float64  `json:"scale"`
	Pos    Position `json:"pos"`
	Align  string   `json:"align,omitempty"`
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

type OutputImage struct {
	FileName string `json:"fileName"`
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

func drawText(img *image.NRGBA, text string, pos image.Point, col color.Color, fontFace font.Face, align string) {
	lines := strings.Split(text, "\\n")
	for i, line := range lines {
		var x fixed.Int26_6
		textWidth := font.MeasureString(fontFace, line).Ceil()
		switch align {
		case "center":
			x = fixed.I(pos.X) - fixed.I(textWidth/2)
		case "right":
			x = fixed.I(pos.X) - fixed.I(textWidth)
		default:
			x = fixed.I(pos.X)
		}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(col),
			Face: fontFace,
			Dot:  fixed.Point26_6{
				X: x,
				Y: fixed.I(pos.Y) + fixed.I(i*int(fontFace.Metrics().Height.Ceil())),
			},
		}
		d.DrawString(line)
	}
}

func main() {
	configFileName := flag.String("conf", "config.json", "Configuration file name")
	flag.Parse()

	configPath := filepath.Join("conf", *configFileName)

	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		fmt.Println("Error decoding config file:", err)
		return
	}

	bgImgPath := filepath.Join("src", config.BgImg.FilePath)
	bgImg, err := loadImage(bgImgPath)
	if err != nil {
		fmt.Println("Error loading background image:", err)
		return
	}

	outputImg := imaging.Clone(bgImg)

	sort.Slice(config.CompositeItemList, func(i, j int) bool {
		return config.CompositeItemList[i].CommonParam.Depth < config.CompositeItemList[j].CommonParam.Depth
	})

	fontCache := make(map[string]font.Face)

	for _, item := range config.CompositeItemList {
		switch item.CommonParam.Type {
		case "image":
			imgPath := filepath.Join("src", item.SpecificParam.FilePath)
			img, err := loadImage(imgPath)
			if err != nil {
				fmt.Printf("Error loading composite image %s: %v\n", item.SpecificParam.FilePath, err)
				continue
			}

			scaledImg := imaging.Resize(img, int(float64(img.Bounds().Dx())*item.CommonParam.Scale), int(float64(img.Bounds().Dy())*item.CommonParam.Scale), imaging.Lanczos)

			bgWidth := outputImg.Bounds().Dx()
			bgHeight := outputImg.Bounds().Dy()
			posX := int(float64(bgWidth) * (item.CommonParam.Pos.X / 100.0))
			posY := int(float64(bgHeight) * (item.CommonParam.Pos.Y / 100.0))

			posX -= scaledImg.Bounds().Dx() / 2
			posY -= scaledImg.Bounds().Dy() / 2

			outputImg = imaging.Overlay(outputImg, scaledImg, image.Pt(posX, posY), 1.0)

		case "text":
			bgWidth := outputImg.Bounds().Dx()
			bgHeight := outputImg.Bounds().Dy()
			posX := int(float64(bgWidth) * (item.CommonParam.Pos.X / 100.0))
			posY := int(float64(bgHeight) * (item.CommonParam.Pos.Y / 100.0))

			col, err := parseHexColor(item.SpecificParam.Color)
			if err != nil {
				fmt.Printf("Error parsing color %s: %v\n", item.SpecificParam.Color, err)
				return
			}

			fontFace, ok := fontCache[item.SpecificParam.Font]
			if !ok {
				var fontDir string
				switch runtime.GOOS {
				case "windows":
					fontDir = `C:\Windows\Fonts`
				case "darwin":
					fontDir = `/Library/Fonts`
				default:
					fmt.Println("Unsupported operating system")
					return
				}
				fontPath := filepath.Join(fontDir, item.SpecificParam.Font+".ttf")

				if _, err := os.Stat(fontPath); os.IsNotExist(err) {
					fmt.Printf("Font file does not exist: %s\n", fontPath)
					return
				}

				fontBytes, err := os.ReadFile(fontPath)
				if err != nil {
					fmt.Printf("Error loading font %s: %v\n", fontPath, err)
					return
				}
				f, err := opentype.Parse(fontBytes)
				if err != nil {
					fmt.Printf("Error parsing font %s: %v\n", fontPath, err)
					return
				}

				const dpi = 72
				fontFace, err = opentype.NewFace(f, &opentype.FaceOptions{
					Size:    float64(item.CommonParam.Scale * 12),
					DPI:     dpi,
					Hinting: font.HintingFull,
				})
				if err != nil {
					fmt.Printf("Error creating font face %s: %v\n", fontPath, err)
					return
				}

				fontCache[item.SpecificParam.Font] = fontFace
			}

			rgbaImg := outputImg
			drawText(rgbaImg, item.SpecificParam.Text, image.Pt(posX, posY), col, fontFace, item.CommonParam.Align)
			outputImg = rgbaImg
		}
	}

	outputDir := "dst"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}
	outputFilePath := filepath.Join(outputDir, config.OutputImg.FileName)

	if err := imaging.Save(outputImg, outputFilePath); err != nil {
		fmt.Println("Error saving output image:", err)
		return
	}

	fmt.Println("Image composite completed successfully!")
}

