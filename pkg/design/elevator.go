package design

/**
 * Design Elevator
 *    OO design:
 *
 *    Scheduler decides which elevator handles which request.
 *    Elevators process requests and move between floors.
 */

type Direction int

const (
    Up Direction = iota
    Down
    Idle
)

func (d Direction) String() string {
    switch d {
    case Up:
        return "Up"
    case Down:
        return "Down"
    default:
        return "Idle"
    }
}

type Request struct {
    Floor     int
    Direction Direction
}

type Button struct {
    Floor     int
    Direction Direction
}

type Elevator struct {
    ID        int
    Current   int
    Direction Direction
    Requests  []int // queue of floors to visit
}

func (e *Elevator) AddRequest(floor int) {
    e.Requests = append(e.Requests, floor)
}

func (e *Elevator) Move() {
    if len(e.Requests) == 0 {
        e.Direction = Idle
        return
    }

    target := e.Requests[0]
    if e.Current < target {
        e.Current++
        e.Direction = Up
    } else if e.Current > target {
        e.Current--
        e.Direction = Down
    }

    // Arrived at target
    if e.Current == target {
        e.Requests = e.Requests[1:]
    }
}

/**
 * Scheduler
 */
type Scheduler interface {
    AssignRequest(req Request, elevators []*Elevator)
}

type NearestScheduler struct{}

func (s NearestScheduler) AssignRequest(req Request, elevators []*Elevator) {
    var best *Elevator
    minDist := 1 << 30

    for _, e := range elevators {
        dist := abs(e.Current - req.Floor)
        if dist < minDist {
            minDist = dist
            best = e
        }
    }

    if best != nil {
        best.AddRequest(req.Floor)
    }
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

type Building struct {
    Elevators []*Elevator
    Scheduler Scheduler
}

func NewBuilding(numElevators int) *Building {
    elevators := make([]*Elevator, numElevators)
    for i := 0; i < numElevators; i++ {
        elevators[i] = &Elevator{ID: i, Current: 0, Direction: Idle}
    }
    return &Building{
        Elevators: elevators,
        Scheduler: NearestScheduler{},
    }
}

func (b *Building) RequestElevator(floor int, direction Direction) {
    req := Request{Floor: floor, Direction: direction}
    b.Scheduler.AssignRequest(req, b.Elevators)
}
