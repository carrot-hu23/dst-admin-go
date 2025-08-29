package service

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

	// 初始化tile颜色映射
	g.tileColors = map[int]Color{
		0:   {42, 42, 42, "地表虚空"},
		1:   {42, 42, 44, "一个虚空"},
		2:   {80, 76, 65, "卵石路"},
		3:   {90, 105, 104, "矿石区"},
		4:   {117, 107, 85, "无地皮"},
		5:   {144, 125, 89, "热带草原地皮"},
		6:   {48, 67, 39, "长草地皮"},
		7:   {40, 47, 18, "森林地皮"},
		8:   {81, 28, 194, "沼泽地皮"},
		9:   {0, 0, 7, "空"},
		10:  {85, 69, 48, "木地板"},
		11:  {64, 75, 116, "地毯地板"},
		12:  {115, 133, 201, "棋盘地板"},
		13:  {139, 131, 115, "鸟粪地皮"},
		14:  {74, 65, 77, "蓝真菌"},
		15:  {67, 70, 32, "黏滑地皮"},
		16:  {75, 75, 73, "洞穴岩石地皮"},
		17:  {66, 49, 30, "泥泞地皮"},
		18:  {115, 113, 107, "远古地皮"},
		19:  {86, 86, 81, "仿远古地皮"},
		20:  {74, 61, 84, "远古瓷砖"},
		21:  {66, 50, 76, "仿远古瓷砖"},
		22:  {39, 38, 39, "远古雕砖"},
		23:  {34, 35, 34, "仿远古雕砖"},
		24:  {70, 44, 43, "红真菌"},
		25:  {62, 76, 61, "绿真菌"},
		26:  {42, 42, 44, "空"},
		27:  {42, 42, 44, "空"},
		28:  {42, 42, 44, "空"},
		29:  {42, 42, 44, "空"},
		30:  {91, 62, 14, "落叶林地皮"},
		31:  {117, 86, 46, "沙漠地皮"},
		32:  {31, 27, 27, "龙鳞地皮"},
		33:  {128, 128, 128, "崩溃"},
		34:  {128, 128, 128, "崩溃"},
		35:  {47, 44, 47, "暴食沼泽"},
		36:  {158, 104, 105, "粉桦树林【暴食】"},
		37:  {137, 113, 113, "粉【暴食】"},
		38:  {81, 97, 100, "蓝长草【暴食】"},
		39:  {66, 58, 49, "耕地地皮"},
		40:  {42, 42, 44, "空"},
		41:  {119, 113, 97, "白【暴食】"},
		42:  {84, 108, 107, "岩石海滩"},
		43:  {67, 133, 142, "月球环形山"},
		44:  {154, 146, 186, "贝壳海滩"},
		45:  {138, 96, 73, "远古石刻"},
		46:  {66, 86, 82, "变异真菌地皮"},
		47:  {61, 57, 46, "耕地地皮"},
		48:  {42, 42, 44, "空"},
		200: {0, 0, 11, "深渊"},
		201: {18, 66, 73, "浅海"},
		202: {18, 66, 73, "浅海海岸"},
		203: {7, 46, 61, "中海"},
		204: {1, 32, 46, "深海"},
		205: {6, 54, 81, "盐矿海岸"},
		206: {6, 54, 81, "盐矿海岸"},
		207: {0, 24, 26, "危险海"},
		208: {189, 193, 198, "水中木"},
		247: {42, 42, 44, "空"},
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

func RestoreTileId(original int, colors map[int]Color) int {
	high := original >> 8
	low := original & 0xFF

	// 1) 直接用 high（正常情形：tile 存为 high<<8）
	if _, ok := colors[high]; ok {
		return high
	}

	// 2) 尝试 high-1
	if high > 0 {
		if _, ok := colors[high-1]; ok {
			return high - 1
		}
	}

	// 3) ocean 等大型 ID 区间的保守尝试（0xC8 = 200）
	//    如果 high 在 0xC8..0xFF 范围，优先把它当作 200+ 值尝试映射
	if high >= 0xC8 && high <= 0xFF {
		cand := 200 + (high - 0xC8)
		if _, ok := colors[cand]; ok {
			return cand
		}
	}

	// 4) 其它情况：如果低字节不为0，有可能这不是简单的 high<<8 编码，
	//    你可以在此加入额外规则；暂时返回 0（或你想要的默认）
	_ = low // 如果以后需要可利用 low 做更精细的判断
	return 0
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
		tileId = RestoreTileId(tileId, g.tileColors)
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
