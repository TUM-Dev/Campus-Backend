#!/usr/bin/env sh
# Updates the vendored Swagger UI assets from the swagger-ui-dist npm package.
# Usage: ./update_swagger-ui.sh [version]   (defaults to the latest release)
set -eu

cd "$(dirname "$0")"
version="${1:-$(pnpm view swagger-ui-dist version)}"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "updating swagger-ui to ${version}"
pnpm --dir "$tmp" add --silent "swagger-ui-dist@${version}"

for f in oauth2-redirect.html oauth2-redirect.js \
	swagger-ui-bundle.js swagger-ui-bundle.js.map \
	swagger-ui.css swagger-ui.css.map; do
	cp "$tmp/node_modules/swagger-ui-dist/$f" .
done
