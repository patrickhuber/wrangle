package edit

import (
	"fmt"
	"strings"

	goconfig "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Service interface {
	AddPackage(r *PackageRequest) error
	AddStore(r *StoreRequest) error
}

type service struct {
	fs            fs.FS
	opsys         os.OS
	path          filepath.Provider
	log           log.Logger
	env           env.Environment
	cli           config.CliContext
	systemDefault config.SystemDefaultProvider
}

func NewService(
	fs fs.FS,
	opsys os.OS,
	path filepath.Provider,
	log log.Logger,
	env env.Environment,
	cli config.CliContext,
	systemDefault config.SystemDefaultProvider,
) Service {
	return &service{
		fs:            fs,
		opsys:         opsys,
		path:          path,
		log:           log,
		env:           env,
		cli:           cli,
		systemDefault: systemDefault,
	}
}

const (
	ScopeNearest = "nearest"
	ScopeLocal   = "local"
	ScopeUser    = "user"
	ScopeSystem  = "system"
)

func ValidateScope(scope string) error {
	switch scope {
	case ScopeNearest, ScopeLocal, ScopeUser, ScopeSystem:
		return nil
	default:
		return fmt.Errorf("invalid scope %q (expected nearest, local, user, or system)", scope)
	}
}

func (s *service) update(scope string, update func(*config.Config) error) (string, error) {
	file, err := s.configFile(scope)
	if err != nil {
		return "", err
	}
	exists, err := s.fs.Exists(file)
	if err != nil {
		return "", fmt.Errorf("unable to check config file: %w", err)
	}

	var cfg config.Config
	if exists {
		cfg, err = config.ReadFile(s.fs, file)
		if err != nil {
			return "", fmt.Errorf("unable to read config file '%s': %w", file, err)
		}
	} else {
		if err := s.fs.MkdirAll(s.path.Dir(file), 0755); err != nil {
			return "", fmt.Errorf("unable to create config directory: %w", err)
		}
		s.log.Infof("creating config file '%s'", file)
		cfg = config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Feeds:       []config.Feed{},
				Stores:      []config.Store{},
				Environment: map[string]string{},
				Packages:    []config.Package{},
			},
		}
	}

	if err := update(&cfg); err != nil {
		return "", err
	}
	if err := config.WriteFile(s.fs, file, cfg); err != nil {
		return "", fmt.Errorf("unable to write config file '%s': %w", file, err)
	}
	return file, nil
}

func (s *service) configFile(scope string) (string, error) {
	if scope == "" {
		scope = ScopeNearest
	}
	if err := ValidateScope(scope); err != nil {
		return "", err
	}

	if scope == ScopeNearest || scope == ScopeLocal {
		dir, err := s.opsys.WorkingDirectory()
		if err != nil {
			return "", fmt.Errorf("unable to get working directory: %w", err)
		}
		workDir := dir
		for {
			parent := s.path.Dir(dir)
			if parent == dir {
				break
			}
			entries, err := s.fs.ReadDir(dir)
			if err != nil {
				return "", fmt.Errorf("unable to read config directory %q: %w", dir, err)
			}
			var matches []string
			for _, entry := range entries {
				name := entry.Name()
				ext := s.path.Ext(name)
				if !entry.IsDir() && strings.HasSuffix(strings.TrimSuffix(name, ext), ".wrangle") &&
					(ext == ".toml" || ext == ".yml" || ext == ".yaml" || ext == ".json" || ext == ".env") {
					matches = append(matches, s.path.Join(dir, name))
				}
			}
			if len(matches) > 1 {
				return "", fmt.Errorf("multiple local config files in %q; cannot choose one", dir)
			}
			if len(matches) == 1 {
				if s.path.Ext(matches[0]) == ".env" {
					return "", fmt.Errorf("cannot edit local config file %q: .env format is not supported", matches[0])
				}
				return matches[0], nil
			}
			dir = parent
		}
		if scope == ScopeLocal {
			return s.path.Join(workDir, global.LocalConfigurationFileName), nil
		}
	}

	providers := []goconfig.Provider{
		config.NewDefaultProvider(s.opsys, s.env, s.path),
		config.NewEnvProvider(s.env, global.EnvPrefix),
		config.NewCliProvider(s.cli),
	}
	base, err := goconfig.NewRoot(providers...).Get(&goconfig.GetContext{})
	if err != nil {
		return "", err
	}
	systemFile, err := dataptr.GetAs[string]("/spec/env/"+global.EnvSystemConfig, base)
	if err != nil {
		return "", err
	}
	if scope == ScopeSystem {
		return systemFile, nil
	}
	systemExists, err := s.fs.Exists(systemFile)
	if err != nil {
		return "", err
	}
	if systemExists {
		providers = append(providers, config.NewSystemProvider(s.fs, s.path, s.systemDefault, true))
		base, err = goconfig.NewRoot(providers...).Get(&goconfig.GetContext{})
		if err != nil {
			return "", err
		}
	}
	userFile, err := dataptr.GetAs[string]("/spec/env/"+global.EnvUserConfig, base)
	if err != nil {
		return "", err
	}
	if scope == ScopeUser {
		return userFile, nil
	}
	userExists, err := s.fs.Exists(userFile)
	if err != nil {
		return "", err
	}
	if userExists {
		return userFile, nil
	}
	if systemExists {
		return systemFile, nil
	}
	return "", fmt.Errorf("no existing local, user, or system config file found")
}
