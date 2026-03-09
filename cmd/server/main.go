// @title           DST Admin Go API
// @version         1.0
// @description     饥荒联机版服务器管理后台 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/carrot-hu23/dst-admin-go
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"dst-admin-go/internal/api"
	"dst-admin-go/internal/config"
	"dst-admin-go/internal/database"
	"fmt"
)

func main() {

	cfg := config.Load()
	db := database.InitDB(cfg)

	route := api.NewRoute(cfg, db)

	err := route.Run(cfg.BindAddress + ":" + cfg.Port)
	if err != nil {
		fmt.Println("启动失败！！！", err)
	}
}
