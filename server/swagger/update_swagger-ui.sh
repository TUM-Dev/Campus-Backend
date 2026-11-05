#!/usr/bin/env sh
# Updates the vendored Swagger UI assets from the swagger-ui-dist npm package.
# Usage: ./update_swagger-ui.sh [version]   (defaults to the latest release)
set -eu

cd "$(dirname "$0")"
version="${1:-$(npm view swagger-ui-dist version)}"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "updating swagger-ui to ${version}"
npm pack --silent --pack-destination "$tmp" "swagger-ui-dist@${version}" >/dev/null
tar -xzf "$tmp"/swagger-ui-dist-*.tgz -C "$tmp"

for f in oauth2-redirect.html oauth2-redirect.js \
	swagger-ui-bundle.js swagger-ui-bundle.js.map \
	swagger-ui.css swagger-ui.css.map; do
	cp "$tmp/package/$f" .
done

echo "updating the favicon"
curl -fsSL https://www.tum.app/assets/img/fav.png --output fav.png
