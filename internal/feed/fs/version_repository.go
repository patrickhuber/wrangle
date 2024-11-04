package fs

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/packages"
	"gopkg.in/yaml.v3"
)

type versionRepository struct {
	fs               fs.FS
	workingDirectory string
	path             filepath.Provider
}

func NewVersionRepository(fs fs.FS, path filepath.Provider, workingDirectory string) feed.VersionRepository {
	return &versionRepository{
		fs:               fs,
		path:             path,
		workingDirectory: workingDirectory,
	}
}

func (r *versionRepository) Save(name string, version *packages.Version) error {
	versionPath := r.GetVersionFolderPath(name, version.Version)
	err := r.fs.MkdirAll(versionPath, 0775)
	if err != nil {
		return err
	}
	versionFile := r.GetVersionFilePath(name, version.Version)
	manifest := version.Manifest

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	return r.fs.WriteFile(versionFile, data, 0644)
}

func (r *versionRepository) Get(name string, version string) (*packages.Version, error) {
	versionFile := r.GetVersionFilePath(name, version)
	data, err := r.fs.ReadFile(versionFile)
	if err != nil {
		return nil, err
	}

	var obj any
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	manifest := &packages.Manifest{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      manifest,
		ErrorUnused: true,
		ErrorUnset:  true,
		DecodeHook:  decodeHook,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(obj)
	if err != nil {
		return nil, fmt.Errorf("unable to decode manifest for package '%s' version '%s' %w", name, version, err)
	}

	if manifest.Package == nil {
		return nil, fmt.Errorf("invalid package %s. manifest.Package is nil", versionFile)
	}
	v := &packages.Version{
		Version:  manifest.Package.Version,
		Manifest: manifest,
	}
	return v, nil
}

func decodeHook(fromType reflect.Type, toType reflect.Type, from any) (any, error) {
	switch toType {
	case reflect.TypeOf((*platform.Platform)(nil)).Elem():
		return platform.Parse(from.(string)), nil
	case reflect.TypeOf((*arch.Arch)(nil)).Elem():
		return arch.Parse(from.(string)), nil
	}
	return from, nil
}

func (r *versionRepository) List(name string) ([]*packages.Version, error) {
	packagePath := r.GetPackageFolderPath(name)
	files, err := r.fs.ReadDir(packagePath)
	if err != nil {
		return nil, err
	}
	versions := []*packages.Version{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		version, err := r.Get(name, file.Name())
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (r *versionRepository) Remove(name string, version string) error {
	versionPath := r.path.Join(r.workingDirectory, name, version)
	return r.fs.RemoveAll(versionPath)
}

func (r *versionRepository) GetPackageFolderPath(name string) string {
	return r.path.Join(r.workingDirectory, name)
}

func (r *versionRepository) GetVersionFilePath(name, version string) string {
	return r.path.Join(r.GetVersionFolderPath(name, version), "package.yml")
}

func (r *versionRepository) GetVersionFolderPath(name, version string) string {
	return r.path.Join(r.GetPackageFolderPath(name), version)
}
