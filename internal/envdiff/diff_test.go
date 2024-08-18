package envdiff_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	type test struct {
		name     string
		prev     map[string]string
		next     map[string]string
		expected []envdiff.Change
	}
	tests := []test{
		{"equal",
			map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2"},
			map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2"},
			nil},
		{"add",
			map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2"},
			map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2",
				"KEY3": "VALUE3"},
			[]envdiff.Change{envdiff.Add{Key: "KEY3", Value: "VALUE3"}}},
		{"remove",
			map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2",
				"KEY3": "VALUE3"},
			map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2"},
			[]envdiff.Change{envdiff.Remove{Key: "KEY3", Previous: "VALUE3"}}},
		{"update",
			map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2",
				"KEY3": "VALUE3"},
			map[string]string{
				"KEY1": "VALUE1",
				"KEY2": "VALUE2",
				"KEY3": "VALUE"},
			[]envdiff.Change{envdiff.Update{Key: "KEY3", Previous: "VALUE3", Value: "VALUE"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			changes := envdiff.Diff(test.prev, test.next)
			require.Equal(t, test.expected, changes)
		})
	}
}

func TestCompress(t *testing.T) {
	changes := []envdiff.Change{
		envdiff.Add{
			Key:   "KEY1",
			Value: "VALUE1",
		},
		envdiff.Remove{
			Key:      "KEY2",
			Previous: "VALUE2",
		},
		envdiff.Update{
			Key:      "KEY3",
			Value:    "VALUE3",
			Previous: "VALUE",
		},
	}
	str, err := envdiff.Encode(changes)
	require.NoError(t, err)
	actual, err := envdiff.Decode(str)
	require.NoError(t, err)
	require.Equal(t, changes, actual)
}
