package handler

import (
	"dst-admin-go/internal/database"
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/login"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	PasswordPath = "./password.txt"
)

type LoginHandler struct {
	loginService *login.LoginService
}

func NewLoginHandler(loginService *login.LoginService) *LoginHandler {
	return &LoginHandler{
		loginService: loginService,
	}
}

func (h *LoginHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/user", h.GetUserInfo)
	router.POST("/api/login", h.Login)
	router.GET("/api/logout", h.Logout)
	router.POST("/api/change/password", h.ChangePassword)
	router.POST("/api/user", h.UpdateUserInfo)
	router.GET("/api/init", h.CheckIsFirst)
	router.POST("/api/init", h.InitFirst)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=login.UserInfo}
// @Router /api/user/info [get]
func (h *LoginHandler) GetUserInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "Init user success",
		Data: h.loginService.GetUserInfo(),
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags user
// @Accept json
// @Produce json
// @Param user body login.UserInfo true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/user/login [post]
func (h *LoginHandler) Login(ctx *gin.Context) {
	var userInfo login.UserInfo
	err := ctx.ShouldBind(&userInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "Invalid request body: " + err.Error(),
			Data: nil,
		})
		return
	}

	response := h.loginService.Login(userInfo, ctx)
	ctx.JSON(http.StatusOK, response)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/user/logout [get]
func (h *LoginHandler) Logout(ctx *gin.Context) {
	h.loginService.Logout(ctx)
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "Logout success",
		Data: nil,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改密码
// @Tags user
// @Accept json
// @Produce json
// @Param password body object true "新密码"
// @Success 200 {object} response.Response
// @Router /api/user/changePassword [post]
func (h *LoginHandler) ChangePassword(ctx *gin.Context) {
	var body struct {
		NewPassword string `json:"newPassword"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "Invalid request body: " + err.Error(),
			Data: nil,
		})
		return
	}
	response := h.loginService.ChangePassword(body.NewPassword)

	ctx.JSON(http.StatusOK, response)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param user body object true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/user/update [post]
func (h *LoginHandler) UpdateUserInfo(ctx *gin.Context) {
	var body struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
		PhotoURL    string `json:"photoURL"`
		Password    string `json:"password"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "Invalid request body: " + err.Error(),
			Data: nil,
		})
		return
	}
	err := fileUtils.WriterLnFile(PasswordPath, []string{
		"username = " + body.Username,
		"password = " + body.Password,
		"displayName=" + body.DisplayName,
		"photoURL=" + body.PhotoURL,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Code: 500,
			Msg:  "修改用户信息失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "Update user info success",
		Data: nil,
	})
}

// InitFirst 初始化首次用户
// @Summary 初始化首次用户
// @Description 初始化系统首次用户信息
// @Tags user
// @Accept json
// @Produce json
// @Param userInfo body login.UserInfo true "用户信息"
// @Success 200 {object} response.Response
// @Router /api/init [post]
func (h *LoginHandler) InitFirst(ctx *gin.Context) {
	db := database.Db
	kv := model.KV{}
	db.Where("key = 'FIRST_INIT'").First(&kv)
	if kv.Value == "TRUE" || fileUtils.Exists("./first") {
		log.Panicln("非法请求")
	}
	var payload struct {
		UserInfo login.UserInfo `json:"userInfo"`
	}
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "Invalid request body: " + err.Error(),
			Data: nil,
		})
	}
	// 事务
	// 记录已经初始化
	// 保存用户信息
	err = db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&model.KV{Key: "FIRST_INIT", Value: "TRUE"})
		h.loginService.InitUserInfo(payload.UserInfo)
		return nil
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Code: 500,
			Msg:  "初始化失败: " + err.Error(),
			Data: nil,
		})
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// CheckIsFirst 检查是否首次初始化
// @Summary 检查是否首次初始化
// @Description 检查系统是否进行了首次初始化
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/init [get]
func (h *LoginHandler) CheckIsFirst(ctx *gin.Context) {

	exist := false
	db := database.Db
	kv := model.KV{}
	db.Where("key = 'FIRST_INIT'").First(&kv)
	if kv.Value == "TRUE" {
		exist = true
	} else {
		exist = fileUtils.Exists("./first")
	}

	code := 200
	msg := "is first"
	if exist {
		code = 400
		msg = "is not first"
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
