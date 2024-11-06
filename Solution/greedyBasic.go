package solution

import (
	"cmp"
	"first/models"
	"slices"
	"sort"
)

func GreedySolution(state State) models.CMaxValue {
	tasks := state.tasks
	workers := state.workers

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].time > tasks[j].time
	})

	for t := range tasks {
		task := &tasks[t]
		minWorkerID := minimumHeuristicWorker(workers, task)
		workers[minWorkerID].addTask(task)
	}
	slices.SortFunc(workers, func(i, j worker) int {
		return cmp.Compare(j.cSum, i.cSum)
	})

	value := workers[0].cSum
	return models.CMaxValue(value)
}

type workerID int

func minimumHeuristicWorker(workers []worker, task *task) workerID {
	cSumCache, addedCache := make([]int, len(workers)), make([]int, len(workers))
	for w := range workers {
		worker := &workers[w]
		cSumCache[worker.id], addedCache[worker.id] = worker.testTask(task)
	}
	minWorker := slices.MinFunc(workers, func(i, j worker) int {
		return cmp.Or(
			cmp.Compare(cSumCache[i.id], cSumCache[j.id]),
			cmp.Compare(addedCache[i.id], addedCache[j.id]),
			cmp.Compare(i.id, j.id))
	})
	return workerID(minWorker.id)
}
