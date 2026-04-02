# Contributing & Releasing

## Repository Layout

```
cmd/                  Cobra command layer and text rendering
cmd/pchome/           Binary entrypoint
cmd/testdata/         Golden fixtures for CLI help and text output
pkg/catalog/          Normalized product models and aggregate service
pkg/config/           Config loading and validation
pkg/i18n/             Language handling and translation catalog
pkg/output/           Shared table rendering
pkg/pchome/           Upstream PChome API clients
pkg/*/testdata/       Package-level test fixtures
```

## Development

```bash
# Fast local build for the current platform
make build

# Unit tests + GoReleaser config validation
make verify

# Simulate a full release locally into ./dist
make release-snapshot
```

On first run, `make verify` and `make release-snapshot` automatically download and cache the pinned GoReleaser version.

Run unit tests directly:

```bash
go test ./...
```

## Build and Release

Recommended maintainer workflow:

```bash
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

After a `v*` tag is pushed, GitHub Actions runs the release workflow and GoReleaser publishes macOS, Linux, and Windows archives for `amd64` and `arm64`, plus `checksums.txt`, to GitHub Releases.

### Homebrew / Scoop Setup

One-time setup for package manager distribution:

- Create an `oliyy/homebrew-tap` repository for `Formula/pchome-cli.rb`
- Create an `oliyy/scoop-bucket` repository for `pchome-cli.json`
- Add a `PACKAGE_REPOS_TOKEN` GitHub Actions secret to `pchome-cli`
- That token needs write access to both repositories

Notes:

- The Homebrew path here is a personal tap, not `homebrew/core`
- GoReleaser updates the tap and bucket on normal release tags; prerelease tags are skipped automatically

## Notes

- Recommendation token precedence is `hermes.token` -> `PCHOME_HERMES_TOKEN` -> bundled fallback token.
- Machine-readable output keeps stable English schema keys regardless of locale, so agent integrations do not break.
- The current schema version is `--schema-version v1`.
