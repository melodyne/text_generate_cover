package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func main() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	config := CoverConfig{
		Width:      630,
		Height:     1200,
		Title:      "Go编程完全指南",
		Subtitle:   "构建高效、可靠的现代后端服务",
		Author:     "AI生成 © 2024",
		FontPath:   "fonts/NotoSansSC-Regular.ttf", // 下载中文字体放到这个路径
		OutputPath: "cover_chinese.png",
	}

	fmt.Println("生成中文封面...")

	// 检查字体文件是否存在
	if _, err := os.Stat(config.FontPath); os.IsNotExist(err) {
		fmt.Println("❌ 找不到字体文件:", config.FontPath)
		fmt.Println("请下载中文字体，例如：")
		fmt.Println("1. 思源黑体: https://github.com/adobe-fonts/source-han-sans")
		fmt.Println("2. 下载后放到 fonts/ 目录")
		fmt.Println("3. 重命名为 NotoSansSC-Regular.ttf")
		return
	}

	err := generateChineseCover(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✓ 封面生成成功:", config.OutputPath)
}

type CoverConfig struct {
	Width      int
	Height     int
	Title      string
	Subtitle   string
	Author     string
	FontPath   string
	OutputPath string
}

func generateChineseCover(config CoverConfig) error {
	// 1. 创建基础图像
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	// 2. 绘制随机渐变背景
	drawRandomGradientBackground(img)

	// 3. 加载中文字体
	fontData, err := ioutil.ReadFile(config.FontPath)
	if err != nil {
		return fmt.Errorf("读取字体文件失败: %v", err)
	}

	fnt, err := opentype.Parse(fontData)
	if err != nil {
		return fmt.Errorf("解析字体失败: %v", err)
	}

	// 4. 添加文本（支持自动换行）
	err = addTextWithFontAndWrap(img, config, fnt)
	if err != nil {
		return err
	}

	// 5. 保存图片（去掉所有装饰，保持简洁）
	return saveImage(img, config.OutputPath)
}

func drawRandomGradientBackground(img *image.RGBA) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// 随机选择两种好看的颜色
	colorPairs := [][2]color.RGBA{
		{{41, 128, 185, 255}, {39, 174, 96, 255}},    // 蓝到绿
		{{142, 68, 173, 255}, {230, 126, 34, 255}},   // 紫到橙
		{{231, 76, 60, 255}, {241, 196, 15, 255}},    // 红到黄
		{{52, 152, 219, 255}, {155, 89, 182, 255}},   // 浅蓝到紫
		{{22, 160, 133, 255}, {39, 174, 96, 255}},    // 深绿到浅绿
		{{230, 126, 34, 255}, {231, 76, 60, 255}},    // 橙到红
		{{85, 98, 112, 255}, {78, 205, 196, 255}},    // 灰到青
		{{253, 121, 168, 255}, {120, 119, 198, 255}}, // 粉到紫
		{{30, 60, 114, 255}, {18, 194, 233, 255}},    // 深蓝到亮蓝
		{{255, 107, 107, 255}, {255, 159, 243, 255}}, // 亮红到粉
	}

	pair := colorPairs[rand.Intn(len(colorPairs))]
	startColor := pair[0]
	endColor := pair[1]

	// 垂直渐变（更自然）
	for y := 0; y < height; y++ {
		// 计算垂直比例
		ratio := float64(y) / float64(height)
		r := uint8(float64(startColor.R)*(1-ratio) + float64(endColor.R)*ratio)
		g := uint8(float64(startColor.G)*(1-ratio) + float64(endColor.G)*ratio)
		b := uint8(float64(startColor.B)*(1-ratio) + float64(endColor.B)*ratio)

		lineColor := color.RGBA{r, g, b, 255}

		// 绘制一行
		for x := 0; x < width; x++ {
			img.Set(x, y, lineColor)
		}
	}
}

