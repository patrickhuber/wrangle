# Wrangle

Wrangle is a solution for managing multiple configuration stores and multiple command line interfaces across multiple environments. It is meant to address team collaboration by using declarative configuration that is committed to source control. 

## Getting Started

### Download 

The latest release can be found on the releases page of github. Download the release for your target platform

https://github.com/patrickhuber/wrangle/releases

### Install

Use tar or 7zip to extract and decompress the package. 

Set the WRANGLE_BIN environment variable. Add WRANGLE_BIN to your path. 

Place the wrangle executable in WRANGLE_BIN

### Environment Variables

For ease of use you can set the WRANGLE_CONFIG and WRANGLE_PACKAGES environment variables. 

* WRANGLE_CONFIG - specifies where the configuration will reside, similar to bosh bootloader's BBL_STATE_DIR, this is set to your working directory. You can also specify this file with the -c flag. 
* WRANGLE_PACKAGES - specifies where package versions will be installed. You can also specify this path with the -p flag. 
* WRANGLE_ROOT - the root directory where wrangle stores artifacts
* WRANGLE_BIN - the location where wrangle stores symlinks (or shims)

### The Configuration File

An example configuration file is located here [config file](examples/example-config.yml)

The default location for the config file is the current working directory

The config file has the following structure:

```
stores:
processes:
imports:
```

#### Stores

Stores are configuration sources that can be cascaded to configure the cli commands. 

The available stores are listed below:

1. file - loads yaml files 
2. credHub - loads creds from credhub
3. env - loads from environment variables
4. meta - contains meta information about the current wrangle config file
5. prompt - prompts the user to enter information over stdin

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

Gpg encrypted files are supported when the path points to a file with 'gpg' extension. The default gnupg keyring is used. These are the search  paths:

> windows

```
%APPDATA%\gnupg\pubring.gpg
%APPDATA%\gnupg\secring.gpg
```

> linux and darwin

```
$HOME/.gnupg/pubring.gpg
$HOME/.gnupg/secring.gpg
```

Golang only support gnupg v1 pubring.gpg and secring.gpg key files. If you have v2, you need to export the secret and public keys using the following commands from the gnupg directory:

> windows

```
cd %APPDATA%\gnupg
gpg --export-secret-keys --output secring.gpg
gpg --export --output pubring.gpg
```

> linux and darwin

```
cd $HOME/.gnupg
gpg --export-secret-keys --output secring.gpg
gpg --export --output pubring.gpg
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

#### META

The meta store provides contextual information about the current config file. This can be useful when you need to know the directory of the config file or config file containing directory for loading other files.

> command

```
wrangle list -s meta
```

> result

```
bin: "/opt/wrangle/bin"
root: "/opt/wrangle"
packges_folder: "/opt/wrangle/packages"
config_file: "/home/abc/source/github.com/org/repo/wrangle.yml"
config_file_folder: "/home/abc/source/github.com/org/repo"
```

#### Packages

Packages are currently just ways of installing self contained binaries. Each package can target multiple platforms. Packages have a download and an extract step. If the package downloaded is a binary, the extract step can be skipped. If the Package is a tarball, tgz or zip file, the extract step can be used to extract the binary. 

Each package platform has an alias that will be used to create a symlink to the fully versioned name of the package. This allows scripts to reference the short name of the executable while allowing multiple CLIs to be installed. 

Packages support the ((version)) variable in the url, out and filter parameters of `extract` and `download`. The variable will be replaced with the version specified as part of the package specification. 

The filter parameter of the extract is a regex filter that will be used to include a single file. This pattern will match one file and quit so it should be as specific as possible in order to extract the proper binary for the package platform. 

This is an example that both downloads and extracts the credhub cli:

```yml
- name: credhub
  version: 1.7.6  
  targets:    
  - platform: linux
    tasks:    
    - download: 
        url: https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/((version))/credhub-linux-((version)).tgz      
        out: credhub-((version))-linux.tgz
    - extract:        
        archive: credhub-((version))-linux.tgz
  - platform: darwin    
    tasks:
    - download:
        url: https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/((version))/credhub-darwin-((version)).tgz      
        out: credhub-((version))-darwin.tgz
    - extract:
        archive: credhub-((version))-darwin.tgz
  - name: windows    
    tasks:
    - download: 
        url: https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/((version))/credhub-windows-((version)).tgz      
        out: credhub-((version))-windows.tgz
    - extract:        
        archive: credhub-((version))-windows.tgz
```

This example simply downloads the binary and doesn't do any extraction

```yml
- name: bosh
  version: 3.0.1  
  targets:
  - platform: linux
    tasks:
    - download:
        url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-linux-amd64
        out: bosh-cli-((version))-linux-amd64
  - platform: windows    
    tasks:
    - download:
        url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-windows-amd64.exe
        out: bosh-cli-((version))-windows-amd64.exe
  - platform: darwin    
    tasks:
    - download:
        url: https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-((version))-darwin-amd64
        out: bosh-cli-((version))-darwin-amd64
```

## sample files

[creds file](doc/example-creds.yml)