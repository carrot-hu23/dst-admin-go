package clusterUtils

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"github.com/gin-gonic/gin"
)

func GetCluster(clusterName string) *model.Cluster {
	db := database.DB
	cluster := &model.Cluster{}
	db.Where("cluster_name=?", clusterName).First(cluster)
	return cluster
}

func GetClusterFromGin(ctx *gin.Context) *model.Cluster {
	clusterName := ctx.GetHeader("cluster")
	db := database.DB
	cluster := &model.Cluster{}
	db.Where("cluster_name=?", clusterName).First(cluster)
	return cluster
}
