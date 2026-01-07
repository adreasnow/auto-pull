# Auto-pull

Janky little go program that iterates over a list of directories and checks for updates.

Uses [github.com/caseymrm/menuet](https://github.com/caseymrm/menuet) for a simple taskbar UI.

## Building

- `make` will build and run the app
- `make build` will build the app without running it

## Running

Config is in `~/.config/auto-pull/config.yaml` and looks like this:

- Note that wildcards are only supported as the last entry
- Only github repos are supported and only via the `https` protocol
- Upon launching, the program will query for a token and save it to the keychain. The token must have contents:read permissions for your repos.

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

Build and run with `cd cmd && make`
