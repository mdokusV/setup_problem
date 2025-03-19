package solution

import (
	"first/models"
	"log"
)

const FILE_DIR = "./m5-vs-N"

type cMaxValue = models.CMaxValue

func Check() {
	tests := prepareFiles()
	for i, test := range tests {
		initialState := parseInput(test.input)
		state, value := parseSolution(test.solution, initialState)
		start_state := removeInOrder(state)

		cMax, _ := GreedySolution(start_state)
		if cMax != value {
			log.Println("not equal on ", i)
		}

	}
}
