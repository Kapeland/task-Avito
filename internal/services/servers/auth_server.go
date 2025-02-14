package servers

import (
	"context"
	"errors"
	structs2 "github.com/Kapeland/task-Avito/internal/services/structs"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"unicode"

	"github.com/Kapeland/task-Avito/internal/models"
	"github.com/Kapeland/task-Avito/internal/models/structs"
)

type AuthServer struct {
	A models.AuthModelManager
}

func isPasswordValid(s string) bool {
	if len(s) < 8 {
		return false
	}
	symbols := 0
	number, upper, lower, special := false, false, false, false
	for _, c := range s {
		switch {
		case unicode.IsDigit(c):
			number = true
			symbols++
		case unicode.IsUpper(c):
			upper = true
			symbols++
		case unicode.IsLower(c):
			lower = true
			symbols++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			symbols++
		case unicode.IsLetter(c):
			symbols++
		}
	}
	if symbols < 8 || !(number && special && upper && lower) {
		return false
	}

	return true
}

func isLoginValid(s string) bool {
	if len(s) < 8 {
		return false
	}
	for _, c := range s {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) {
			return false
		}
	}

	return true
}

func IsLoginPswdValid(login string, password string) bool {
	return isLoginValid(login) && isPasswordValid(password)
}

func (s *AuthServer) Register(c *gin.Context) {
	lgr := logger.GetLogger()

	var regInfo structs2.RegisterReqBody

	if err := c.ShouldBindBodyWithJSON(&regInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	login := regInfo.Username
	pswd := regInfo.Password

	if !IsLoginPswdValid(login, pswd) { // bad password or login
		lgr.Info("Bad pass or login", "authServer", "Register", "IsLoginPswdValid")

		c.JSON(http.StatusBadRequest, gin.H{"errors": "Bad pass or login"})
		return
	}

	userInfo := structs.RegisterUserInfo{
		Login: login,
		Pswd:  pswd,
	}
	tokStr, status := s.register(c.Request.Context(), userInfo)

	if status == http.StatusBadRequest {
		lgr.Info("item already exists", "authServer", "Register", "register")
		//TODO: Это не ошибка. Это значит, что уже регался и просто нужно вернуть JWT
		c.JSON(status, gin.H{"errors": "item already exists"})
		return
	}
	if status == http.StatusInternalServerError {
		lgr.Error("internal server error", "authServer", "Register", "register")
		c.JSON(status, gin.H{"errors": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokStr})
}

func (s *AuthServer) register(ctx context.Context, info structs.RegisterUserInfo) (string, int) {
	lgr := logger.GetLogger()

	tokStr, err := s.A.RegisterUser(ctx, info)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			return "", http.StatusBadRequest
		}

		lgr.Error(err.Error(), "authServer", "register", "RegisterUser")

		return "", http.StatusInternalServerError
	}

	return tokStr, http.StatusOK
}
