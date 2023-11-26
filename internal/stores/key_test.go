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
		key  *stores.Key
	}
	tests := []test{
		{"name",
			"name", &stores.Key{
				Data: &stores.Data{
					Name: "name",
					Version: stores.Version{
						Latest: true,
						Value:  "",
					},
				},
				Path: &dataptr.DataPointer{},
			}},
		{"name version",
			"name@v1.0.0", &stores.Key{
				Data: &stores.Data{
					Name: "name",
					Version: stores.Version{
						Value:  "v1.0.0",
						Latest: false,
					},
				},
				Path: &dataptr.DataPointer{},
			}},
		{"name path", "name/test", &stores.Key{
			Data: &stores.Data{
				Name: "name",
				Version: stores.Version{
					Latest: true,
					Value:  "",
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
		{"name version path", "name@v1.0.0/test", &stores.Key{
			Data: &stores.Data{
				Name: "name",
				Version: stores.Version{
					Value:  "v1.0.0",
					Latest: false,
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
			p, err := stores.ParseKey(test.str)
			require.NoError(t, err)
			require.NotNil(t, p)
			require.Equal(t, test.key, p)
		})
	}
}
