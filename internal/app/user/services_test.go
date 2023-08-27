package user

import (
	"database/sql/driver"
	"fmt"
	"net/http"
	"spiralmatrix/config"
	"spiralmatrix/internal/models"
	"spiralmatrix/internal/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type createUserTestStruct struct {
	input            models.CreateUser
	expectedResponse models.User
	expectedError    utils.ErrorWrapper
	dbErr            error
}

type findAllUsersTestStruct struct {
	expectedResponse []models.User
	expectedError    utils.ErrorWrapper
	mockRows         *sqlmock.Rows
	dbErr            error
}

type findByIdTestStruct struct {
	id               int
	expectedResponse models.User
	expectedError    utils.ErrorWrapper
	mockRows         *sqlmock.Rows
	dbErr            error
}

type changePasswordTestStruct struct {
	id            int
	passwordInput models.ChangePassword
	expectedError utils.ErrorWrapper
	mockRows      *sqlmock.Rows
	dbErr         error
	updateResult  driver.Result
	expectCommit  bool
}

type deleteUserStruct struct {
	id            int
	expectedError utils.ErrorWrapper
	mockRows      *sqlmock.Rows
	expectCommit  bool
	dbErr         error
}

func getMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	return db, mock
}

func TestCreateUser(t *testing.T) {
	db, mock := getMockDB()
	service := NewUserHandler(db)
	normalUser := models.User{
		Name:     "test",
		Password: "af0426e71dd57c0fdf93f23f6f191a4aa0578ad7d71897e936746028b8ffd31d",
	}
	normalUser.ID = 1
	testCases := []createUserTestStruct{
		{
			input: models.CreateUser{
				Name:     "test",
				Password: "testPass",
			},
			expectedResponse: normalUser,
		},
		{
			input: models.CreateUser{
				Name:     "test",
				Password: "testPass",
			},
			dbErr:            fmt.Errorf("db fail"),
			expectedResponse: models.User{},
			expectedError:    utils.NewErrorWrapper(config.CREATE_ERROR, 500, fmt.Errorf("db fail")),
		},
	}
	t.Run("create user test", func(t *testing.T) {

		for _, tc := range testCases {
			rows := sqlmock.NewRows([]string{"Name", "ID"}).AddRow("test", 1)
			mock.ExpectBegin().WillReturnError(tc.dbErr)
			mock.ExpectQuery("INSERT").WillReturnRows(rows)
			mock.ExpectCommit()
			resp, errWrapper := service.createUser(tc.input)
			if diff := cmp.Diff(resp, tc.expectedResponse, cmpopts.IgnoreFields(models.User{}, "CreatedAt", "UpdatedAt")); diff != "" {
				t.Errorf(diff)
			}
			if diff := cmp.Diff(errWrapper, tc.expectedError, cmpopts.IgnoreFields(utils.ErrorWrapper{}, "Error")); diff != "" {
				t.Errorf(diff)
			}
			if errWrapper.Error != tc.expectedError.Error {
				if errWrapper.Error.Error() != tc.expectedError.Error.Error() {
					t.Errorf("got %v, want %v", errWrapper.Error, tc.expectedError.Error)
				}

			}

		}
	})
}
func TestFindAllUsers(t *testing.T) {
	db, mock := getMockDB()
	service := NewUserHandler(db)
	normalUser := models.User{
		Name:     "test",
		Password: "af0426e71dd57c0fdf93f23f6f191a4aa0578ad7d71897e936746028b8ffd31d",
	}
	normalUser.ID = 1
	testCases := []findAllUsersTestStruct{
		{
			expectedResponse: []models.User{
				normalUser,
			},
			mockRows: sqlmock.NewRows([]string{"Name", "ID", "Password"}).AddRow(normalUser.Name, normalUser.ID, normalUser.Password),
		},
		{
			dbErr:            fmt.Errorf("db fail"),
			mockRows:         &sqlmock.Rows{},
			expectedResponse: []models.User{},
			expectedError:    utils.NewErrorWrapper(config.FIND_ERROR, 500, fmt.Errorf("db fail")),
		},
	}
	t.Run("find all users test", func(t *testing.T) {

		for _, tc := range testCases {
			mock.ExpectQuery("SELECT").WillReturnRows(tc.mockRows).WillReturnError(tc.dbErr)
			resp, errWrapper := service.findAllUsers()
			if diff := cmp.Diff(resp, tc.expectedResponse); diff != "" {
				t.Errorf(diff)
			}
			if diff := cmp.Diff(errWrapper, tc.expectedError, cmpopts.IgnoreFields(utils.ErrorWrapper{}, "Error")); diff != "" {
				t.Errorf(diff)
			}
			if errWrapper.Error != tc.expectedError.Error {
				if errWrapper.Error.Error() != tc.expectedError.Error.Error() {
					t.Errorf("got %v, want %v", errWrapper.Error, tc.expectedError.Error)
				}

			}

		}
	})
}

