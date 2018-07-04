# CLI Manager

Template and simplify command line operations for several CLIs in the cloud foundry ecosystem.

## usage

```
NAME:
   cli-mgr - a cli management tool

USAGE:
   cli-mgr [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     run, r              run a command
     print, p            print command environemnt variables
     environments, e     prints the list of environments in the config file
     packages, k         prints the list of packages and versions in the config file
     install-package, i  installs the package with the given `NAME` for the current platform
     help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE (default: "/home/patrick/.cli-mgr/config.yml") [$CLI_MGR_CONFIG]
   --help, -h              show help
   --version, -v           print the version
```

### run command

```bash
cli-mgr run --help
```

```
NAME:
   cli-mgr run - run a command

USAGE:
   cli-mgr run [command options] [arguments...]

OPTIONS:
   --name NAME, -n NAME                       Execute command named NAME
   --environment ENVIRONMENT, -e ENVIRONMENT  Use environment named ENVIRONMENT
```

### print command

```bash
cli-mgr print --help
```

```
NAME:
   cli-mgr print - print command environemnt variables

USAGE:
   cli-mgr print [command options] [arguments...]

OPTIONS:
   --name NAME, -n NAME                       process named NAME
   --environment ENVIRONMENT, -e ENVIRONMENT  Use environment named ENVIRONMENT
```

### environments command

```bash
cli-mgr environments --help
```

```
NAME:
   cli-mgr environments - prints the list of environments in the config file

USAGE:
   cli-mgr environments [arguments...]
```

### packages command

```bash
cli-mgr packages --help
```

```
NAME:
   cli-mgr packages - prints the list of packages and versions in the config file

USAGE:
   cli-mgr packages [arguments...]
```

### install-package command

```bash
cli-mgr install-package --help
```

```
NAME:
   cli-mgr install-package - installs the package with the given `NAME` for the current platform

USAGE:
   cli-mgr install-package [command options] [arguments...]

OPTIONS:
   --name NAME, -n NAME    package named NAME
   --path value, -p value  the package install path [$CLI_MGR_PACKAGE_INSTALL_PATH]
```

## building

to restore packages (requires dep)

```
make restore
```

to perform a build

```
make build
```

## testing

to run unit tests

```
make unit
```

## sample files

[config file](doc/example-config.yml)

[creds file](doc/example-creds.yml)