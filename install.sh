#!/bin/sh
# Totion installer — https://github.com/milankatira/totion
#
#   curl -fsSL https://raw.githubusercontent.com/milankatira/totion/main/install.sh | sh
#
# Downloads the latest release binary for your platform into ~/.local/bin.

set -eu

REPO="milankatira/totion"
INSTALL_DIR="${TOTION_INSTALL_DIR:-$HOME/.local/bin}"

main() {
    os=$(uname -s)
    case "$os" in
        Darwin) os="darwin" ;;
        Linux) os="linux" ;;
        *)
            echo "error: unsupported operating system: $os" >&2
            echo "On Windows, download totion_windows_*.exe from:" >&2
            echo "  https://github.com/$REPO/releases/latest" >&2
            exit 1
            ;;
    esac

    arch=$(uname -m)
    case "$arch" in
        x86_64 | amd64) arch="amd64" ;;
        arm64 | aarch64) arch="arm64" ;;
        *)
            echo "error: unsupported architecture: $arch" >&2
            exit 1
            ;;
    esac

    url="https://github.com/$REPO/releases/latest/download/totion_${os}_${arch}"

    echo "Downloading totion (${os}/${arch})..."
    mkdir -p "$INSTALL_DIR"
    curl -fSL --progress-bar "$url" -o "$INSTALL_DIR/totion"
    chmod +x "$INSTALL_DIR/totion"

    echo "Installed: $INSTALL_DIR/totion"

    case ":$PATH:" in
        *":$INSTALL_DIR:"*)
            echo "Run 'totion' to get started."
            ;;
        *)
            echo
            echo "Note: $INSTALL_DIR is not on your PATH. Add it with:"
            echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
            ;;
    esac
}

main
