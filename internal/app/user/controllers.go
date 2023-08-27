package user

import (
	"fmt"
	"net/http"
	"spiralmatrix/config"
	"spiralmatrix/internal/models"
	"spiralmatrix/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func getUserId(r *http.Request) (int, utils.ErrorWrapper) {

	userIdStr := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId < 0 {
		return 0, utils.NewErrorWrapper(config.GET_ID, http.StatusBadRequest, fmt.Errorf("%s is not a valid number", userIdStr))
	}
	return userId, utils.ErrorWrapper{}
}

func (u *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var userCreate models.CreateUser
	userCreate, errWrapper := utils.GetBody(r, userCreate)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	resp, errWrapper := u.createUser(userCreate)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	resp.Password = ""
	utils.CreateResponse("User created", resp, w)
}

func (u *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	resp, errWrapper := u.findAllUsers()
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	utils.CreateResponse("Success", resp, w)
}

func (u *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userId, errWrapper := getUserId(r)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	resp, errWrapper := u.findById(userId)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	utils.CreateResponse("Success", resp, w)
}

func (u *UserHandler) HandleChangePass(w http.ResponseWriter, r *http.Request) {
	userId, errWrapper := getUserId(r)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	var changePassword models.ChangePassword
	changePassword, errWrapper = utils.GetBody(r, changePassword)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	if errWrapper = u.changePassword(userId, changePassword); errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	utils.CreateResponse("Success", "", w)
}

func (u *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, errWrapper := getUserId(r)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	if errWrapper = u.deleteUser(userId); errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	utils.CreateResponse("Success", "", w)
}
