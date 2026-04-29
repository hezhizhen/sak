# Release Guide

This project uses [GoReleaser](https://goreleaser.com/) with GitHub Actions for automated releases.

## Quick Start

To create a new release:

```bash
# Create and push a tag
git tag v0.1.0
git push origin v0.1.0
```

GitHub Actions will automatically:
- Build binaries for multiple platforms (Linux/macOS/Windows, amd64/arm64)
- Generate checksums
- Create a GitHub Release with changelog
- Upload all artifacts

## Testing Locally

Before creating a real release, test the build locally:

```bash
# Install goreleaser (if not already installed)
brew install goreleaser

# Test release build (no publishing)
make release-snapshot

# Check the generated files
ls -la dist/
```

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- `v1.0.0` - Major release (breaking changes)
- `v0.1.0` - Minor release (new features)
- `v0.0.1` - Patch release (bug fixes)

## Release Artifacts

Each release includes:
- Binaries for: `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64`, `windows_amd64`, `windows_arm64`
- Compressed archives (`.tar.gz` for Unix, `.zip` for Windows)
- `checksums.txt` for verification
- Auto-generated changelog

## Manual Release (if needed)

If you need to create a release manually:

```bash
# Using gh CLI
gh release create v0.1.0 --generate-notes

# Or via GitHub web interface
# Visit: https://github.com/hezhizhen/sak/releases/new
```
