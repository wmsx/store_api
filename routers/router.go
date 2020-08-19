package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	mygin "github.com/wmsx/pkg/gin"
	"github.com/wmsx/store_api/handler"
)

/**
 * 初始化路由信息
 */
func InitRouter(c client.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	storeHandler := handler.NewStoreHandler(c)
	storeRouter := r.Group("/store")

	storeRouter.POST("/upload/avatar", mygin.AuthWrapper(storeHandler.UploadAvatar))
	storeRouter.POST("/upload/file", mygin.AuthWrapper(storeHandler.UploadFile))

	return r
}
