package config_test

import (
	"testing"

	cfgpkg "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

func TestFlagProvider(t *testing.T) {
	// arrange
	expected := map[string]string{
		global.EnvRoot:         "/opt/wrangle",
		global.EnvUserConfig:   "/home/fake/.wrangle/config.yml",
		global.EnvSystemConfig: "/opt/wrangle/config/config.yml",
		global.EnvPackages:     "/opt/wrangle/packages",
		global.EnvBin:          "/opt/wrangle/bin",
		global.EnvLogLevel:     "debug",
	}
	flagProvider := config.NewFlagProvider([]string{
		"--" + global.FlagBin, expected[global.EnvBin],
		"--" + global.FlagPackages, expected[global.EnvPackages],
		"--" + global.FlagRoot, expected[global.EnvRoot],
		"--" + global.FlagUserConfig, expected[global.EnvUserConfig],
		"--" + global.FlagSystemConfig, expected[global.EnvSystemConfig],
		"--" + global.FlagLogLevel, expected[global.EnvLogLevel],
	})

	// act
	cfg, err := flagProvider.Get(&cfgpkg.GetContext{
		MergedConfiguration: map[string]any{
			"spec": map[string]any{
				"flags": map[string]any{
					"testFlag": "testValue",
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
	for envVar, expectedValue := range expected {
		value, err := dataptr.GetAs[string]("/spec/env/"+envVar, cfg)
		if err != nil {
			t.Fatal(err)
		}
		if value != expectedValue {
			t.Fatalf("expected %s to be '%s', got '%s'", envVar, expectedValue, value)
		}
	}
}
