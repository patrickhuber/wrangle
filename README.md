# Wrangle 

A tool for managing devops environments

## Features

* Package Management
* Secrets Access
* Cascading Variables

## Getting Started

These snippets will run the scripts directly from the internet. You can optionally download the install files and run them by hand.

```bash
curl https://raw.githubusercontent.com/patrickhuber/wrangle/main/install.sh | bash
```

```powershell
# verbose expression
Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/patrickhuber/wrangle/main/install.ps1')

# or short expression
iwr -useb 'https://raw.githubusercontent.com/patrickhuber/wrangle/main/install.ps1' | iex
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