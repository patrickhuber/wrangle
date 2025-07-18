package stores_test

import (
	"testing"

	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/stretchr/testify/require"
)

func TestKeyParse(t *testing.T) {
	type test struct {
		name string
		str  string
		key  stores.Key
	}
	tests := []test{
		{
			"name",
			"name", stores.Key{
				Data: stores.Data{
					Name: "name",
					Version: stores.Version{
						Latest: true, // latest is true if the value is empty
						Value:  "",
					},
				},
				Path: dataptr.DataPointer{},
			}},
		{
			"name underscore",
			"name_part",
			stores.Key{
				Data: stores.Data{
					Name: "name_part",
					Version: stores.Version{
						Latest: true, // latest is true if the value is empty
						Value:  "",
					},
				},
			}},
		{
			"name version",
			"name@v1.0.0", stores.Key{
				Data: stores.Data{
					Name: "name",
					Version: stores.Version{
						Value:  "v1.0.0",
						Latest: false,
					},
				},
				Path: dataptr.DataPointer{},
			}},
		{
			"name path",
			"name/test", stores.Key{
				Data: stores.Data{
					Name: "name",
					Version: stores.Version{
						Latest: true, // latest is true if the value is empty
						Value:  "",
					},
				},
				Path: dataptr.DataPointer{
					Segments: []dataptr.Segment{
						dataptr.Key{
							Key: "test",
						},
					},
				},
			}},
		{
			"name version path",
			"name@v1.0.0/test", stores.Key{
				Data: stores.Data{
					Name: "name",
					Version: stores.Version{
						Value:  "v1.0.0",
						Latest: false,
					},
				},
				Path: dataptr.DataPointer{
					Segments: []dataptr.Segment{
						dataptr.Key{
							Key: "test",
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
