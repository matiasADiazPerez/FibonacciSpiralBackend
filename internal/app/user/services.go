package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"spiralmatrix/config"
	"spiralmatrix/internal/models"
	"spiralmatrix/internal/utils"

	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) UserHandler {
	return UserHandler{
		db: db,
	}
}

func hashPassword(password string) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil

}

func (u *UserHandler) createUser(input models.CreateUser) (models.User, utils.ErrorWrapper) {
	newPass, err := hashPassword(input.Password)
	if err != nil {
		return models.User{}, utils.NewErrorWrapper(config.CREATE_ERROR, 0, err)
	}
	newUser := models.User{
		Name:     input.Name,
		Password: newPass,
	}
	result := u.db.Create(&newUser)
	if result.Error != nil {
		return models.User{}, utils.NewErrorWrapper(config.CREATE_ERROR, 0, result.Error)
	}
	return newUser, utils.ErrorWrapper{}
}

func (u *UserHandler) findAllUsers() ([]models.User, utils.ErrorWrapper) {
	users := []models.User{}
	result := u.db.Omit("Password").Find(&users).Where("Where DeletedAt is NULL")
	if result.Error != nil {
		return []models.User{}, utils.NewErrorWrapper(config.FIND_ERROR, 0, result.Error)
	}
	return users, utils.ErrorWrapper{}
}

func (u *UserHandler) findById(id int) (models.User, utils.ErrorWrapper) {
	user := models.User{}
	code := 0
	result := u.db.Omit("Password").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		}
		return models.User{}, utils.NewErrorWrapper(config.FIND_ERROR, code, result.Error)
	}
	return user, utils.ErrorWrapper{}
}

func (u *UserHandler) changePassword(id int, passwordInput models.ChangePassword) utils.ErrorWrapper {
	user := models.User{}
	code := 0
	result := u.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		}
		return utils.NewErrorWrapper(config.CHANGE_PASS, code, result.Error)
	}

	currentPass, err := hashPassword(passwordInput.CurrentPassword)
	if err != nil {
		return utils.NewErrorWrapper(config.CHANGE_PASS, 0, err)
	}
	if currentPass != user.Password {
		return utils.NewErrorWrapper(config.CHANGE_PASS, http.StatusUnauthorized, fmt.Errorf("Wrong password"))
	}

	newPass, err := hashPassword(passwordInput.NewPassword)
	if err != nil {
		return utils.NewErrorWrapper(config.CHANGE_PASS, 0, err)
	}
	user.Password = newPass
	u.db.Save(user)
	return utils.ErrorWrapper{}
}

func (u *UserHandler) deleteUser(id int) utils.ErrorWrapper {
	user := models.User{}
	code := 0
	result := u.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
		}
		return utils.NewErrorWrapper(config.CHANGE_PASS, code, result.Error)
	}
	u.db.Delete(&user)
	return utils.ErrorWrapper{}
}
