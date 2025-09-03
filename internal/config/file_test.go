package config_test

import (
	"path/filepath"
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/internal/config"
)

func TestFileConfig(t *testing.T) {
	t.Run("can write and read yaml file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		fileConfig := config.NewFileConfig(fs)

		cfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Environment: map[string]string{
					"TEST": "value",
				},
				Feeds: []config.Feed{
					{
						Name: "test",
						Type: "fs",
						URI:  "/tmp/test",
					},
				},
			},
		}

		filePath := "/tmp/config.yml"
		err := fileConfig.Write(filePath, cfg)
		require.NoError(t, err)

		readCfg, err := fileConfig.Read(filePath)
		require.NoError(t, err)
		require.Equal(t, cfg.ApiVersion, readCfg.ApiVersion)
		require.Equal(t, cfg.Kind, readCfg.Kind)
		require.Equal(t, cfg.Spec.Environment, readCfg.Spec.Environment)
		require.Equal(t, len(cfg.Spec.Feeds), len(readCfg.Spec.Feeds))
		require.Equal(t, cfg.Spec.Feeds[0].Name, readCfg.Spec.Feeds[0].Name)
	})

	t.Run("can write and read json file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		fileConfig := config.NewFileConfig(fs)

		cfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Environment: map[string]string{
					"TEST": "value",
				},
				Stores: []config.Store{
					{
						Name: "test-store",
						Type: "memory",
						Properties: map[string]string{
							"key": "value",
						},
					},
				},
			},
		}

		filePath := "/tmp/config.json"
		err := fileConfig.Write(filePath, cfg)
		require.NoError(t, err)

		readCfg, err := fileConfig.Read(filePath)
		require.NoError(t, err)
		require.Equal(t, cfg.ApiVersion, readCfg.ApiVersion)
		require.Equal(t, cfg.Kind, readCfg.Kind)
		require.Equal(t, cfg.Spec.Environment, readCfg.Spec.Environment)
		require.Equal(t, len(cfg.Spec.Stores), len(readCfg.Spec.Stores))
		require.Equal(t, cfg.Spec.Stores[0].Name, readCfg.Spec.Stores[0].Name)
	})

	t.Run("write if not exists creates file when it doesn't exist", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		fileConfig := config.NewFileConfig(fs)

		filePath := "/tmp/new_config.yml"
		defaultFactory := func() config.Config {
			return config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
				Spec: config.Spec{
					Environment: map[string]string{
						"DEFAULT": "value",
					},
				},
			}
		}

		err := fileConfig.WriteIfNotExists(filePath, defaultFactory)
		require.NoError(t, err)

		exists, err := fs.Exists(filePath)
		require.NoError(t, err)
		require.True(t, exists)

		readCfg, err := fileConfig.Read(filePath)
		require.NoError(t, err)
		require.Equal(t, config.ApiVersion, readCfg.ApiVersion)
		require.Equal(t, "value", readCfg.Spec.Environment["DEFAULT"])
	})

	t.Run("write if not exists doesn't overwrite existing file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		fileConfig := config.NewFileConfig(fs)

		filePath := "/tmp/existing_config.yml"
		originalCfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Environment: map[string]string{
					"ORIGINAL": "value",
				},
			},
		}

		// Write original file
		err := fileConfig.Write(filePath, originalCfg)
		require.NoError(t, err)

		// Try to write with default factory
		defaultFactory := func() config.Config {
			return config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
				Spec: config.Spec{
					Environment: map[string]string{
						"DEFAULT": "value",
					},
				},
			}
		}

		err = fileConfig.WriteIfNotExists(filePath, defaultFactory)
		require.NoError(t, err)

		// Verify original content is preserved
		readCfg, err := fileConfig.Read(filePath)
		require.NoError(t, err)
		require.Equal(t, "value", readCfg.Spec.Environment["ORIGINAL"])
		require.NotContains(t, readCfg.Spec.Environment, "DEFAULT")
	})

	t.Run("read or create creates file when it doesn't exist", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		fileConfig := config.NewFileConfig(fs)

		filePath := "/tmp/read_or_create.yml"
		defaultFactory := func() config.Config {
			return config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
				Spec: config.Spec{
					Environment: map[string]string{
						"CREATED": "value",
					},
				},
			}
		}

		readCfg, err := fileConfig.ReadOrCreate(filePath, defaultFactory)
		require.NoError(t, err)
		require.Equal(t, config.ApiVersion, readCfg.ApiVersion)
		require.Equal(t, "value", readCfg.Spec.Environment["CREATED"])

		// Verify file was created
		exists, err := fs.Exists(filePath)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("read or create reads existing file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		fileConfig := config.NewFileConfig(fs)

		filePath := "/tmp/existing_read_or_create.yml"
		originalCfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Environment: map[string]string{
					"EXISTING": "value",
				},
			},
		}

		// Write original file
		err := fileConfig.Write(filePath, originalCfg)
		require.NoError(t, err)

		defaultFactory := func() config.Config {
			return config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
				Spec: config.Spec{
					Environment: map[string]string{
						"DEFAULT": "value",
					},
				},
			}
		}

		readCfg, err := fileConfig.ReadOrCreate(filePath, defaultFactory)
		require.NoError(t, err)
		require.Equal(t, "value", readCfg.Spec.Environment["EXISTING"])
		require.NotContains(t, readCfg.Spec.Environment, "DEFAULT")
	})

	t.Run("returns error for unsupported file extension", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		fileConfig := config.NewFileConfig(fs)

		cfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
		}

		filePath := "/tmp/config.txt"
		err := fileConfig.Write(filePath, cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unable to determine encoding")
	})
}

