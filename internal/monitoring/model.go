package monitoring

type Task struct {
	URL    string
	UserId uint32
}

type Result struct {
	Task       *Task
	StatusCode int
	Body       string
}
