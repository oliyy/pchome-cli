#!/usr/bin/env sh

set -eu

VERSION="${GORELEASER_VERSION:-v2.15.1}"
RAW_VERSION="${VERSION#v}"
CACHE_ROOT="${XDG_CACHE_HOME:-$HOME/.cache}/pchome-cli/tools/goreleaser"
CACHE_DIR="$CACHE_ROOT/$RAW_VERSION"
BIN="$CACHE_DIR/goreleaser"

if [ ! -x "$BIN" ]; then
  case "$(uname -s)" in
    Darwin) OS="Darwin" ;;
    Linux) OS="Linux" ;;
    *)
      echo "unsupported OS for GoReleaser bootstrap: $(uname -s)" >&2
      exit 1
      ;;
  esac

  case "$(uname -m)" in
    x86_64|amd64) ARCH="x86_64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
      echo "unsupported architecture for GoReleaser bootstrap: $(uname -m)" >&2
      exit 1
      ;;
  esac

  ARCHIVE="goreleaser_${OS}_${ARCH}.tar.gz"
  BASE_URL="https://github.com/goreleaser/goreleaser/releases/download/$VERSION"
  TMP_DIR="$(mktemp -d)"

  cleanup() {
    rm -rf "$TMP_DIR"
  }
  trap cleanup EXIT

  mkdir -p "$CACHE_DIR"

  curl -fsSL "$BASE_URL/$ARCHIVE" -o "$TMP_DIR/$ARCHIVE"
  curl -fsSL "$BASE_URL/checksums.txt" -o "$TMP_DIR/checksums.txt"

  if command -v shasum >/dev/null 2>&1; then
    (
      cd "$TMP_DIR"
      grep "  $ARCHIVE$" checksums.txt | shasum -a 256 -c -
    )
  elif command -v sha256sum >/dev/null 2>&1; then
    (
      cd "$TMP_DIR"
      grep "  $ARCHIVE$" checksums.txt | sha256sum -c -
    )
  else
    echo "unable to verify GoReleaser download: missing shasum or sha256sum" >&2
    exit 1
  fi

  tar -xzf "$TMP_DIR/$ARCHIVE" -C "$TMP_DIR" goreleaser
  install -m 0755 "$TMP_DIR/goreleaser" "$BIN"
fi

exec "$BIN" "$@"
