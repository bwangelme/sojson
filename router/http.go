package router

import (
	"net/http"

	"sojson/controller"

	"github.com/gin-gonic/gin"
)

// indexHandler 主页处理器
func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

// RegisterHTTP 设置所有路由
func RegisterHTTP(engine *gin.Engine) *gin.Engine {
	// 主页路由
	engine.GET("/", indexHandler)

	// API路由组
	api := engine.Group("/api")
	{
		api.POST("/unescape", controller.JSONController.UnescapeJSON)
		api.POST("/format", controller.JSONController.FormatJSON)
		api.POST("/process", controller.JSONController.ProcessJSON)
		api.POST("/validate", controller.JSONController.ValidateJSON)
	}

	return engine
}
