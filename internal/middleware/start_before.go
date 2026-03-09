package middleware

import (
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/levelConfig"
	"log"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func StartBeforeMiddleware(archive *archive.PathResolver, levelConfigUtils *levelConfig.LevelConfigUtils) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		makeRunVersion(ctx, archive, levelConfigUtils)
		customcommandsFile(ctx, archive, levelConfigUtils)
	}
}

func copyOsFile() {

}

func customcommandsFile(ctx *gin.Context, archive *archive.PathResolver, levelConfigUtils *levelConfig.LevelConfigUtils) {
	clusterName := context.GetClusterName(ctx)
	config, _ := levelConfigUtils.GetLevelConfig(clusterName)
	for _, item := range config.LevelList {
		levelName := item.File
		levelPath := archive.LevelPath(clusterName, levelName)
		path := filepath.Join(levelPath, "customcommands.lua")
		_ = fileUtils.CreateFileIfNotExists(path)
		customcommands, _ := fileUtils.ReadFile("./static/customcommands.lua")
		fileUtils.WriterTXT(path, customcommands)
	}
}

func makeRunVersion(ctx *gin.Context, archive *archive.PathResolver, levelConfigUtils *levelConfig.LevelConfigUtils) {
	clusterName := context.GetClusterName(ctx)
	version, _ := archive.GetLocalDstVersion(clusterName)
	urlPath := ctx.Request.URL.Path
	config, err := levelConfigUtils.GetLevelConfig(clusterName)
	if strings.Contains(urlPath, "all") {
		if err == nil {
			for i := range config.LevelList {
				config.LevelList[i].RunVersion = version
				config.LevelList[i].Version = version
			}
		}
	} else {
		levelName := ctx.Query("levelName")
		if err == nil {
			for i := range config.LevelList {
				if config.LevelList[i].File == levelName {
					config.LevelList[i].RunVersion = version
					config.LevelList[i].Version = version
				}
			}
		}
	}
	err = levelConfigUtils.SaveLevelConfig(clusterName, config)
	if err != nil {
		log.Println(err)
	}
}
