package controller

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"sojson/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router 路由管理器
type Router struct {
	jsonController *JSONController
}

// NewRouter 创建路由管理器
func NewRouter(staticFiles embed.FS, templateFiles embed.FS) *Router {
	jsonService := service.NewJSONProcessorService()
	jsonController := NewJSONController(jsonService)

	return &Router{
		jsonController: jsonController,
	}
}

// SetupRoutes 设置所有路由
func (r *Router) SetupRoutes(staticFiles embed.FS, templateFiles embed.FS) *gin.Engine {
	// 设置为发布模式
	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	engine.Use(cors.New(config))

	// 使用嵌入的静态文件
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	engine.StaticFS("/static", http.FS(staticFS))

	// 使用嵌入的模板文件
	templ := template.Must(template.New("").ParseFS(templateFiles, "templates/*.html"))
	engine.SetHTMLTemplate(templ)

	// 主页路由
	engine.GET("/", r.indexHandler)

	// API路由组
	api := engine.Group("/api")
	{
		api.POST("/unescape", r.jsonController.UnescapeJSON)
		api.POST("/format", r.jsonController.FormatJSON)
		api.POST("/process", r.jsonController.ProcessJSON)
		api.POST("/validate", r.jsonController.ValidateJSON)
	}

	return engine
}

// indexHandler 主页处理器
func (r *Router) indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "SoJSON - JSON工具",
	})
}
