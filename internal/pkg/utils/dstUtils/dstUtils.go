package dstUtils

import (
	"bytes"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/dstConfig"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	textTemplate "text/template"
)

func EscapePath(path string) string {
	if runtime.GOOS == "windows" {
		return path
	}
	// 在这里添加需要转义的特殊字符
	escapedChars := []string{" ", "'", "(", ")"}
	for _, char := range escapedChars {
		path = strings.ReplaceAll(path, char, "\\"+char)
	}
	return path
}

func WorkshopIds(content string) []string {
	var workshopIds []string

	re := regexp.MustCompile("\"workshop-\\w[-\\w+]*\"")
	workshops := re.FindAllString(content, -1)

	for _, workshop := range workshops {
		workshop = strings.Replace(workshop, "\"", "", -1)
		split := strings.Split(workshop, "-")
		workshopId := strings.TrimSpace(split[1])
		workshopIds = append(workshopIds, workshopId)
	}
	return workshopIds
}

func DedicatedServerModsSetup(dstConfig dstConfig.DstConfig, modConfig string) error {
	return DedicatedServerModsSetupExact(dstConfig, []string{modConfig})
}

func DedicatedServerModsSetupExact(dstConfig dstConfig.DstConfig, modConfigs []string) error {
	modSetupPath := GetModSetup2(dstConfig)
	mods, err := fileUtils.ReadLnFile(modSetupPath)
	if err != nil {
		return err
	}

	seen := map[string]bool{}
	var serverModSetup []string
	for _, modConfig := range modConfigs {
		if modConfig == "" {
			continue
		}
		workshopIds := WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			if workshopId == "" || seen[workshopId] {
				continue
			}
			seen[workshopId] = true
			serverModSetup = append(serverModSetup, "ServerModSetup(\""+workshopId+"\")")
		}
	}

	serverModSetupPattern := regexp.MustCompile(`^\s*ServerModSetup\s*\(`)
	var preserved []string
	for _, line := range mods {
		if serverModSetupPattern.MatchString(line) {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		preserved = append(preserved, line)
	}

	newServerModSetup := append(serverModSetup, preserved...)
	return fileUtils.WriterLnFile(modSetupPath, newServerModSetup)
}

func GetModSetup2(dstConfig dstConfig.DstConfig) string {
	return filepath.Join(dstConfig.Force_install_dir, "mods", "dedicated_server_mods_setup.lua")
}

func ParseTemplate(templatePath string, data interface{}) string {

	// 读取文件内容
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		panic(err)
	}

	// 创建模板对象
	tmpl, err := textTemplate.New("myTemplate").Parse(string(content))
	if err != nil {
		panic(err)
	}

	// 执行模板并保存结果到字符串
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()

}
