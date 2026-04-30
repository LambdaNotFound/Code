package interview

import "fmt"

type ITask interface {
	ID() string
	Complete()
	IsCompleted() bool
	Dependencies() []ITask
}

type Task struct {
	id        string
	completed bool
	deps      []ITask
}

/*
 * // variadic — caller passes args directly
 * NewTask("A")                        // deps = []ITask{}
 * NewTask("B", taskA)                 // deps = []ITask{taskA}
 * NewTask("C", taskA, taskB, taskC)   // deps = []ITask{taskA, taskB, taskC}
 *
 * // slice — caller must construct the slice explicitly
 * NewTask("A", []ITask{})
 * NewTask("B", []ITask{taskA})
 * NewTask("C", []ITask{taskA, taskB, taskC})
 */
func NewTask(id string, deps ...ITask) *Task {
	return &Task{id: id, deps: deps}
}

func (t *Task) ID() string            { return t.id }
func (t *Task) Complete()             { t.completed = true }
func (t *Task) IsCompleted() bool     { return t.completed }
func (t *Task) Dependencies() []ITask { return t.deps }

type TaskScheduler struct {
	taskList []ITask

	// Version III
	indegree   map[ITask]int
	dependents map[ITask][]ITask
	ready      []ITask
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		indegree:   make(map[ITask]int),
		dependents: make(map[ITask][]ITask),
	}
}

func (s *TaskScheduler) Add(t ITask) ITask {
	// count how many of t's dependencies are not yet completed
	pending := 0
	for _, dependency := range t.Dependencies() {
		if !dependency.IsCompleted() {
			pending++
			// register t as a dependent of d so we can notify on completion
			s.dependents[dependency] = append(s.dependents[dependency], t)
		}
	}
	s.indegree[t] = pending

	if pending == 0 {
		s.ready = append(s.ready, t)
	}

	return t
}

// Next pops a ready task
func (s *TaskScheduler) Next() ITask {
	for len(s.ready) > 0 {
		t := s.ready[0]
		s.ready = s.ready[1:]
		if !t.IsCompleted() {
			return t
		}
	}
	return nil
}

// Notify is called by the caller after completing a task
// It decrements indegree of dependents and moves them to ready
func (s *TaskScheduler) Notify(completed ITask) {
	for _, dependent := range s.dependents[completed] {
		s.indegree[dependent]--
		if s.indegree[dependent] == 0 {
			s.ready = append(s.ready, dependent)
		}
	}
	delete(s.dependents, completed)
}

// PrintOrder returns all tasks in topological order
func (s *TaskScheduler) PrintOrder(tasks []ITask) ([]ITask, error) {
	// build indegree map
	indegree := map[ITask]int{}
	dependents := map[ITask][]ITask{}
	for _, t := range tasks {
		indegree[t] = 0
	}
	for _, t := range tasks {
		for _, d := range t.Dependencies() { // d -> t
			indegree[t]++
			dependents[d] = append(dependents[d], t)
		}
	}

	// seed queue with indegree-0 tasks
	ready := []ITask{}
	for _, t := range tasks {
		if indegree[t] == 0 {
			ready = append(ready, t)
		}
	}

	// Kahn's algorithm
	order := []ITask{}
	for len(ready) > 0 {
		t := ready[0]
		ready = ready[1:]
		order = append(order, t)

		for _, dep := range dependents[t] {
			indegree[dep]--
			if indegree[dep] == 0 {
				ready = append(ready, dep)
			}
		}
	}

	// cycle detection: if any task has indegree > 0 after processing, cycle exists
	if len(order) != len(tasks) {
		return nil, fmt.Errorf("cycle detected in dependency graph")
	}

	return order, nil
}