func addTextWithFontAndWrap(img *image.RGBA, config CoverConfig, fnt *opentype.Font) error {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// 计算可用宽度（留出边距）
	usableWidth := width - 120

	// 1. 绘制标题（自动换行）
	titleFontSize := 58.0 // 稍微增大一点
	if len(config.Title) > 20 {
		titleFontSize = 50.0 // 标题较长时缩小字体
	}

	// 分割标题为多行
	titleLines := wrapText(config.Title, titleFontSize, usableWidth, fnt)
	titleY := height / 3

	// 绘制每行标题
	for i, line := range titleLines {
		y := titleY + i*int(titleFontSize*1.6)
		drawChineseText(img, fnt, titleFontSize, line,
			width/2, y,
			color.RGBA{255, 255, 255, 255}, true)
	}

	// 2. 绘制副标题 - 使用更清晰的颜色
	subtitleY := titleY + len(titleLines)*int(titleFontSize*1.6) + 40
	subtitleLines := wrapText(config.Subtitle, 34.0, usableWidth, fnt)

	for i, line := range subtitleLines {
		lineSpacing := 34.0 * 1.4
		y := subtitleY + i*int(lineSpacing)
		// 使用更白的颜色，alpha值提高
		drawChineseText(img, fnt, 34.0, line,
			width/2, y,
			color.RGBA{245, 245, 245, 230}, true) // 更白，更不透明
	}

	// 3. 绘制作者 - 避免被遮挡
	// 计算作者文本宽度
	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    22.0,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err == nil {
		defer face.Close()

		authorWidth := 0
		for _, ch := range config.Author {
			advance, ok := face.GlyphAdvance(ch)
			if ok {
				authorWidth += advance.Ceil()
			}
		}

		// 确保作者信息不会被截断
		authorX := width - authorWidth - 30
		if authorX < 30 {
			authorX = 30
		}

		drawChineseText(img, fnt, 22.0, config.Author,
			authorX, height-50,
			color.RGBA{255, 255, 255, 200}, false)
	}

	return nil
}

// wrapText 将文本分割为多行以适应宽度
func wrapText(text string, fontSize float64, maxWidth int, fnt *opentype.Font) []string {
	// 创建临时字体面来计算宽度
	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return []string{text}
	}
	defer face.Close()

	var lines []string
	var currentLine strings.Builder
	var currentWidth int

	// 按字符分割（支持中文字符）
	for _, ch := range text {
		runeStr := string(ch)
		advance, ok := face.GlyphAdvance(ch)
		if !ok {
			continue
		}

		charWidth := advance.Ceil()

		// 如果添加当前字符会超出宽度，开始新行
		if currentWidth+charWidth > maxWidth && currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentWidth = 0
		}

		// 如果是标点符号，尽量不换行
		if currentLine.Len() > 0 && isPunctuation(ch) && currentWidth+charWidth > maxWidth*9/10 {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentWidth = 0
		}

		currentLine.WriteString(runeStr)
		currentWidth += charWidth
	}

	// 添加最后一行
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	// 如果没有换行但文本太长，强制分割
	if len(lines) == 1 && len([]rune(text)) > 15 {
		textRunes := []rune(text)
		mid := len(textRunes) / 2
		// 寻找合适的分割点（空格或标点）
		for i := mid; i < len(textRunes); i++ {
			ch := textRunes[i]
			if isPunctuation(ch) || ch == ' ' || ch == '、' || ch == '，' || ch == '。' {
				return []string{string(textRunes[:i+1]), string(textRunes[i+1:])}
			}
		}
		// 没有找到合适的分割点，直接从中间分割
		return []string{string(textRunes[:mid]), string(textRunes[mid:])}
	}

	return lines
}

func isPunctuation(ch rune) bool {
	punctuation := "，。；：！？、（）《》【】「」"
	for _, p := range punctuation {
		if ch == p {
			return true
		}
	}
	return false
}

func drawChineseText(img *image.RGBA, fnt *opentype.Font, size float64,
	text string, x, y int, col color.Color, center bool) {

	// 创建字体面
	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Printf("创建字体面失败: %v", err)
		return
	}
	defer face.Close()

	// 计算文本宽度（如果是居中）
	textWidth := 0
	if center {
		for _, ch := range text {
			advance, ok := face.GlyphAdvance(ch)
			if ok {
				textWidth += advance.Ceil()
			}
		}
		x -= textWidth / 2
	}

	// 绘制文本
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(text)
}

func saveImage(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
