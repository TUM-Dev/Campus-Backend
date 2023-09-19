#!/usr/bin/sh

echo updating the swagger-ui
pnpm i swagger-ui

rsync node_modules/swagger-ui/dist/oauth2-redirect.html .
rsync node_modules/swagger-ui/dist/swagger-ui-bundle.js .
rsync node_modules/swagger-ui/dist/swagger-ui-bundle.js.map .
rsync node_modules/swagger-ui/dist/swagger-ui.css .
rsync node_modules/swagger-ui/dist/swagger-ui.css.map .

echo cleaning up temporary files
rm -fr node_modules
rm -fr package*.json

echo updating the favicon
curl https://www.tum.app/assets/img/fav.png --output fav.png
