package spiral

import (
	"fmt"
	"net/http"
	"spiralmatrix/config"
	"spiralmatrix/internal/models"
	"spiralmatrix/internal/utils"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type spiralTestStruct struct {
	cols             int
	rows             int
	expectedResponse [][]string
	expectedError    utils.ErrorWrapper
}

func TestCreateUser(t *testing.T) {
	service := NewSpiralHandler()
	testCases := []spiralTestStruct{
		{
			cols:             3,
			rows:             3,
			expectedResponse: [][]string{{"0", "1", "1"}, {"13", "21", "2"}, {"8", "5", "3"}},
			expectedError:    utils.ErrorWrapper{},
		},
		{
			cols:             3,
			rows:             1,
			expectedResponse: [][]string{{"0", "1", "1"}},
			expectedError:    utils.ErrorWrapper{},
		},
		{
			cols:             1,
			rows:             3,
			expectedResponse: [][]string{{"0"}, {"1"}, {"1"}},
			expectedError:    utils.ErrorWrapper{},
		},
		{
			cols:             0,
			rows:             3,
			expectedResponse: [][]string{{}, {}, {}},
			expectedError:    utils.ErrorWrapper{},
		},
		{
			cols:             0,
			rows:             0,
			expectedResponse: [][]string{},
			expectedError:    utils.ErrorWrapper{},
		},
		{
			cols:             3,
			rows:             -1,
			expectedResponse: [][]string{},
			expectedError:    utils.NewErrorWrapper(config.SPIRAL, http.StatusBadRequest, fmt.Errorf("Invalid values for rows: %v and cols: %v", -1, 3)),
		},
	}
	t.Run("create user test", func(t *testing.T) {

		for _, tc := range testCases {
			resp, errWrapper := service.Spiral(tc.cols, tc.rows)
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
