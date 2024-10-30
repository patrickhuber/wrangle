package services

import (
	"bytes"
	"fmt"
	"slices"
	"text/template"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Shim interface {
	Execute(r *ShimRequest) error
}

type ShimRequest struct {
	Shell       string
	Package     string
	Version     string
	Executables []string
}

type ShimResponse struct {
	Body []byte
	Path string
}

func NewShim(
	fs fs.FS,
	path *filepath.Processor,
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
	path          *filepath.Processor
	configuration Configuration
	log           log.Logger
}

const bashTemplate = `#!/usr/bin/bash
exec {{.Executable}} "$@"`

const powershellTemplate = `{{.Executable}} @args`

var executableExtensions = []string{
	".bat",
	".cmd",
	".exe",
}

func (s *shim) Execute(req *ShimRequest) error {
	if !slices.Contains(s.shells, req.Shell) {
		return fmt.Errorf("expected req.Shell to be one of %v", s.shells)
	}

	// load the configuration
	cfg, err := s.configuration.Get()
	if err != nil {
		return err
	}

	executables, err := s.getPackageVersionExecutables(cfg, req.Package, req.Version)
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
	for _, executable := range executables {
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

func (s *shim) getPackageVersionExecutables(cfg config.Config, packageName string, packageVersion string) ([]string, error) {

	// check package root
	packagesRoot := cfg.Spec.Environment[global.EnvPackages]
	packageRootPath := s.path.Join(packagesRoot, packageName)
	_, err := s.fs.Stat(packageRootPath)
	if err != nil {
		return nil, fmt.Errorf("unable to find installed package %s %w", packageName, err)
	}

	// check package version
	packageVersionPath := s.path.Join(packageRootPath, packageVersion)
	_, err = s.fs.Stat(packageVersionPath)
	if err != nil {
		return nil, fmt.Errorf("unable to find installed package %s %s %w", packageName, packageVersion, err)
	}

	// loop over all files looking for executable flag to be set or known executable extensions
	// TODO promote this to a 'executable' property on the package target
	files, err := s.fs.ReadDir(packageVersionPath)
	if err != nil {
		return nil, err
	}
	var executables []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := s.path.Join(packageVersionPath, file.Name())
		stat, err := s.fs.Stat(filePath)
		if err != nil {
			return nil, err
		}
		ext := s.path.Ext(file.Name())
		if stat.Mode()&0111 == 0 && !slices.Contains(executableExtensions, ext) {
			continue
		}
		executables = append(executables, s.path.Join(packageVersionPath, file.Name()))
	}
	return executables, nil
}

func (s *shim) getShimPath(cfg config.Config, executable string) string {
	fileName := s.path.Base(executable)
	return s.path.Join(cfg.Spec.Environment[global.EnvBin], fileName)
}
