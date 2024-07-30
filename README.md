<h1 align="center">
<strong>Mangal 4 ☄️</strong>
</h1>

<p align="center">
    <img alt="Linux" src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
    <img alt="macOS" src="https://img.shields.io/badge/mac%20os-000000?style=for-the-badge&logo=macos&logoColor=F0F0F0">
    <img alt="Windows" src="https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white">
    <img alt="Termux" src="https://img.shields.io/badge/Termux-000000?style=for-the-badge&logo=GNOME%20Terminal&logoColor=white">
</p>

<h3 align="center">
    The most advanced CLI manga downloader in the entire universe!
</h3>

<p align="center">
    <img alt="Mangal 4 TUI" src="assets/tui.gif">
</p>


## Try it!

```shell
curl -sSL mangal.metafates.one/run | sh
```

> **Note** This script does not install anything, it just downloads, verifies and runs Mangal.
> Not available on Windows.

## Table of contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Custom scrapers](#custom-scrapers)
- [Anilist](#anilist)
- [Honorable mentions](#honorable-mentions)

## Features

- __Lua Scrapers!!!__ You can add any source you want by creating your own _(or using someone's else)_ scraper with
  __Lua 5.1__. See [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)
- __4 Built-in sources__ - [Mangadex](https://mangadex.org), [Manganelo](https://m.manganelo.com/wwww), [Manganato](https://manganato.com) & [Mangapill](https://mangapill.com)
- __Download & Read Manga__ - I mean, it would be strange if you couldn't, right?
- __Caching__ - Mangal will cache as much data as possible, so you don't have to wait for it to download the same data over and over again. 
- __4 Different export formats__ - PDF, CBZ, ZIP and plain images
- __TUI ✨__ - You already know how to use it! (ﾉ>ω<)ﾉ :｡･::･ﾟ’★,｡･:･ﾟ’☆
- __Scriptable__ - You can use Mangal in your scripts, it's just a CLI app after all. [Examples](https://github.com/metafates/mangal/wiki/Inline-mode)
- __History__ - Resume your reading from where you left off!
- __Fast?__ - YES.
- __Monolith__ - ZERO runtime dependencies. Even Lua is built in. Easy to install and use.
- __Cross-Platform__ - Linux, macOS, Windows, Termux, even your toaster. (¬‿¬ )
- __Anilist integration__ - Mangal will collect additional data from Anilist and use it to improve your reading experience. It can also sync your progress!

## Installation

### Script (Linux, MacOS, Termux)

Install using [this shell script](https://github.com/metafates/mangal/blob/main/scripts/install)

```shell
curl -sSL mangal.metafates.one/install | sh
```

This script will automatically detect OS & Distro and use the best option available.
For example, on macOS it will try to use Homebrew, on Ubuntu it will install the `.deb` package and so on...

### Arch Linux

[AUR package](https://aur.archlinux.org/packages/mangal-bin) (maintained by [@balajsra](https://github.com/balajsra),
thank you)

### MacOS

Install using [Homebrew](https://brew.sh/)

    brew tap metafates/mangal
    brew install mangal

### Windows

Install using [Scoop](https://scoop.sh/) (thanks to [@SonaliBendre](https://github.com/SonaliBendre) for adding it to
the official bucket)

    scoop bucket add extras
    scoop install mangal

### Termux

Thanks to [@T-Dynamos](https://github.com/T-Dynamos) for adding it to the [termux-packages](https://github.com/termux/termux-packages)

```shell
pkg install mangal
```

### Gentoo

Install using third-party overlay [raiagent](https://github.com/leycec/raiagent). Thanks to [@leycec](https://github.com/leycec) for maintaining it.

```shell
eselect repository enable raiagent
emerge --sync raiagent
emerge mangal
```

### Nix 

Install using [Nix](https://nixos.org/download.html#download-nix). Thanks to [@bertof](https://github.com/bertof) for adding it to the [nixpkgs](https://github.com/NixOS/nixpkgs)

```shell
# NixOS
nix-env -iA nixos.mangal

# Non NixOS
nix-env -iA nixpkgs.mangal
```

### Docker

Install using Docker. (thanks to [@ArabCoders](https://github.com/ArabCoders) for reference)

    docker pull metafates/mangal

To run

```shell
docker run --rm -ti -e "TERM=xterm-256color" -v $(PWD)/mangal/downloads:/downloads -v $(PWD)/mangal/config:/config metafates/mangal
```

### From source

Visit this link to install [Go](https://go.dev/doc/install).

Clone the repo
```shell
git clone --depth 1 https://github.com/metafates/mangal.git
cd mangal
```

GNU Make **(Recommended)**
```shell
make install # if you want to compile and install mangal to path
make build # if you want to just build the binary
```

<details>
<summary>If you don't have GNU Make use this</summary>


```shell
# To build
go build -ldflags "-X 'github.com/metafates/mangal/constant.BuiltAt=$(date -u)' -X 'github.com/metafates/mangal/constant.BuiltBy=$(whoami)' -X 'github.com/metafates/mangal/constant.Revision=$(git rev-parse --short HEAD)' -s -w"

# To install
go install -ldflags "-X 'github.com/metafates/mangal/constant.BuiltAt=$(date -u)' -X 'github.com/metafates/mangal/constant.BuiltBy=$(whoami)' -X 'github.com/metafates/mangal/constant.Revision=$(git rev-parse --short HEAD)' -s -w"
```

</details>

If you want to build mangal for other architecture, say ARM, you'll have to set env variables `GOOS` and `GOARCH`

```shell
GOOS=linux GOARCH=arm64 make build
```

[Available GOOS and GOARCH combinations](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)

### Pre-compiled

Download the pre-compiled binaries from the [releases page](https://github.com/metafates/mangal/releases/latest)
and copy them to the desired location.

## Usage

### TUI

Just run `mangal` and you're ready to go.

<details>
<summary>Keybinds</summary>

| Bind                                                        | Description                          |
|-------------------------------------------------------------|--------------------------------------|
| <kbd>?</kbd>                                                | Show help                            |
| <kbd>↑/j</kbd> <kbd>↓/k</kbd> <kbd>→/l</kbd> <kbd>←/h</kbd> | Navigate                             |
| <kbd>g</kbd>                                                | Go to first                          |
| <kbd>G</kbd>                                                | Go to last                           |
| <kbd>/</kbd>                                                | Filter                               |
| <kbd>esc</kbd>                                              | Back                                 |
| <kbd>space</kbd>                                            | Select one                           |
| <kbd>tab</kbd>                                              | Select all                           |
| <kbd>v</kbd>                                                | Select volume                        |
| <kbd>backspace</kbd>                                        | Unselect all                         |
| <kbd>enter</kbd>                                            | Confirm                              |
| <kbd>o</kbd>                                                | Open URL                             |
| <kbd>r</kbd>                                                | Read                                 |
| <kbd>q</kbd>                                                | Quit                                 |
| <kbd>ctrl+c</kbd>                                           | Force quit                           |
| <kbd>a</kbd>                                                | Select Anilist manga (chapters list) |
| <kbd>d</kbd>                                                | Delete single history entry          |

</details>

![TUI](https://user-images.githubusercontent.com/62389790/198830334-fd85c74f-cf3b-4e56-9262-5d62f7f829f4.png)

> If you wonder what those icons mean - `D` stands for "downloaded", `*` shows that chapter is marked to be downloaded.
> You can choose different icons, e.g. nerd font ones - just run mangal with `--icons nerd`.
> Available options are `nerd`, `emoji`, `kaomoji` and `squares`

### Mini

Mini mode tries to mimic [ani-cli](https://github.com/pystardust/ani-cli)

To run: `mangal mini`

![mini](https://user-images.githubusercontent.com/62389790/198830544-f2005ec4-c206-4fe0-bd08-862ffd08320e.png)

### Inline

Inline mode is intended for use with other scripts.

Type `mangal help inline` for more information.

See [Wiki](https://github.com/metafates/mangal/wiki/Inline-mode) for more examples.

<p align="center">
    <img alt="Mangal 4 Inline" src="assets/inline.gif">
</p>

### Other

See `mangal help` for more information

## Configuration

Mangal uses [TOML](https://toml.io) format for configuration under the `mangal.toml` filename.
Config path depends on the OS.
To find yours, use `mangal where --config`.
For example, on __Linux__ it would be `~/.config/mangal/mangal.toml`.

Use env variable `MANGAL_CONFIG_PATH` to set custom config path.
> See `mangal env` to show all available env variables.

| Command               | Description                                      |
|-----------------------|--------------------------------------------------|
| `mangal config get`   | Get config value for specific key                |
| `mangal config set`   | Set config value for specific key                |
| `mangal config reset` | Reset config value for specific key              |
| `mangal config info`  | List all config fields with description for each |
| `mangal config write` | Write current config to a file                   |

## Custom scrapers

TLDR; To browse and install a custom scraper
from [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers) run

    mangal sources install

Mangal has a Lua5.1 VM built-in + some useful libraries, such as headless chrome, http client, html parser and so on...

Check the [defined modules](https://github.com/metafates/mangal-lua-libs) for more information.

For scrapers examples, check the [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers)

### Creating a custom scraper

This command will create `example.lua` file in the `mangal where --sources` directory.

    mangal sources gen --name example --url https://example.com

Open the file and edit it as you wish.
Take a look at the comments for more information.
See [mangal-scrapers repository](https://github.com/metafates/mangal-scrapers) for examples.

You can test it by running `mangal run <filepath>`

It should automatically appear in the list of available scrapers.

> New to Lua? [Quick start guide](https://learnxinyminutes.com/docs/lua/)

## Anilist

Mangal also supports integration with anilist.

Besides fetching metadata for each manga when downloading,
mangal can also mark chapters as read on your Anilist profile when you read them inside mangal.

For more information see [wiki](https://github.com/metafates/mangal/wiki/Anilist-Integration)

## Honorable mentions

### Projects using mangal

- [kaizoku](https://github.com/oae/kaizoku) - Self-hosted manga downloader with mangal as its core 🚀

### Similar Projects

- [mangadesk](https://github.com/darylhjd/mangadesk) - Terminal client for MangaDex
- [ani-cli](https://github.com/pystardust/ani-cli) - A cli tool to browse and play anime
- [manga-py](https://github.com/manga-py/manga-py) - Universal manga downloader
- [animdl](https://github.com/justfoolingaround/animdl) - A highly efficient, fast, powerful and light-weight anime
  downloader and streamer
- [tachiyomi](https://github.com/tachiyomiorg/tachiyomi) - Free and open source manga reader for Android

### Libraries

- [bubbletea](https://github.com/charmbracelet/bubbletea), [bubbles](https://github.com/charmbracelet/bubbles)
  & [lipgloss](https://github.com/charmbracelet/lipgloss) - Made mangal shine! The best TUI libraries ever ✨
- [gopher-lua](https://github.com/yuin/gopher-lua) - Made it possible to write custom scrapers with Lua ❤️
- [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) - Responsible for the awesome CLI
  & config experience 🛠
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - Fast pdf processor in pure go 📄
- _And many others!_

### Contributors

And of course, thanks to all contributors! You are awesome!

<p align="center">
<a href="https://github.com/metafates/mangal/graphs/contributors">
  <img alt="Contributors" src="https://contrib.rocks/image?repo=metafates/mangal" />
</a>
</p>

---

<p align="center">
If you find this project useful or want to say thank you,
please consider starring it, that would mean a lot to me ⭐
</p>

<p align="center">
<a href="https://star-history.com/#metafates/mangal&Date">
<img alt="Star History" src="https://api.star-history.com/svg?repos=metafates/mangal&type=Date"/>
</a>
</p>
