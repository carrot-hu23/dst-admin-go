package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant"
	"dst-admin-go/model"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserApi struct {
}

func (u *UserApi) QueryUserList(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}

	db := database.DB
	db2 := database.DB
	if name, isExist := ctx.GetQuery("username"); isExist {
		db = db.Where("username LIKE ?", "%"+name+"%")
		db2 = db2.Where("username LIKE ?", "%"+name+"%")
	}
	if displayName, isExist := ctx.GetQuery("displayName"); isExist {
		db = db.Where("display_name LIKE ?", "%"+displayName+"%")
		db2 = db2.Where("display_name LIKE ?", "%"+displayName+"%")
	}
	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)
	users := make([]model.User, 0)

	if err := db.Find(&users).Error; err != nil {
		fmt.Println(err.Error())
	}

	var total int64
	db2.Model(&model.User{}).Count(&total)

	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: vo.Page{
			Data:       users,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (u *UserApi) CreateUser(ctx *gin.Context) {

	user := model.User{}
	err := ctx.ShouldBind(&user)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%v", user)

	if user.Username == "" {
		log.Panicln("create user is error, Username is null")
	}
	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}
	if user.Password == "" {
		log.Panicln("create user is error, Password is null")
	}

	// 找到管理员
	admin, err := fileUtils.ReadLnFile(constant.PASSWORD_PATH)
	if err != nil {
		log.Panicln("Not find password file error: " + err.Error())
	}

	adminUsername := strings.TrimSpace(strings.Split(admin[0], "=")[1])
	if adminUsername == user.Username {
		log.Panicln("username not same admin username")
	}

	db1 := database.DB
	oldUser := model.User{}
	db1.Where("username=?", user.Username).First(&oldUser)
	if oldUser.ID != 0 {
		log.Panicln("create user is error, user is existed")
	}
	db := database.DB
	db.Create(&user)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (u *UserApi) DeleteUser(ctx *gin.Context) {

	id := ctx.Query("id")
	db := database.DB
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	User := model.User{}
	result := tx.Where("id = ?", id).Unscoped().Delete(&User)
	// 同时删除已经分配的集群
	if result.Error != nil {
		log.Panicln(result.Error)
	}
	var userClusters []model.UserCluster
	tx.Where("user_id = ?", User.ID).Delete(userClusters)
	tx.Commit()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (u *UserApi) UpdateUser(ctx *gin.Context) {

	user := model.User{}
	err := ctx.ShouldBind(&user)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%v", user)

	db := database.DB
	oldUser := &model.User{}
	db.Where("ID = ?", user.ID).First(oldUser)
	if oldUser.ID == 0 {
		log.Panicln("not find user")
	}

	if user.Username != "" {
		oldUser.Username = user.Username
	}
	if user.Description != "" {
		oldUser.Description = user.Description
	}
	if user.DisplayName != "" {
		oldUser.DisplayName = user.DisplayName
	}
	if user.Password != "" {
		oldUser.Password = user.Password
	}
	oldUser.PhotoURL = user.PhotoURL
	db.Updates(oldUser)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

type UserClusterVO struct {
	ID                    int
	Name                  string `json:"name"`
	ClusterName           string `json:"clusterName"`
	AllowAddLevel         bool   `json:"allowAddLevel"`
	AllowEditingServerIni bool   `json:"allowEditingServerIni"`
}

func (u *UserApi) GetUserClusterList(ctx *gin.Context) {

	userId := ctx.Query("userId")
	db := database.DB
	var userClusterVOList []UserClusterVO
	db.Raw("select uc.id as ID, c.name as name, c.cluster_name, uc.allow_add_level, uc.allow_editing_server_ini from user_clusters uc join clusters as c on c.id = uc.cluster_id where uc.deleted_at is null and uc.user_id=?", userId).Scan(&userClusterVOList)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: userClusterVOList,
	})
}

func (u *UserApi) GetUserCluster(ctx *gin.Context) {
	session := sessions.Default(ctx)
	role := session.Get("role")
	userId := session.Get("userId")
	clusterName := ctx.Query("clusterName")
	log.Println("role", role, "userId", userId, "clusterName", clusterName)
	userClusterVO := UserClusterVO{}
	if role == "admin" {
		userClusterVO.AllowAddLevel = true
		userClusterVO.AllowEditingServerIni = true
	} else {
		db := database.DB
		db.Raw("select uc.id as ID, c.name as name, c.cluster_name, uc.allow_add_level, uc.allow_editing_server_ini from user_clusters uc join clusters as c on c.id = uc.cluster_id where uc.deleted_at is null and uc.user_id= ? and c.cluster_name = ?", userId, clusterName).Scan(&userClusterVO)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: userClusterVO,
	})
}

type AddUserCluster struct {
	UserId                int  `json:"userId"`
	ClusterId             int  `json:"clusterId"`
	AllowAddLevel         bool `json:"allowAddLevel"`
	AllowEditingServerIni bool `json:"allowEditingServerIni"`
}

func (u *UserApi) AddUserCluster(ctx *gin.Context) {

	addUserCluster := AddUserCluster{}
	err := ctx.ShouldBind(&addUserCluster)
	if err != nil {
		log.Panicln(err)
	}

	if addUserCluster.UserId != 0 && addUserCluster.ClusterId != 0 {
		db := database.DB
		userCluster := model.UserCluster{
			UserId:                addUserCluster.UserId,
			ClusterId:             addUserCluster.ClusterId,
			AllowAddLevel:         addUserCluster.AllowAddLevel,
			AllowEditingServerIni: addUserCluster.AllowEditingServerIni,
		}
		db.Create(&userCluster)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (u *UserApi) RemoveUserCluster(ctx *gin.Context) {

	id := ctx.Query("id")
	db := database.DB
	var userCluster model.UserCluster
	db.Where("id = ?", id).Delete(&userCluster)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (u *UserApi) UpdateUserAllow(ctx *gin.Context) {

	var payload struct {
		ID                    int  `json:"ID"`
		AllowAddLevel         bool `json:"allowAddLevel"`
		AllowEditingServerIni bool `json:"allowEditingServerIni"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln("参数解析失败", err)
	}

	db := database.DB
	userCluster := model.UserCluster{}
	db.Where("id = ?", payload.ID).First(&userCluster)
	log.Println("userCluster", userCluster)

	userCluster.AllowAddLevel = payload.AllowAddLevel
	userCluster.AllowEditingServerIni = payload.AllowEditingServerIni

	db2 := database.DB
	db2.Save(&userCluster)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: userCluster,
	})

}
