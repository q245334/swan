package executor

// Executor is responsible for creating execution environment for given workload.
// It returns Task handle when workload started gracefully.
// Workload is executed asynchronously.
type Executor interface {
	Execute(command string) (TaskHandle, error)
	GetIsolators() string
}
