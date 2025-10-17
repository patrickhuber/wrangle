package install_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/actions"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/feed"
	feedmemory "github.com/patrickhuber/wrangle/internal/feed/memory"
	"github.com/patrickhuber/wrangle/internal/fixtures"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/install"
	"github.com/patrickhuber/wrangle/internal/packages"
	"github.com/patrickhuber/wrangle/internal/shim"
	"github.com/stretchr/testify/require"
)

func TestInstall(t *testing.T) {
	type packageTest struct {
		name    string
		version string
	}
	type fileTest struct {
		platform platform.Platform
	}
	files := []fileTest{
		{
			platform: platform.Windows,
		},
		{
			platform: platform.Linux,
		},
		{
			platform: platform.Darwin,
		},
	}
	packages := []packageTest{
		{
			name:    "test",
			version: "latest",
		},
		{
			name:    "test",
			version: "1.0.0",
		},
	}
	for _, f := range files {
		for _, p := range packages {
			t.Run(fmt.Sprintf("%s_%s_%s", f.platform.String(), p.name, p.version), func(t *testing.T) {
				RunInstallTest(t, p.name, p.version, f.platform)
			})
		}
	}
}

func RunInstallTest(t *testing.T,
	packageName string,
	packageVersion string,
	plat platform.Platform) {

	target := cross.NewTest(plat, arch.AMD64)
	err := fixtures.Apply(target.OS(), target.FS(), target.Env())
	require.NoError(t, err)

	root, err := config.GetRoot(target.Env(), target.Path(), plat)
	require.NoError(t, err)

	packagesDir := config.GetDefaultPackagesPath(target.Path(), root)
	appName, err := config.GetAppName("test", plat)
	require.NoError(t, err)

	actualPackageVersion := packageVersion
	if actualPackageVersion == "latest" {
		actualPackageVersion = "1.0.0"
	}

	cfg := config.Config{
		Spec: config.Spec{
			Environment: map[string]string{
				global.EnvPackages: packagesDir,
			},
			Feeds: []config.Feed{
				{
					Name: feedmemory.ProviderType,
					Type: feedmemory.ProviderType,
				},
			},
		},
	}

	metadataProvider := actions.NewMetadataProvider(target.Path())
	logger := log.Default(log.WithLevel(log.DebugLevel))
	configuration := config.NewMock(cfg)

	service := install.NewService(
		target.FS(),
		feed.NewServiceFactory(feedmemory.NewProvider(logger, []*feed.Item{
			{
				State: &feed.State{
					LatestVersion: actualPackageVersion,
				},
				Package: &packages.Package{
					Name: packageName,
					Versions: []*packages.Version{
						{
							Version: actualPackageVersion,
							Manifest: &packages.Manifest{
								Package: &packages.ManifestPackage{
									Name:    packageName,
									Version: actualPackageVersion,
									Targets: []*packages.ManifestTarget{
										{
											Platform:     plat,
											Architecture: arch.AMD64,
											Steps: []*packages.ManifestStep{
												{
													Action: "move",
													With: map[string]any{
														"source":      appName,
														"destination": appName,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}...)),
		actions.NewRunner(
			actions.NewFactory(
				actions.NewMoveProvider(target.FS(), target.Path(), logger),
			)),
		target.OS(),
		configuration,
		metadataProvider,
		target.Path(),
		shim.NewService(target.FS(), target.Path(), configuration, logger),
		logger)

	req := &install.Request{
		Package: packageName,
		Version: packageVersion,
	}

	// write out package file
	metadata := metadataProvider.Get(&cfg, packageName, actualPackageVersion)
	packageVersionFileLocation := target.Path().Join(metadata.PackageVersionPath, appName)
	target.FS().WriteFile(packageVersionFileLocation, []byte("test"), 0644)

	err = service.Execute(req)
	require.NoError(t, err)

	fs := target.FS()

	ok, err := fs.Exists(packageVersionFileLocation)
	require.NoError(t, err)
	require.True(t, ok, "file '%s' not found", packageVersionFileLocation)
}
