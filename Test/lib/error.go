package lib

import "github.com/gin-gonic/gin"

var (
	OK                    = gin.H{"message": "成功"}
	BadRequest            = gin.H{"message": "参数错误"}
	Unauthorized          = gin.H{"message": "未登录或非法访问"}
	InternalServerError   = gin.H{"message": "系统错误"}
	InsufficientBalance   = gin.H{"message": "余额不足"}
	NoFoundSku            = gin.H{"message": "暂无商品信息"}
	InvalidSignature      = gin.H{"message": "签名验证未通过"}
	InvalidCharacter      = gin.H{"message": "当前对象已存在"}
	InvalidAuthentication = gin.H{"message": "用户名或密码错误"}
	ErrAdminRoleExists    = gin.H{"message": "重复分配"}
	EmptyResponse         = gin.H{}
)
