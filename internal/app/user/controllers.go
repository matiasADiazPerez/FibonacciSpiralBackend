package user

import (
	"fmt"
	"net/http"
	"spiralmatrix/config"
	_ "spiralmatrix/docs"
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

//	HandleCreateUser godoc
//
//	@Summary		Create a User
//	@Description	Store a new User in the database, this endpoint doesn't need authentication
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			models.CreateUser	body		models.CreateUser	true	"Add User"
//	@Success		200					{object}	models.User
//	@Failure		400					{object}	utils.ErrorWrapper
//	@Failure		500					{object}	utils.ErrorWrapper
//	@Router			/public/user [post]
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

//	HandleGetAllUsers godoc
//
//	@Summary		List all users
//	@Description	List all non deleted users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.User
//	@Failure		400	{object}	utils.ErrorWrapper
//	@Failure		404	{object}	utils.ErrorWrapper
//	@Failure		500	{object}	utils.ErrorWrapper
//	@Router			/user [get]
func (u *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	resp, errWrapper := u.findAllUsers()
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
		return
	}
	utils.CreateResponse("Success", resp, w)
}

//	HandleGetUser godoc
//
//	@Summary		Retrieves a user
//	@Description	Retrieves a non deleted user by id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int	true	"Get user"
//	@Success		200		{object}	models.User
//	@Failure		400		{object}	utils.ErrorWrapper
//	@Failure		404		{object}	utils.ErrorWrapper
//	@Failure		500		{object}	utils.ErrorWrapper
//	@Router			/user/{userId} [get]
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

//	HandleChangePass godoc
//
//	@Summary		Change Password of user
//	@Description	Verifies the current password and then changes the password of a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			models.ChangePassword	body		models.ChangePassword	true	"Change Password"
//	@Param			userId					path		int						true	"Get user"
//	@Success		200						string		Success
//	@Failure		400						{object}	utils.ErrorWrapper
//	@Failure		404						{object}	utils.ErrorWrapper
//	@Failure		500						{object}	utils.ErrorWrapper
//	@Router			/user/{userId} [patch]
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

//	HandleDeleteUser godoc
//
//	@Summary		Deletes a User
//	@Description	Soft deletes a User
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int	true	"Get user"
//	@Success		200		string		Success
//	@Failure		400		{object}	utils.ErrorWrapper
//	@Failure		500		{object}	utils.ErrorWrapper
//	@Router			/user [delete]
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
