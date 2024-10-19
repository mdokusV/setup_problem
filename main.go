package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

const FILE_DIR = "./vs-m"

func main() {
	tests := prepareFiles()

	solutions := parallelRun(tests)
	sort.Slice(solutions, func(i, j int) bool {
		iID := solutions[i].id
		jID := solutions[j].id
		return slices.Compare(iID[:], jID[:]) < 0
	})
	for _, s := range solutions {
		fmt.Printf("id: %v, name: %v, result: %v, time: %v\n", s.id, s.name, s.result, s.time)
	}
}

type testFile struct {
	input  string
	output string
}
type testSolution struct {
	id     [4]int
	name   string
	result int
	time   int
}

func (tf testFile) Equal() {
	inputName := strings.Split(tf.input, ".")[0]
	outputName := strings.Split(tf.output, ".")[0]
	if inputName != outputName {
		log.Fatal("input and output names do not match")
	}
}

func prepareFiles() []testFile {
	files, err := os.ReadDir(FILE_DIR)
	if err != nil {
		log.Fatal(err)
	}

	if len(files)%3 != 0 {
		log.Fatal("number of files not divisible by 3")
	}

	fileTriples := make([]testFile, len(files)/3)

	for i, f := range files {
		tripletSliceIndex := i / 3
		if i%3 == 1 {
			continue
		}

		if i%3 == 0 {
			fileTriples[tripletSliceIndex].input = f.Name()
		} else if i%3 == 2 {
			fileTriples[tripletSliceIndex].output = f.Name()
		}
	}

	for _, tf := range fileTriples {
		tf.Equal()
	}
	return fileTriples
}

type initialState struct {
	workerNumber int
	tasks        []task
	setups       []setup
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

func scanNonEmpty(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			return scanner.Text()
		}
	}
	log.Panic("empty scanner")
	return ""
}

func parseInput(in string) initialState {
	f, err := os.Open(FILE_DIR + "/" + in)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	taskNumber, err := strconv.Atoi(strings.Fields(scanNonEmpty(scanner))[1])
	if err != nil {
		log.Panic(err)
	}

	setupNumber, err := strconv.Atoi(strings.Fields(scanNonEmpty(scanner))[1])
	if err != nil {
		log.Panic(err)
	}

	workerNumber, err := strconv.Atoi(strings.Fields(scanNonEmpty(scanner))[1])
	if err != nil {
		log.Panic(err)
	}

	tasks := make([]task, 0, taskNumber)
	for i, v := range strings.Fields(scanNonEmpty(scanner)) {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Panic(err)
		}
		tasks = append(tasks, task{i, val, nil})
	}

	setups := make([]setup, 0, setupNumber)
	for i, v := range strings.Fields(scanNonEmpty(scanner)) {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Panic(err)
		}
		setups = append(setups, setup{i, val, nil})
	}

	for taskID := range tasks {
		task := &tasks[taskID]
		for _, val := range strings.Fields(scanNonEmpty(scanner))[1:] {
			setupID, err := strconv.Atoi(val)
			setup := &setups[setupID]
			if err != nil {
				log.Panic(err)
			}
			task.setups = append(task.setups, setup)
			setup.tasks = append(setup.tasks, task)
		}
	}
	return initialState{workerNumber, tasks, setups}
}

func parallelRun(tests []testFile) []testSolution {
	const numWorkers = 10
	numJobs := len(tests)

	var wg sync.WaitGroup
	results := make([]testSolution, numJobs)

	pool, _ := ants.NewPoolWithFunc(numWorkers, func(i interface{}) {
		defer wg.Done()
		index, test := i.(struct {
			index int
			test  testFile
		}).index, i.(struct {
			index int
			test  testFile
		}).test
		state := parseInput(test.input)

		result, time := runTest(state)
		noSuffix, _ := strings.CutSuffix(test.input, ".in")
		dashSplit := strings.Split(noSuffix, "-")
		exampleNumber, _ := strconv.Atoi(dashSplit[len(dashSplit)-1])
		results[index] = testSolution{id: [4]int{len(state.tasks), len(state.setups), state.workerNumber, exampleNumber}, name: test.input, result: result, time: time}
	}, ants.WithPanicHandler(func(i interface{}) {
		log.Printf("%v\n%s", i, debug.Stack())
	}))
	defer pool.Release()

	for i, t := range tests {
		wg.Add(1)
		_ = pool.Invoke(struct {
			index int
			test  testFile
		}{index: i, test: t})
	}
	wg.Wait()
	return results
}

func runTest(state initialState) (int, int) {
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	value := state.workerNumber
	elapsed := time.Since(start)
	return value, int(elapsed.Milliseconds())
}
