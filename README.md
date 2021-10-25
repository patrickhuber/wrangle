# Wrangle 

A tool for managing devops environments

## Features

* Package Management
* Secrets Access
* Cascading Variables

## Getting Started

These snippets will run the scripts directly from the internet. You can optionally download the install files and run them by hand.

```bash
bash <(curl https://github.com/patrickhuber/wrangle/master/install.sh)
```

```powershell
# verbose expression
Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://github.com/patrickhuber/wrangle/master/install.ps1')

# or short expression
iwr -useb 'https://github.com/patrickhuber/wrangle/master/install.ps1' | iex
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