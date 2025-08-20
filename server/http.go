package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"sojson/env"
	"sojson/router"
	"sojson/static"
	"sojson/zlog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func newStaticHandler(staticFiles embed.FS) gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := strings.TrimPrefix(c.Param("filepath"), "/")

		// 忽略 source map 请求以减少 404 错误 - 返回 204 No Content
		if strings.HasSuffix(filePath, ".map") {
			c.Status(http.StatusNoContent)
			return
		}

		// 尝试从嵌入的文件系统读取文件
		fileData, err := staticFiles.ReadFile(filePath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		// 使用文件扩展名让 Gin 自动设置正确的 Content-Type
		ext := filepath.Ext(filePath)
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		c.Data(http.StatusOK, contentType, fileData)
	}
}
func newGinEngine(ctx context.Context, staticFiles embed.FS, templateFiles embed.FS) *gin.Engine {
	// 设置为发布模式
	if env.IsTest() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	// 使用嵌入的模板文件
	templ := template.Must(template.New("").ParseFS(templateFiles, "templates/*.html"))
	engine.SetHTMLTemplate(templ)

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	engine.Use(cors.New(config))

	// 使用自定义的静态文件处理器以确保正确的 MIME 类型
	engine.GET("/static/*filepath", newStaticHandler(staticFiles))

	return engine
}

// RunHTTPServer 运行服务器
func RunHTTPServer(ctx *cli.Context) error {
	host := ctx.String("host")
	port := ctx.Int("port")
	address := fmt.Sprintf("%s:%d", host, port)

	// 创建路由
	engine := newGinEngine(ctx.Context, static.StaticFiles, static.TemplateFiles)

	// 处理Chrome DevTools的特定请求，返回204避免404错误
	engine.GET("/.well-known/appspecific/com.chrome.devtools.json", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	router.RegisterHTTP(engine)

	// 启动服务器
	zlog.Infof(ctx.Context, "🚀 SoJSON 服务器启动在 http://%s\n", address)
	zlog.Infof(ctx.Context, "📝 访问主页: http://%s\n", address)
	zlog.Infof(ctx.Context, "🔧 API 文档: http://%s/api\n", address)

	return engine.Run(address)
}
