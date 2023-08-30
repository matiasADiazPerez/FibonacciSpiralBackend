package spiral

import (
	"fmt"
	"net/http"
	"spiralmatrix/config"
	"spiralmatrix/internal/utils"
	"strconv"
)

// HandleSpiral godoc
//
//	@Summary		Generate Fibonacci Spiral
//	@Description	Generates a Matrix of size cols x rows filled with fibonacci numbers and ordered like a spiral
//	@Tags			spiral
//	@Accept			json
//	@Produce		json
//	@Param			cols	query		int	true	"columns"
//	@Param			rows	query		int	true	"rows"
//	@Success		200		{object}	[][]string{}
//	@Failure		400		{object}	utils.ErrorWrapper
//	@Failure		500		{object}	utils.ErrorWrapper
//	@Router			/spiral [get]
func (s *SpiralHandler) HandleSpiral(w http.ResponseWriter, r *http.Request) {
	colsParam := r.URL.Query().Get("cols")
	rowsParam := r.URL.Query().Get("rows")

	cols, err := strconv.Atoi(colsParam)
	if err != nil {
		utils.HandleError(utils.NewErrorWrapper(config.SPIRAL, http.StatusBadRequest, fmt.Errorf("%v is not a valid number for columns", colsParam)), w)
		return
	}
	rows, err := strconv.Atoi(rowsParam)
	if err != nil {
		utils.HandleError(utils.NewErrorWrapper(config.SPIRAL, http.StatusBadRequest, fmt.Errorf("%v is not a valid number for rows", rowsParam)), w)
		return
	}

	resp, errWrapper := s.Spiral(cols, rows)
	if errWrapper.Error != nil {
		utils.HandleError(errWrapper, w)
	}
	utils.CreateResponse("Success", resp, w)
}
