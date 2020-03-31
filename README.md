
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
    <a href="#options">Options</a> ❘
    <a href="#how-it-works">How It Works</a> ❘
    <a href="#contributing">Contributing</a> ❘
    <a href="#credits-and-license">Credits & License</a>
</p>

![screenshot](https://github.com/unleashable/apker/raw/master/.github/images/demo.gif)


## Installation

[⚠ **WARNING**]: Apker under development and its core features are not completed yet. <br>
[⚠ **WARNING**]: Please do not use this in production until v1 stable, there may be breaking changes.

You can install Apker via go, or download pre-compiled versions.

#### Compiled:

Download latest version from [Releases](https://github.com/Unleashable/apker/releases).

Then move the binary file to your bin directory:

```bash
sudo mv ./apker /usr/bin/apker
sudo chmod +x /usr/bin/apker
```

### Via Go Mod:


```bash
go get https://github.com/Unleashable/apker
```
Apker will be installed automatically into your `$GOPATH/bin`

### Manual Install:

```bash
git clone https://github.com/Unleashable/apker /tmp/apker
cd /tmp/apker
make install
```
Node: this requires golang.

## Usage

WIP!

You can try the demo project:

Open your terminal and export your provider name and it's API key (at the moment Apker supports only DO):

```bash
export APKER_PROVIDER=digitalocen
export APKER_KEY=YOUR_DIGITALOCEAN_API_KEY_HERE
```

Then run deploy command like this:
```bash
apker deploy --url https://github.com/melbahja/apker-demo
```
If your private ssh key protected with passphrase just add `--passphrase` flag to the end of the command.


## Options

WIP!

## How It Works

WIP!

## Contributing

PRs, issues, and feedback from ninja gophers are very welcomed.

## Credits and License

#### Credits:
Models Used: [go.mod](https://github.com/Unleashable/apker/blob/master/go.mod) <br>
Icon By: <a href="https://thenounproject.com/smodgekar/">Sachin Modgekar</a>


#### License:

Apker is provided under the [MIT License](https://github.com/Unleashable/apker/blob/master/LICENSE)
