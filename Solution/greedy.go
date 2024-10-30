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
		minWorker := minimumHeuristicWorker(workers, task)
		minWorker.addTask(task)
	}
	sort.Slice(workers, func(i, j int) bool {
		return workers[i].cSum > workers[j].cSum
	})

	value := workers[0].cSum
	return models.CMaxValue(value)
}

func minimumHeuristicWorker(workers []worker, task *task) *worker {
	cacheHeuristic := make(map[int]int, len(workers))
	for w := range workers {
		worker := &workers[w]
		cacheHeuristic[worker.id] = worker.testTask(task)
	}
	minWorker := slices.MinFunc(workers, func(i, j worker) int {
		return cmp.Compare(cacheHeuristic[i.id], cacheHeuristic[j.id])
	})
	return &workers[minWorker.id]
}
