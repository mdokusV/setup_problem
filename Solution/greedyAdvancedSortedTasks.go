package solution

import (
	"cmp"
	"first/models"
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
)

func GreedyWithAdvancedSortedTasks(state State) models.CMaxValue {
	tasks := state.tasks
	workers := state.workers

	for t := range tasks {
		task := &tasks[t]
		task.sumTimeWithSetups()
	}

	setupsUsedSet := mapset.NewSet[*setup]()
	// Find initial task to put

	bigFirst := slices.MaxFunc(tasks, func(i, j task) int {
		return cmp.Or(
			cmp.Compare(len(i.setups), len(j.setups)),
			cmp.Compare(i.timeWithSetup, j.timeWithSetup),
			cmp.Compare(i.time, j.time))
	})

	// put initial task
	workers[0].addTask(bigFirst)
	setupsUsedSet.Append(bigFirst.setups...)
	tasks = slices.DeleteFunc(tasks, func(i task) bool { return i.id == bigFirst.id })

	// fill first task in workers
	for {
		maxTask := slices.MaxFunc(tasks, func(i, j task) int {
			return cmp.Or(
				cmp.Compare(i.setupsSet.Difference(setupsUsedSet).Cardinality(), j.setupsSet.Difference(setupsUsedSet).Cardinality()),
				cmp.Compare(i.timeWithSetup, j.timeWithSetup),
				cmp.Compare(i.time, j.time))
		})
		if maxTask.setupsSet.IsSubset(setupsUsedSet) {
			break
		}

		minWorker := minimumHeuristicWorker(workers, &maxTask)
		minWorker.addTask(maxTask)
		setupsUsedSet.Append(maxTask.setups...)
		tasks = slices.DeleteFunc(tasks, func(i task) bool { return i.id == maxTask.id })
	}

	slices.SortFunc(tasks, func(i, j task) int {
		return cmp.Or(cmp.Compare(j.timeWithSetup, i.timeWithSetup),
			cmp.Compare(j.time, i.time))
	})

	for t := range tasks {
		task := &tasks[t]
		minWorkerID := minimumHeuristicWorker(workers, task)
		minWorkerID.addTask(*task)
	}
	slices.SortFunc(workers, func(i, j worker) int {
		return cmp.Compare(j.cSum, i.cSum)
	})

	value := workers[0].cSum
	return models.CMaxValue(value)
}
