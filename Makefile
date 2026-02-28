.PHONY: generate-api dev release build invalidate foldermap

OPENAPI_SPEC = api/openapi.yaml

generate-api: 
	@echo "Generating API code from OpenAPI spec..."
	@oapi-codegen -config api/oapi_codegen.yaml $(OPENAPI_SPEC)
	@echo "✓ All API code generated successfully"

dev:
	cd ./external && npm run dev --turbo

release:
	cd ./external && npm run build && npm run start

build:
	cd ./external && npm run build --debug-prerender

invalidate:
	cd ./external && powershell -Command "Remove-Item -Recurse -Force .next, node_modules, package-lock.json"
	cd ./external && npm install

foldermap:
	py folder_map.py -fcg --hide-empty --no-format -r external -o folder_map.txt --match "^(?!package-lock|README)"
