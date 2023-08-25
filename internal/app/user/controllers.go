package user

import (
	"encoding/json"
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
	if err != nil {
		return 0, utils.NewErrorWrapper(config.GET_ID, http.StatusBadRequest, fmt.Errorf("%s is not a valid number", userIdStr))
	}
	return userId, utils.ErrorWrapper{}
}

func (u *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var userCreate models.CreateUser
	err := json.NewDecoder(r.Body).Decode(&userCreate)
	if err != nil {
		utils.HandleError(utils.NewErrorWrapper(config.CREATE_ERROR, http.StatusInternalServerError, err), w)
		return
	}
	resp, errWrapper := u.createUser(userCreate)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
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
	err := json.NewDecoder(r.Body).Decode(&changePassword)
	if err != nil {
		utils.HandleError(utils.NewErrorWrapper(config.CHANGE_PASS, http.StatusInternalServerError, err), w)
		return
	}
	errWrapper = u.changePassword(userId, changePassword)
	if errWrapper.Error != nil {
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
	errWrapper = u.deleteUser(userId)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	utils.CreateResponse("Success", "", w)
}
