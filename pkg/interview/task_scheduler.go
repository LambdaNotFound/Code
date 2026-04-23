package interview

import "fmt"

type Task interface {
	ID() string
	Complete()
	IsCompleted() bool
	Dependencies() []Task
}

type SimpleTask struct {
	id        string
	completed bool
	deps      []Task
}

func NewTask(id string, deps ...Task) *SimpleTask {
	return &SimpleTask{id: id, deps: deps}
}

func (t *SimpleTask) ID() string           { return t.id }
func (t *SimpleTask) Complete()            { t.completed = true }
func (t *SimpleTask) IsCompleted() bool    { return t.completed }
func (t *SimpleTask) Dependencies() []Task { return t.deps }

type TaskScheduler struct {
	indegree   map[Task]int
	dependents map[Task][]Task
	ready      []Task
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		indegree:   map[Task]int{},
		dependents: map[Task][]Task{},
	}
}

func (s *TaskScheduler) Add(t Task) {
	// count how many of t's dependencies are not yet completed
	pending := 0
	for _, d := range t.Dependencies() {
		if !d.IsCompleted() {
			pending++
			// register t as a dependent of d so we can notify on completion
			s.dependents[d] = append(s.dependents[d], t)
		}
	}
	s.indegree[t] = pending

	if pending == 0 {
		s.ready = append(s.ready, t)
	}
}

// Next pops a ready task
func (s *TaskScheduler) Next() Task {
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
func (s *TaskScheduler) Notify(completed Task) {
	for _, dependent := range s.dependents[completed] {
		s.indegree[dependent]--
		if s.indegree[dependent] == 0 {
			s.ready = append(s.ready, dependent)
		}
	}
	delete(s.dependents, completed)
}

// PrintOrder returns all tasks in topological order
func (s *TaskScheduler) PrintOrder(tasks []Task) ([]Task, error) {
	// build indegree map
	indegree := map[Task]int{}
	dependents := map[Task][]Task{}
	for _, t := range tasks {
		indegree[t] = 0
	}
	for _, t := range tasks {
		for _, d := range t.Dependencies() {
			indegree[t]++
			dependents[d] = append(dependents[d], t)
		}
	}

	// seed queue with indegree-0 tasks
	ready := []Task{}
	for _, t := range tasks {
		if indegree[t] == 0 {
			ready = append(ready, t)
		}
	}

	// Kahn's algorithm
	order := []Task{}
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

func testTaskScheduler() {
	s := NewTaskScheduler()

	a := NewTask("A")
	b := NewTask("B")
	c := NewTask("C", a, b)
	d := NewTask("D", c)

	s.Add(a)
	s.Add(b)
	s.Add(c)
	s.Add(d)

	for {
		t := s.Next()
		if t == nil {
			break
		}
		fmt.Printf("Processing %s\n", t.ID())
		t.Complete()
		s.Notify(t) // tell scheduler this task finished
	}
}
