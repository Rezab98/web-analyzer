package slicetools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name          string
		input         []int
		shouldInclude func(int) bool
		expected      []int
	}{
		{
			name:  "Even numbers",
			input: []int{1, 2, 3, 4, 5, 6},
			shouldInclude: func(n int) bool {
				return n%2 == 0
			},
			expected: []int{2, 4, 6},
		},
		{
			name:  "Odd numbers",
			input: []int{1, 2, 3, 4, 5, 6},
			shouldInclude: func(n int) bool {
				return n%2 != 0
			},
			expected: []int{1, 3, 5},
		},
		{
			name:  "Greater than 3",
			input: []int{1, 2, 3, 4, 5, 6},
			shouldInclude: func(n int) bool {
				return n > 3
			},
			expected: []int{4, 5, 6},
		},
		{
			name:  "Empty input",
			input: []int{},
			shouldInclude: func(n int) bool {
				return n%2 == 0
			},
			expected: []int{},
		},
		{
			name:  "No matches",
			input: []int{1, 3, 5},
			shouldInclude: func(n int) bool {
				return n%2 == 0
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.input, tt.shouldInclude)
			assert.Equal(t, tt.expected, result)
		})
	}
}
