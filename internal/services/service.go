package services

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Kapeland/task-Avito/internal/models"
	"github.com/Kapeland/task-Avito/internal/services/servers"
	"github.com/Kapeland/task-Avito/internal/utils/config"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/pkg/errors"
)

type Service struct {
	um models.UsersModelManager
	am models.AuthModelManager
}

func NewService(um models.UsersModelManager, am models.AuthModelManager) Service {
	return Service{um: um, am: am}
}

func (s Service) Launch(cfg *config.Config, lgr *logger.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	implAuth := servers.AuthServer{A: s.am}
	implShop := servers.ShopServer{U: s.um, A: s.am}

	restAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)

	restServer := servers.CreateRESTServer(implAuth, implShop, restAddr)

	go func() {
		lgr.InfoMsg(fmt.Sprintf("REST server is running on %s", restAddr))
		if err := restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lgr.ErrorMsg("Failed running REST server")
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	statusServer := servers.CreateStatusServer(cfg, isReady)

	go func() {
		statusAdrr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)
		lgr.InfoMsg(fmt.Sprintf("Status server is running on %s", statusAdrr))

		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lgr.ErrorMsg("Failed running status server")
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		lgr.InfoMsg("The Service successfully launched")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		lgr.InfoMsg(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		lgr.InfoMsg(fmt.Sprintf("ctx.Done: %v", done))
	}

	isReady.Store(false)

	if err := restServer.Shutdown(ctx); err != nil {
		lgr.Error("Failed shutting down REST server", "Service", "Launch", "restServer.Shutdown")
	} else {
		lgr.InfoMsg("REST server shut down successfully")
	}

	if err := statusServer.Shutdown(ctx); err != nil {
		lgr.Error("Failed shutting down Status server", "Service", "Launch", "statusServer.Shutdown")
	} else {
		lgr.InfoMsg("Status server shut down successfully")
	}

	return nil
}
