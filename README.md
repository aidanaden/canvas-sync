<!-- omit in toc -->
# canvas-sync

CLI to download (course files, videos, etc) & view (deadlines, events, announcements) from [Canvas](https://www.instructure.com/canvas)

<!-- omit in toc -->
## Contents

- [Install](#install)
  - [Brew (mac/linux/wsl)](#brew-maclinuxwsl)
  - [Scoop (windows)](#scoop-windows)
  - [Autocomplete (mac/linux/wsl)](#autocomplete-maclinuxwsl)
    - [zsh](#zsh)
    - [bash](#bash)
    - [fish (not necessary if you installed fish via homebrew)](#fish-not-necessary-if-you-installed-fish-via-homebrew)
  - [Updating](#updating)
    - [Brew](#brew)
    - [Scoop](#scoop)
- [Set-up](#set-up)
- [Config](#config)
- [Commands](#commands)
  - [Init](#init)
  - [Pull](#pull)
    - [Pull Files](#pull-files)
    - [Pull Videos](#pull-videos)
  - [Update](#update)
    - [Update Files](#update-files)
    - [Update Videos](#update-videos)
  - [View](#view)
    - [View Deadlines (assignments)](#view-deadlines-assignments)
    - [View Events (Announcements/lectures/tutorials)](#view-events-announcementslecturestutorials)
    - [View People (from a given course)](#view-people-from-a-given-course)
- [FAQ](#faq)
- [LICENSE](#license)

## Install

### Brew (mac/linux/wsl)

Brew is a package manager for macOS (or linux) that helps you install packages easily - more info [here](https://brew.sh/)

To install Brew:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

To install canvas-sync with Brew:

```bash
brew install aidanaden/tools/canvas-sync
```

### Scoop (windows)

Scoop is a package manager for windows that helps you install programs from the command line (Brew but for windows) - more info [here](https://scoop.sh/)

To install Scoop, launch powershell and run:

```bash
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
```

To install canvas-sync with Scoop:

```bash
scoop bucket add scoop-bucket https://github.com/aidanaden/scoop-bucket.git
scoop install canvas-sync
```

You can also download directly from the [releases](https://github.com/aidanaden/canvas-sync/releases) page

### Autocomplete (mac/linux/wsl)

**Warning: skip if you don't know what zsh is**

Run the following code block to add autocompletion for canvas-sync to your shell.

#### zsh

<details>
  <summary>
    Code
  </summary>

  ```bash
  echo "\n\nif type brew &>/dev/null
  then
    FPATH="$(brew --prefix)/share/zsh/site-functions:${FPATH}"

    autoload -Uz compinit
    compinit
  fi" >> ~/.zshrc && source ~/.zshrc

  ```

</details>

#### bash

<details>
  <summary>
    Code
  </summary>

  ```bash
  echo "if type brew &>/dev/null
  then
    HOMEBREW_PREFIX="$(brew --prefix)"
    if [[ -r "${HOMEBREW_PREFIX}/etc/profile.d/bash_completion.sh" ]]
    then
      source "${HOMEBREW_PREFIX}/etc/profile.d/bash_completion.sh"
    else
      for COMPLETION in "${HOMEBREW_PREFIX}/etc/bash_completion.d/"*
      do
        [[ -r "${COMPLETION}" ]] && source "${COMPLETION}"
      done
    fi
  fi" >> ~/.bash_profile && source ~/.bash_profile
  ```

</details>

#### fish (not necessary if you installed fish via homebrew)

<details>
  <summary>
    Code
  </summary>

  ```bash
  echo "if test -d (brew --prefix)"/share/fish/completions"
      set -gx fish_complete_path $fish_complete_path (brew --prefix)/share/fish/completions
  end

  if test -d (brew --prefix)"/share/fish/vendor_completions.d"
      set -gx fish_complete_path $fish_complete_path (brew --prefix)/share/fish/vendor_completions.d
  end" >> ~/.config/fish/config.fish && source ~/.config/fish/config.fish
  ```

</details>

### Updating

#### Brew

If installed using brew, simply run:

```bash
canvas-sync upgrade
```

#### Scoop

If installed using scoop, simply run:

```bash
scoop update; scoop update canvas-sync
```

You can also download the latest version directly from the [releases](https://github.com/aidanaden/canvas-sync/releases) page

## Set-up

To set up canvas-sync, run the `init` command:

```bash
canvas-sync init
```

![init completed](examples/init/complete.png)

1. First, enter the directory to store all downloaded canvas data (files, videos, etc) e.g. `$HOME/Desktop/canvas`, if left blank all downloaded data will be stored in `$HOME/canvas-sync/data`

2. Next, enter your school's canvas website url, if left blank it'll be set to `https://canvas.nus.edu.sg` (i'm from nus after all)

3. You'll be asked for your username and password to log in to canvas **(required for video downloads)**

4. After logging in, your config will be successfully created and all other commands will work, check them out [here](#commands)

## Config

All configuration is done in the `$HOME/canvas-sync/config.yaml` file.

4 values can be configured:

- data_dir: directory to store downloaded canvas data, defaults to `$HOME/canvas-sync/data`  
- canvas_url: URL of your target canvas site, defaults to `https://canvas.nus.edu.sg`
- canvas_username: your canvas site username
- canvas_password: your canvas site password
- access_token **(DO NOT EDIT)**: token generated by the `init` command to download from canvas directly, if not filled you'll need to run `canvas-sync init`

To create a config file, run `canvas-sync init`

## Commands

### Init

Creates a new config file in the default directory `$HOME/canvas-sync`

![init](examples/init/init.gif)

### Pull

Downloads data (files, videos, etc) from canvas, overwrites all existing data

#### Pull Files

![pull files demo](examples/pull_files/run.gif)

View documentation via `pull files -h`

#### Pull Videos

![pull videos demo](examples/pull_videos/run.gif)

View documentation via `pull videos -h`

### Update

Updates downloaded data (files, videos, etc) from canvas

#### Update Files

![update files demo](examples/update_files/run.gif)

View documentation via `update files -h`

#### Update Videos

![update videos demo](examples/update_videos/run.gif)

View documentation via `update videos -h`

### View

Display data from canvas (deadlines, events, announcements, etc)

#### View Deadlines (assignments)

Display past/future assignment deadlines

![view deadlines demo](examples/view_deadlines/run.gif)

#### View Events (Announcements/lectures/tutorials)

Display past/future lectures/announcements

![view events demo](examples/view_events/run.gif)

#### View People (from a given course)

Display people from a given course code

![view people demo](examples/view_people/run.gif)

## FAQ

<details>
  <summary>
    What is this $HOME thing?
  </summary>
  It's the main directory of the current user. For mac/linux it'll be `/Users/<your user name>`, for windows it'll be `C:\Users\<your username>`
</details>

<details>
  <summary>
    It doesn't work for my university
  </summary>
  This tool was built by an NUS student (me), hence i'm not able to test it with canvas sites from any other university. Please create an issue at https://github.com/aidanaden/canvas-sync/issues regarding any problems you face with your school's canvas website. Thank you! :)
</details>

<details>
  <summary>
    Is my username/password stored anywhere?
  </summary>
  No, we do not store any credentials. When running the `init` command, the tool logs in to your canvas website and creates an access token that allows the tool to access your canvas data on your behalf.
</details>

If there are any other questions, please create an issue [here](https://github.com/aidanaden/canvas-sync/issues), if it's a common enough issue i'll add it to the FAQ section here :)

## LICENSE

MIT
