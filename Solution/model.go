package solution

import mapset "github.com/deckarep/golang-set/v2"

type State struct {
	workers []worker
	tasks   []task
	setups  []setup
}

type worker struct {
	id         int
	setups     mapset.Set[int]
	tasks      []*task
	predictSum float64
	cSum       int
}

func (w *worker) addTaskWithPredict(t *task) {
	w.tasks = append(w.tasks, t)
	w.cSum += t.time
	w.predictSum += t.setupPopularity + float64(t.time)
	for _, s := range t.setups {
		if !w.setups.ContainsOne(s.id) {
			w.setups.Add(s.id)
			w.cSum += s.time
			w.predictSum += float64(s.time)
		}
	}
}

func (w *worker) testTask(t *task) int {
	sum := w.cSum
	sum += t.time

	for _, s := range t.setups {
		if !w.setups.ContainsOne(s.id) {
			sum += s.time
		}
	}
	return sum

}

func (w *worker) addTask(t *task) {
	w.tasks = append(w.tasks, t)
	w.cSum += t.time
	for _, s := range t.setups {
		if !w.setups.ContainsOne(s.id) {
			w.setups.Add(s.id)
			w.cSum += s.time
		}
	}
}

type task struct {
	id              int
	time            int
	setupPopularity float64
	setups          []*setup
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
	tasks       []*task
}

func (s *setup) alphaAverageTasks(alpha float64) {
	s.taskAverage = 0
	setupsValues := make([]int, 0, len(s.tasks))
	for _, s := range s.tasks {
		setupsValues = append(setupsValues, s.time)
	}
	s.taskAverage = alphaAverage(setupsValues, alpha)
}
