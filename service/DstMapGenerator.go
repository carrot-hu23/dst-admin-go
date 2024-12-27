package service

//import (
//	"encoding/base64"
//	"fmt"
//	"image"
//	"image/color"
//	"image/png"
//	"io/ioutil"
//	"math"
//	"os"
//	"regexp"
//)
//
//// Color 定义一个RGB颜色
//type Color struct {
//	R, G, B uint8
//	Name    string
//}
//
//// DSTMapGenerator 存储地图生成器的状态
//type DSTMapGenerator struct {
//	tileColors   map[int]Color
//	entityColors map[string]color.Color
//}
//
//// NewDSTMapGenerator 创建新的地图生成器实例
//func NewDSTMapGenerator() *DSTMapGenerator {
//	g := &DSTMapGenerator{
//		tileColors:   make(map[int]Color),
//		entityColors: make(map[string]color.Color),
//	}
//
//	// 初始化地形颜色映射
//	g.tileColors = map[int]Color{
//		0:  {255, 252, 146, "沙漠地皮"},
//		1:  {34, 139, 34, "平原"},
//		2:  {81, 29, 195, "沼泽"},
//		3:  {210, 180, 140, "沙地"},
//		4:  {152, 251, 152, "水中木"},
//		5:  {169, 169, 169, "岩石"},
//		6:  {36, 213, 75, "草地"},
//		7:  {160, 82, 215, "沼泽------"},
//		8:  {42, 42, 42, "海边缘"},
//		9:  {143, 140, 21, "草原"},
//		10: {218, 165, 32, "平原"},
//		11: {0, 206, 209, "浅水域"},
//		12: {17, 81, 165, "中水域"},
//		13: {148, 0, 211, "魔法地形"},
//		14: {0, 85, 252, "浅水区"},
//		15: {255, 0, 0, "危险地形"},
//		16: {10, 85, 252, "中水域"},
//		17: {72, 191, 30, "草地"},
//		18: {0, 21, 165, "深水域"},
//		19: {152, 251, 152, "草甸"},
//		20: {0, 21, 165, "深海"},
//		21: {88, 217, 212, "水中木"},
//		22: {189, 183, 107, "荒漠"},
//		23: {247, 171, 94, "桦树林"},
//		24: {255, 153, 71, "森林"},
//		25: {82, 153, 71, "森林"},
//		26: {200, 238, 200, "月岛"},
//		27: {0, 85, 252, "浅海"},
//		28: {125, 121, 121, "海礁石"},
//		29: {222, 184, 135, "荒原"},
//		30: {234, 152, 84, "红树林"},
//		31: {112, 128, 144, "山地"},
//		32: {176, 196, 222, "高山"},
//		33: {176, 196, 222, "高山"},
//		34: {255, 196, 222, "高山"},
//	}
//
//	// 初始化实体颜色映射
//	g.entityColors = map[string]color.Color{
//		"evergreen":     color.RGBA{0, 100, 0, 255},
//		"deciduoustree": color.RGBA{34, 139, 34, 255},
//		"rock1":         color.RGBA{128, 128, 128, 255},
//		"rock2":         color.RGBA{169, 169, 169, 255},
//		"sapling":       color.RGBA{144, 238, 144, 255},
//		"grass":         color.RGBA{124, 252, 0, 255},
//		"berrybush":     color.RGBA{255, 0, 0, 255},
//		"berrybush2":    color.RGBA{220, 20, 60, 255},
//		"spiderden":     color.RGBA{128, 0, 128, 255},
//		"pond":          color.RGBA{0, 191, 255, 255},
//		"rabbithole":    color.RGBA{139, 69, 19, 255},
//		"cave":          color.RGBA{105, 105, 105, 255},
//		"reeds":         color.RGBA{189, 183, 107, 255},
//		"marsh_tree":    color.RGBA{85, 107, 47, 255},
//		"marsh_bush":    color.RGBA{107, 142, 35, 255},
//		"carrot":        color.RGBA{255, 140, 0, 255},
//		"flower":        color.RGBA{255, 192, 203, 255},
//		"wormhole":      color.RGBA{138, 43, 226, 255},
//		"pighouse":      color.RGBA{255, 160, 122, 255},
//		"mound":         color.RGBA{160, 82, 45, 255},
//		"ruins":         color.RGBA{119, 136, 153, 255},
//		"fireflies":     color.RGBA{255, 255, 0, 255},
//	}
//
//	return g
//}
//
//// ReadSaveFile 读取存档文件
//func (g *DSTMapGenerator) ReadSaveFile(filePath string) (string, error) {
//	content, err := ioutil.ReadFile(filePath)
//	if err != nil {
//		return "", fmt.Errorf("读取文件失败: %v", err)
//	}
//
//	// 使用正则表达式提取地图数据
//	re := regexp.MustCompile(`tiles="([^"]+)"`)
//	matches := re.FindSubmatch(content)
//	if len(matches) < 2 {
//		return "", fmt.Errorf("无法在存档中找到地图数据")
//	}
//
//	base64Data := string(matches[1])
//	fmt.Printf("提取的base64数据长度: %d\n", len(base64Data))
//	return base64Data, nil
//}
//
//// DecodeMapData 解码地图数据
//func (g *DSTMapGenerator) DecodeMapData(tilesBase64 string) ([]int, error) {
//	tileBytes, err := base64.StdEncoding.DecodeString(tilesBase64)
//	if err != nil {
//		return nil, fmt.Errorf("base64解码失败: %v", err)
//	}
//
//	// 处理文件头
//	dataStart := 0
//	if len(tileBytes) > 5 && string(tileBytes[:5]) == "VRSTN" {
//		dataStart = 5
//		for dataStart < len(tileBytes) && tileBytes[dataStart] == 0 {
//			dataStart++
//		}
//	}
//
//	tileBytes = tileBytes[dataStart:]
//	if len(tileBytes)%2 != 0 {
//		tileBytes = tileBytes[:len(tileBytes)-1]
//	}
//
//	// 解码tile IDs
//	tileIds := make([]int, 0, len(tileBytes)/2)
//	for i := 0; i < len(tileBytes); i += 2 {
//		if i+1 >= len(tileBytes) {
//			break
//		}
//		tileId := (int(tileBytes[i+1]) << 8) | int(tileBytes[i])
//		tileId = tileId % 31 // 映射到0-33范围
//		tileIds = append(tileIds, tileId)
//	}
//
//	// 打印分布情况
//	tileCounts := make(map[int]int)
//	for _, id := range tileIds {
//		tileCounts[id]++
//	}
//
//	total := len(tileIds)
//	fmt.Println("\n瓦片分布情况:")
//	for id, count := range tileCounts {
//		percentage := float64(count) / float64(total) * 100
//		if percentage > 0.01 {
//			fmt.Printf("瓦片ID %d: %d 个 (%.2f%%)\n", id, count, percentage)
//		}
//	}
//
//	return tileIds, nil
//}
//
//// CreateMapImage 创建地图图像
//func (g *DSTMapGenerator) CreateMapImage(tileIds []int, width, height int) *image.RGBA {
//	img := image.NewRGBA(image.Rect(0, 0, width, height))
//
//	// 填充地形颜色
//	for y := 0; y < height; y++ {
//		for x := 0; x < width; x++ {
//			idx := y*width + x
//			if idx >= len(tileIds) {
//				continue
//			}
//
//			tileId := tileIds[idx]
//			tileColor, exists := g.tileColors[tileId]
//			if !exists {
//				// 使用默认颜色 (浅绿色)
//				tileColor = Color{144, 238, 144, "未知地形"}
//			}
//
//			// 在Go中，图像坐标是从左到右的，所以需要翻转X坐标
//			flippedX := width - x - 1
//			img.Set(flippedX, y, color.RGBA{
//				R: tileColor.R,
//				G: tileColor.G,
//				B: tileColor.B,
//				A: 255,
//			})
//		}
//	}
//
//	return img
//}
//
//// GenerateMap 生成完整的地图
//func (g *DSTMapGenerator) GenerateMap(saveFilePath, outputPath string) error {
//	// 读取并解码地图数据
//	tilesBase64, err := g.ReadSaveFile(saveFilePath)
//	if err != nil {
//		return err
//	}
//
//	tileIds, err := g.DecodeMapData(tilesBase64)
//	if err != nil {
//		return err
//	}
//
//	// 计算地图尺寸
//	width := int(math.Sqrt(float64(len(tileIds))))
//	if width%2 != 0 {
//		width--
//	}
//	height := len(tileIds) / width
//
//	fmt.Printf("计算得到的地图尺寸: %dx%d\n", width, height)
//	fmt.Printf("总瓦片数: %d\n", len(tileIds))
//	fmt.Printf("实际使用面积: %d\n", width*height)
//
//	// 创建地图图像
//	img := g.CreateMapImage(tileIds, width, height)
//
//	// 保存图像
//	f, err := os.Create(outputPath)
//	if err != nil {
//		return fmt.Errorf("创建输出文件失败: %v", err)
//	}
//	defer f.Close()
//
//	if err := png.Encode(f, img); err != nil {
//		return fmt.Errorf("保存图像失败: %v", err)
//	}
//
//	fmt.Printf("地图已保存到: %s\n", outputPath)
//	return nil
//}
//
//// WalrusHut_Plains
//func main() {
//	generator := NewDSTMapGenerator()
//	err := generator.GenerateMap(
//		"save/session/C7B8F320B11E4E9E/0000000023",
//		"dst_map.png",
//	)
//	if err != nil {
//		fmt.Printf("生成地图时出错: %v\n", err)
//		os.Exit(1)
//	}
//}

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
)

