package config

import (
	"fmt"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/wrangle/internal/global"
)

func NewTestDefaultProvider(os os.OS, env env.Environment, path filepath.Provider) config.Provider {
	return &testDefaultProvider{
		os:   os,
		env:  env,
		path: path,
	}
}

type testDefaultProvider struct {
	path filepath.Provider
	os   os.OS
	env  env.Environment
}

// Get implements config.Provider.
func (t *testDefaultProvider) Get(ctx *config.GetContext) (any, error) {

	root := "/opt/wrangle"
	if platform.IsWindows(t.os.Platform()) {
		programData := t.env.Get("ProgramData")
		if len(programData) == 0 {
			return nil, fmt.Errorf("failed to get ProgramData environment variable")
		}
		root = t.path.Join(programData, "wrangle")
	}
	home, err := t.os.Home()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"spec": map[string]any{
			"env": map[string]any{
				global.EnvRoot:         root,
				global.EnvBin:          t.path.Join(root, "bin"),
				global.EnvPackages:     t.path.Join(root, "packages"),
				global.EnvSystemConfig: t.path.Join(root, "config", "config.yml"),
				global.EnvUserConfig:   t.path.Join(home, ".wrangle", "config.yml"),
			},
		},
	}, nil
}
