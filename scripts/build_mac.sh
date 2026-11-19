#!/bin/bash
# Build EDA-Game.app, a native macOS window around the static reader in site/, and install it.
#
#     python3 scripts/build_site.py     # if site/ is stale
#     scripts/build_mac.sh              # builds build/mac/EDA-Game.app, copies it to /Applications
#     scripts/build_mac.sh --no-install
#
# The app works offline: marked (mac/vendor) and KaTeX (web/next-app/node_modules) are bundled
# and the page's CDN links are rewritten to the bundled copies.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/build/mac"
APP="$OUT/EDA-Game.app"
KATEX="$ROOT/web/next-app/node_modules/katex/dist"

rm -rf "$APP"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"

swiftc -O -o "$APP/Contents/MacOS/EDA-Game" "$ROOT/mac/main.swift" -framework Cocoa -framework WebKit

# The reader, minus files only the website needs.
rsync -a --exclude '.nojekyll' "$ROOT/site/" "$APP/Contents/Resources/site/"
VENDOR="$APP/Contents/Resources/site/vendor"
mkdir -p "$VENDOR/fonts"
cp "$ROOT/mac/vendor/marked.min.js" "$KATEX/katex.min.js" "$KATEX/katex.min.css" "$VENDOR/"
cp "$KATEX"/fonts/*.woff2 "$VENDOR/fonts/"
sed -i '' -E \
  -e 's#https://cdn.jsdelivr.net/npm/marked@[^"]*/marked.min.js#vendor/marked.min.js#' \
  -e 's#https://cdn.jsdelivr.net/npm/katex@[^"]*/dist/katex.min.(js|css)#vendor/katex.min.\1#g' \
  "$APP/Contents/Resources/site/index.html"
if grep -q 'cdn.jsdelivr' "$APP/Contents/Resources/site/index.html"; then
  echo "index.html still loads from a CDN" >&2; exit 1
fi

# Icon.
ICONSET="$OUT/AppIcon.iconset"
rm -rf "$ICONSET"; mkdir -p "$ICONSET"
swift "$ROOT/mac/icon.swift" "$OUT/icon.png"
for s in 16 32 128 256 512; do
  sips -z $s $s "$OUT/icon.png" --out "$ICONSET/icon_${s}x${s}.png" >/dev/null
  sips -z $((s * 2)) $((s * 2)) "$OUT/icon.png" --out "$ICONSET/icon_${s}x${s}@2x.png" >/dev/null
done
iconutil -c icns "$ICONSET" -o "$APP/Contents/Resources/AppIcon.icns"

cat > "$APP/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key><string>EDA-Game</string>
  <key>CFBundleDisplayName</key><string>EDA-Game</string>
  <key>CFBundleIdentifier</key><string>io.github.adityanagachandra.eda-game</string>
  <key>CFBundleExecutable</key><string>EDA-Game</string>
  <key>CFBundleIconFile</key><string>AppIcon</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>CFBundleShortVersionString</key><string>1.0</string>
  <key>CFBundleVersion</key><string>$(git -C "$ROOT" rev-list --count HEAD)</string>
  <key>LSMinimumSystemVersion</key><string>13.0</string>
  <key>LSApplicationCategoryType</key><string>public.app-category.education</string>
  <key>NSHighResolutionCapable</key><true/>
  <key>NSHumanReadableCopyright</key><string>Lessons from EDA-Game, built on mathua.</string>
</dict>
</plist>
EOF

# Ad-hoc signature, so macOS runs it on this Mac without a developer certificate.
codesign --force --deep --sign - "$APP"
echo "built $APP ($(du -sh "$APP" | cut -f1))"

if [[ "${1:-}" != "--no-install" ]]; then
  rm -rf /Applications/EDA-Game.app
  ditto "$APP" /Applications/EDA-Game.app
  echo "installed /Applications/EDA-Game.app"
fi
