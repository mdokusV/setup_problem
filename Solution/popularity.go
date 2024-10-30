package solution

import (
	"first/models"
	"sort"
)

func PopularitySolution(state State) models.CMaxValue {
	workers := state.workers
	tasks := state.tasks
	setups := state.setups

	for s := range setups {
		setup := &setups[s]
		setup.alphaAverageTasks(1)
	}

	for t := range tasks {
		task := &tasks[t]
		task.alphaAverageSetups(1)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].setupPopularity+float64(tasks[i].time) < tasks[j].setupPopularity+float64(tasks[i].time)
	})

	for t := range tasks {
		task := &tasks[t]
		sort.Slice(workers, func(i, j int) bool {
			return workers[i].predictSum < workers[j].predictSum
		})
		minWorker := &workers[0]
		minWorker.addTaskWithPredict(task)
	}
	sort.Slice(workers, func(i, j int) bool {
		return workers[i].cSum > workers[j].cSum
	})
	value := workers[0].cSum

	return models.CMaxValue(value)
}
