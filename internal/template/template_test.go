package template_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/template"
	"github.com/stretchr/testify/require"
)

func TestEvaluate(t *testing.T) {

	type test struct {
		name     string
		data     any
		m        map[string]any
		expected any
	}
	tests := []test{
		{
			name:     "string",
			data:     "((key))",
			m:        map[string]any{"key": "value"},
			expected: "value",
		},
		{
			name:     "slice",
			data:     []any{"((key))", "1", "2"},
			m:        map[string]any{"key": "0"},
			expected: []any{"0", "1", "2"},
		},
		{
			name:     "map",
			data:     map[any]any{"a": "A", "b": "B", "c": "((key))"},
			m:        map[string]any{"key": "C"},
			expected: map[any]any{"a": "A", "b": "B", "c": "C"},
		},
		{
			name:     "map key",
			data:     map[string]string{"a": "A", "((b))": "((B))", "c": "C"},
			m:        map[string]any{"b": "b", "B": "B"},
			expected: map[string]string{"a": "A", "b": "B", "c": "C"},
		},
		{
			name:     "int",
			data:     int(1),
			m:        map[string]any{},
			expected: int(1),
		},
		{
			name:     "int64",
			data:     int64(1),
			m:        map[string]any{},
			expected: int64(1),
		},
		{
			name:     "float64",
			data:     float64(1),
			m:        map[string]any{},
			expected: float64(1),
		},
		{
			name:     "multiple inline",
			data:     "((one)),((two)),((three))",
			m:        map[string]any{"one": "1", "two": "2", "three": "3"},
			expected: "1,2,3",
		},
		{
			name:     "itoa",
			data:     "((one))",
			m:        map[string]any{"one": 1},
			expected: "1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mp = template.MapProvider{}
			for k, v := range test.m {
				mp[k] = v
			}
			tmp := template.New(test.data, template.WithProvider(mp))
			result, err := tmp.Evaluate()
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Empty(t, result.Unresolved)
			require.Equal(t, test.expected, result.Value)
		})
	}
}

func TestEvaluateFail(t *testing.T) {

	type test struct {
		name     string
		data     any
		varNames []string
	}
	tests := []test{
		{
			name:     "string",
			data:     "((key))",
			varNames: []string{"key"},
		},
		{
			name:     "slice",
			data:     []any{"((key))", "1", "2"},
			varNames: []string{"key"},
		},
		{
			name:     "map",
			data:     map[any]any{"a": "A", "b": "B", "c": "((key))"},
			varNames: []string{"key"},
		},
		{
			name:     "multi",
			data:     "((a))((b))((c))",
			varNames: []string{"a", "b", "c"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mp = template.MapProvider{}
			tmp := template.New(test.data, template.WithProvider(mp))
			result, err := tmp.Evaluate()
			require.NoError(t, err)
			require.Equal(t, len(test.varNames), len(result.Unresolved))
		})
	}
}
