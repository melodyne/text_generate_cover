# text_generate_cover

go语言的文字生成封面图
## 配置
```go
config := CoverConfig{
		Width:      630,
		Height:     1200,
		Title:      "Go编程完全指南：从入门到精通，掌握高性能并发编程与最佳实践",
		Subtitle:   "构建高效、可靠的现代后端服务",
		Author:     "AI生成 © 2024",
		FontPath:   "fonts/NotoSansSC-Regular.ttf", // 下载中文字体放到这个路径
		OutputPath: "cover_chinese.png",
	}
```
## 运行
```shell
go mod tidy
go run .
```
## 特性
* 随机渐变背景
* 支持中文字体

## 案例
<img src="cover_chinese.png" alt="描述" width="300">
