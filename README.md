<!-- omit in toc -->
# canvas-sync

CLI to download (course files, videos, etc) & view (deadlines, events, announcements) from [Canvas](https://www.instructure.com/canvas)

<!-- omit in toc -->
## Contents

- [Install](#install)
  - [Brew](#brew)
  - [Autocomplete](#autocomplete)
    - [zsh](#zsh)
    - [bash](#bash)
    - [fish (not necessary if you installed fish via homebrew)](#fish-not-necessary-if-you-installed-fish-via-homebrew)
- [Config](#config)
- [Commands](#commands)
  - [Pull](#pull)
    - [Pull Files](#pull-files)
    - [Pull Videos](#pull-videos)
  - [Update](#update)
    - [Update Files](#update-files)
    - [Update Videos](#update-videos)
  - [View](#view)
    - [View Deadlines](#view-deadlines)
    - [View Events (lectures/tutorials)](#view-events-lecturestutorials)
    - [View Announcements](#view-announcements)
- [LICENSE](#license)

## Install

### Brew

```bash
brew install aidanaden/tools/canvas-sync
```

You can also download directly from the [releases](https://github.com/aidanaden/canvas-sync/releases) page

### Autocomplete

Add autocompletion for canvas-sync in your terminal:

#### zsh

<details>
  <summary>
    Code for zsh autocomplete
  </summary>

  ```bash
  # replace '~/.zshrc' with wherever your zsh config file is
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
    Code for bash autocomplete
  </summary>

  ```bash
  # replace '~/.bash_profile' with wherever your bash config file is
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
    Code for fish autocomplete
  </summary>

  ```bash
  # replace '~/.config/fish/config.fish' with wherever your fish config file is
  echo "if test -d (brew --prefix)"/share/fish/completions"
      set -gx fish_complete_path $fish_complete_path (brew --prefix)/share/fish/completions
  end

  if test -d (brew --prefix)"/share/fish/vendor_completions.d"
      set -gx fish_complete_path $fish_complete_path (brew --prefix)/share/fish/vendor_completions.d
  end" >> ~/.config/fish/config.fish && source ~/.config/fish/config.fish
  ```

</details>

## Config

All configuration is done in the `$HOME/.canvas-sync/config.yaml` file.

3 values can be configured:

- `access_token`: token generated by the canvas user (you) to download from canvas directly, if not filled canvas-sync will prompt you to log in
- `data_dir`: directory to store downloaded canvas data, defaults to `$HOME/.canvas-sync/data`  
- `canvas_url`: URL of your target canvas site, defaults to `https://canvas.nus.edu.sg`

## Commands

### Pull

Downloads data (files, videos, etc) from canvas, overwrites all existing data

#### Pull Files

View documentation via `pull files -h`:

![pull files help](examples/pull_files_help.gif)

![pull files demo](examples/pull_files_all.gif)

#### Pull Videos

WIP

### Update

Updates downloaded data (files, videos, etc) from canvas

#### Update Files

View documentation via `update files -h`:

![update files help](examples/update_files_help.gif)

![update files demo](examples/update_files_all.gif)

#### Update Videos

WIP

### View

Display data from canvas (deadlines, events, announcements, etc)

#### View Deadlines

WIP

#### View Events (lectures/tutorials)

WIP

#### View Announcements

WIP

## LICENSE

MIT
