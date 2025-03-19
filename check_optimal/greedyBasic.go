package solution

import (
	"cmp"
	"first/helpers"
	"first/models"
	"slices"
)

func GreedySolution(state State) (models.CMaxValue, State) {
	tasks := state.tasks
	workers := state.workers

	for t := range tasks {
		task := &tasks[t]
		minWorker := minimumHeuristicWorker(workers, task)
		minWorker.addTask(*task)
	}
	max := slices.MaxFunc(workers, func(i, j worker) int {
		return cmp.Compare(i.cSum, j.cSum)
	})

	return models.CMaxValue(max.cSum), state
}

func minimumHeuristicWorker(workers []worker, task *task) *worker {
	cSumCache, addedCache := make([]int, len(workers)), make([]int, len(workers))
	for w := range workers {
		worker := &workers[w]
		cSumCache[worker.id], addedCache[worker.id] = worker.testTask(task)
	}
	minWorker, _ := helpers.FindBestMatch(workers, func(i, j worker) int {
		return cmp.Or(
			cmp.Compare(cSumCache[i.id], cSumCache[j.id]),
			cmp.Compare(addedCache[i.id], addedCache[j.id]),
			cmp.Compare(i.id, j.id))
	})
	return minWorker
}
