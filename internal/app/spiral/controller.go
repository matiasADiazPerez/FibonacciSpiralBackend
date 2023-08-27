package spiral

import (
	"fmt"
	"net/http"
	"spiralmatrix/config"
	"spiralmatrix/internal/utils"
	"strconv"
)

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
