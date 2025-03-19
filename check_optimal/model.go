package solution

import (
	"cmp"
	"first/models"
	"fmt"
	"slices"
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
)

type State struct {
	workers []worker
	tasks   []task
	setups  []setup
}

func prepareStartState(initState models.InitialState) State {
	workers := make([]worker, initState.WorkerNumber)
	for i := range workers {
		workers[i].id = i
		workers[i].setups = mapset.NewSet[*setup]()
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
		tasks[t.ID].setupsSet = mapset.NewSet[*setup](tasks[t.ID].setups...)

	}

	return State{workers, tasks, setups}
}

func (state State) IsCorrect() {
	allTasksOnce := make([]bool, len(state.tasks))
	for _, t := range state.tasks {
		if allTasksOnce[t.id] {
			panic("task used twice")
		}
		allTasksOnce[t.id] = true
	}
	if len(allTasksOnce) != len(state.tasks) {
		panic("not all tasks used")
	}
}

func (state *State) checkSolution() {
	workers := state.workers
	tasks := state.tasks
	// all tasks used by workers
	usedTasks := make(map[int]bool, len(workers))
	usedTasksNum := 0
	for _, w := range workers {
		for _, t := range w.tasks {
			if usedTasks[t.id] {
				panic("task already used")
			}

			//all setups assigned
			for _, s := range t.setups {
				if !w.setups.ContainsOne(s) {
					panic("missing setup")
				}
			}

			usedTasksNum++
			usedTasks[t.id] = true
		}

		// correct cSum
		sum := 0
		for _, t := range w.tasks {
			sum += t.time
		}
		for _, s := range w.setups.ToSlice() {
			sum += s.time
		}
		if sum != w.cSum {
			panic("wrong cSum")
		}
	}

	if usedTasksNum != len(tasks) {
		panic("not all tasks used")
	}
}

func (state *State) Print() {
	slices.SortFunc(state.workers, func(i, j worker) int {
		return cmp.Compare(i.id, j.id)
	})
	for _, w := range state.workers {
		w.print()
	}
}

type worker struct {
	id         int
	setups     mapset.Set[*setup]
	tasks      []task
	predictSum float64
	cSum       int
}

func (w *worker) Reset() {
	w.setups.Clear()
	w.tasks = nil
	w.predictSum = 0
	w.cSum = 0
}

func (w *worker) print() {
	sort.Slice(w.tasks, func(i, j int) bool {
		return w.tasks[i].id < w.tasks[j].id
	})

	fmt.Printf("worker: %1v cSum: %4v | setups ", w.id, w.cSum)
	setups := w.setups.ToSlice()
	sort.Slice(setups, func(i, j int) bool {
		return setups[i].id < setups[j].id
	})
	for _, s := range setups {
		s.print()
	}
	fmt.Printf(" | tasks ")
	sort.Slice(w.tasks, func(i, j int) bool {
		return w.tasks[i].id < w.tasks[j].id
	})

	for _, t := range w.tasks {
		t.print()
	}

	fmt.Printf("\n")
}

func (w worker) copy() worker {
	tasks := make([]task, len(w.tasks))
	copy(tasks, w.tasks)
	clone := worker{
		id:         w.id,
		setups:     w.setups.Clone(),
		tasks:      tasks,
		predictSum: w.predictSum,
		cSum:       w.cSum,
	}
	return clone
}

func (w *worker) addTaskWithPredict(t task) {
	w.tasks = append(w.tasks, t)
	w.cSum += t.time
	w.predictSum += t.setupPopularity + float64(t.time)
	for _, s := range t.setups {
		if !w.setups.ContainsOne(s) {
			w.setups.Add(s)
			w.cSum += s.time
			w.predictSum += float64(s.time)
		}
	}
}

func (w *worker) testTask(t *task) (int, int) {
	sum := w.cSum
	sum += t.time

	added := t.time

	for _, s := range t.setups {
		if !w.setups.ContainsOne(s) {
			added += s.time
			sum += s.time
		}
	}
	return sum, added

}

func (w *worker) addTask(t task) {
	w.tasks = append(w.tasks, t)
	w.cSum += t.time
	for _, s := range t.setups {
		if w.setups.Add(s) {
			w.cSum += s.time
		}
	}
}

type task struct {
	id              int
	time            int
	timeWithSetup   int
	setupPopularity float64
	setups          []*setup
	setupsSet       mapset.Set[*setup]
}

func (t *task) print() {
	fmt.Printf("%2v(%2v),", t.id, t.time)
}

func (t *task) printSetups() {
	fmt.Printf("%2v:", t.id)
	for _, s := range t.setups {
		fmt.Printf("%2v,", s.id)
	}
	fmt.Printf("\n")
}

func (t *task) sumTimeWithSetups() {
	t.timeWithSetup = t.time
	for _, s := range t.setups {
		t.timeWithSetup += s.time
	}
}

type setup struct {
	id          int
	time        int
	taskAverage float64
	tasks       []*task
}

func (s *setup) print() {
	fmt.Printf("%2v(%2v),", s.id, s.time)
}
