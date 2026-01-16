# sak

Swiss Army Knife - A personal CLI toolset.

## Installation

```bash
go install github.com/hezhizhen/sak/cmd/sak@latest
```

## Commands

### version

Show version information.

```bash
sak version
sak version --json
sak version --short
```

### worktime

Analyze work time data from `worktime.csv` in current directory.

```bash
sak worktime                # Show current period statistics
sak worktime -c             # Include comparison with previous periods
```

### compare

Compare files or directories between current directory and home directory using VS Code diff.

```bash
sak compare .bashrc         # Compare ./.bashrc with ~/.bashrc
sak compare config/         # Compare all files in ./config/ with ~/config/
```

### brew

Search Homebrew packages and display their info.

```bash
sak brew git                # Search packages matching "git"
```

### ccusage

Run ccusage tools to check AI coding assistant usage.

```bash
sak ccusage claude          # Run npx ccusage@latest
sak ccusage amp             # Run npx @ccusage/amp@latest
sak ccusage opencode        # Run npx @ccusage/opencode@latest
sak ccusage codex           # Run npx @ccusage/codex@latest
```

## Global Flags

- `--verbose` - Enable debug output
