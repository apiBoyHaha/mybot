package main

import (
	"log"
	"mybot/api/routes"
	"mybot/internal/database"
	"mybot/pkg/config"
	"net/http"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 设置路由
	router := routes.SetupRoutes(db, cfg)

	log.Printf("服务器启动在 %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
