package config_test

import (
	"testing"

	cfgpkg "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

func TestEnvProvider(t *testing.T) {
	// arrange
	target := cross.NewTest(platform.Linux, arch.AMD64)
	e := target.Env()
	e.Set(global.EnvRoot, "/opt/wrangle")
	e.Set(global.EnvUserConfig, "/home/fake/.wrangle/config.yml")
	e.Set(global.EnvSystemConfig, "/opt/wrangle/config/config.yml")
	e.Set(global.EnvBin, "/opt/wrangle/bin")
	e.Set(global.EnvPackages, "/opt/wrangle/packages")
	e.Set("TEST", "test")
	envProvider := config.NewEnvProvider(target.Env())

	// act
	cfg, err := envProvider.Get(&cfgpkg.GetContext{})

	// assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected a configuration, got nil")
	}
	envConfig, err := dataptr.GetAs[map[string]any]("/spec/env", cfg)
	if err != nil {
		t.Fatal(err)
	}
	if envConfig == nil {
		t.Fatal("expected env configuration, got nil")
	}
	if len(envConfig) != 5 {
		t.Fatalf("expected 5 environment variables, got %d", len(envConfig))
	}
}
