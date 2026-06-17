package design

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Direction_String(t *testing.T) {
	assert.Equal(t, "Up", Up.String())
	assert.Equal(t, "Down", Down.String())
	assert.Equal(t, "Idle", Idle.String())
	assert.Equal(t, "Idle", Direction(99).String())
}

func Test_Elevator_AddRequest(t *testing.T) {
	e := &Elevator{ID: 0, Current: 0, Direction: Idle}
	e.AddRequest(3)
	e.AddRequest(7)
	assert.Equal(t, []int{3, 7}, e.Requests)
}

func Test_Elevator_Move(t *testing.T) {
	t.Run("idle_when_no_requests", func(t *testing.T) {
		e := &Elevator{Current: 2, Direction: Up}
		e.Move()
		assert.Equal(t, Idle, e.Direction)
		assert.Equal(t, 2, e.Current)
	})

	t.Run("moves_up_toward_target", func(t *testing.T) {
		e := &Elevator{Current: 1, Requests: []int{3}}
		e.Move()
		assert.Equal(t, 2, e.Current)
		assert.Equal(t, Up, e.Direction)
		assert.Equal(t, []int{3}, e.Requests) // not arrived yet
	})

	t.Run("moves_down_toward_target", func(t *testing.T) {
		e := &Elevator{Current: 5, Requests: []int{3}}
		e.Move()
		assert.Equal(t, 4, e.Current)
		assert.Equal(t, Down, e.Direction)
	})

	t.Run("pops_request_on_arrival", func(t *testing.T) {
		e := &Elevator{Current: 2, Requests: []int{3, 6}}
		e.Move() // moves to 3, arrives
		assert.Equal(t, 3, e.Current)
		assert.Equal(t, []int{6}, e.Requests)
	})

	t.Run("already_at_target_pops_immediately", func(t *testing.T) {
		e := &Elevator{Current: 4, Requests: []int{4}}
		e.Move()
		assert.Equal(t, 4, e.Current)
		assert.Empty(t, e.Requests)
	})
}

func Test_NearestScheduler(t *testing.T) {
	sched := NearestScheduler{}

	t.Run("assigns_to_nearest_elevator", func(t *testing.T) {
		elevators := []*Elevator{
			{ID: 0, Current: 10},
			{ID: 1, Current: 3},
			{ID: 2, Current: 8},
		}
		sched.AssignRequest(Request{Floor: 5, Direction: Up}, elevators)
		// elevator 1 at floor 3 is nearest (dist 2)
		assert.Equal(t, []int{5}, elevators[1].Requests)
		assert.Empty(t, elevators[0].Requests)
		assert.Empty(t, elevators[2].Requests)
	})

	t.Run("no_elevators_no_panic", func(t *testing.T) {
		sched.AssignRequest(Request{Floor: 5, Direction: Up}, []*Elevator{})
	})
}

func Test_Building(t *testing.T) {
	t.Run("creates_correct_number_of_elevators", func(t *testing.T) {
		b := NewBuilding(3)
		assert.Len(t, b.Elevators, 3)
		for i, e := range b.Elevators {
			assert.Equal(t, i, e.ID)
			assert.Equal(t, 0, e.Current)
			assert.Equal(t, Idle, e.Direction)
		}
	})

	t.Run("request_elevator_assigns_to_nearest", func(t *testing.T) {
		b := NewBuilding(2)
		b.Elevators[0].Current = 1
		b.Elevators[1].Current = 8
		b.RequestElevator(7, Up)
		// elevator 1 at floor 8 is nearest floor 7 (dist 1 vs dist 6)
		assert.Equal(t, []int{7}, b.Elevators[1].Requests)
	})
}
