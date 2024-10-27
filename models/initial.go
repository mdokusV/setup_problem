package models

// initialState represents the initial state of the problem
type InitialState struct {
	WorkerNumber int
	Tasks        []InitialTask
	Setups       []InitialSetup
}

// InitialTask represents a task
type InitialTask struct {
	ID     int
	Time   int
	Setups []*InitialSetup
}

// InitialSetup represents a setup
type InitialSetup struct {
	ID    int
	Time  int
	Tasks []*InitialTask
}

type CMaxValue int
