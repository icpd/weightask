package priortask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrioritySlice_Remove(t *testing.T) {
	t.Run("Remove value from slice", func(t *testing.T) {
		numbers := PrioritySlice{1, 5, 3, 6, 2, 4}
		numbers.Remove(3)
		assert.Equal(t, PrioritySlice{1, 5, 6, 2, 4}, numbers)
		numbers.Remove(1)
		assert.Equal(t, PrioritySlice{5, 6, 2, 4}, numbers)
		numbers.Remove(4)
		assert.Equal(t, PrioritySlice{5, 6, 2}, numbers)
	})

	t.Run("Remove value does not exist in slice", func(t *testing.T) {
		numbers := PrioritySlice{1, 5, 6, 2, 4}
		numbers.Remove(10)
		assert.Equal(t, PrioritySlice{1, 5, 6, 2, 4}, numbers)
	})
}

func TestPrioritySlice_Sort(t *testing.T) {
	numbers := PrioritySlice{1, 5, 3, 6, 2, 4}
	numbers.Sort()
	assert.Equal(t, PrioritySlice{6, 5, 4, 3, 2, 1}, numbers)
}
