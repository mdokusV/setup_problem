package main

import (
	solution "first/Solution"
	"fmt"
	"testing"
)

func Benchmark_parallelRun(t *testing.B) {
	files := prepareFiles()
	for i := 0; i < t.N; i++ {
		parallelRun(files, solution.GreedySolution)
	}
}

func TestGA(t *testing.T) {
	files := prepareFiles()
	state := parseInput(files[2].input)
	val, time, bestSolution := solution.RunSolution(state, solution.GreedyGAinitStateSolution)
	bestSolution.Print()
	fmt.Println(val, time)
}
