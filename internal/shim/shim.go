package shim

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Service interface {
	Execute(r *Request) error
}

type Request struct {
	Shell       string
	Executables []string
}

func NewService(
	fs fs.FS,
	path filepath.Provider,
	configuration config.Configuration,
	log log.Logger) Service {
	shells := []string{shellhook.Bash, shellhook.Powershell}
	return &service{
		shells:        shells,
		log:           log,
		configuration: configuration,
		fs:            fs,
		path:          path,
	}
}

type service struct {
	shells        []string
	fs            fs.FS
	path          filepath.Provider
	configuration config.Configuration
	log           log.Logger
}

const bashTemplate = `#!/usr/bin/bash
exec {{.Executable}} "$@"`

const powershellTemplate = `{{.Executable}} @args`

func (s *service) Execute(req *Request) error {
	if !slices.Contains(s.shells, req.Shell) {
		return fmt.Errorf("expected req.Shell to be one of %v", s.shells)
	}

	// load the configuration
	cfg, err := s.configuration.Get()
	if err != nil {
		return fmt.Errorf("failed to get configuration: %w", err)
	}

	content, err := s.getTemplate(req.Shell)
	if err != nil {
		return fmt.Errorf("failed to get template for shell %s: %w", req.Shell, err)
	}

	temp, err := template.New("").Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template for shell %s: %w", req.Shell, err)
	}

	var buf bytes.Buffer
	for _, executable := range req.Executables {
		data := map[string]any{
			"Executable": executable,
		}
		err = temp.Execute(&buf, data)
		if err != nil {
			return fmt.Errorf("failed to execute template for %s: %w", executable, err)
		}

		shimPath := s.getShimPath(cfg, req.Shell, executable)

		err = s.fs.WriteFile(shimPath, buf.Bytes(), 0775)
		if err != nil {
			return fmt.Errorf("failed to write shim file %s: %w", shimPath, err)
		}
	}
	return nil
}

func (s *service) getTemplate(shell string) (string, error) {
	switch shell {
	case shellhook.Bash:
		return bashTemplate, nil
	case shellhook.Powershell:
		return powershellTemplate, nil
	}
	return "", fmt.Errorf("unable to find shim template for shell %s", shell)
}

func (s *service) getShimPath(cfg config.Config, shell string, executable string) string {
	fileName := s.path.Base(executable)
	if shell == shellhook.Powershell {
		ext := s.path.Ext(fileName)
		fileNameWithoutExtension := strings.TrimSuffix(fileName, ext)
		fileName = fileNameWithoutExtension + ".ps1"
	}
	return s.path.Join(cfg.Spec.Environment[global.EnvBin], fileName)
}
