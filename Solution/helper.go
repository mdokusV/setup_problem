package solution

import (
	"first/models"

	mapset "github.com/deckarep/golang-set/v2"
)

func transformInitialState(initState models.InitialState) state {
	workers := make([]worker, initState.WorkerNumber)
	for i := range workers {
		workers[i].setups = mapset.NewSet[int]()
	}

	// generate tasks
	tasks := make([]task, len(initState.Tasks))
	setups := make([]setup, len(initState.Setups))
	for i, t := range initState.Tasks {
		tasks[i] = task{id: t.ID, time: t.Time, setups: make([]*setup, len(t.Setups))}
	}

	// generate setups
	for i, s := range initState.Setups {
		setups[i] = setup{id: s.ID, time: s.Time, tasks: make([]*task, len(s.Tasks))}
	}

	//populate tasks
	for _, s := range initState.Setups {
		for j, t := range s.Tasks {
			setups[s.ID].tasks[j] = &tasks[t.ID]
		}
	}

	// populate setups
	for _, t := range initState.Tasks {
		for j, s := range t.Setups {
			tasks[t.ID].setups[j] = &setups[s.ID]
		}
	}

	return state{
		workers: workers,
		tasks:   tasks,
		setups:  setups,
	}
}
