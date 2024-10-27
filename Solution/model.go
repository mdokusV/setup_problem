package solution

import mapset "github.com/deckarep/golang-set/v2"

type state struct {
	workers []worker
	tasks   []task
	setups  []setup
}

type worker struct {
	setups mapset.Set[int]
	tasks  []*task
	cSum   int
}

type task struct {
	id     int
	time   int
	setups []*setup
}

type setup struct {
	id    int
	time  int
	tasks []*task
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
