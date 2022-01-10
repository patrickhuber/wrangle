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
curl https://raw.githubusercontent.com/patrickhuber/wrangle/main/install.sh | bash
```

```powershell
# verbose expression
Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/patrickhuber/wrangle/main/install.ps1')

# or short expression
iwr -useb 'https://raw.githubusercontent.com/patrickhuber/wrangle/main/install.ps1' | iex
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

## Package Management

Concepts

* Feeds
* Packages
* Versions

Commands

* Install
* Upgrade

### Feeds

### Packages

### Versions