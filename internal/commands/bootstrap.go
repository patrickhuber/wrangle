package commands

import (
	"io"
	"os"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/env"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/operatingsystem"
	"github.com/urfave/cli/v2"
)

type BootstrapOptions struct {
	FileSystem      filesystem.FileSystem
	OperatingSystem operatingsystem.OS
	Environment     env.Environment
	Config          *config.Config
	ApplicationName string
	GlobalPath      string
	Force           bool
}

func Bootstrap(ctx *cli.Context) error {
	return BootstrapInternal(&BootstrapOptions{
		Config: &config.Config{
			PackagePath: ctx.String("packages"),
			BinPath:     ctx.String("bin"),
			RootPath:    ctx.String("root"),
		},
		FileSystem:      ctx.App.Metadata["fileSystem"].(filesystem.FileSystem),
		OperatingSystem: ctx.App.Metadata["os"].(operatingsystem.OS),
		ApplicationName: ctx.App.Name,
		GlobalPath:      ctx.String("global"),
		Force:           ctx.Bool("force"),
	})
}

func BootstrapInternal(opt *BootstrapOptions) error {

	binary := opt.ApplicationName

	// get the current working directory from the operating system
	workingDirectory, err := opt.OperatingSystem.WorkingDirectory()
	if err != nil {
		return err
	}

	// TODO: install the required packages for wrangle
	// copy the currently running cli to the bin path
	targetFilePath := crosspath.Join(opt.Config.BinPath, opt.ApplicationName)
	sourceFilePath := crosspath.Join(workingDirectory, binary)

	// writer from target file
	writer, err := opt.FileSystem.OpenFile(targetFilePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()

	// reader from source file
	reader, err := opt.FileSystem.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// copy the reader to the writer
	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}

	// get the home directory
	globalPath := opt.GlobalPath
	globalProvider := config.NewFileProvider(opt.FileSystem, globalPath)

	// check that the global config exists
	ok, err := opt.FileSystem.Exists(globalPath)
	if err != nil {
		return err
	}

	// if the file exists or force is disabled, return
	if ok && !opt.Force {
		return nil
	}

	// set the global config
	return globalProvider.Set(opt.Config)
}
