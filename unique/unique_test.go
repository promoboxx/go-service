package unique

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInt32Slice(t *testing.T) {

	type testCase struct {
		name     string
		input    []int32
		expected []int32
		validate func(t *testing.T, result []int32, expected []int32)
	}

	tests := []testCase{
		{
			name:     "produces unique slice of int32",
			input:    []int32{1, 1, 2, 2, 3},
			expected: []int32{1, 2, 3},
			validate: func(t *testing.T, result []int32, expected []int32) {
				assert.Equal(t, result, expected)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Int32Slice(tc.input)
			tc.validate(t, result, tc.expected)
		})
	}
}

func TestStringSlice(t *testing.T) {

	type testCase struct {
		name     string
		input    []string
		expected []string
		validate func(t *testing.T, result []string, expected []string)
	}

	tests := []testCase{
		{
			name:     "produces unique slice of string",
			input:    []string{"foo", "foo", "bar"},
			expected: []string{"foo", "bar"},
			validate: func(t *testing.T, result []string, expected []string) {
				assert.Equal(t, result, expected)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := StringSlice(tc.input)
			tc.validate(t, result, tc.expected)
		})
	}
}
