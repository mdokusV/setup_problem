package models

type InitialState struct {
	WorkerNumber int
	Tasks        []InitialTask
	Setups       []InitialSetup
}

type InitialTask struct {
	ID     int
	Time   int
	Setups []*InitialSetup
}

type InitialSetup struct {
	ID    int
	Time  int
	Tasks []*InitialTask
}

type CMaxValue int
