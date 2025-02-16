package solution

import (
	"first/models"
	"math"
	"math/rand"
	"slices"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

func randomInts(k uint, min, max int, rng *rand.Rand) []int {
	var ints = make([]int, k)
	for i := 0; i < int(k); i++ {
		ints[i] = i + min
	}
	for i := int(k); i < max-min; i++ {
		var j = rng.Intn(i + 1)
		if j < int(k) {
			ints[j] = i + min
		}
	}
	return ints
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

func alphaAverage[T int | float64](items []T, alpha float64) float64 {
	if alpha == math.Inf(1) {
		return float64(slices.Max(items))
	}
	if alpha == math.Inf(-1) {
		return float64(slices.Min(items))
	}

	sum := 0.0
	for _, v := range items {
		sum += math.Pow(float64(v), alpha)
	}
	sum /= float64(len(items))
	sum = math.Pow(sum, 1/alpha)
	return sum
}

func RunSolution(initialState models.InitialState, solution func(State) (models.CMaxValue, State)) (models.CMaxValue, time.Duration, State) {
	state := prepareStartState(initialState)
	start := time.Now()

	value, bestSolution := solution(state)

	elapsed := time.Since(start)

	bestSolution.checkSolution()

	return value, elapsed, bestSolution
}

func findBestMatch[S ~[]E, E any](x S, cmp func(a, b E) int) (*E, int) {
	if len(x) < 1 {
		panic("slices.MinFunc: empty list")
	}
	m := &x[0]
	ind := 0
	for i := 1; i < len(x); i++ {
		if cmp(x[i], *m) < 0 {
			m = &x[i]
			ind = i
		}
	}
	return m, ind
}

func removeNoOrder[T any](s []T, i int) []T {
	s[i], s[len(s)-1] = s[len(s)-1], s[i]
	return s[:len(s)-1]
}
