package design

/**
 * Design Parking Lot
 */

type VehicleSize int

const (
    Small VehicleSize = iota
    Medium
    Large
)

// Optional: implement Stringer so it's human-readable
func (vs VehicleSize) String() string {
    switch vs {
    case Small:
        return "Small"
    case Medium:
        return "Medium"
    case Large:
        return "Large"
    default:
        return "Unknown"
    }
}

type Vehicle interface {
    GetSize() VehicleSize
    GetID() string
}

type Motorcycle struct {
    ID string
}

func (m Motorcycle) GetSize() VehicleSize { return Small }
func (m Motorcycle) GetID() string        { return m.ID }

type Car struct {
    ID string
}

func (c Car) GetSize() VehicleSize { return Medium }
func (c Car) GetID() string        { return c.ID }

type Bus struct {
    ID string
}

func (b Bus) GetSize() VehicleSize { return Large }
func (b Bus) GetID() string        { return b.ID }