// Color 定义一个RGB颜色
type Color struct {
	R, G, B uint8
	Name    string
}

// DSTMapGenerator 存储地图生成器的状态
type DSTMapGenerator struct {
	tileColors   map[int]Color
	entityColors map[string]color.Color
}

// NewDSTMapGenerator 创建新的地图生成器实例
func NewDSTMapGenerator() *DSTMapGenerator {
	g := &DSTMapGenerator{
		tileColors:   make(map[int]Color),
		entityColors: make(map[string]color.Color),
	}

	// 初始化地形颜色映射
	g.tileColors = map[int]Color{
		0:  {255, 252, 146, "沙漠地皮"},
		1:  {34, 139, 34, "平原"},
		2:  {81, 29, 195, "沼泽"},
		3:  {210, 180, 140, "沙地"},
		4:  {152, 251, 152, "水中木"},
		5:  {169, 169, 169, "岩石"},
		6:  {36, 213, 75, "草地"},
		7:  {160, 82, 215, "沼泽------"},
		8:  {42, 42, 42, "海边缘"},
		9:  {143, 140, 21, "草原"},
		10: {218, 165, 32, "平原"},
		11: {0, 206, 209, "浅水域"},
		12: {17, 81, 165, "中水域"},
		13: {148, 0, 211, "魔法地形"},
		14: {0, 85, 252, "浅水区"},
		15: {255, 0, 0, "危险地形"},
		16: {10, 85, 252, "中水域"},
		17: {72, 191, 30, "草地"},
		18: {0, 21, 165, "深水域"},
		19: {152, 251, 152, "草甸"},
		20: {0, 21, 165, "深海"},
		21: {88, 217, 212, "水中木"},
		22: {189, 183, 107, "荒漠"},
		23: {247, 171, 94, "桦树林"},
		24: {255, 153, 71, "森林"},
		25: {82, 153, 71, "森林"},
		26: {200, 238, 200, "月岛"},
		27: {0, 85, 252, "浅海"},
		28: {125, 121, 121, "海礁石"},
		29: {222, 184, 135, "荒原"},
		30: {234, 152, 84, "红树林"},
		31: {112, 128, 144, "山地"},
		32: {176, 196, 222, "高山"},
		33: {176, 196, 222, "高山"},
		34: {255, 196, 222, "高山"},
	}

	// 初始化实体颜色映射
	g.entityColors = map[string]color.Color{
		"evergreen":     color.RGBA{0, 100, 0, 255},
		"deciduoustree": color.RGBA{34, 139, 34, 255},
		"rock1":         color.RGBA{128, 128, 128, 255},
		"rock2":         color.RGBA{169, 169, 169, 255},
		"sapling":       color.RGBA{144, 238, 144, 255},
		"grass":         color.RGBA{124, 252, 0, 255},
		"berrybush":     color.RGBA{255, 0, 0, 255},
		"berrybush2":    color.RGBA{220, 20, 60, 255},
		"spiderden":     color.RGBA{128, 0, 128, 255},
		"pond":          color.RGBA{0, 191, 255, 255},
		"rabbithole":    color.RGBA{139, 69, 19, 255},
		"cave":          color.RGBA{105, 105, 105, 255},
		"reeds":         color.RGBA{189, 183, 107, 255},
		"marsh_tree":    color.RGBA{85, 107, 47, 255},
		"marsh_bush":    color.RGBA{107, 142, 35, 255},
		"carrot":        color.RGBA{255, 140, 0, 255},
		"flower":        color.RGBA{255, 192, 203, 255},
		"wormhole":      color.RGBA{138, 43, 226, 255},
		"pighouse":      color.RGBA{255, 160, 122, 255},
		"mound":         color.RGBA{160, 82, 45, 255},
		"ruins":         color.RGBA{119, 136, 153, 255},
		"fireflies":     color.RGBA{255, 255, 0, 255},
	}

	return g
}

