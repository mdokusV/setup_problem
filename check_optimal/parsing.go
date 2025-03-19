package solution

import (
	"bufio"
	"first/models"
	"log"
	"os"
	"strconv"
	"strings"
)

type testFile struct {
	input    string
	solution string
}

func (tf testFile) Equal() {
	inputName := strings.Split(tf.input, ".")[0]
	outputName := strings.Split(tf.solution, ".")[0]
	if inputName != outputName {
		log.Fatal("input and output names do not match")
	}
}

func parseSolution(out string, initialState models.InitialState) (State, cMaxValue) {
	f, err := os.Open(FILE_DIR + "/" + out)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	state := prepareStartState(initialState)
	cMax := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "jobs") {
			// jobs in set 0:
			parts := strings.Fields(line)
			last_word := parts[len(parts)-1]
			number_str := last_word[:len(last_word)-1]
			number, err := strconv.Atoi(number_str)
			if err != nil {
				log.Panic(err)
			}
			worker := &state.workers[number]

			// 0 9
			scanner.Scan()
			line = scanner.Text()
			parts = strings.Fields(line)
			for _, p := range parts {
				job_id, err := strconv.Atoi(p)
				if err != nil {
					log.Panic(err)
				}
				task := state.tasks[job_id]
				worker.addTask(task)
			}

		} else if strings.Contains(line, "IP2position") {
			break
		} else if strings.Contains(line, "Cmax") {
			parts := strings.Fields(line)
			cMax_str := parts[1]
			cMax, err = strconv.Atoi(cMax_str)
			if err != nil {
				log.Panic(err)
			}
		}

	}
	return state, cMaxValue(cMax)
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
		if i%3 == 2 {
			continue
		}

		if i%3 == 0 {
			fileTriples[tripletSliceIndex].input = f.Name()
		} else if i%3 == 1 {
			fileTriples[tripletSliceIndex].solution = f.Name()
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
