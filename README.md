# Wrangle 

A tool for managing devops environments

## Features

* Package Management
* Secrets Access
* Cascading Variables

## Getting Started

There are three methods to getting wrangle on a PC. You can do `go install`, a manual download and install or run one of the scripts below. 

### Go Install

```
go install github.com/patrickhuber/wrangle/cmd/wrangle@0.10.0
```

OR

```
git clone https://github.com/patrickhuber/wrangle
cd wrangle
go install cmd/wrangle@0.10.0
```

### Scripted Install

These snippets will run the scripts directly from the internet. You can also download the scripts and run them by hand. 

```bash
curl https://raw.githubusercontent.com/patrickhuber/wrangle/main/scripts/install.sh | bash
```

```powershell
# verbose expression
Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/patrickhuber/wrangle/main/scripts/install.ps1')

# or short expression
iwr -useb 'https://raw.githubusercontent.com/patrickhuber/wrangle/main/scripts/install.ps1' | iex
```

### Manual Install

The scripts above download the latest version of wrangle and then run a `wrangle bootstrap` command. The manual install involves the same steps run by hand.

Linux

```bash
export VERSION=0.10.0
export ARCHIVE="wrangle-${VERSION}-linux-amd64.tar.gz"
wget https://github.com/patrickhuber/wrangle/releases/download/${VERSION}/${ARCHIVE}
tar xfz ${ARCHIVE}
rm ${ARCHIVE}
chmod +x wrangle
wrangle bootstrap
rm wrangle
rm README.md
```

Darwin

```bash
export VERSION=0.10.0
export ARCHIVE="wrangle-${VERSION}-darwin-amd64.tar.gz"
wget https://github.com/patrickhuber/wrangle/releases/download/${VERSION}/${ARCHIVE}
tar xfz ${ARCHIVE}
rm ${ARCHIVE}
chmod +x wrangle
wrangle bootstrap
rm wrangle
rm README.md
```

Windows (Powershell)

```powershell
$version = "0.10.0"
$archive = "wrangle-$version-windows-amd64.zip"
iwr -Uri "https://github.com/patrickhuber/wrangle/releases/download/$version/$archive" -OutFile $archive
Extract-Archive $archive -DestinationPath .
Remove-Item $archive
.\wrangle.exe bootstrap
Remove-Item wrangle.exe
Remove-Item README.md
```

## Usage

Once wrangle is installed you can install the latest package with the install command.

```bash
wrangle install yq
```

or a specific version

```
wrangle install yq@4.31.1
```

You can also install packages by creating a .wrangle(.yml|.json) file in the directory and run the `wrangle restore` command.

> .wrangle.yml

```yaml
apiVersion: wrangle/v1
kind: Config
spec:
  packages:
  - name: jq
    version: 4.31.1

```

> .wrangle.json

```json
{
    "apiVersion": "wrangle/v1",
    "kind": "Config",
    "spec": {
        "packages": [
            {
                "name": "jq",
                "version": "4.31.1"
            }
        ]
    }
}
```

```bash
wrangle restore
```

## Shell Integration

Wrangle can integrate into your shell to enable environment variable injection. 

> bash

add the following to the ~/.bashrc file

```bash
eval "$(wrangle hook bash)"
```

> powershell

add the following to end of the $PROFILE file

```powershell
iex $(wrangle hook powershell | Out-String)
```

## Variable Replacement

In instances when you need variable replacement for secrets or other configuraiton, wrangle supports external stores. Variables are surrounded with double parenthesis `((<VariablePath>))`, where <VariablePath> is the path to the variable in a store. 

The following stores are supported:

### Azure Key Vault 

```yaml
stores:
- name: default
  type: azure.keyvault  
  properties:
    uri: {key vault uri} // (required)
```

| property | description |
| -------- | ----------- |
| uri      | the uri to the key vault | https://quickstart-kv.vault.azure.net |

### Keyring

```yaml
stores:
- name: default
  type: keyring
  properties:
    service: test
```

| property | description |
| -------- | ----------- |
| service  | the service under which to store secrets |

## Package Management

Wrangle is a simple package manager much like [arkade](https://github.com/alexellis/arkade). Arkade tends to focus on the lastest package version while Wrangle is fully version aware. Arkade also embeds its package feed into the execuable while Wrangle utilizes external package feeds.

Wrangle can find packages in one or more feeds. The default feed is https://github.com/patrickhuber/wrangle-packages.

When wrangle bootstrap is called, wrangle will go to default feed (or the overrided feed) and install itself and a shim execuable.

### Feeds

The default feed is located at https://github.com/patrickhuber/wrangle-packages. A feed is a git repository with a top level `/feed` folder. 

In the feed folder, each package is created as a folder under /feed. For example, wrangle itself is located under `/feed/wrangle`. 

Under the package folder, each version has its own folder as well. For example, wrangle version 0.9.0 is located under `/feed/wrangle/0.9.0`.

### Packages

### Versions

## Logging and Debugging

To enable logging, set the WRANGLE_LOG_LEVEL environment variable. 

The following values are accepted:

| level     | description |
| --------- | ----------- |
| debug     | all, verbose line level | 
| info      | informational, warnings and errors |
| warn      | warnings and errors |
| error     | errors only (default) |