// ReadSaveFile 读取存档文件
func (g *DSTMapGenerator) ReadSaveFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 使用正则表达式提取地图数据
	re := regexp.MustCompile(`tiles="([^"]+)"`)
	matches := re.FindSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法在存档中找到地图数据")
	}

	base64Data := string(matches[1])
	fmt.Printf("提取的base64数据长度: %d\n", len(base64Data))
	return base64Data, nil
}

// DecodeMapData 解码地图数据
func (g *DSTMapGenerator) DecodeMapData(tilesBase64 string) ([]int, error) {
	tileBytes, err := base64.StdEncoding.DecodeString(tilesBase64)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %v", err)
	}

	// 处理文件头
	dataStart := 0
	if len(tileBytes) > 5 && string(tileBytes[:5]) == "VRSTN" {
		dataStart = 5
		for dataStart < len(tileBytes) && tileBytes[dataStart] == 0 {
			dataStart++
		}
	}

	tileBytes = tileBytes[dataStart:]
	if len(tileBytes)%2 != 0 {
		tileBytes = tileBytes[:len(tileBytes)-1]
	}

	// 解码tile IDs
	tileIds := make([]int, 0, len(tileBytes)/2)
	for i := 0; i < len(tileBytes); i += 2 {
		if i+1 >= len(tileBytes) {
			break
		}
		tileId := (int(tileBytes[i+1]) << 8) | int(tileBytes[i])
		tileId = tileId % 31 // 映射到0-33范围
		tileIds = append(tileIds, tileId)
	}

	// 打印分布情况
	tileCounts := make(map[int]int)
	for _, id := range tileIds {
		tileCounts[id]++
	}

	total := len(tileIds)
	fmt.Println("\n瓦片分布情况:")
	for id, count := range tileCounts {
		percentage := float64(count) / float64(total) * 100
		if percentage > 0.01 {
			fmt.Printf("瓦片ID %d: %d 个 (%.2f%%)\n", id, count, percentage)
		}
	}

	return tileIds, nil
}

