package solution

import mapset "github.com/deckarep/golang-set/v2"

type State struct {
	workers []worker
	tasks   []task
	setups  []setup
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

type worker struct {
	id         int
	setups     mapset.Set[*setup]
	tasks      []task
	predictSum float64
	cSum       int
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

func (t *task) sumTimeWithSetups() {
	t.timeWithSetup = t.time
	for _, s := range t.setups {
		t.timeWithSetup += s.time
	}
}

func (t *task) alphaAverageSetups(alpha float64) {
	t.setupPopularity = 0
	setupsValues := make([]float64, 0, len(t.setups))
	for _, s := range t.setups {
		setupsValues = append(setupsValues, s.taskAverage)
	}
	t.setupPopularity = alphaAverage(setupsValues, alpha)
}

type setup struct {
	id          int
	time        int
	taskAverage float64
	tasks       []task
}

func (s *setup) alphaAverageTasks(alpha float64) {
	s.taskAverage = 0
	setupsValues := make([]int, 0, len(s.tasks))
	for _, s := range s.tasks {
		setupsValues = append(setupsValues, s.time)
	}
	s.taskAverage = alphaAverage(setupsValues, alpha)
}
