# Auto-pull

A macOS menu bar app that monitors a list of git directories and automatically pulls updates.

Uses [github.com/caseymrm/menuet](https://github.com/caseymrm/menuet) for the menu bar UI.

## Install

### Homebrew

```sh
brew install --cask adreasnow/tap/auto-pull
```

### Manual

Download the latest `.app` zip for your architecture from the [releases page](https://github.com/adreasnow/auto-pull/releases), unzip it, and move `Auto Pull.app` to `/Applications`.

> **Note:** The app is not codesigned. On first launch macOS may block it — right-click the app and choose **Open**, or run:
>
> ```sh
> xattr -dr com.apple.quarantine /Applications/Auto\ Pull.app
> ```

## Configuration

Config lives at `~/.config/auto-pull/config.yaml`:

```yaml
directories:
  - /path/to/directory1
  - ~/path/to/directory2
  - ~/path/to/directory3/*
refreshSeconds: 30 # defaults to 60
notifications: # defaults to true
  failed: true
  fetchedNoPull: true
  pulled: true
```

- Wildcards (`*`) are only supported as the last path component
- Only GitHub repos over `https` are supported
- On first launch the app will prompt for a GitHub token and store it in the keychain — the token needs `contents:read` permission for your repos

Logs are written to `~/Library/Logs/com.github.adreasnow.auto-pull/`.

## Building

### With GoReleaser (recommended)

```sh
goreleaser build --snapshot --clean
```

The resulting `.app` bundle will be in `dist/`.

### With Make

```sh
make build   # build the .app bundle
make         # build and run
```
