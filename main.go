package main

import (
	"fmt"
	"log"
	"os"

	"sojson/env"
	"sojson/server"
	"sojson/zlog"

	"github.com/urfave/cli/v2"
)

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
		Before: initApp,
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
				Action: server.RunHTTPServer,
			},
		},
		DefaultCommand: "server",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// initApp 初始化日志系统
func initApp(ctx *cli.Context) error {
	// 初始化日志，输出到 stdout
	if err := zlog.InitLogger("info", ""); err != nil {
		return fmt.Errorf("初始化日志失败: %v", err)
	}

	env.Init()

	zlog.Info(ctx.Context, "系统初始化成功")
	return nil
}
