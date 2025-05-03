package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MedianFinder(t *testing.T) {
    medianFinder := MedianFinderConstructor()

    medianFinder.AddNum(1)
    medianFinder.AddNum(2)
    value := medianFinder.FindMedian()
    assert.Equal(t, 1.50000, value)
    medianFinder.AddNum(3)
    value = medianFinder.FindMedian()
    assert.Equal(t, 2.00000, value)
}