func TestFindById(t *testing.T) {
	db, mock := getMockDB()
	service := NewUserHandler(db)
	normalUser := models.User{
		Name:     "test",
		Password: "af0426e71dd57c0fdf93f23f6f191a4aa0578ad7d71897e936746028b8ffd31d",
	}
	normalUser.ID = 1
	testCases := []findByIdTestStruct{
		{
			id:               1,
			expectedResponse: normalUser,
			mockRows:         sqlmock.NewRows([]string{"Name", "ID", "Password"}).AddRow(normalUser.Name, normalUser.ID, normalUser.Password),
		},
		{
			dbErr:            fmt.Errorf("db fail"),
			mockRows:         &sqlmock.Rows{},
			expectedResponse: models.User{},
			expectedError:    utils.NewErrorWrapper(config.FIND_ERROR, 500, fmt.Errorf("db fail")),
		},
		{
			dbErr:            gorm.ErrRecordNotFound,
			mockRows:         &sqlmock.Rows{},
			expectedResponse: models.User{},
			expectedError:    utils.NewErrorWrapper(config.FIND_ERROR, http.StatusNotFound, gorm.ErrRecordNotFound),
		},
	}
	t.Run("find user by test", func(t *testing.T) {

		for _, tc := range testCases {
			mock.ExpectQuery("SELECT").WillReturnRows(tc.mockRows).WillReturnError(tc.dbErr)
			resp, errWrapper := service.findById(tc.id)
			if diff := cmp.Diff(resp, tc.expectedResponse); diff != "" {
				t.Errorf(diff)
			}
			if diff := cmp.Diff(errWrapper, tc.expectedError, cmpopts.IgnoreFields(utils.ErrorWrapper{}, "Error")); diff != "" {
				t.Errorf(diff)
			}
			if errWrapper.Error != tc.expectedError.Error {
				if errWrapper.Error.Error() != tc.expectedError.Error.Error() {
					t.Errorf("got %v, want %v", errWrapper.Error, tc.expectedError.Error)
				}

			}

		}
	})
}
func TestChangePassword(t *testing.T) {
	normalUser := models.User{
		Name:     "test",
		Password: "af0426e71dd57c0fdf93f23f6f191a4aa0578ad7d71897e936746028b8ffd31d",
	}
	normalUser.ID = 1
	testCases := []changePasswordTestStruct{
		{
			id: 1,
			passwordInput: models.ChangePassword{
				NewPassword:     "changedPasswrd",
				CurrentPassword: "testPass",
			},
			mockRows:     sqlmock.NewRows([]string{"Name", "ID", "Password"}).AddRow(normalUser.Name, normalUser.ID, normalUser.Password),
			expectCommit: true,
			updateResult: sqlmock.NewResult(1, 1),
		},
		{
			id: 1,
			passwordInput: models.ChangePassword{
				NewPassword:     "changedPasswrd",
				CurrentPassword: "wrongPassword",
			},
			mockRows:      sqlmock.NewRows([]string{"Name", "ID", "Password"}).AddRow(normalUser.Name, normalUser.ID, normalUser.Password),
			expectedError: utils.NewErrorWrapper(config.CHANGE_PASS, http.StatusUnauthorized, fmt.Errorf("Wrong password")),
			updateResult:  sqlmock.NewResult(1, 1),
		},
		{
			dbErr:         fmt.Errorf("db fail"),
			mockRows:      &sqlmock.Rows{},
			expectedError: utils.NewErrorWrapper(config.CHANGE_PASS, 0, fmt.Errorf("db fail")),
		},
		{
			dbErr:         gorm.ErrRecordNotFound,
			mockRows:      &sqlmock.Rows{},
			expectedError: utils.NewErrorWrapper(config.CHANGE_PASS, http.StatusNotFound, gorm.ErrRecordNotFound),
		},
	}
	t.Run("change password test", func(t *testing.T) {

		for _, tc := range testCases {
			db, mock := getMockDB()
			service := NewUserHandler(db)
			mock.ExpectQuery("SELECT").WillReturnRows(tc.mockRows).WillReturnError(tc.dbErr)
			if tc.expectCommit {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(tc.updateResult)
				mock.ExpectCommit()
			}
			errWrapper := service.changePassword(tc.id, tc.passwordInput)
			if diff := cmp.Diff(errWrapper, tc.expectedError, cmpopts.IgnoreFields(utils.ErrorWrapper{}, "Error")); diff != "" {
				t.Errorf(diff)
			}
			if errWrapper.Error != nil && tc.expectedError.Error == nil {
				t.Errorf("unexpected error: %v", errWrapper.Error.Error())
			} else if errWrapper.Error == nil && tc.expectedError.Error != nil {
				t.Errorf("was expecting error: %v", tc.expectedError.Error.Error())

			} else if errWrapper.Error != tc.expectedError.Error {
				if errWrapper.Error.Error() != tc.expectedError.Error.Error() {
					t.Errorf("got %v, want %v", errWrapper.Error, tc.expectedError.Error)
				}

			}

		}
	})
}
func TestDeleteUser(t *testing.T) {
	db, mock := getMockDB()
	service := NewUserHandler(db)
	normalUser := models.User{
		Name:     "test",
		Password: "af0426e71dd57c0fdf93f23f6f191a4aa0578ad7d71897e936746028b8ffd31d",
	}
	normalUser.ID = 1
	testCases := []deleteUserStruct{
		{
			id:           1,
			mockRows:     sqlmock.NewRows([]string{"Name", "ID", "Password"}).AddRow(normalUser.Name, normalUser.ID, normalUser.Password),
			expectCommit: true,
		},
		{
			dbErr:         fmt.Errorf("db fail"),
			mockRows:      &sqlmock.Rows{},
			expectedError: utils.NewErrorWrapper(config.DELETE_USER, 500, fmt.Errorf("db fail")),
		},
		{
			dbErr:         gorm.ErrRecordNotFound,
			mockRows:      &sqlmock.Rows{},
			expectedError: utils.NewErrorWrapper(config.DELETE_USER, http.StatusNotFound, gorm.ErrRecordNotFound),
		},
	}
	t.Run("find user by test", func(t *testing.T) {

		for _, tc := range testCases {
			mock.ExpectQuery("SELECT").WillReturnRows(tc.mockRows).WillReturnError(tc.dbErr)
			errWrapper := service.deleteUser(tc.id)

			if diff := cmp.Diff(errWrapper, tc.expectedError, cmpopts.IgnoreFields(utils.ErrorWrapper{}, "Error")); diff != "" {
				t.Errorf(diff)
			}
			if errWrapper.Error != nil && tc.expectedError.Error == nil {
				t.Errorf("unexpected error: %v", errWrapper.Error.Error())
			} else if errWrapper.Error == nil && tc.expectedError.Error != nil {
				t.Errorf("was expecting error: %v", tc.expectedError.Error.Error())

			} else if errWrapper.Error != tc.expectedError.Error {
				if errWrapper.Error.Error() != tc.expectedError.Error.Error() {
					t.Errorf("got %v, want %v", errWrapper.Error, tc.expectedError.Error)
				}

			}
		}
	})
}
