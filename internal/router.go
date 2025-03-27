package internal

import (
	"github.com/gin-contrib/graceful"
	"gorm.io/gorm"

	"github.com/zgsm-ai/gatewayctl/internal/handler"
)

func InitRouter(db *gorm.DB) {
	router, err := graceful.Default()
	if err != nil {
		panic(err)
	}

	router.POST("/plugins/register", handler.RegisterPlugins)
	router.POST("/plugins/unregister", handler.UnregisterPlugins)
	router.GET("/plugins/list", handler.ListPlugins)
	router.GET("/plugins/:id", handler.GetPlugin)

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	router.Run(":8081")
}
