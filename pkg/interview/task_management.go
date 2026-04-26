package interview

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

type Task struct {
	id         string
	name       string
	priority   int
	counter    int    // original insertion order
	timestamp  int    // when added
	userId     string // assigned user, "" if unassigned
	startTime  int    // when assigned
	finishTime int    // when assignment ends
	completed  bool
}

type User struct {
	id    string
	quota int
	tasks map[string]*Task // taskId -> task (active assignments)
}

type TaskManager struct {
	tasks   map[string]*Task
	users   map[string]*User
	counter int
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
		users: make(map[string]*User),
	}
}

// ==================== LEVEL 1 ====================

func (tm *TaskManager) AddTask(timestamp int, name string, priority int) string {
	tm.counter++
	id := "task_id_" + strconv.Itoa(tm.counter)
	tm.tasks[id] = &Task{
		id:        id,
		name:      name,
		priority:  priority,
		counter:   tm.counter,
		timestamp: timestamp,
	}
	return id
}

func (tm *TaskManager) UpdateTask(timestamp int, taskId, name string, priority int) bool {
	t, ok := tm.tasks[taskId]
	if !ok {
		return false
	}
	t.name = name
	t.priority = priority
	return true
}

func (tm *TaskManager) GetTask(taskId string) string {
	t, ok := tm.tasks[taskId]
	if !ok {
		return ""
	}
	data := map[string]interface{}{
		"name":     t.name,
		"priority": t.priority,
	}
	b, _ := json.Marshal(data)
	return string(b)
}

// ==================== LEVEL 2 ====================

func (tm *TaskManager) SearchTask(timestamp int, nameFilter string, limit int) []string {
	var matched []*Task
	for _, t := range tm.tasks {
		if strings.Contains(t.name, nameFilter) {
			matched = append(matched, t)
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		if matched[i].priority != matched[j].priority {
			return matched[i].priority > matched[j].priority // desc priority
		}
		return matched[i].counter < matched[j].counter // asc task ID (counter = insertion order)
	})

	if limit > len(matched) {
		limit = len(matched)
	}
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = matched[i].name
	}
	return result
}

func (tm *TaskManager) TaskList(limit int) []string {
	all := make([]*Task, 0, len(tm.tasks))
	for _, t := range tm.tasks {
		all = append(all, t)
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].counter != all[j].counter {
			return all[i].counter < all[j].counter // asc task ID
		}
		return all[i].timestamp < all[j].timestamp // asc timestamp tie-break
	})

	if limit > len(all) {
		limit = len(all)
	}
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = all[i].name
	}
	return result
}

// ==================== LEVEL 3 ====================

func (tm *TaskManager) AddUser(timestamp int, userId string, quota int) bool {
	if _, exists := tm.users[userId]; exists {
		return false
	}
	tm.users[userId] = &User{
		id:    userId,
		quota: quota,
		tasks: make(map[string]*Task),
	}
	return true
}

func (tm *TaskManager) AssignTask(timestamp int, taskId, userId string, finishTime int) bool {
	t, ok := tm.tasks[taskId]
	if !ok {
		return false
	}
	u, ok := tm.users[userId]
	if !ok {
		return false
	}

	// Count active (unfinished) tasks — free quota for completed ones
	active := 0
	for _, assigned := range u.tasks {
		if !assigned.completed && timestamp < assigned.finishTime {
			active++
		}
	}
	if active >= u.quota {
		return false
	}

	t.userId = userId
	t.startTime = timestamp
	t.finishTime = finishTime
	u.tasks[taskId] = t
	return true
}

func (tm *TaskManager) GetTaskUser(timestamp int, userId string) []string {
	u, ok := tm.users[userId]
	if !ok {
		return []string{}
	}

	var result []string
	for _, t := range u.tasks {
		// startTime <= timestamp <= finishTime
		if t.startTime <= timestamp && timestamp <= t.finishTime {
			result = append(result, t.name)
		}
	}
	return result
}
