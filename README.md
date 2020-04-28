<div align="center">
    <a href="https://github.com/unleashable/apker">
        <img src="https://github.com/unleashable/apker/raw/master/.github/images/icon.png" width="200">
    </a>
    <h1>APKER</h1>
</div>

<h4 align="center">
    Deploy and manage your custom images to your cloud provider in seconds.
</h4>

<p align="center">
    <a href="#installation">Installation</a> ❘
    <a href="#usage">Usage</a> ❘
    <a href="#contributing">Contributing</a> ❘
    <a href="#license">License</a>
</p>

[![Apker](https://github.com/unleashable/apker/raw/master/.github/images/demo.png)](https://asciinema.org/a/TGa2tfGhFfVmtriuE51xwqFYM)


## Installation

You can install Apker via go, manually, or download pre-compiled versions.

#### Pre-Compiled:

```bash
# go to tmp dir.
cd /tmp

# download latest version.
curl -sfL https://git.io/apker | sh

# make the binary executable.
chmod +x /tmp/bin/apker

# move the binary to your PATH
sudo mv /tmp/bin/apker /usr/bin/apker
```

### Go:


```bash
go get https://github.com/Unleashable/apker
```
Apker will be installed automatically into your `$GOPATH/bin`

### Manually:

```bash
git clone https://github.com/Unleashable/apker /tmp/apker
cd /tmp/apker
make
sudo make install
```
Note: this requires golang.

## Usage

### Quick Start:

Try Apker Demo: [https://github.com/melbahja/apker-demo](https://github.com/melbahja/apker-demo/)

### Cli Tool:

Type `--help` to see all available options and flags.

```bash
# global help
apker --help

# apker deploy help
apker deploy --help

# apker run help
apker run --help
```

### Apker File:

To be able to deploy with Apker your project most have a `apker.yaml` file, In `apker.yaml` file you can specify the deploy steps to build your image, and actions.

A simple example of `apker.yaml`:

```yaml
version: v1
name: my-nginx-demo
image:
  size: small
  from: centos-8-x64

provider:
  name: {{Env "APKER_PROVIDER"}}
  credentials:
    API_KEY: {{Env "APKER_KEY"}}

deploy:

  env:
    MY_VAR: {{GetOr "myvar" "HELLO WORLD"}}

  setup:
    - dnf update -y
    - dnf install rsync -y
    - dnf install git -y

  steps:
    - run dnf install nginx -y
    - copy public/ /usr/share/nginx/html/
    - run echo $MY_VAR > /usr/share/nginx/html/hi.html
    - run chown nginx:nginx /usr/share/nginx/html/ -R
    - run systemctl enable nginx
    - reboot

actions:
  status: systemctl status nginx
  restart: systemctl restart nginx
  reboot: reboot &

events:
  success: echo "success event executed"
  failure: echo "failure event executed"
```


#### Apker File Properties:

| Name | Type | Descriptin | Required |
|------|:----:|------------|:--------:|
| `version`              | string | Just for apker.                                        | YES |
| `name`                 | string | Your image name                                        | YES |
| `image.size`           | string | Image size: `small, medium, large`                     | NO  |
| `image.from`           | string | Base distro name or remote `.qcow2` url [1]            | YES |
| `provider.name`        | string | The cloud provider name: `digitalocean`, `aws`[2]      | YES |
| `provider.credentials` | key: value | The cloud provider credentials like api keys.      | YES |
| `deploy.env`           | key: value | The deployment env vars                            | NO  |
| `deploy.setup`         | list of commands | Required `git` and `rsync` install commands. | YES |
| `deploy.steps`         | List of deploy steps | Deployment steps                         | YES |
| `actions`              | key: value | Actions to run later via `apker run`               | NO  |
| `events.success` | bash command | Command to run on **host** machine after successful deployment. | NO |
| `events.failure` | bash command | Command to run on **host** machine after deployment failure.      | NO |

[1]: You can use remote url to build from your own custom images like: `https://cloud.centos.org/centos/8/x86_64/images/CentOS-8-ec2-8.1.1911-20200113.3.x86_64.qcow2`

[2]: AWS not supported yet, you can deploy only to Digitalocean. But you can deploy to a Custom Provider.

#### Deployment Steps:
This is the: `deploy.steps` that you can use to build your image.

| Name | Description | Example |
|------|-------------|:--------|
| `run`    | Run a shell command.               | `run: apt-get -y install nginx` |
| `dir`    | Create a directory.                | `dir: /var/www/myapp/public` [1]  |
| `copy`   | Copy file or directory.             |  `copy: . /var/www/myapp`      |
| `reboot` | Reboot the machine.                | `reboot`                    |

[1]: Create a new directory equivalent to `mkdir -p`

[2]: Copy all repo content to `/var/www/myapp`

### Deployment:

Apker currently only supports digitalocean, to deploy your project you must export these env vars before running deploy command:
```bash
export APKER_PROVIDER=digitalocean
export APKER_KEY=your_do_api_key_here
```

then use the `deploy` subcommand:
```bash
apker deploy --url https://github.com/username/repo
```
`username/repo` must have a valid `apker.yaml` file, and your public key must be in digitalocean keys.

#### Deploy To A Custom Provider:
If you want to deploy a project to unsupported cloud provider for example aws, just create a new instance based on the project distro `name` in the `apker.yaml` file, add your public ssh key to it and run the following command:

```bash
apker deploy --url https://github.com/username/repo --ip 127.0.0.1
```

replace `127.0.0.1` with your instance public ip address.


WIP!

[**⚠ WARNING**]: Apker is under development and its core features are not completed yet, please do not use this in production until v1 stable, there may be BREAKING CHANGES.


## Contributing

PRs, issues, and feedback from ninja gophers are very welcomed.

## License

Apker is provided under the [MIT License](https://github.com/Unleashable/apker/blob/master/LICENSE).
