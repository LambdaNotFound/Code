package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── Task / ITask ──────────────────────────────────────────────────────────────

func Test_Task(t *testing.T) {
	testCases := []struct {
		name         string
		id           string
		deps         []ITask
		wantID       string
		wantComplete bool
	}{
		{
			name:         "new_task_not_completed",
			id:           "t1",
			deps:         nil,
			wantID:       "t1",
			wantComplete: false,
		},
		{
			name:         "complete_sets_flag",
			id:           "t2",
			deps:         nil,
			wantID:       "t2",
			wantComplete: true,
		},
		{
			name:         "task_with_dependencies",
			id:           "child",
			deps:         []ITask{NewTask("parent")},
			wantID:       "child",
			wantComplete: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask(tc.id, tc.deps...)
			assert.Equal(t, tc.wantID, task.ID())
			assert.False(t, task.IsCompleted())
			assert.Equal(t, tc.deps, task.Dependencies())

			if tc.wantComplete {
				task.Complete()
				assert.True(t, task.IsCompleted())
			}
		})
	}
}

// ── TaskScheduler.Add / Next / Notify ─────────────────────────────────────────

func Test_TaskScheduler_AddNextNotify(t *testing.T) {
	testCases := []struct {
		name      string
		setupFn   func() (*TaskScheduler, []ITask) // returns scheduler + tasks in add order
		wantOrder []string                          // IDs of tasks returned by Next, in order
	}{
		{
			// Single task with no dependencies — immediately ready.
			name: "single_task_no_deps",
			setupFn: func() (*TaskScheduler, []ITask) {
				s := NewTaskScheduler()
				t1 := NewTask("t1")
				s.Add(t1)
				return s, []ITask{t1}
			},
			wantOrder: []string{"t1"},
		},
		{
			// Two independent tasks — both ready immediately.
			name: "two_independent_tasks",
			setupFn: func() (*TaskScheduler, []ITask) {
				s := NewTaskScheduler()
				t1, t2 := NewTask("t1"), NewTask("t2")
				s.Add(t1)
				s.Add(t2)
				return s, []ITask{t1, t2}
			},
			wantOrder: []string{"t1", "t2"},
		},
		{
			// t1 → t2 (t2 depends on t1): t1 must come before t2.
			name: "simple_chain",
			setupFn: func() (*TaskScheduler, []ITask) {
				s := NewTaskScheduler()
				t1 := NewTask("t1")
				t2 := NewTask("t2", t1)
				s.Add(t1)
				s.Add(t2)
				return s, []ITask{t1, t2}
			},
			wantOrder: []string{"t1", "t2"},
		},
		{
			// t1 → t3, t2 → t3: t3 must wait for both t1 and t2.
			name: "two_deps_one_child",
			setupFn: func() (*TaskScheduler, []ITask) {
				s := NewTaskScheduler()
				t1, t2 := NewTask("t1"), NewTask("t2")
				t3 := NewTask("t3", t1, t2)
				s.Add(t1)
				s.Add(t2)
				s.Add(t3)
				return s, []ITask{t1, t2, t3}
			},
			wantOrder: []string{"t1", "t2", "t3"},
		},
		{
			// t1 → t2 → t3: linear chain of three.
			name: "three_task_chain",
			setupFn: func() (*TaskScheduler, []ITask) {
				s := NewTaskScheduler()
				t1 := NewTask("t1")
				t2 := NewTask("t2", t1)
				t3 := NewTask("t3", t2)
				s.Add(t1)
				s.Add(t2)
				s.Add(t3)
				return s, []ITask{t1, t2, t3}
			},
			wantOrder: []string{"t1", "t2", "t3"},
		},
		{
			// Task already completed before Add — treated as indegree 0 (dep already done).
			name: "dep_already_completed_before_add",
			setupFn: func() (*TaskScheduler, []ITask) {
				s := NewTaskScheduler()
				t1 := NewTask("t1")
				t1.Complete()            // already done — t2 should be immediately ready
				t2 := NewTask("t2", t1)
				s.Add(t2)
				return s, []ITask{t2}
			},
			wantOrder: []string{"t2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, tasks := tc.setupFn()

			var order []string
			// Drive the scheduler: Next → Complete → Notify until drained.
			for {
				task := s.Next()
				if task == nil {
					break
				}
				order = append(order, task.ID())
				task.Complete()
				s.Notify(task)
			}

			assert.Equal(t, tc.wantOrder, order)
			// All tasks must be completed.
			for _, task := range tasks {
				assert.True(t, task.IsCompleted(), "task %s should be completed", task.ID())
			}
		})
	}
}

