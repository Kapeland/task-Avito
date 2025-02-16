package servers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Kapeland/task-Avito/internal/models"
	"github.com/Kapeland/task-Avito/internal/services/servers/middleware"
	"github.com/Kapeland/task-Avito/internal/services/structs"
	"github.com/Kapeland/task-Avito/internal/storage"
	"github.com/Kapeland/task-Avito/internal/storage/repository/postgresql/auth"
	"github.com/Kapeland/task-Avito/internal/storage/repository/postgresql/users"
	"github.com/Kapeland/task-Avito/internal/utils/config"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(implAuth AuthServer, implShop ShopServer, lgr *logger.Logger) *gin.Engine {
	router := gin.Default()
	storeGr := router.Group("/api", middleware.CheckJWT(implAuth.A, lgr))
	{
		storeGr.GET("/info", implShop.Info)
	}
	authGR := router.Group("/api")
	{
		authGR.POST("/auth", implAuth.Register)
	}

	operGr := router.Group("/api", middleware.CheckJWT(implAuth.A, lgr))
	{
		operGr.POST("/sendCoin", implShop.SendCoin)
		operGr.GET("/buy/:item", implShop.BuyItem)
	}
	return router
}

func initServer() (*gin.Engine, error) {
	ctx := context.Background()
	if err := config.ReadLocalConfigYAML(); err != nil {
		return nil, err
	}

	cfg := config.GetConfig()
	lgr := logger.CreateLogger(&cfg)

	dbStor, err := storage.NewPostgresStorage(ctx)
	if err != nil {
		return nil, err
	}

	usersRepo := users.New(dbStor.DB)
	authRepo := auth.New(dbStor.DB)

	authStorage := storage.NewAuthStorage(authRepo)
	usersStorage := storage.NewUsersStorage(usersRepo)

	umdl := models.NewModelUsers(&usersStorage)
	amdl := models.NewModelAuth(&authStorage, &usersStorage)

	implAuth := AuthServer{A: &amdl}
	implShop := ShopServer{U: &umdl, A: &amdl}

	tmp := setupRouter(implAuth, implShop, &lgr)
	return tmp, nil
}

func TestShopServer_NoAuth(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/buy/pen", nil)
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		t.Log(w.Body.String())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/info", nil)
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		t.Log(w.Body.String())
	}

	sendReq := structs.SendCoinReqBody{
		To:     "user1user2",
		Amount: 2,
	}

	sendReqJson, err := json.Marshal(sendReq)
	if err != nil {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/sendCoin", strings.NewReader(string(sendReqJson)))
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		t.Log(w.Body.String())
	}

}

func TestShopServer_Buy_ExistingItem(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/buy/pen", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3MzQ3NTUsInNJRCI6IjVmMzRmMGRlLTRhODQtNDc0ZC05NTBkLTc2NTcyYmQ0N2E3ZSIsInN1YiI6InVzZXIxdXNlcjEifQ.KsDL5GBzuk4jz6_RK9nqH9jeeYQigGxVSZHRuvr8uwU")
	router.ServeHTTP(w, req)

	if !assert.Equal(t, 200, w.Code) {
		t.Log(w.Body.String())
	}

}

func TestShopServer_Buy_NotExistingItem(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/buy/car", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3MzQ3NTUsInNJRCI6IjVmMzRmMGRlLTRhODQtNDc0ZC05NTBkLTc2NTcyYmQ0N2E3ZSIsInN1YiI6InVzZXIxdXNlcjEifQ.KsDL5GBzuk4jz6_RK9nqH9jeeYQigGxVSZHRuvr8uwU")
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		t.Log(w.Body.String())
	}
}

func TestShopServer_SendCoin_BothExist(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}

	sendReq := structs.SendCoinReqBody{
		To:     "user1user2",
		Amount: 2,
	}

	sendReqJson, err := json.Marshal(sendReq)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(string(sendReqJson)))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3MzQ3NTUsInNJRCI6IjVmMzRmMGRlLTRhODQtNDc0ZC05NTBkLTc2NTcyYmQ0N2E3ZSIsInN1YiI6InVzZXIxdXNlcjEifQ.KsDL5GBzuk4jz6_RK9nqH9jeeYQigGxVSZHRuvr8uwU")
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusOK, w.Code) {
		t.Log(w.Body.String())
	}
}

