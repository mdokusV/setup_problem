package solution

import (
	"first/models"
	"log"
	"math/rand"
	"slices"

	"github.com/MaxHalford/eaopt"
)

func GreedyGAinitStateSolution(state State) (models.CMaxValue, State) {
	ga, err := eaopt.NewDefaultGAConfig().NewGA()
	ga.PopSize = 40
	ga.NGenerations = 40
	ga.NPops = 2
	if err != nil {
		log.Fatal(err)
	}
	ga.Minimize(state.makeGAState)
	return models.CMaxValue(ga.HallOfFame[0].Fitness), State((ga.HallOfFame[0].Genome).(stateGA))
}

func (state State) makeGAState(rng *rand.Rand) eaopt.Genome {
	cloneState := stateGA(state).Clone()
	tasks := State(cloneState.(stateGA)).tasks

	rng.Shuffle(len(tasks), func(i, j int) { tasks[i], tasks[j] = tasks[j], tasks[i] })
	return cloneState
}

type stateGA State
type taskGA []task

func (t taskGA) At(i int) interface{} {
	return t[i]
}
func (t taskGA) ID(i int) int {
	return t[i].id
}

// Set method from Slice
func (t taskGA) Set(i int, v interface{}) {
	t[i] = v.(task)
}

// Len method from Slice
func (t taskGA) Len() int {
	return len(t)
}

// Swap method from Slice
func (t taskGA) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Slice method from Slice
func (t taskGA) Slice(a, b int) eaopt.Slice {
	return t[a:b]
}

// Split method from Slice
func (t taskGA) Split(k int) (eaopt.Slice, eaopt.Slice) {
	return t[:k], t[k:]
}

// Append method from Slice
func (t taskGA) Append(q eaopt.Slice) eaopt.Slice {
	return append(t, q.(taskGA)...)
}

// Replace method from Slice
func (p taskGA) Replace(q eaopt.Slice) {
	copy(p, q.(taskGA))
}

// Copy method from Slice
func (p taskGA) Copy() eaopt.Slice {
	var clone = make(taskGA, len(p))
	copy(clone, p)
	return clone
}

func (s stateGA) Evaluate() (mismatches float64, err error) {
	s.clear()
	eval, state := GreedySolution(State(s))
	state.CheckSolution(int(eval))

	return float64(eval), nil
}

func (s stateGA) clear() {
	for i := range s.workers {
		s.workers[i].Reset()
	}
}

func (s stateGA) Mutate(rng *rand.Rand) {
	for i := 0; i < 3; i++ {
		// Choose two points on the genome
		var points = randomInts(2, 0, len(s.tasks), rng)
		s.tasks[points[0]], s.tasks[points[1]] = s.tasks[points[1]], s.tasks[points[0]]
	}
	State(s).IsCorrect()
}

func (s stateGA) Crossover(q eaopt.Genome, rng *rand.Rand) {
	var indexes = randomInts(2, 1, len(s.tasks), rng)
	slices.Sort(indexes)
	a := indexes[0]
	b := indexes[1]
	p1 := taskGA(s.tasks)
	p2 := taskGA(q.(stateGA).tasks)
	var (
		n  = p1.Len()
		o1 = p1.Copy()
		o2 = p2.Copy()
	)
	// Create lookup maps to quickly see if a gene has been visited
	var (
		p1Visited, p2Visited = make([]bool, n), make([]bool, n)
		o1Visited, o2Visited = make([]bool, n), make([]bool, n)
	)
	for i := a; i < b; i++ {
		p1Visited[p1.ID(i)] = true
		p2Visited[p2.ID(i)] = true
		o1Visited[i] = true
		o2Visited[i] = true
	}
	for i := a; i < b; i++ {
		// Find the element in the second parent that has not been copied in the first offspring
		if !p1Visited[p2.ID(i)] {
			var j = i
			for o1Visited[j] {
				s := (o1).(taskGA)
				j = slices.IndexFunc(p2, func(t task) bool { return t.id == s.ID(j) })
				if j == -1 {
					log.Fatal("j == -1")
				}

			}
			o1.Set(j, p2.At(i))
			o1Visited[j] = true
		}
		// Find the element in the first parent that has not been copied in the second offspring
		if !p2Visited[p1.ID(i)] {
			var j = i
			for o2Visited[j] {
				s := (o2).(taskGA)
				j = slices.IndexFunc(p1, func(t task) bool { return t.id == s.ID(j) })
				if j == -1 {
					log.Fatal("j == -1")
				}
			}
			o2.Set(j, p1.At(i))
			o2Visited[j] = true
		}
	}
	// Fill in the offspring's missing values with the opposite parent's values
	for i := 0; i < n; i++ {
		if !o1Visited[i] {
			o1.Set(i, p2.At(i))
		}
		if !o2Visited[i] {
			o2.Set(i, p1.At(i))
		}
	}
	p1.Replace(o1)
	p2.Replace(o2)
}

func (s stateGA) Clone() eaopt.Genome {
	var clone stateGA

	clone.tasks = make([]task, len(s.tasks))
	copy(clone.tasks, s.tasks)

	clone.workers = make([]worker, len(s.workers))
	for i := range clone.workers {
		clone.workers[i] = s.workers[i].copy()
	}

	clone.setups = s.setups
	return clone
}
