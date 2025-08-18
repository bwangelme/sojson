package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"

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

	// 使用嵌入的静态文件
	engine.StaticFS("/static", http.FS(staticFiles))

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
	controller.RegisterRoutes(engine)

	// 启动服务器
	zlog.Infof(ctx.Context, "🚀 SoJSON 服务器启动在 http://%s\n", address)
	zlog.Infof(ctx.Context, "📝 访问主页: http://%s\n", address)
	zlog.Infof(ctx.Context, "🔧 API 文档: http://%s/api\n", address)

	return engine.Run(address)
}
