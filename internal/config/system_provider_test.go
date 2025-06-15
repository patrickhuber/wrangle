package config_test

import (
	"testing"

	cfgpkg "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

func TestSystemProvider(t *testing.T) {

	// arrange
	target := cross.NewTest(platform.Linux, arch.AMD64)
	fileSystem := target.FS()
	fakeSystemConfigPath := "/opt/wrangle/config/config.yml"
	err := config.WriteFile(fileSystem, fakeSystemConfigPath, config.Config{})
	if err != nil {
		t.Fatal(err)
	}
	systemProvider := config.NewSystemProvider(fileSystem)

	// act
	cfg, err := systemProvider.Get(&cfgpkg.GetContext{
		MergedConfiguration: map[string]any{
			"spec": map[string]any{
				"env": map[string]any{
					global.EnvSystemConfig: fakeSystemConfigPath,
				},
			},
		},
	})

	// assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected a configuration, got nil")
	}
}
