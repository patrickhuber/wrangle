package commands

import (
	"os"
	"testing"

	"github.com/patrickhuber/wrangle/global"
)

func TestEnv(t *testing.T) {
	os.Unsetenv(global.PackagePathKey)
	os.Unsetenv(global.ConfigFileKey)

}
