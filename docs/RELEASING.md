# Qredex Go SDK Release Guide

This document describes the release process for the Qredex Go SDK.

## Versioning

The SDK follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html):

- **MAJOR.MINOR.PATCH** (e.g., `1.2.3`)
- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

### Pre-release Versions

Pre-release versions use hyphen suffix:

- `0.1.0-alpha.1` — First alpha release
- `0.1.0-beta.1` — First beta release
- `0.1.0-rc.1` — First release candidate

## Release Checklist

### 1. Prepare the Release

- [ ] Review all commits since last release: `git log --oneline v0.1.0..HEAD`
- [ ] Ensure all changes are documented in `CHANGELOG.md`
- [ ] Verify all tests pass: `go test ./...`
- [ ] Run `go vet ./...` to check for issues
- [ ] Update version in `config.go` (`SDKVersion` constant)
- [ ] Update `go.mod` if dependencies changed

### 2. Update CHANGELOG.md

Edit `CHANGELOG.md` to:

1. Move `[Unreleased]` section to new version section
2. Add release date: `## [0.2.0] - 2026-01-20`
3. Add new `[Unreleased]` section at top
4. Ensure all changes are categorized (Added, Changed, Fixed, etc.)

Example:

```markdown
# Changelog

## [Unreleased]

## [0.2.0] - 2026-01-20

### Added
- New feature X
- New resource Y

### Fixed
- Bug fix for issue Z
```

### 3. Verify Artifacts

Run these commands to verify the SDK is release-ready:

```bash
# Build the SDK
go build ./...

# Run all tests
go test ./...

# Run vet
go vet ./...

# Check formatting
go fmt ./...

# Verify module
go mod verify
```

### 4. Commit Changes

Create a release commit:

```bash
git add -A
git commit -m "chore: prepare release v0.2.0"
```

### 5. Create Git Tag

Create an annotated tag:

```bash
git tag -a v0.2.0 -m "Release v0.2.0"
```

### 6. Push to Remote

Push commits and tags:

```bash
git push origin main
git push origin v0.2.0
```

### 7. Create GitHub Release

1. Go to https://github.com/qredex/sdk-go/releases
2. Click **Create a new release**
3. Select the tag you just pushed (`v0.2.0`)
4. Copy the CHANGELOG.md section for this version as the release description
5. Click **Publish release**

### 8. Notify Packagist (if applicable)

If the SDK is published on Packagist:

1. Ensure the GitHub webhook is configured in Packagist
2. The release should be automatically picked up
3. Verify the package is updated: https://packagist.org/packages/qredex/sdk-go

### 9. Announce the Release

- [ ] Post release announcement in internal channels
- [ ] Update documentation links if necessary
- [ ] Notify stakeholders of the new release

## Hotfix Releases

For critical bug fixes:

1. Create a hotfix branch from the release tag: `git checkout -b hotfix/v0.1.1 v0.1.0`
2. Apply the fix
3. Update `CHANGELOG.md` with the fix
4. Bump PATCH version: `0.1.1`
5. Follow the release checklist above
6. Merge hotfix back to main: `git checkout main && git merge hotfix/v0.1.1`

## Major Version Releases

For breaking changes (MAJOR version bump):

1. **Deprecation Cycle**: Deprecate old API in a MINOR release first
2. **Migration Guide**: Document migration path in CHANGELOG.md
3. **Extended Testing**: Ensure thorough testing of breaking changes
4. **Communication**: Announce breaking changes well in advance

## Release Notes Template

Use this template for GitHub release descriptions:

```markdown
## What's New

[Brief summary of key changes]

## Changes

### Added
- [List new features]

### Changed
- [List changes to existing behavior]

### Fixed
- [List bug fixes]

### Security
- [List security improvements]

## Upgrade Guide

[Instructions for upgrading from previous version]

## Contributors

Thanks to [@contributor1, @contributor2] for contributions!

## Full Changelog

https://github.com/qredex/sdk-go/compare/v0.1.0...v0.2.0
```

## Troubleshooting

### Tag Already Exists

If the tag already exists:

```bash
# Delete local tag
git tag -d v0.2.0

# Delete remote tag
git push origin :refs/tags/v0.2.0

# Recreate tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

### Tests Failing

If tests fail during release preparation:

1. Fix the failing tests
2. Update CHANGELOG.md if the fix is user-visible
3. Re-run the release checklist

### Module Issues

If there are Go module issues:

```bash
# Clean module cache
go clean -modcache

# Tidy dependencies
go mod tidy

# Verify module
go mod verify
```

## Contact

For questions about the release process, contact:

- **Email**: os@qredex.com
- **GitHub**: https://github.com/qredex/sdk-go/issues
