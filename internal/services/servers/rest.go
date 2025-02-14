package servers

import (
	"net/http"

	"github.com/Kapeland/task-Avito/internal/services/servers/middleware"
	"github.com/Kapeland/task-Avito/internal/utils/config"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/gin-gonic/gin"
)

func CreateRESTServer(implAuth AuthServer, implShop ShopServer, restAddr string) *http.Server {
	cfg := config.GetConfig()

	if !cfg.Project.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	lgr := logger.GetLogger()
	router := gin.New()
	router.HandleMethodNotAllowed = true // Обрабатывает 405 код
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	storeGr := router.Group("/api", middleware.CheckJWT(implAuth.A, &lgr))
	{
		storeGr.GET("/info", implShop.Info)
	}
	authGR := router.Group("/api")
	{
		authGR.POST("/auth", implAuth.Register)
	}

	operGr := router.Group("/api", middleware.CheckJWT(implAuth.A, &lgr))
	{
		operGr.POST("/sendCoin", implShop.SendCoin)
		operGr.GET("/buy/:item", implShop.BuyItem)

	}
	restServer := &http.Server{
		Addr:    restAddr,
		Handler: router.Handler(),
	}

	return restServer
}

//TODO: проверить все ли коды возврата
