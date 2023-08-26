package auth

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"spiralmatrix/config"
	"spiralmatrix/internal/models"
	"spiralmatrix/internal/utils"
	"strings"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) AuthHandler {
	return AuthHandler{
		db: db,
	}
}

func getSecret() (ed25519.PrivateKey, utils.ErrorWrapper) {
	var secret ed25519.PrivateKey
	secretString := os.Getenv("JWT_SECRET")
	if secretString == "" {
		return secret, utils.NewErrorWrapper(config.LOGIN, http.StatusUnauthorized, fmt.Errorf("environment var missing"))
	}
	seed64, err := utils.HashPassword(secretString)
	if secretString == "" {
		return secret, utils.NewErrorWrapper(config.LOGIN, http.StatusUnauthorized, err)
	}
	seed32 := []byte(seed64)[:32]
	secret = ed25519.NewKeyFromSeed(seed32)
	return secret, utils.ErrorWrapper{}
}

func (a *AuthHandler) verifyUser(authUser models.AuthUser) utils.ErrorWrapper {
	user := models.User{}
	result := a.db.First(&user, authUser.Id)
	if result.Error != nil {
		return utils.NewErrorWrapper(config.LOGIN, http.StatusUnauthorized, result.Error)
	}
	hashPass, err := utils.HashPassword(authUser.Password)
	if err != nil {
		return utils.NewErrorWrapper(config.LOGIN, http.StatusUnauthorized, result.Error)
	}
	if hashPass != user.Password {
		return utils.NewErrorWrapper(config.LOGIN, http.StatusUnauthorized, fmt.Errorf("Failed login"))
	}
	return utils.ErrorWrapper{}
}

func createJWT(id int) (string, utils.ErrorWrapper) {
	newJWT := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"user": id,
	})
	secret, errWrapper := getSecret()
	if errWrapper.Error != nil {
		return "", errWrapper
	}
	signedJWT, err := newJWT.SignedString(secret)
	if err != nil {

		return "", utils.NewErrorWrapper(config.LOGIN, http.StatusUnauthorized, err)
	}
	return signedJWT, utils.ErrorWrapper{}
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var authUser models.AuthUser
	err := json.NewDecoder(r.Body).Decode(&authUser)
	if err != nil {
		utils.HandleError(utils.NewErrorWrapper(config.LOGIN, http.StatusUnauthorized, err), w)
		return
	}
	errWrapper := a.verifyUser(authUser)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}

	token, errWrapper := createJWT(authUser.Id)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	utils.CreateResponse("Success", token, w)
}

func (a *AuthHandler) AuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqToken := r.Header.Get("Authorization")

			splitToken := strings.Split(reqToken, " ")
			if splitToken[0] != "Bearer" || len(splitToken) != 2 {
				utils.HandleError(utils.NewErrorWrapper(config.AUTH, http.StatusUnauthorized, fmt.Errorf("Malformed token")), w)
				return
			}
			tokenStr := splitToken[1]
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				secret, errWrapper := getSecret()
				pubKey := secret.Public()
				return pubKey, errWrapper.Error
			})
			if err != nil {
				utils.HandleError(utils.NewErrorWrapper(config.AUTH, http.StatusUnauthorized, err), w)
				return
			}
			if !token.Valid {
				utils.HandleError(utils.NewErrorWrapper(config.AUTH, http.StatusUnauthorized, fmt.Errorf("Invalid token")), w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
