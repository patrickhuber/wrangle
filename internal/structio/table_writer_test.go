package structio_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/wrangle/internal/structio"
	"github.com/stretchr/testify/require"
)

func TestTableWriter(t *testing.T) {
	t.Run("slice", func(t *testing.T) {
		type Person struct {
			FirstName string
			LastName  string
			Age       int
		}
		people := []Person{
			{"John", "Doe", 20},
			{"Some", "Guy", 30},
			{"Some", "Gal", 40},
		}
		buf := &bytes.Buffer{}
		writer := structio.NewTableWriter(buf)
		err := writer.Write(people)
		require.NoError(t, err)
		result := buf.String()
		t.Log(result)
	})
}
