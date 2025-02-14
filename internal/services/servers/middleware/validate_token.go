package middleware

import (
	"fmt"
	"net/http"

	"github.com/Kapeland/task-Avito/internal/models"
	structs2 "github.com/Kapeland/task-Avito/internal/models/structs"
	"github.com/Kapeland/task-Avito/internal/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

func CheckJWT(a models.AuthModelManager, lgr *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			//TODO: Здесь ещё должен быть нужный метод, чтобы передать управление регистрации + сама регистрация
		}
		tokenString = tokenString[7:] // Strip the "Bearer " prefix from the token.

		claims, err := getUnverifiedTokenClaims(tokenString, lgr)
		if err != nil {
			lgr.Error(err.Error(), "validate_token", "CheckJWT", "getUnverifiedTokenClaims")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		login := claims["sub"].(string)
		sessionID := claims["sID"].(string)

		userSecret, err := a.GetUserSecretByLoginAndSession(c.Request.Context(), structs2.UserSecret{Login: login, SessionID: sessionID})
		if err != nil {
			lgr.Error(err.Error(), "validate_token", "CheckJWT", "GetUserSecretByLoginAndSession")
			//TODO: точно ли этот код?
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(userSecret.Secret), nil
		})

		switch {
		case token.Valid:
			break
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			lgr.Error(err.Error(), "validate_token", "CheckJWT", "Parse")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		default:
			lgr.Error(err.Error(), "validate_token", "CheckJWT", "Parse")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": "Invalid token"})
			return
		}
		c.Set("login", login)
		c.Next()
	}
}

func getUnverifiedTokenClaims(tokenStr string, lgr *logger.Logger) (jwt.MapClaims, error) {
	parser := jwt.Parser{}
	unverToken, _, err := parser.ParseUnverified(tokenStr,
		jwt.MapClaims{
			"sub": "",
			"sID": "",
			"exp": "",
		})
	if err != nil {
		lgr.Error(err.Error(), "validate_token", "getUnverifiedTokenClaims", "ParseUnverified")

		return nil, models.ErrInvalidToken
	}

	claims, ok := unverToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, models.ErrInvalidToken
	}
	return claims, nil
}
