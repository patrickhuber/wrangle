# Wrangle

Wrangle is a solution for managing multiple configuration stores and multiple command line interfaces across multiple environments. It is meant to address team collaboration by using declarative configuration that is committed to source control. 

## Getting Started

### Download 

The latest release can be found on the releases page of github. Download the release for your target platform

https://github.com/patrickhuber/wrangle/releases


### Install

Use tar or 7zip to extract and decompress the package. 

Add move the binary to a folder under your PATH environment variable. 

### Environment Variables

For ease of use you can set the WRANGLE_CONFIG_FILE and WRANGLE_PACKAGE_PATH environment variables. 

* WRANGLE_CONFIG_FILE - specifies where the configuration will reside, similar to bosh bootloader's BBL_STATE_DIR, this is set to your working directory. You can also specify this file with the -c flag. 
* WRANGLE_PACKAGE_PATH - specifies where packages will be extracted and linked. Make sure this is in your PATH environment variable as each CLI is stored here and a symlink is created for the most recentlly installed. You can also specify this path with the -p flag. 


### The Configuration File

An example configuration file is located here [config file](doc/example-config.yml)

The default location for the configuration file is in your user directory under:

> mac & linux

`~/.wrangle/config.yml`

> windows

`%userprofile%\.wrangle\config.yml`

The config file has the following structure:

```
stores:
environments:
packages:
```

#### Stores

Stores are configuration sources that can be cascaded to configure the cli commands. There are two config sources by default:

1. file - loads yaml files 
2. credHub - loads creds from credhub
3. env - loads from environment variables

Additional config sources could be other key managers like Vault, LastPass, Amazon Key Management service etc. 

Stores can receive configuration from other stores through their "stores: " list. 

##### FILE

this is an example file config

```yml
stores:

- name: bosh-lab-yaml
  type: file
  params:
    path: state/vars/director-vars-store.yml
```

##### CREDHUB

This is an example credhub config. You can see this config references the bosh-lab-yaml configuration where it will read any variables defined in the params.  

```yml
stores:
- name: bosh-lab-credhub
  type: credhub
  stores: 
  - bosh-lab-yaml
  params:
    client_id: credhub_admin
    client_secret: ((credhub_admin_client_secret))
    server: https://192.168.3.11:4343
    ca_cert: ((credhub_ca.certificate))
    skip_tls_validation: false
```

> note

If the `CREDHUB_PROXY` environment variable is set, wrangle will use it to make the connection. This is the default behavior of the credhub cli and wrangle imports that cli as a library.

##### ENV

This is an example env config. You can load any enivironment variable as a config variable. The environment variable must be present, or the lookup will fail. If the variable is defined but is empty the lookup will succeed.

```yml
stores:
- name: environment
  type: env
  params:
    some_variable: SOME_VARIABLE 
```

With the configuration example above, `some_variable` is now available as a variable in consumers of this store.

#### Environments

Environments allow for different parameters to be passed to CLIs that may share the same name. For example, a customer may have several credhubs across lab and production environments as well as several in a control plane and PCF install. Environments provide a easy grouping to avoid name conflicts when attempting to run a cli. 

Each environment has a list of processes. Processes are comprised of the path to the process as well as arguments, environment variables and a list of configurations. Arguments and Environment Variables can contain variables which can be looked up in configurations. 

The configurations are evaluated in the order they are specified. Variables can be cascaded and will be resolved in the order the configuraitons are listed. 

Here is an example process that run fly login:

```yml
environments:

- name: lab
  processes:
  
  - name: fly
    stores: 
    - bosh-lab-credhub
    path: fly
    args:
    - -t
    - main
    - login
    - -u
    - ((/bosh-lab/concourse/atc_basic_auth.username))
    - -p 
    - ((/bosh-lab/concourse/atc_basic_auth.password))
```

It assumes the cli is the environment PATH, if you placed your WRANGLE_PACKAGE_PATH environment variable in the PATH, the above will resolve once you install the fly package. 

This is an example of running the command above using the wrangle:

```bash
wrangle run -e lab -n fly
```

If you would like to see the output that would be executed you can use the `print` command

```bash
wrangle print -e lab -n fly
```

