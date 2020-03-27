package grtp

/**
exe service interface
*/
type Executor interface {
	Execute() //oop
}

type Task struct {
	run func() //like run() in java
}

func NewTask(f func()) *Task {
	t := Task{run: f}
	return &t
}

func (t *Task) Execute() {
	t.run() //only call
}
