package router

import (
	"gin-gorm-oj/docs"
	"gin-gorm-oj/middlewares"
	"gin-gorm-oj/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 问题
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)
	//用户
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/login", service.Login)
	r.POST("/send_code", service.SendCode)
	r.POST("/register", service.Register)
	// 排行榜
	r.GET("/rank-list", service.GetRankList)

	// 提交
	r.GET("/submit-list", service.GetSubmitList)

	/**
	管理员私有方法
	*/
	authAdmin := r.Group("admin", middlewares.AuthAdminCheck())
	// 问题创建
	authAdmin.POST("problem-create", service.ProblemCreate)
	// 问题修改
	authAdmin.PUT("problem-modify", service.ProblemModify)
	// 分类列表
	authAdmin.GET("category-list", service.GetCategoryList)
	// 分类创建
	authAdmin.POST("category-create", service.CategoryCreate)
	// 分类修改
	authAdmin.PUT("category-modify", service.CategoryModify)
	// 分类删除
	authAdmin.DELETE("category-delete", service.CategoryDelete)
	/**
	用户私有方法
	*/
	authUser := r.Group("user", middlewares.AuthUserCheck())
	// 代码提交
	authUser.POST("submit", service.Submit)

	return r
}
