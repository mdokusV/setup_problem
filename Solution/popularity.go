package solution

import (
	"first/models"
	"time"
)

func popularitySolution(initialState models.InitialState) (models.CMaxValue, time.Duration) {
	start := time.Now()
	state := transformInitialState(initialState)

	elapsed := time.Since(start)
	return models.CMaxValue(len(state.workers)), elapsed
}
