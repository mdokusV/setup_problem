package solution

import (
	"first/models"
	"sort"
	"time"
)

func GreedySolution(initialState models.InitialState) (models.CMaxValue, time.Duration) {
	start := time.Now()
	state := transformInitialState(initialState)
	workers := state.workers
	tasks := state.tasks

	for t := range tasks {
		task := &tasks[t]
		sort.Slice(workers, func(i, j int) bool {
			return workers[i].cSum < workers[j].cSum
		})
		minWorker := &workers[0]
		minWorker.addTask(task)
	}
	sort.Slice(workers, func(i, j int) bool {
		return workers[i].cSum > workers[j].cSum
	})

	value := workers[0].cSum

	elapsed := time.Since(start)
	return models.CMaxValue(value), elapsed
}
