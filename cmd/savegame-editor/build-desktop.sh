#!/usr/bin/env bash
#
# Builds the savegame editor as a native desktop application (Wails v2,
# build tag "desktop") for the current platform, without the Wails CLI.
#
#   macOS:   build/bin/X-COM Savegame Editor.app
#   Linux:   build/bin/savegame-editor
#   Windows: build/bin/savegame-editor.exe   (run from Git Bash / MSYS2)
#
# Usage: ./build-desktop.sh [-s] [-o OUTDIR]
#
#   -s         skip the frontend build (reuse frontend/dist)
#   -o OUTDIR  output directory (default: build/bin)
#
# The version defaults to `git describe`; override it with VERSION=...
# The build needs cgo and the platform's webview (WebKit on macOS, WebKitGTK
# on Linux, WebView2 on Windows), so it cannot be cross-compiled.

set -euo pipefail

cd "$(dirname "$0")"

APP_NAME="X-COM Savegame Editor"
BINARY="savegame-editor"
BUNDLE_ID="de.redtoad.xcom-savegame-editor"

skip_frontend=0
outdir="build/bin"

usage() {
	sed -n '9,13p' "$0" | sed 's/^# \{0,1\}//'
}

while getopts ":so:h" opt; do
	case "$opt" in
	s) skip_frontend=1 ;;
	o) outdir="$OPTARG" ;;
	h) usage; exit 0 ;;
	*) usage >&2; exit 2 ;;
	esac
done

log() { printf '==> %s\n' "$*"; }

for tool in go npm; do
	command -v "$tool" >/dev/null || { echo "error: $tool not found" >&2; exit 1; }
done

version="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
commit="$(git rev-parse --short HEAD 2>/dev/null || echo none)"
date="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
ldflags="-s -w -X main.version=${version} -X main.commit=${commit} -X main.date=${date} -X main.builtBy=build-desktop.sh"
tags="desktop,production"

if [ "$skip_frontend" -eq 0 ]; then
	log "Building frontend"
	if [ -f frontend/package-lock.json ]; then
		npm --prefix frontend ci
	else
		npm --prefix frontend install
	fi
	npm --prefix frontend run build
elif [ ! -d frontend/dist ]; then
	echo "error: frontend/dist missing, run without -s" >&2
	exit 1
fi

export CGO_ENABLED=1
mkdir -p "$outdir"

case "$(uname -s)" in
Darwin)
	# minimum macOS version of the Go runtime (Go 1.27: macOS 13)
	min_macos="${MACOSX_DEPLOYMENT_TARGET:-13.0}"
	export CGO_CFLAGS="${CGO_CFLAGS:-} -mmacosx-version-min=${min_macos}"
	export CGO_LDFLAGS="${CGO_LDFLAGS:-} -mmacosx-version-min=${min_macos} -framework UniformTypeIdentifiers"

	app="$outdir/$APP_NAME.app"
	log "Building $app"
	rm -rf "$app"
	mkdir -p "$app/Contents/MacOS" "$app/Contents/Resources"
	go build -tags "$tags" -ldflags "$ldflags" -o "$app/Contents/MacOS/$BINARY" .

	# app icon (.icns) from build/appicon.png
	if [ -f build/appicon.png ] && command -v iconutil >/dev/null; then
		iconset="$(mktemp -d)/iconfile.iconset"
		mkdir -p "$iconset"
		for size in 16 32 128 256 512; do
			sips -z "$size" "$size" build/appicon.png --out "$iconset/icon_${size}x${size}.png" >/dev/null
			double=$((size * 2))
			sips -z "$double" "$double" build/appicon.png --out "$iconset/icon_${size}x${size}@2x.png" >/dev/null
		done
		iconutil -c icns "$iconset" -o "$app/Contents/Resources/iconfile.icns"
		rm -rf "$(dirname "$iconset")"
	fi

	plist_version="${version#v}"
	cat >"$app/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundlePackageType</key><string>APPL</string>
	<key>CFBundleName</key><string>${APP_NAME}</string>
	<key>CFBundleDisplayName</key><string>${APP_NAME}</string>
	<key>CFBundleExecutable</key><string>${BINARY}</string>
	<key>CFBundleIdentifier</key><string>${BUNDLE_ID}</string>
	<key>CFBundleVersion</key><string>${plist_version}</string>
	<key>CFBundleShortVersionString</key><string>${plist_version}</string>
	<key>CFBundleIconFile</key><string>iconfile</string>
	<key>CFBundleGetInfoString</key><string>Savegame editor for X-COM: UFO Defense</string>
	<key>LSMinimumSystemVersion</key><string>${min_macos}</string>
	<key>NSHighResolutionCapable</key><true/>
	<key>NSHumanReadableCopyright</key><string>Sebastian Rahlf</string>
</dict>
</plist>
EOF

	# ad-hoc signature, required to run on Apple Silicon
	if command -v codesign >/dev/null; then
		codesign --force --deep --sign - "$app"
	fi
	log "Done: $app"
	;;

Linux)
	# Wails supports WebKitGTK 4.1 (newer distributions) and 4.0.
	if command -v pkg-config >/dev/null && pkg-config --exists webkit2gtk-4.1; then
		tags="$tags,webkit2_41"
	fi
	log "Building $outdir/$BINARY"
	go build -tags "$tags" -ldflags "$ldflags" -o "$outdir/$BINARY" .
	log "Done: $outdir/$BINARY"
	;;

MINGW* | MSYS* | CYGWIN*)
	log "Building $outdir/$BINARY.exe"
	go build -tags "$tags" -ldflags "$ldflags -H windowsgui" -o "$outdir/$BINARY.exe" .
	log "Done: $outdir/$BINARY.exe"
	;;

*)
	echo "error: unsupported platform $(uname -s)" >&2
	exit 1
	;;
esac
