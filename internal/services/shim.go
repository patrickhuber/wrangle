package services

import (
	"bytes"
	"fmt"
	"slices"
	"text/template"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Shim interface {
	Execute(r *ShimRequest) error
}

type ShimRequest struct {
	Shell       string
	Executables []string
}

type ShimResponse struct {
	Body []byte
	Path string
}

func NewShim(
	fs fs.FS,
	path filepath.Provider,
	configuration Configuration,
	log log.Logger) Shim {
	shells := []string{"bash", "powershell"}
	return &shim{
		shells:        shells,
		log:           log,
		configuration: configuration,
		fs:            fs,
		path:          path,
	}
}

type shim struct {
	shells        []string
	fs            fs.FS
	path          filepath.Provider
	configuration Configuration
	log           log.Logger
}

const bashTemplate = `#!/usr/bin/bash
exec {{.Executable}} "$@"`

const powershellTemplate = `{{.Executable}} @args`

func (s *shim) Execute(req *ShimRequest) error {
	if !slices.Contains(s.shells, req.Shell) {
		return fmt.Errorf("expected req.Shell to be one of %v", s.shells)
	}

	// load the configuration
	cfg, err := s.configuration.Get()
	if err != nil {
		return err
	}

	content, err := s.getTemplate(req.Shell)
	if err != nil {
		return err
	}

	temp, err := template.New("").Parse(content)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, executable := range req.Executables {
		data := map[string]any{
			"Executable": executable,
		}
		err = temp.Execute(&buf, data)
		if err != nil {
			return err
		}

		shimPath := s.getShimPath(cfg, executable)

		err = s.fs.WriteFile(shimPath, buf.Bytes(), 0775)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *shim) getTemplate(shell string) (string, error) {
	switch shell {
	case "bash":
		return bashTemplate, nil
	case "powershell":
		return powershellTemplate, nil
	}
	return "", fmt.Errorf("unable to find shim template for shell %s", shell)
}

func (s *shim) getShimPath(cfg config.Config, executable string) string {
	fileName := s.path.Base(executable)
	return s.path.Join(cfg.Spec.Environment[global.EnvBin], fileName)
}
