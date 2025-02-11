package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"

	"github.com/Kapeland/task-Avito/internal/utils/config"
)

func CreateStatusServer(cfg *config.Config, isReady *atomic.Value) *http.Server {
	statusAddr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)

	if !cfg.Project.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.HandleMethodNotAllowed = true // Обрабатывает 405 код
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET(cfg.Status.LivenessPath, livenessHandler)
	router.GET(cfg.Status.ReadinessPath, readinessHandler(isReady))
	router.GET(cfg.Status.VersionPath, versionHandler(cfg))

	statusServer := &http.Server{
		Addr:    statusAddr,
		Handler: router.Handler(),
	}

	return statusServer
}

func livenessHandler(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func readinessHandler(isReady *atomic.Value) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isReady == nil || !isReady.Load().(bool) {
			c.String(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))
			c.Abort()
			return
		}
		c.Writer.WriteHeader(http.StatusOK)
	}
}

func versionHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := map[string]interface{}{
			"debug": cfg.Project.Debug,
		}
		c.JSON(http.StatusOK, data)
	}
}
