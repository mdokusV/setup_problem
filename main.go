package main

import (
	"bufio"
	solution "first/Solution"
	"first/models"
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

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/panjf2000/ants/v2"
)

const FILE_DIR = "./vs-m"

func main() {
	tests := prepareFiles()

	solutions := parallelRun(tests, solution.GreedySolution)
	for _, s := range solutions {
		fmt.Printf("id: %-2v name: %-58v result: %-5v time: %-10v IPsolVal: %-5v GreedyVal: %-5v\n",
			s.id, s.name, s.cMax, s.time, s.IPsolVal, s.GreedyVal)
	}
}

type testFile struct {
	input  string
	output string
}
type testSolution struct {
	id        [4]int
	name      string
	cMax      models.CMaxValue
	time      time.Duration
	IPsolVal  cMaxValue
	GreedyVal cMaxValue
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

func scanNonEmpty(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			return scanner.Text()
		}
	}
	log.Panic("empty scanner")
	return ""
}

func parseInput(in string) models.InitialState {
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

	tasks := make([]models.InitialTask, 0, taskNumber)
	for i, v := range strings.Fields(scanNonEmpty(scanner)) {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Panic(err)
		}
		tasks = append(tasks, models.InitialTask{i, val, nil})
	}

	setups := make([]models.InitialSetup, 0, setupNumber)
	for i, v := range strings.Fields(scanNonEmpty(scanner)) {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Panic(err)
		}
		setups = append(setups, models.InitialSetup{i, val, nil})
	}

	for taskID := range tasks {
		task := &tasks[taskID]
		for _, val := range strings.Fields(scanNonEmpty(scanner))[1:] {
			setupID, err := strconv.Atoi(val)
			setup := &setups[setupID]
			if err != nil {
				log.Panic(err)
			}
			task.Setups = append(task.Setups, setup)
			setup.Tasks = append(setup.Tasks, task)
		}
	}
	return models.InitialState{workerNumber, tasks, setups}
}

func parseOutput(out string) (cMaxValue, cMaxValue) {
	f, err := os.Open(FILE_DIR + "/" + out)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var ip1anchorStr, greedyStr string
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "IP2position") {
			parts := strings.Fields(line)
			ip1anchorStr = parts[1]
		}
		if strings.Contains(line, "Greedy ") {
			parts := strings.Fields(line)
			greedyStr = parts[1]
		}
	}

	ip1anchor, _ := strconv.Atoi(ip1anchorStr)
	greedy, _ := strconv.Atoi(greedyStr)
	return cMaxValue(ip1anchor), cMaxValue(greedy)
}

func parallelRun(tests []testFile, runTest func(models.InitialState) (models.CMaxValue, time.Duration)) []testSolution {
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
		IPsolVal, GreedyVal := parseOutput(test.output)

		cMax, duration := runTest(state)
		noSuffix, _ := strings.CutSuffix(test.input, ".in")
		dashSplit := strings.Split(noSuffix, "-")
		exampleNumber, _ := strconv.Atoi(dashSplit[len(dashSplit)-1])
		results[index] = testSolution{id: [4]int{len(state.Tasks), len(state.Setups), state.WorkerNumber, exampleNumber}, name: test.input, cMax: cMax, time: duration, IPsolVal: IPsolVal, GreedyVal: GreedyVal}
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
	sort.Slice(results, func(i, j int) bool {
		iID := results[i].id
		jID := results[j].id
		return slices.Compare(iID[:], jID[:]) < 0
	})
	return results
}

func testTest(state initialState) (cMaxValue, time.Duration) {
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	value := state.workerNumber
	elapsed := time.Since(start)
	return cMaxValue(value), elapsed
}

type initialState struct {
	workerNumber int
	tasks        []initialTask
	setups       []initialSetup
}

type initialTask struct {
	id     int
	time   int
	setups []*initialSetup
}
type initialSetup struct {
	id    int
	time  int
	tasks []*initialTask
}

type worker struct {
	setups mapset.Set[int]
	tasks  []*initialTask
	cSum   int
}

func (w *worker) addTask(t *initialTask) {
	w.tasks = append(w.tasks, t)
	w.cSum += t.time
	for _, s := range t.setups {
		if !w.setups.ContainsOne(s.id) {
			w.setups.Add(s.id)
			w.cSum += s.time
		}
	}
}

type cMaxValue int

func greedySolution(state initialState) (cMaxValue, time.Duration) {
	start := time.Now()
	workers := make([]worker, state.workerNumber)
	for i := range workers {
		workers[i].setups = mapset.NewSet[int]()
	}
	for t := range state.tasks {
		task := &state.tasks[t]
		sort.Slice(workers, func(i, j int) bool {
			return workers[i].cSum < workers[j].cSum
		})
		minWorker := &workers[0]
		minWorker.addTask(task)
	}
	sort.Slice(workers, func(i, j int) bool {
		return workers[i].cSum > workers[j].cSum
	})

	value := workers[0].cSum

	elapsed := time.Since(start)
	return cMaxValue(value), elapsed
}

func popularitySolution(state initialState) (cMaxValue, time.Duration) {
	start := time.Now()
	workers := make([]worker, state.workerNumber)
	for i := range workers {
		workers[i].setups = mapset.NewSet[int]()
	}

	elapsed := time.Since(start)
	return cMaxValue(1), elapsed
}