// CreateMapImage 创建地图图像
func (g *DSTMapGenerator) CreateMapImage(tileIds []int, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充地形颜色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			if idx >= len(tileIds) {
				continue
			}

			tileId := tileIds[idx]
			tileColor, exists := g.tileColors[tileId]
			if !exists {
				// 使用默认颜色 (黑色)
				tileColor = Color{0, 0, 0, "未知地形"}
			}

			// 在Go中，图像坐标是从左到右的，所以需要翻转X坐标
			flippedX := width - x - 1
			img.Set(flippedX, y, color.RGBA{
				R: tileColor.R,
				G: tileColor.G,
				B: tileColor.B,
				A: 255,
			})
		}
	}

	return img
}

// GenerateMap 生成完整的地图
func (g *DSTMapGenerator) GenerateMap(saveFilePath, outputPath string, width, height int) error {
	// 读取并解码地图数据
	tilesBase64, err := g.ReadSaveFile(saveFilePath)
	if err != nil {
		return err
	}

	tileIds, err := g.DecodeMapData(tilesBase64)
	if err != nil {
		return err
	}

	// 确保生成的图像为指定的宽度和高度
	fmt.Printf("生成的地图尺寸: %dx%d\n", width, height)

	// 创建地图图像
	img := g.CreateMapImage(tileIds, width, height)

	// 保存图像
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("保存图像失败: %v", err)
	}

	fmt.Printf("地图已保存到: %s\n", outputPath)
	return nil
}

// ExtractDimensions 从文件中读取并提取 height 和 width
func ExtractDimensions(filePath string) (int, int, error) {
	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, 0, err
	}

	// 将文件内容转换为字符串
	input := string(data)

	// 定义正则表达式
	heightRegex := regexp.MustCompile(`height=(\d+)`)
	widthRegex := regexp.MustCompile(`width=(\d+)`)

	// 查找匹配
	heightMatch := heightRegex.FindStringSubmatch(input)
	widthMatch := widthRegex.FindStringSubmatch(input)

	// 返回结果
	var height, width int
	if len(heightMatch) > 1 {
		fmt.Sscanf(heightMatch[1], "%d", &height)
	} else {
		height = 0 // 未找到时返回 0
	}

	if len(widthMatch) > 1 {
		fmt.Sscanf(widthMatch[1], "%d", &width)
	} else {
		width = 0 // 未找到时返回 0
	}

	return height, width, nil
}

// WalrusHut_Plains
func main() {
	filePath := "save/session/50D0753A78BF681E/0000000004"
	height, width, err := ExtractDimensions(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	generator := NewDSTMapGenerator()
	err = generator.GenerateMap(
		filePath,
		"dst_map.png",
		height, // 指定宽度
		width,  // 指定高度
	)
	if err != nil {
		fmt.Printf("生成地图时出错: %v\n", err)
		os.Exit(1)
	}
}