func TestShopServer_SendCoin_MoneyNotEnough(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}

	sendReq := structs.SendCoinReqBody{
		To:     "user1user2",
		Amount: 200000,
	}

	sendReqJson, err := json.Marshal(sendReq)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(string(sendReqJson)))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3MzQ3NTUsInNJRCI6IjVmMzRmMGRlLTRhODQtNDc0ZC05NTBkLTc2NTcyYmQ0N2E3ZSIsInN1YiI6InVzZXIxdXNlcjEifQ.KsDL5GBzuk4jz6_RK9nqH9jeeYQigGxVSZHRuvr8uwU")
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		t.Log(w.Body.String())
	}
}

func TestShopServer_SendCoin_RecipientNotExist(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}

	sendReq := structs.SendCoinReqBody{
		To:     "user1user2132132",
		Amount: 2,
	}

	sendReqJson, err := json.Marshal(sendReq)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/sendCoin", strings.NewReader(string(sendReqJson)))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3MzQ3NTUsInNJRCI6IjVmMzRmMGRlLTRhODQtNDc0ZC05NTBkLTc2NTcyYmQ0N2E3ZSIsInN1YiI6InVzZXIxdXNlcjEifQ.KsDL5GBzuk4jz6_RK9nqH9jeeYQigGxVSZHRuvr8uwU")
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		t.Log(w.Body.String())
	}
}

func TestAuthServer_RegisterExisting(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}

	authReq := structs.RegisterReqBody{
		Username: "user1user1",
		Password: "Lhjxb[eq1",
	}

	authReqJson, err := json.Marshal(authReq)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(authReqJson)))
	router.ServeHTTP(w, req)

	if !assert.Equal(t, 200, w.Code) {
		t.Log(w.Body.String())
	}
	t.Log(w.Body.String())

}

func TestAuthServer_RegisterNew(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}
	tmpNumb := strconv.Itoa(rand.Intn(1000) + 1000)

	authReq := structs.RegisterReqBody{
		Username: "user1user" + tmpNumb,
		Password: "Lhjxb[eq" + tmpNumb,
	}

	authReqJson, err := json.Marshal(authReq)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(authReqJson)))
	router.ServeHTTP(w, req)

	if !assert.Equal(t, 200, w.Code) {
		t.Log(w.Body.String())
	}
	t.Log(w.Body.String())

}

func TestAuthServer_Register_BadPass(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}
	tmpNumb := strconv.Itoa(rand.Intn(1000) + 1000)

	authReq := structs.RegisterReqBody{
		Username: "user1user" + tmpNumb,
		Password: "aaa",
	}

	authReqJson, err := json.Marshal(authReq)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(authReqJson)))
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		t.Log(w.Body.String())
	}
	t.Log(w.Body.String())

}

func TestAuthServer_Register_BadLogin(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}
	tmpNumb := strconv.Itoa(rand.Intn(1000) + 1000)

	authReq := structs.RegisterReqBody{
		Username: "aaaa",
		Password: "Lhjxb[eq" + tmpNumb,
	}

	authReqJson, err := json.Marshal(authReq)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(authReqJson)))
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusBadRequest, w.Code) {
		t.Log(w.Body.String())
	}
	t.Log(w.Body.String())

}

func TestAuthServer_Register_WrongPass(t *testing.T) {
	router, err := initServer()
	if err != nil {
		t.Error(err)
	}
	tmpNumb := strconv.Itoa(rand.Intn(1000) + 1000)

	authReq := structs.RegisterReqBody{
		Username: "user1user1",
		Password: "Lhjxb[eq" + tmpNumb,
	}

	authReqJson, err := json.Marshal(authReq)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(authReqJson)))
	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusUnauthorized, w.Code) {
		t.Log(w.Body.String())
	}
	t.Log(w.Body.String())
}
