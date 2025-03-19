package solution

import (
	"cmp"
	"container/heap"
	"log"
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
)

type shortTask struct {
	id   int
	time int
}

type shortWorker struct {
	id                 int
	tasks              []shortTask
	latestLastTaskTime int
}

type PriorityQueue []*shortWorker

func (pq PriorityQueue) Len() int {
	return len(pq)
}
func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest
	comparison := cmp.Or(cmp.Compare(pq[j].latestLastTaskTime, pq[i].latestLastTaskTime), cmp.Compare(pq[j].id, pq[i].id))
	if comparison == 0 {
		log.Fatalln("comparison is equal")
	}
	return comparison < 0
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*shortWorker)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func removeInOrder(state State) State {
	// generate worker state that have combined times of setups to tasks
	state.Print()

	shortWorkers := make([]shortWorker, 0, len(state.workers))
	for _, worker := range state.workers {
		shortTasks := make([]shortTask, 0, len(worker.tasks))
		usedSetups := mapset.NewSetWithSize[*setup](len(worker.tasks))
		usedTime := 0
		for _, task := range worker.tasks {
			unusedSetups := task.setupsSet.Difference(usedSetups)
			unusedSetupsTimeSum := 0
			for _, uS := range unusedSetups.ToSlice() {
				unusedSetupsTimeSum += uS.time
			}
			newTask := shortTask{task.id, task.time + unusedSetupsTimeSum}
			shortTasks = append(shortTasks, newTask)
			usedTime += newTask.time
			usedSetups.Append(unusedSetups.ToSlice()...)
		}
		shortWorkers = append(shortWorkers, shortWorker{worker.id, shortTasks, usedTime - shortTasks[len(shortTasks)-1].time})
	}

	// remove tasks in order of latest start time of last available shortTask
	pq := make(PriorityQueue, 0, len(state.workers))
	for _, worker := range shortWorkers {
		pq = append(pq, &worker)
	}
	shortTasks := make([]shortTask, 0, len(state.tasks))

	heap.Init(&pq)
	for pq.Len() > 0 {
		worker := heap.Pop(&pq).(*shortWorker)

		shortTasks = append(shortTasks, worker.tasks[len(worker.tasks)-1])
		worker.tasks = worker.tasks[:len(worker.tasks)-1]
		if len(worker.tasks) > 0 {
			worker.latestLastTaskTime -= worker.tasks[len(worker.tasks)-1].time
		}

		if len(worker.tasks) != 0 {
			heap.Push(&pq, worker)
		}
	}

	slices.Reverse(shortTasks)

	newTaskOrder := make([]task, len(shortTasks))
	for i, t := range shortTasks {
		newTaskOrder[i] = state.tasks[t.id]
	}
	state.tasks = newTaskOrder

	for i := range state.workers {
		state.workers[i].Reset()
	}

	return state
}
