package design

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VehicleSize_String(t *testing.T) {
	assert.Equal(t, "Small", Small.String())
	assert.Equal(t, "Medium", Medium.String())
	assert.Equal(t, "Large", Large.String())
	assert.Equal(t, "Unknown", VehicleSize(99).String())
}

func Test_Vehicle_GetSize(t *testing.T) {
	assert.Equal(t, Small, Motorcycle{ID: "m1"}.GetSize())
	assert.Equal(t, Medium, Car{ID: "c1"}.GetSize())
	assert.Equal(t, Large, Bus{ID: "b1"}.GetSize())
}

func Test_Vehicle_GetID(t *testing.T) {
	assert.Equal(t, "m1", Motorcycle{ID: "m1"}.GetID())
	assert.Equal(t, "c1", Car{ID: "c1"}.GetID())
	assert.Equal(t, "b1", Bus{ID: "b1"}.GetID())
}
