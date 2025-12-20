# Auto-pull

Janky little go program that iterates over a list of directories and checks for updates.

Uses [github.com/caseymrm/menuet](https://github.com/caseymrm/menuet) for a simple taskbar UI.

## Building



## Running

Config is in `~/.config/auto-pull/config.yaml` and looks like this:

- Note that wildcards are only supported as the last entry
- Only github repos are supported and only via the `https` protocol

```yaml
directories:
  - /path/to/directory1
  - ~/path/to/directory2
  - ~/path/to/directory3/*
refreshSeconds: 30
```

Build and run with `cd cmd && make`