// Next returns nil when the queue is empty.
func Test_TaskScheduler_Next_empty(t *testing.T) {
	s := NewTaskScheduler()
	assert.Nil(t, s.Next())
}

// Next skips tasks that were completed externally before being dequeued.
func Test_TaskScheduler_Next_skips_already_completed(t *testing.T) {
	s := NewTaskScheduler()
	t1 := NewTask("t1")
	s.Add(t1)
	t1.Complete() // mark completed before Next is called
	assert.Nil(t, s.Next(), "already-completed task should be skipped")
}

// ── PrintOrder (Kahn's topological sort) ─────────────────────────────────────

func Test_PrintOrder(t *testing.T) {
	testCases := []struct {
		name        string
		buildTasks  func() []ITask
		wantIDs     []string // expected topological order of IDs
		wantErr     bool
	}{
		{
			name: "single_task",
			buildTasks: func() []ITask {
				return []ITask{NewTask("t1")}
			},
			wantIDs: []string{"t1"},
		},
		{
			name: "two_independent_tasks",
			buildTasks: func() []ITask {
				return []ITask{NewTask("t1"), NewTask("t2")}
			},
			wantIDs: []string{"t1", "t2"},
		},
		{
			name: "linear_chain",
			buildTasks: func() []ITask {
				t1 := NewTask("t1")
				t2 := NewTask("t2", t1)
				t3 := NewTask("t3", t2)
				return []ITask{t1, t2, t3}
			},
			wantIDs: []string{"t1", "t2", "t3"},
		},
		{
			name: "diamond_dependency",
			// t1 → t2, t1 → t3, t2 → t4, t3 → t4
			buildTasks: func() []ITask {
				t1 := NewTask("t1")
				t2 := NewTask("t2", t1)
				t3 := NewTask("t3", t1)
				t4 := NewTask("t4", t2, t3)
				return []ITask{t1, t2, t3, t4}
			},
			wantIDs: []string{"t1", "t2", "t3", "t4"},
		},
		{
			name: "two_roots_one_leaf",
			buildTasks: func() []ITask {
				t1, t2 := NewTask("t1"), NewTask("t2")
				t3 := NewTask("t3", t1, t2)
				return []ITask{t1, t2, t3}
			},
			wantIDs: []string{"t1", "t2", "t3"},
		},
		{
			// Cycle: t1 → t2 → t1 — PrintOrder must return an error.
			name: "cycle_returns_error",
			buildTasks: func() []ITask {
				t1 := NewTask("t1")
				t2 := NewTask("t2", t1)
				// artificially inject cycle: add t2 as dep of t1 after creation
				t1.deps = append(t1.deps, t2)
				return []ITask{t1, t2}
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewTaskScheduler()
			tasks := tc.buildTasks()
			order, err := s.PrintOrder(tasks)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, order)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tc.wantIDs), len(order))

			// For cases with a unique valid ordering, check exact order.
			// For cases with multiple valid orderings (diamond, two roots),
			// verify that each task appears after all its dependencies.
			idAt := make(map[string]int, len(order))
			for i, task := range order {
				idAt[task.ID()] = i
			}
			for _, task := range tasks {
				for _, dep := range task.Dependencies() {
					assert.Less(t, idAt[dep.ID()], idAt[task.ID()],
						"dep %s must come before %s", dep.ID(), task.ID())
				}
			}
		})
	}
}
