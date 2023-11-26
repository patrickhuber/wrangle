package stores_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/dataptr"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	type test struct {
		name string
		str  string
		key  *stores.Value
	}
	tests := []test{
		{"name",
			"name", &stores.Value{
				Secret: &stores.Secret{
					Name:    "name",
					Version: stores.Version{Latest: true},
				},
				Path: &dataptr.DataPointer{},
			}},
		{"name version",
			"name@v1.0.0", &stores.Value{
				Secret: &stores.Secret{
					Name: "name",
					Version: stores.Version{
						Major:    1,
						Minor:    0,
						Revision: 0,
					},
				},
				Path: &dataptr.DataPointer{},
			}},
		{"name path", "name/test", &stores.Value{
			Secret: &stores.Secret{
				Name:    "name",
				Version: stores.Version{Latest: true},
			},
			Path: &dataptr.DataPointer{
				Segments: []dataptr.Segment{
					dataptr.Element{
						Name: "test",
					},
				},
			},
		}},
		{"name version path", "name@v1.0.0/test", &stores.Value{
			Secret: &stores.Secret{
				Name: "name",
				Version: stores.Version{
					Major:    1,
					Minor:    0,
					Revision: 0,
				},
			},
			Path: &dataptr.DataPointer{
				Segments: []dataptr.Segment{
					dataptr.Element{
						Name: "test",
					},
				},
			},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := stores.Parse(test.str)
			require.NoError(t, err)
			require.NotNil(t, p)
			require.Equal(t, test.key, p)
		})
	}
}
