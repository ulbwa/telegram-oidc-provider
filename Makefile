.PHONY: generate-api

OPENAPI_SPEC = api/openapi.yaml
GENERATED_DIR = internal/transport/http/generated

generate-api:
	@echo "Generating Go server code from OpenAPI spec..."
	@mkdir -p $(GENERATED_DIR)
	@oapi-codegen -package generated -generate types $(OPENAPI_SPEC) > $(GENERATED_DIR)/types.go
	@echo "✓ Types generated: $(GENERATED_DIR)/types.go"
	@oapi-codegen -package generated -generate server $(OPENAPI_SPEC) > $(GENERATED_DIR)/server.go
	@echo "✓ Server interfaces generated: $(GENERATED_DIR)/server.go"
	@oapi-codegen -package generated -generate spec $(OPENAPI_SPEC) > $(GENERATED_DIR)/spec.go
	@echo "✓ OpenAPI spec embedded: $(GENERATED_DIR)/spec.go"