func TestReadFile(t *testing.T) {
	t.Run("can read yaml file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()

		// Write a test file first
		cfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Environment: map[string]string{
					"TEST": "value",
				},
			},
		}

		filePath := "/tmp/test_read.yml"
		err := config.WriteFile(fs, filePath, cfg)
		require.NoError(t, err)

		readCfg, err := config.ReadFile(fs, filePath)
		require.NoError(t, err)
		require.Equal(t, cfg.ApiVersion, readCfg.ApiVersion)
		require.Equal(t, cfg.Kind, readCfg.Kind)
		require.Equal(t, cfg.Spec.Environment, readCfg.Spec.Environment)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()

		_, err := config.ReadFile(fs, "/tmp/non_existent.yml")
		require.Error(t, err)
	})
}

func TestWriteFile(t *testing.T) {
	t.Run("can write yaml file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()

		cfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Packages: []config.Package{
					{
						Name:    "test-package",
						Version: "1.0.0",
					},
				},
			},
		}

		filePath := "/tmp/test_write.yml"
		err := config.WriteFile(fs, filePath, cfg)
		require.NoError(t, err)

		exists, err := fs.Exists(filePath)
		require.NoError(t, err)
		require.True(t, exists)

		// Verify content
		readCfg, err := config.ReadFile(fs, filePath)
		require.NoError(t, err)
		require.Equal(t, cfg.ApiVersion, readCfg.ApiVersion)
		require.Equal(t, len(cfg.Spec.Packages), len(readCfg.Spec.Packages))
		require.Equal(t, cfg.Spec.Packages[0].Name, readCfg.Spec.Packages[0].Name)
	})

	t.Run("can write json file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()

		cfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Packages: []config.Package{
					{
						Name:    "test-package",
						Version: "1.0.0",
					},
				},
			},
		}

		filePath := "/tmp/test_write.json"
		err := config.WriteFile(fs, filePath, cfg)
		require.NoError(t, err)

		exists, err := fs.Exists(filePath)
		require.NoError(t, err)
		require.True(t, exists)

		// Verify content
		readCfg, err := config.ReadFile(fs, filePath)
		require.NoError(t, err)
		require.Equal(t, cfg.ApiVersion, readCfg.ApiVersion)
		require.Equal(t, len(cfg.Spec.Packages), len(readCfg.Spec.Packages))
		require.Equal(t, cfg.Spec.Packages[0].Name, readCfg.Spec.Packages[0].Name)
	})
}

