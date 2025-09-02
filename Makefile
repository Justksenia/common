generate_kafka_schema:
	@protoc  \
    		   --go_out=schema/kafka/gen --go_opt=paths=source_relative \
    		   --go-grpc_out=schema/kafka/gen --go-grpc_opt=paths=source_relative\
    		   --proto_path schema/kafka/proto \
    		   schema/kafka/proto/*

lint:
	@golangci-lint run ./... --timeout 20s

lint-fix:
	@golangci-lint run ./... --fix --timeout 20s

lint++:
	golangci-lint run ./...  --issues-exit-code 0 --out-format code-climate | \
	tee gl-code-quality-report.json | \
	jq -r '.[] | "\(.location.path):\(.location.lines.begin) \(.description)"'

gen-openapi:
	@echo "Generating openapi..."
	@go run -mod=mod github.com/ogen-go/ogen/cmd/ogen@v0.82.0 --target crm/client/v1/gen -package gen ./crm/api/openapi.yaml