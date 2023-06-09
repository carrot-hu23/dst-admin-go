package api

import (
	"dst-admin-go/schedule"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SteamApi struct {
}

func (s *SteamApi) DstNews(ctx *gin.Context) {
	imageSchedule := schedule.ImageSchedule{}
	tasks := imageSchedule.StartSchedule()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: tasks,
	})
}