If you would just like to print the environment variables you can use the `print-env` command


```bash
wrangle print-env -e lab -n fly
```

#### Packages

Packages allow the cli manager tool to actually manage CLIs. Each package can target multiple platforms. Packages have a download and an extract step. If the package downloaded is a binary, the extract step can be skipped. If the Package is a tarball, tgz or zip file, the extract step can be used to extract the binary. 

Each package platform has an alias that will be used to create a symlink to the fully versioned name of the package. This allows scripts to reference the short name of the executable while allowing multiple CLIs to be installed. 

Packages support the ((version)) variable in the url, out and filter parameters of `extract` and `download`. The variable will be replaced with the version specified as part of the package specification. 

The filter parameter of the extract is a regex filter that will be used to include a single file. This pattern will match one file and quit so it should be as specific as possible in order to extract the proper binary for the package platform. 

This is an example that both downloads and extracts the credhub cli:

```yml
- name: credhub
  version: 1.7.6  
  platforms:    
  - name: linux
    alias: credhub
    download: 
      url: https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/((version))/credhub-linux-((version)).tgz      
      out: credhub-((version))-linux.tgz
    extract:
      filter: credhub
      out: credhub-((version))-linux
  - name: darwin
    alias: credhub
    download: 
      url: https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/((version))/credhub-darwin-((version)).tgz      
      out: credhub-((version))-darwin.tgz
    extract:
      filter: credhub
      out: credhub-((version))-darwin
  - name: windows
    alias: credhub.exe
    download: 
      url: https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/((version))/credhub-windows-((version)).tgz      
      out: credhub-((version))-windows.tgz
    extract:
      filter: credhub
      out: credhub-((version))-windows.exe
```

This example simply downloads the binary and doesn't do any extraction

```yml
- name: bosh
  version: 3.0.1  
  platforms:
  - name: linux
    alias: bosh
    download:
      url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-linux-amd64
      out: bosh-cli-((version))-linux-amd64
  - name: windows
    alias: bosh.exe
    download:
      url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-windows-amd64.exe
      out: bosh-cli-((version))-windows-amd64.exe
  - name: darwin
    alias: bosh
    download:
      url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-darwin-amd64
      out: bosh-cli-((version))-darwin-amd64
```

## usage

```
NAME:
   wrangle - a cli management tool

USAGE:
   wrangle [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     run, r              run a command
     print, p            print command environemnt variables
     environments, e     prints the list of environments in the config file
     packages, k         prints the list of packages and versions in the config file
     install, i  installs the package with the given `NAME` for the current platform
     help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE (default: "/home/patrick/.wrangle/config.yml") [$WRANGLE_CONFIG]
   --help, -h              show help
   --version, -v           print the version
```

### run command

```bash
wrangle run --help
```

```
NAME:
   wrangle run - run a command

USAGE:
   wrangle run [command options] [arguments...]

OPTIONS:
   --name NAME, -n NAME                       Execute command named NAME
   --environment ENVIRONMENT, -e ENVIRONMENT  Use environment named ENVIRONMENT
```

### print command

```bash
wrangle print --help
```

```
NAME:
   wrangle print - print command environemnt variables

USAGE:
   wrangle print [command options] [arguments...]

OPTIONS:
   --name NAME, -n NAME                       process named NAME
   --environment ENVIRONMENT, -e ENVIRONMENT  Use environment named ENVIRONMENT
```

### environments command

```bash
wrangle environments --help
```

```
NAME:
   wrangle environments - prints the list of environments in the config file

USAGE:
   wrangle environments [arguments...]
```

### packages command

```bash
wrangle packages --help
```

```
NAME:
   wrangle packages - prints the list of packages and versions in the config file

USAGE:
   wrangle packages [arguments...]
```

### install-package command

```bash
wrangle install-package --help
```

```
NAME:
   wrangle install-package - installs the package with the given `NAME` for the current platform

USAGE:
   wrangle install-package [command options] [arguments...]

OPTIONS:
   --name NAME, -n NAME    package named NAME
   --path value, -p value  the package install path [$WRANGLE_PACKAGE_PATH]
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

[creds file](doc/example-creds.yml)