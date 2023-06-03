package service

import (
	"bytes"
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"fmt"
	"html/template"
)

type DstHelper struct {
}

func (dst *DstHelper) ParseTemplate(templatePath string, data interface{}) string {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, data)
	fmt.Println("解析文本模板")
	fmt.Printf("buf.String():\n%v\n", buf.String())
	return buf.String()
}

func (dst *DstHelper) DedicatedServerModsSetup(modConfig string) {
	if modConfig != "" {
		var serverModSetup = ""
		workshopIds := WorkshopIds(modConfig)
		for _, workshopId := range workshopIds {
			serverModSetup += "ServerModSetup(\"" + workshopId + "\")\n"
		}
		fileUtils.WriterTXT(constant.GET_DST_MOD_SETUP_PATH(), serverModSetup)
	}

}
