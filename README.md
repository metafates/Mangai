<h1 align="center">Mangal</h1>
<p align="center">
    <img width="200" src="assets/logo.png" alt="logo">
</p>

<h3 align="center">Manga Browser & Downloader</h3>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/metafates/mangal">
    <img src="https://goreportcard.com/badge/github.com/metafates/mangal">
  </a>

  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg">
  </a>
</p>

https://user-images.githubusercontent.com/62389790/174501320-119474c3-c745-4f95-8e7d-fbf5bd40920b.mov

## Table of Contents

- [About](#about)
- [Examples](#examples)
- [Config](#config)
- [Commands](#commands)
- [Install](#install)
- [Update](#update)
- [Uninstall](#uninstall)
- [Build](#build)
- [Limitations](#limitations)

## About

✨ __Mangal__ is feature rich, configurable manga browser & downloader
written in Go with support for different formats

⚙️ One of the most important features of Mangal is that it supports user defined scrapers
that can be added with just a few lines of config file (see [config](#config) & [limitations](#limitations))

🦎 Works in both modes - TUI & Inline. Use it as a standalone app or integrate with scripts

🍿 This app is inspired by __awesome__ [ani-cli](https://github.com/pystardust/ani-cli). Check it out!

Currently, Mangal supports these formats
- PDF
- Epub
- CBZ
- Zip
- Plain (just images)

> Type `mangal formats` for more info

## Examples

### TUI usage example

https://user-images.githubusercontent.com/62389790/174574562-011f9c30-db6f-45a9-9ce2-03973564ace0.mov

### Inline mode usage example

> For more information about inline mode type `mangal inline --help`

```bash
# Search manga. Returns a list of found manga
mangal inline --query "death note"

# Search manga. Returns a JSON list of found manga
mangal inline --query "death note" --json

# Get chapters of the first manga in the list
mangal inline --query "death note" --manga 1

# Download first chapter of the first manga in the list
mangal inline --query "death note" --manga 1 --chapter 1
```

## Config

> TLDR: Use `mangal config where` to show where config should be located
> and `mangal config init` to create default config


<details>
<summary>
Config is located at the OS default config directory.
</summary>

- __Unix__ - `$XDG_CONFIG_HOME/mangal/config.toml` if `$XDG_CONFIG_HOME` exists, else `$HOME/.config/mangal/config.toml`
- __Darwin__ (macOS) - `$HOME/Library/Application\ Support/mangal/config.toml`
- __Windows__ - `%AppData%\mangal\config.toml`
</details>


You can load config from custom path by using `--config` flag

`mangal --config /user/configs/config.toml`

By default, Mangal uses [manganelo](https://ww5.manganelo.tv) as a source

<details>
<summary>Click here to show config example</summary>

```toml
# Which sources to use. You can use several sources, it won't affect perfomance'
use = ['manganelo']

# Available options: pdf, epub, cbz, zip, plain (just images)
format = "pdf"

# If false, then OS default pdf reader will be used
use_custom_reader = false
custom_reader = "zathura"

# Custom download path, can be either relative (to the current directory) or absolute
download_path = '.'

# Add images to cache
# If set to true mangal could crash when trying to redownload something quickly
# Usually happens on slow machines
cache_images = false

[ui]
# Fullscreen mode
fullscreen = true

# Input prompt icon
prompt = ">"

# Input placeholder
placeholder = "What shall we look for?"

# Selected chapter mark
mark = "▼"

# Search window title
title = "Mangal"

[sources]
[sources.manganelo]
# Base url
base = 'https://ww5.manganelo.tv'

# Search endpoint. Put %s where the query should be
search = 'https://ww5.manganelo.tv/search/%s'

# Selector of entry anchor (<a></a>) on search page
manga_anchor = '.search-story-item a.item-title'

# Selector of entry title on search page
manga_title = '.search-story-item a.item-title'

# Manga chapters anchors selector
chapter_anchor = 'li.a-h a.chapter-name'

# Manga chapters titles selector
chapter_title = 'li.a-h a.chapter-name'

# Reader page images selector
reader_page = '.container-chapter-reader img'

# Random delay between requests
random_delay_ms = 500 # ms

# Are chapters listed in reversed order on that source?
# reversed order -> from newest chapter to oldest
reversed_chapters_order = true
```
</details>

## Commands

```
The ultimate manga downloader multitool

Usage:
  mangal [flags]
  mangal [command]

Available Commands:
  cleanup     Remove cached and temp files
  completion  Generate the autocompletion script for the specified shell
  config      Config actions
  formats     Information about available formats
  help        Help about any command
  inline      Search & Download manga in inline mode
  version     Show version

Flags:
  -c, --config string   use config from path
  -f, --format string   use custom format
  -h, --help            help for mangal

Use "mangal [command] --help" for more information about a command.
```

## Install

- [Go](#go)
- [From source](#from-source)
- [MacOS](#macos)
- [Windows](#windows)
- [Linux](#linux)


### Go

You will need [Go installed](https://go.dev/doc/install)

```bash
go install -ldflags="-s -w" github.com/metafates/mangal@latest
```

#### Update

```bash
go install github.com/metafates/mangal@latest
```

#### Uninstall

- Bash / zsh - `rm $(which mangal)`
- Fish - `rm (which mangal)`
- Powershell - `rm $(which.exe mangal)`

### From source
```bash
git clone https://github.com/metafates/mangal.git
cd mangal
go install -ldflags="-s -w"
```

### MacOS

Install using [Homebrew](https://brew.sh/)

```bash
brew tap metafates/mangal
brew install mangal
```

<details>
<summary>Update & Uninstall</summary>

#### Update

```bash
brew upgrade mangal
```

#### Uninstall

```bash
brew uninstall mangal
```

</details>

### Windows

```powershell
scoop install https://raw.githubusercontent.com/metafates/scoop-mangal/main/mangal.json
```

<details>
<summary>Update & Uninstall</summary>

#### Update

```powershell
scoop update mangal
```

#### Uninstall

```powershell
scoop uninstall mangal
```

### Linux

> Soon... 😴

## Build

```bash
git clone https://github.com/metafates/mangal.git
cd mangal
go build
```


> You can also cross build for windows, linux & macos
> by running `cross-compile.py` (you will need Python 3)
> 
> Built binaries and generated packages
> will be stored in the `bin` folder

## Limitations

Even though many manga sites will work,
there exist some (serious) limitations to which sites could be added

- Navigation layout should follow this model
    - Each manga have a separate page
    - Manga page should have some form of chapters list (not lazy loaded)
    - Each chapter should have a separate reader page with all the images


Some sites that work well

- https://manganato.com
- https://ww3.mangakakalot.tv
- https://ww5.manganelo.tv


I'm planning to make a more advanced scraper creation system
to overcome this roadblocks somewhere in the future

---

Logo taken from [here](https://www.flaticon.com/free-icon/parchment_1391306)
