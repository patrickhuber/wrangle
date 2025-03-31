package parse_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/dataptr/ast"
	"github.com/patrickhuber/wrangle/internal/dataptr/parse"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	type test struct {
		name string
		str  string
		ptr  ast.DataPointer
	}
	tests := []test{

		{"name", "name", ast.DataPointer{
			Segments: []ast.Segment{
				ast.Element{
					Name: "name",
				},
			},
		}},
		{"index", "0", ast.DataPointer{
			Segments: []ast.Segment{
				ast.Index{
					Index: 0,
				},
			},
		}},
		{"constraint", "key=value", ast.DataPointer{
			Segments: []ast.Segment{
				ast.Constraint{
					Key:   "key",
					Value: "value",
				},
			},
		}},
		{"multi name", "parent/child", ast.DataPointer{
			Segments: []ast.Segment{
				ast.Element{
					Name: "parent",
				},
				ast.Element{
					Name: "child",
				},
			},
		}},
		{"name constraint", "name/key=value", ast.DataPointer{
			Segments: []ast.Segment{
				ast.Element{
					Name: "name",
				},
				ast.Constraint{
					Key:   "key",
					Value: "value",
				},
			},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parse.Parse(test.str)
			require.NoError(t, err)
			require.Equal(t, test.ptr, actual)
		})
	}
}
