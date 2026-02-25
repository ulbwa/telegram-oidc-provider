.PHONY: generate-api

OPENAPI_SPEC = api/openapi.yaml

generate-api: 
	@echo "Generating API code from OpenAPI spec..."
	@oapi-codegen -config api/oapi_codegen.yaml $(OPENAPI_SPEC)
	@echo "✓ All API code generated successfully"