func TestWriteFileIfNotExists(t *testing.T) {
	t.Run("creates file when it doesn't exist", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()

		filePath := "/tmp/write_if_not_exists.yml"
		defaultFactory := func() (config.Config, error) {
			return config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
				Spec: config.Spec{
					Environment: map[string]string{
						"NEW": "value",
					},
				},
			}, nil
		}

		err := config.WriteFileIfNotExists(fs, filePath, defaultFactory)
		require.NoError(t, err)

		exists, err := fs.Exists(filePath)
		require.NoError(t, err)
		require.True(t, exists)

		readCfg, err := config.ReadFile(fs, filePath)
		require.NoError(t, err)
		require.Equal(t, "value", readCfg.Spec.Environment["NEW"])
	})

	t.Run("doesn't overwrite existing file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()

		filePath := "/tmp/existing_write_if_not_exists.yml"
		originalCfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Environment: map[string]string{
					"ORIGINAL": "value",
				},
			},
		}

		// Write original file
		err := config.WriteFile(fs, filePath, originalCfg)
		require.NoError(t, err)

		defaultFactory := func() (config.Config, error) {
			return config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
				Spec: config.Spec{
					Environment: map[string]string{
						"DEFAULT": "value",
					},
				},
			}, nil
		}

		err = config.WriteFileIfNotExists(fs, filePath, defaultFactory)
		require.NoError(t, err)

		// Verify original content is preserved
		readCfg, err := config.ReadFile(fs, filePath)
		require.NoError(t, err)
		require.Equal(t, "value", readCfg.Spec.Environment["ORIGINAL"])
		require.NotContains(t, readCfg.Spec.Environment, "DEFAULT")
	})
}

func TestReadOrCreateFile(t *testing.T) {
	t.Run("creates and reads file when it doesn't exist", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()

		filePath := "/tmp/read_or_create_file.yml"
		defaultFactory := func() (config.Config, error) {
			return config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
				Spec: config.Spec{
					Environment: map[string]string{
						"CREATED": "value",
					},
				},
			}, nil
		}

		readCfg, err := config.ReadOrCreateFile(fs, filePath, defaultFactory)
		require.NoError(t, err)
		require.Equal(t, config.ApiVersion, readCfg.ApiVersion)
		require.Equal(t, "value", readCfg.Spec.Environment["CREATED"])

		// Verify file was created
		exists, err := fs.Exists(filePath)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("reads existing file", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()

		filePath := "/tmp/existing_read_or_create_file.yml"
		originalCfg := config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Environment: map[string]string{
					"EXISTING": "value",
				},
			},
		}

		// Write original file
		err := config.WriteFile(fs, filePath, originalCfg)
		require.NoError(t, err)

		defaultFactory := func() (config.Config, error) {
			return config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
				Spec: config.Spec{
					Environment: map[string]string{
						"DEFAULT": "value",
					},
				},
			}, nil
		}

		readCfg, err := config.ReadOrCreateFile(fs, filePath, defaultFactory)
		require.NoError(t, err)
		require.Equal(t, "value", readCfg.Spec.Environment["EXISTING"])
		require.NotContains(t, readCfg.Spec.Environment, "DEFAULT")
	})
}

func TestGetEncoding(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected config.Encoding
		hasError bool
	}{
		{
			name:     "yaml extension",
			filename: "config.yaml",
			expected: config.Yaml,
			hasError: false,
		},
		{
			name:     "yml extension",
			filename: "config.yml",
			expected: config.Yaml,
			hasError: false,
		},
		{
			name:     "json extension",
			filename: "config.json",
			expected: config.Json,
			hasError: false,
		},
		{
			name:     "uppercase yaml",
			filename: "config.YAML",
			expected: config.Yaml,
			hasError: false,
		},
		{
			name:     "uppercase json",
			filename: "config.JSON",
			expected: config.Json,
			hasError: false,
		},
		{
			name:     "unsupported extension",
			filename: "config.txt",
			expected: "",
			hasError: true,
		},
		{
			name:     "no extension",
			filename: "config",
			expected: "",
			hasError: true,
		},
		{
			name:     "with path",
			filename: filepath.Join("path", "to", "config.yml"),
			expected: config.Yaml,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: getEncoding is not exported, so we test it indirectly through file operations
			h := cross.NewTest(platform.Linux, arch.AMD64)
			fs := h.FS()

			cfg := config.Config{
				ApiVersion: config.ApiVersion,
				Kind:       config.Kind,
			}

			err := config.WriteFile(fs, "/tmp/"+tt.filename, cfg)
			if tt.hasError {
				require.Error(t, err)
				if !tt.hasError {
					require.Contains(t, err.Error(), "unable to determine encoding")
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
