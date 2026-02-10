dev:
	npm run dev --turbo

release:
	npm run build
	npm run start

build:
	npm run build

invalidate:
	powershell -Command "Remove-Item -Recurse -Force .next, node_modules, package-lock.json"
	npm install

foldermap:
	py folder_map.py -fcg --hide-empty --no-format -o folder_map.txt --match "^(?!package-lock|README)"