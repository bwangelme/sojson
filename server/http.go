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

	"sojson/controller"
	"sojson/env"
	"sojson/static"
	"sojson/zlog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func newGinEngine(ctx context.Context, staticFiles embed.FS, templateFiles embed.FS) *gin.Engine {
	// 设置为发布模式
	if env.IsTest() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	engine.Use(cors.New(config))

	// 使用自定义的静态文件处理器以确保正确的 MIME 类型
	engine.GET("/static/*filepath", func(c *gin.Context) {
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

		// 设置正确的 MIME 类型
		ext := filepath.Ext(filePath)
		var contentType string
		switch ext {
		case ".js":
			contentType = "application/javascript; charset=utf-8"
		case ".css":
			contentType = "text/css; charset=utf-8"
		case ".wasm":
			contentType = "application/wasm"
		case ".ttf":
			contentType = "font/ttf"
		case ".woff":
			contentType = "font/woff"
		case ".woff2":
			contentType = "font/woff2"
		default:
			// 使用默认的 MIME 类型检测
			mimeType := mime.TypeByExtension(ext)
			if mimeType != "" {
				contentType = mimeType
			} else {
				contentType = "application/octet-stream"
			}
		}

		// 强制清除缓存以避免 MIME 类型问题
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("Content-Type", contentType)

		// 返回文件内容
		c.Data(http.StatusOK, contentType, fileData)
	})

	// 使用嵌入的模板文件
	templ := template.Must(template.New("").ParseFS(templateFiles, "templates/*.html"))
	engine.SetHTMLTemplate(templ)

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

	controller.RegisterRoutes(engine)

	// 启动服务器
	zlog.Infof(ctx.Context, "🚀 SoJSON 服务器启动在 http://%s\n", address)
	zlog.Infof(ctx.Context, "📝 访问主页: http://%s\n", address)
	zlog.Infof(ctx.Context, "🔧 API 文档: http://%s/api\n", address)

	return engine.Run(address)
}
