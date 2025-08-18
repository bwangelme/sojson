package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// indexHandler 主页处理器
func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "SoJSON - JSON工具",
	})
}

// RegisterRoutes 设置所有路由
func RegisterRoutes(engine *gin.Engine) *gin.Engine {
	// 主页路由
	engine.GET("/", indexHandler)

	// API路由组
	api := engine.Group("/api")
	{
		api.POST("/unescape", JSONController.UnescapeJSON)
		api.POST("/format", JSONController.FormatJSON)
		api.POST("/process", JSONController.ProcessJSON)
		api.POST("/validate", JSONController.ValidateJSON)
	}

	return engine
}
