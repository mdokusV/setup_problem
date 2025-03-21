package main

import (
	"bufio"
	solution "first/Solution"
	"first/models"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type cMaxValue = models.CMaxValue

// const FILE_DIR = "./vs-m"

const FILE_DIR = "./m5-vs-N"

// var NUM_WORKERS = 1

var NUM_WORKERS = runtime.NumCPU()

func main() {
	tests := prepareFiles()

	solutions := parallelRun(tests, solution.GreedyGAinitStateSolution)
	showSummary(solutions)
	saveSummary(solutions)
	showStatistics(solutions)
}

func showStatistics(results []testSolution) {
	sum, min, max := 0.0, math.Inf(1), 0.0
	for _, r := range results {
		resToGreedy := float64(r.cMax) / float64(r.GreedyVal)
		sum += resToGreedy
		if resToGreedy > max {
			max = resToGreedy
		}
		if resToGreedy < min {
			min = resToGreedy
		}
	}

	fmt.Printf("avg: %-5v min: %-5v max: %-5v \n", sum/float64(len(results)), min, max)
}

type testFile struct {
	input  string
	output string
}
type testSolution struct {
	id        [4]int
	name      string
	cMax      cMaxValue
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
		tasks = append(tasks, models.InitialTask{ID: i, Time: val, Setups: nil})
	}

	setups := make([]models.InitialSetup, 0, setupNumber)
	for i, v := range strings.Fields(scanNonEmpty(scanner)) {
		val, err := strconv.Atoi(v)
		if err != nil {
			log.Panic(err)
		}
		setups = append(setups, models.InitialSetup{ID: i, Time: val, Tasks: nil})
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
	return models.InitialState{WorkerNumber: workerNumber, Tasks: tasks, Setups: setups}
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

func parallelRun(tests []testFile, solutionFunc func(solution.State) (models.CMaxValue, solution.State)) []testSolution {
	numJobs := len(tests)

	var wg sync.WaitGroup
	results := make([]testSolution, numJobs)

	runtime.GC()

	pool, err := ants.NewPoolWithFunc(NUM_WORKERS, func(i interface{}) {
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

		noSuffix, _ := strings.CutSuffix(test.input, ".in")
		dashSplit := strings.Split(noSuffix, "-")
		exampleNumber, _ := strconv.Atoi(dashSplit[len(dashSplit)-1])
		id := [4]int{len(state.Tasks), len(state.Setups), state.WorkerNumber, exampleNumber}

		cMax, duration, _ := solution.RunSolution(state, solutionFunc)
		results[index] = testSolution{id: id, name: test.input, cMax: cMax, time: duration, IPsolVal: IPsolVal, GreedyVal: GreedyVal}
		// fmt.Println(id)
		// bestState.Print()
		// fmt.Println("____________________________")
	}, ants.WithPanicHandler(func(i interface{}) {
		log.Printf("caught panic: %v", i)
	}), ants.WithExpiryDuration(time.Minute), ants.WithPreAlloc(true))
	defer pool.Release()

	if err != nil {
		log.Fatal(err)
	}
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

func saveSummary(results []testSolution) {
	f, err := os.Create("summary.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, r := range results {
		resToIPos := float64(r.cMax) / float64(r.IPsolVal)
		resToGreedy := float64(r.cMax) / float64(r.GreedyVal)
		fmt.Fprintf(f, "id: %-2v  result: %-5v time: %-13v | IPsolVal: %-5v resToIPos: %-6.3f | GreedyVal: %-5v resToGreedy: %-6.3f \n",
			r.id, r.cMax, r.time, r.IPsolVal, resToIPos, r.GreedyVal, resToGreedy)
	}
}

func showSummary(results []testSolution) {
	for _, r := range results {
		resToIPos := float64(r.cMax) / float64(r.IPsolVal)
		resToGreedy := float64(r.cMax) / float64(r.GreedyVal)
		fmt.Printf("id: %-2v  result: %-5v time: %-13v | IPsolVal: %-5v resToIPos: %-6.3f | GreedyVal: %-5v resToGreedy: %-6.3f \n",
			r.id, r.cMax, r.time, r.IPsolVal, resToIPos, r.GreedyVal, resToGreedy)
	}
}
