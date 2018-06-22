# Automation Manager CLI

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
     run, r   run a command
     env, e   print command environemnt variables
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE (default: "C:\\Users\\patri\\.cli-mgr\\config.yml") [%CLI_MGR_CONFIG%]
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

### env command

```bash
cli-mgr env --help
```

```
NAME:
   cli-mgr env - print command environemnt variables

USAGE:
   cli-mgr env [command options] [arguments...]

OPTIONS:
   --name NAME, -n NAME                       Execute command named NAME
   --environment ENVIRONMENT, -e ENVIRONMENT  Use environment named ENVIRONMENT
```

## building

clone submodules

```
git submodule update --init --recursive
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

