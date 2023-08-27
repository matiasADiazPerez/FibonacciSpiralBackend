package spiral

import (
	"fmt"
	"math/big"
	"net/http"
	"spiralmatrix/config"
	"spiralmatrix/internal/utils"
)

type SpiralHandler struct {
}

func NewSpiralHandler() SpiralHandler {
	return SpiralHandler{}
}

func (s *SpiralHandler) Spiral(cols, rows int) ([][]string, utils.ErrorWrapper) {
	if cols < 0 || rows < 0 {
		return [][]string{}, utils.NewErrorWrapper(config.SPIRAL, http.StatusBadRequest, fmt.Errorf("Invalid values for rows: %v and cols: %v", rows, cols))
	}
	var left, top, right, bottom int = 0, 0, cols - 1, rows - 1
	ans := make([][]string, rows)
	for i := range ans {
		ans[i] = make([]string, cols)
	}

	current, next, sum := big.NewInt(0), big.NewInt(1), big.NewInt(1)

	for left <= right && top <= bottom {
		for i := left; i <= right; i++ {
			ans[top][i] = current.String()
			current.Set(next)
			next.Set(sum)
			sum.Add(current, next)
		}
		top++
		if top > bottom {
			break
		}
		for i := top; i <= bottom; i++ {
			ans[i][right] = current.String()
			current.Set(next)
			next.Set(sum)
			sum.Add(current, next)
		}
		right--
		if left > right {
			break
		}
		for i := right; i >= left; i-- {
			ans[bottom][i] = current.String()
			current.Set(next)
			next.Set(sum)
			sum.Add(current, next)
		}
		bottom--
		for i := bottom; i >= top; i-- {
			ans[i][left] = current.String()
			current.Set(next)
			next.Set(sum)
			sum.Add(current, next)
		}
		left++
	}
	return ans, utils.ErrorWrapper{}
}
