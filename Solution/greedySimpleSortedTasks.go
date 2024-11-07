package solution

import (
	"cmp"
	"first/models"
	"slices"
)

func GreedyWithSimpleSortedTasks(state State) models.CMaxValue {
	tasks := state.tasks
	workers := state.workers

	for t := range tasks {
		task := &tasks[t]
		task.sumTimeWithSetups()
	}

	slices.SortFunc(tasks, func(i, j task) int {
		return cmp.Or(cmp.Compare(j.timeWithSetup, i.timeWithSetup), cmp.Compare(j.time, i.time))
	})

	for t := range tasks {
		task := &tasks[t]
		minWorker := minimumHeuristicWorker(workers, task)
		minWorker.addTask(*task)
	}
	slices.SortFunc(workers, func(i, j worker) int {
		return cmp.Compare(j.cSum, i.cSum)
	})

	value := workers[0].cSum
	return models.CMaxValue(value)
}
