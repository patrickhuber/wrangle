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
go install github.com/patrickhuber/wrangle/cmd/wrangle@v0.10.8
```

OR

```
git clone https://github.com/patrickhuber/wrangle
cd wrangle
go install cmd/wrangle@v0.10.8
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
export VERSION=v0.10.8
export ARCHIVE="wrangle-${VERSION}-linux-amd64.tar.gz"
wget https://github.com/patrickhuber/wrangle/releases/download/v${VERSION}/${ARCHIVE}
tar xfz ${ARCHIVE}
rm ${ARCHIVE}
chmod +x wrangle
wrangle bootstrap
rm wrangle
rm README.md
```

Darwin

```bash
export VERSION=v0.10.8
export ARCHIVE="wrangle-${VERSION}-darwin-amd64.tar.gz"
wget https://github.com/patrickhuber/wrangle/releases/download/v${VERSION}/${ARCHIVE}
tar xfz ${ARCHIVE}
rm ${ARCHIVE}
chmod +x wrangle
wrangle bootstrap
rm wrangle
rm README.md
```

Windows (Powershell)

```powershell
$version = "v0.10.8"
$archive = "wrangle-$version-windows-amd64.zip"
iwr -Uri "https://github.com/patrickhuber/wrangle/releases/download/v$version/$archive" -OutFile $archive
Expand-Archive $archive -DestinationPath .
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

> .wrangle.toml

```toml
apiVersion = "wrangle/v1"
kind = "Config"

[spec.packages]
name = "jq"
version = "4.31.1"
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

| property | description | values |
| -------- | ----------- | ------ |
| uri      | the uri to the key vault | https://quickstart-kv.vault.azure.net |

### Keyring

> file

```yaml 
stores:
- name: default
  type: keyring
  properties:
    service: test
    allowed_backends: 
    - file
    file.directory: ~/
    file.password: abc123    
```

> pass

```yaml
- name: default
  type: keyring
  properties:
    service: test
    allowed_backends: 
    - pass
    pass.directory: ~/
    pass.command: /usr/bin/pass
    pass.prefix: ""
```

| property         | description                              | reqired | values |
| ---------------- | ---------------------------------------- | ------- | ------ |
| service          | the service under which to store secrets | yes     | |
| allowed_backends | the backends allowed                     | no      | file, secret-service, keychain, keyctl, kwallet, wincred, file, pass |
| file.directory   | file backend directory                   | no      | |
| file.password    | file backend password                    | no      | |
| pass.directory   | pass backend directory                   | no      | | 
| pass.command     | path to pass command                     | no      | |
| pass.prefix      | key prefix                               | no      | |

### HashiCorp Vault

#### Token Authentication

```yaml
stores:
- name: default
  type: vault
  properties:
    address: http://127.0.0.1:8200  # (required)
    token: ((vault_token))  # (optional)
    path: secret  # (optional, defaults to "secret")
variables:
- name: vault_token
  type: password
```

#### AppRole Authentication

```yaml
stores:
- name: default
  type: vault
  properties:
    address: http://127.0.0.1:8200  # (required)
    role_id: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx  # (required for AppRole)
    secret_id: ((vault_secret_id))  # (required for AppRole)
    path: secret  # (optional, defaults to "secret")
variables:
- name: vault_secret_id
  type: password
```

| property  | description |
| --------- | ----------- |
| address   | the address of the Vault server (e.g., http://127.0.0.1:8200) |
| token     | the authentication token (optional, can use VAULT_TOKEN environment variable) |
| role_id   | the AppRole role ID for AppRole authentication (optional) |
| secret_id | the AppRole secret ID for AppRole authentication (optional) |
| path      | the KV v2 secrets engine mount path (optional, defaults to "secret") |

**Authentication Methods** (in order of precedence):
1. **AppRole**: Provide both `role_id` and `secret_id`
2. **Token**: Provide `token`
3. **Environment Variables**: Use `VAULT_TOKEN` and `VAULT_ADDR`

## Variables

Wrangle supports defining variables for secret generation following the BOSH variable types specification. Variables define how secrets should be generated, specifying only the properties needed for generation—not the actual values.

### Variable Types

Variables can be defined in your configuration file and listed using the `wrangle list variables` command.

#### Certificate

Generate TLS/SSL certificates with customizable properties.

```yaml
variables:
  - name: root-ca
    type: certificate
    options:
      common_name: "Root CA"
      is_ca: true
      key_length: 4096
  - name: server-cert
    type: certificate
    options:
      ca: root-ca  # Sign with root-ca
      common_name: "server.example.com"
      alternative_names:
        - www.example.com
        - api.example.com
      extended_key_usage:
        - server_auth
      duration: 365
```

**Options:**

| Option | Description | Required | Default |
| ------ | ----------- | -------- | ------- |
| `common_name` | Common name for the certificate | Yes | |
| `ca` | Name of a CA variable to sign the certificate | No | |
| `organization` | Organization name | No | |
| `alternative_names` | Subject alternative names (array) | No | |
| `is_ca` | Whether to generate a CA certificate | No | false |
| `extended_key_usage` | Extended key usage extensions (array) | No | |
| `duration` | Certificate duration in days | No | 365 |
| `key_length` | RSA key length | No | 2048 |

#### Password

Generate random passwords with specified length.

```yaml
variables:
  - name: db-password
    type: password
    options:
      length: 32
  - name: api-key
    type: password
    options:
      length: 64
```

**Options:**

| Option | Description | Required | Default |
| ------ | ----------- | -------- | ------- |
| `length` | Password length | No | 20 |

#### RSA

Generate RSA key pairs.

```yaml
variables:
  - name: encryption-key
    type: rsa
    options:
      key_length: 4096
```

**Options:**

| Option | Description | Required | Default |
| ------ | ----------- | -------- | ------- |
| `key_length` | RSA key length (bits) | No | 2048 |

#### SSH

Generate SSH key pairs.

```yaml
variables:
  - name: deploy-key
    type: ssh
    options:
      comment: "deploy@example.com"
```

**Options:**

| Option | Description | Required | Default |
| ------ | ----------- | -------- | ------- |
| `comment` | Comment for the SSH key | No | |

### Listing Variables

To view all variables defined in your configuration:

```bash
# Table format (default)
wrangle list variables

# JSON format
wrangle list variables --output json

# YAML format
wrangle list variables --output yaml
```

### Complete Configuration Example

```yaml
apiVersion: wrangle/v1
kind: Config
spec:
  variables:
    - name: root-ca
      type: certificate
      options:
        common_name: "Root CA"
        is_ca: true
    - name: server-cert
      type: certificate
      options:
        ca: root-ca
        common_name: "server.example.com"
        extended_key_usage: ["server_auth"]
    - name: db-password
      type: password
      options:
        length: 32
    - name: deploy-key
      type: ssh
      options:
        comment: "deploy@prod"
  stores:
    - name: vault
      type: vault
      properties:
        address: http://127.0.0.1:8200
        token: ((vault-token))
  packages:
    - name: yq
      version: 4.31.1
```

## Package Management

Wrangle is a simple package manager much like [arkade](https://github.com/alexellis/arkade). Arkade tends to focus on the lastest package version while Wrangle is fully version aware. Arkade also embeds its package feed into the execuable while Wrangle utilizes external package feeds.

Wrangle can find packages in one or more feeds. The default feed is https://github.com/patrickhuber/wrangle-packages.

When wrangle bootstrap is called, wrangle will go to default feed (or the overrided feed) and install itself 'wrangle@latest'.

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
