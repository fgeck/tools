# Release Process

This document describes how to create a new release of the `tools` CLI and publish it to Homebrew.

## Prerequisites

1. **Create the Homebrew Tap Repository**
   - Go to https://github.com/new
   - Create a **public** repository named: `homebrew-tools`
   - Leave it empty (GoReleaser will populate it automatically)

2. **No GitHub Token Needed!**
   - The GitHub Actions workflow uses `secrets.GITHUB_TOKEN` which is automatically provided
   - The token has permissions to:
     - Create releases in the `tools` repository
     - Push to the `homebrew-tools` repository

## Release Steps

### 1. Commit All Changes

```bash
git add .
git commit -m "Prepare for release v0.1.0"
git push origin main
```

### 2. Create and Push a Git Tag

```bash
# Create an annotated tag
git tag -a v0.1.0 -m "Release version 0.1.0"

# Push the tag to trigger the release workflow
git push origin v0.1.0
```

### 3. What Happens Automatically

The GitHub Actions pipeline will:

1. **Run Tests**
   - Unit tests with coverage
   - Integration tests with coverage

2. **GoReleaser Builds**
   - Compile binaries for:
     - Linux (amd64, arm64)
     - macOS/Darwin (amd64, arm64)
     - Windows (amd64, arm64)
   - Generate checksums
   - Create changelog from git commits

3. **Create GitHub Release**
   - Upload all binaries as release assets
   - Include checksums and changelog

4. **Update Homebrew Formula**
   - Automatically create/update `Formula/tools.rb` in `homebrew-tools` repo
   - Calculate SHA256 checksums
   - Include installation instructions

5. **Build and Push Docker Image**
   - Tag with version number
   - Tag as `latest`
   - Push to GitHub Container Registry

### 4. Verify the Release

1. Check GitHub Releases:
   ```
   https://github.com/fgeck/tools/releases
   ```

2. Check Homebrew Formula:
   ```
   https://github.com/fgeck/homebrew-tools/blob/main/Formula/tools.rb
   ```

3. Test Installation:
   ```bash
   brew tap fgeck/tools
   brew install tools
   tools --help
   ```

## Versioning

Follow Semantic Versioning (https://semver.org/):

- **MAJOR** version (v1.0.0) - Incompatible API changes
- **MINOR** version (v0.1.0) - New functionality, backwards compatible
- **PATCH** version (v0.0.1) - Bug fixes, backwards compatible

## Testing a Release Locally (Before Pushing Tag)

```bash
# Test GoReleaser locally without publishing
goreleaser release --snapshot --clean --skip=publish

# Check generated files
ls -la dist/

# Test a binary
./dist/tools_darwin_arm64/tools --help
```

## Rolling Back a Release

If you need to delete a bad release:

```bash
# Delete the remote tag
git push --delete origin v0.1.0

# Delete the local tag
git tag -d v0.1.0

# Delete the GitHub release manually at:
# https://github.com/fgeck/tools/releases
```

## Example Release Commands

```bash
# Patch release (bug fixes)
git tag -a v0.1.1 -m "Fix: resolve TUI input issue"
git push origin v0.1.1

# Minor release (new features)
git tag -a v0.2.0 -m "Add: export/import functionality"
git push origin v0.2.0

# Major release (breaking changes)
git tag -a v1.0.0 -m "Release: stable v1.0.0"
git push origin v1.0.0
```

## Troubleshooting

### Release Workflow Not Triggering

Check that:
- The tag follows the pattern `v*.*.*` (e.g., v0.1.0)
- The tag is pushed to GitHub: `git push origin v0.1.0`
- All tests pass before the release job runs

### Homebrew Formula Not Created

Check that:
- The `homebrew-tools` repository exists and is public
- The GitHub Actions bot has permissions to push
- Check the GoReleaser logs in Actions tab

### Users Can't Install via Brew

Users might need to:
```bash
# Update Homebrew
brew update

# Untap and re-tap
brew untap fgeck/tools
brew tap fgeck/tools

# Clear cache
rm -rf $(brew --cache)

# Install
brew install tools
```

## Post-Release Checklist

- [ ] Verify GitHub release created with assets
- [ ] Verify Homebrew formula in `homebrew-tools` repo
- [ ] Test Homebrew installation on clean system
- [ ] Update README.md with new version number
- [ ] Announce release (if applicable)

## Next Steps

After v1.0.0 is stable, consider:
- Submitting to Homebrew Core (official tap)
- Publishing to other package managers (apt, yum, snap, etc.)
- Creating installation scripts for unsupported platforms
