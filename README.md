# Auto-pull

Janky little go program that iterates over a list of directories and checks for updates.

Uses [github.com/caseymrm/menuet](https://github.com/caseymrm/menuet) for a simple taskbar UI.

Config is in `~/.config/auto-pull/config.yaml` and looks like this:

```yaml
directories:
  - /path/to/directory1
  - ~/path/to/directory2
refreshSeconds: 30
```

Build and run with `cd cmd && make`
