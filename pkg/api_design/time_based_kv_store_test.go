package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TimeMap(t *testing.T) {
    timeMap := ConstructorTimeMap()

    timeMap.Set("foo", "bar", 1)
    val := timeMap.Get("foo", 1)
    assert.Equal(t, "bar", val)

    val = timeMap.Get("foo", 3)
    assert.Equal(t, "bar", val)

    timeMap.Set("foo", "bar2", 4)
    val = timeMap.Get("foo", 4)
    assert.Equal(t, "bar2", val)

    val = timeMap.Get("foo", 5)
    assert.Equal(t, "bar2", val)
}

func Test_TimeMap_2(t *testing.T) {
    timeMap := ConstructorTimeMap()

    timeMap.Set("love", "high", 10)
    timeMap.Set("love", "low", 20)
    val := timeMap.Get("love", 5)
    assert.Equal(t, "", val)

    val = timeMap.Get("love", 10)
    assert.Equal(t, "high", val)

    val = timeMap.Get("love", 15)
    assert.Equal(t, "high", val)

    val = timeMap.Get("love", 20)
    assert.Equal(t, "low", val)

    val = timeMap.Get("love", 25)
    assert.Equal(t, "low", val)
}
