# Wrangle 

A tool for managing devops environments

## Features

* Package Management
* Secrets Access
* Cascading Variables

## Getting Started

There are two methods to getting wrangle on a PC. You can do manual download and install or run one of the scripts below. 

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

```
export VERSION=0.10.0
export ARCHIVE=wrangle-darwin-amd64.tgz
wget https://github.com/patrickhuber/wrangle/releases/download/${VERSION}/${ARCHIVE}
tar -xfz ${ARCHIVE}
rm ${ARCHIVE}
chmod +x wrangle
wrangle bootstrap
rm wrangle
```

Darwin

```
export VERSION=0.10.0
export ARCHIVE=wrangle-darwin-amd64.tgz
wget https://github.com/patrickhuber/wrangle/releases/download/${VERSION}/${ARCHIVE}
tar -xfz ${ARCHIVE}
rm ${ARCHIVE}
chmod +x wrangle
wrangle bootstrap
rm wrangle
```

Windows (Powershell)

```
$version = "0.10.0"
$archive = "wrangle-windows-amd64.zip"
iwr -Uri "https://github.com/patrickhuber/wrangle/releases/download/$VERSION/$archive" -OutFile $archive
Extract-Archive $archive
Remove-Item $archive
.\wrangle.exe bootstrap
Remove-Item wrangle.exe
```

### Usage

Once wrangle is installed you can install a package with the install command.

```bash
wrangle install yq
```

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