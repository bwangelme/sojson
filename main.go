package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"sojson/controller"

	"github.com/urfave/cli/v2"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

func main() {
	app := &cli.App{
		Name:    "sojson",
		Usage:   "JSON 去除转义和格式化工具",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name: "SoJSON Team",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "启动 HTTP 服务器",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "host",
						Aliases: []string{"H"},
						Value:   "0.0.0.0",
						Usage:   "服务器监听地址",
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   2378,
						Usage:   "服务器监听端口",
					},
				},
				Action: runServer,
			},
		},
		DefaultCommand: "server",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// runServer 运行服务器
func runServer(ctx *cli.Context) error {
	host := ctx.String("host")
	port := ctx.Int("port")
	address := fmt.Sprintf("%s:%d", host, port)

	// 创建路由
	router := controller.NewRouter(staticFiles, templateFiles)
	engine := router.SetupRoutes(staticFiles, templateFiles)

	// 启动服务器
	fmt.Printf("🚀 SoJSON 服务器启动在 http://%s\n", address)
	fmt.Printf("📝 访问主页: http://%s\n", address)
	fmt.Printf("🔧 API 文档: http://%s/api\n", address)

	return engine.Run(address)
}
