package servers

import (
	"context"
	"net/http"

	"github.com/Kapeland/task-Avito/internal/models"
	"github.com/Kapeland/task-Avito/internal/models/structs"
	svStruct "github.com/Kapeland/task-Avito/internal/services/structs"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/gin-gonic/gin"
)

type ShopServer struct {
	U models.UsersModelManager
	A models.AuthModelManager
}

func (s *ShopServer) SendCoin(c *gin.Context) {
	lgr := logger.GetLogger()

	var operation svStruct.SendCoinReqBody

	if err := c.ShouldBindBodyWithJSON(&operation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	login := c.Keys["login"].(string) // Получаем из JWT middleware

	err := s.sendCoin(c.Request.Context(), operation, login)
	if err != nil {
		lgr.Error(err.Error(), "ShopServer", "SendCoin", "sendCoin")
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
	}
	c.Status(http.StatusOK)
}

func (s *ShopServer) sendCoin(ctx context.Context, operation svStruct.SendCoinReqBody, fromLogin string) error {
	lgr := logger.GetLogger()

	err := s.U.SendCoin(ctx, structs.SendCoinInfo{
		From:   fromLogin,
		To:     operation.To,
		Amount: operation.Amount,
	})
	if err != nil {
		lgr.Error(err.Error(), "ShopServer", "sendCoin", "SendCoin")
	}
	return err
}

func (s *ShopServer) BuyItem(c *gin.Context) {
	lgr := logger.GetLogger()

	item := c.Param("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "item can't be empty"})
		return
	}

	login := c.Keys["login"].(string) // Получаем из JWT middleware

	err := s.buyItem(c.Request.Context(), item, login)
	if err != nil {
		lgr.Error(err.Error(), "ShopServer", "BuyItem", "buyItem")
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
	}
	c.Status(http.StatusOK)
}

func (s *ShopServer) buyItem(ctx context.Context, item string, login string) error {
	lgr := logger.GetLogger()
	err := s.U.BuyItem(ctx, item, login)
	if err != nil {
		lgr.Error(err.Error(), "ShopServer", "buyItem", "BuyItem")
	}
	return err
}

func (s *ShopServer) Info(c *gin.Context) {
	lgr := logger.GetLogger()

	login := c.Keys["login"].(string) // Получаем из JWT middleware

	accInfo, err := s.info(c.Request.Context(), login)
	if err != nil {
		lgr.Error(err.Error(), "ShopServer", "Info", "info")
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
	}
	c.JSON(http.StatusOK, accInfo)
}

func (s *ShopServer) info(ctx context.Context, login string) (structs.AccInfo, error) {
	lgr := logger.GetLogger()
	accInfo, err := s.U.Info(ctx, login)
	if err != nil {
		lgr.Error(err.Error(), "ShopServer", "info", "Info")
	}
	return accInfo, err
}

//TODO: проверить правильные ли коды возврата
