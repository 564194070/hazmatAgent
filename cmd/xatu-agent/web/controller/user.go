package controller

import (
	"agentFrame/cmd/xatu-agent/utils/tool"
	"agentFrame/cmd/xatu-agent/web/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	users *model.MySQLUser
}

func NewUserController(users *model.MySQLUser) *UserController {
	return &UserController{users: users}
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 注册
func (h *UserController) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	hashPwd, err := tool.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: hashPwd,
		Nickname: req.Nickname,
	}
	err = h.users.GetUserConnection().Create(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "用户名已存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}

// 登录
func (h *UserController) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}
	var user model.User
	h.users.GetUserConnection().Where("username = ?", req.Username).First(&user)
	if user.ID == 0 || !tool.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "账号密码错误"})
		return
	}
	token, err := tool.GenToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// 获取个人信息（需要登录）
func (h *UserController) GetInfo(c *gin.Context) {
	uid, _ := c.Get("user_id")
	var user model.User
	h.users.GetUserConnection().Select("id,username,nickname,email,status,role").First(&user, uid)
	c.JSON(http.StatusOK, gin.H{"data": user})
}
