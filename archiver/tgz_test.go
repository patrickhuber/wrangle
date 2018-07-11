package archiver_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/archiver"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestRoundTripTgz(t *testing.T) {
	r := require.New(t)

	fileSystem := filesystem.NewMemoryMappedFsWrapper(afero.NewMemMapFs())

	err := afero.WriteFile(fileSystem, "/tmp/test", []byte("this is a test"), 0666)
	r.Nil(err)

	output, err := fileSystem.Create("/tmp/temp.tgz")
	r.Nil(err)
	defer output.Close()

	a := archiver.NewTargzArchiver(fileSystem)
	err = a.Write(output, []string{"/tmp/test"})
	r.Nil(err)

	source, err := fileSystem.Open("/tmp/temp.tgz")
	r.Nil(err)
	defer source.Close()
	err = a.Read(source, "/tmp")
	r.Nil(err)
}
