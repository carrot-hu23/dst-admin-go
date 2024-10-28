package api

import (
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
)

type FileApi struct {
}

func (d *FileApi) UploadUgcMods(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Panicln(err)
	}
	files := form.File["files"]
	filePaths := form.Value["filePaths"]

	dstConfig := dstConfigUtils.GetDstConfig()
	ugcModPath := ""
	if dstConfig.Ugc_directory != "" {
		ugcModPath = filepath.Join(dstConfig.Ugc_directory, "content", "322330")
	} else {
		ugcModPath = filepath.Join(dstConfig.Force_install_dir, "ugc_mods", dstConfig.Cluster, "Master", "content", "322330")
	}

	log.Println("上传ugc模组路径 ", ugcModPath)
	for i, file := range files {
		dir := filepath.Join(ugcModPath, filepath.Dir(filePaths[i]))
		fileUtils.CreateDirIfNotExists(dir)
		p := filepath.Join(ugcModPath, filepath.Dir(filePaths[i]), file.Filename)
		log.Println(p)
		if err := ctx.SaveUploadedFile(file, p); err != nil {
			log.Panicln(err)
		}
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}
