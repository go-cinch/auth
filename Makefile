SHELL := /usr/bin/env bash

SERVICE_NAME ?= auth

GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
GRPC_PORT:=$(shell grep -A 3 'grpc:' configs/server.yaml | grep 'addr:' | awk -F':' '{print $$3}' | tr -d ' ')
HTTP_PORT:=$(shell grep -A 3 'http:' configs/server.yaml | grep 'addr:' | awk -F':' '{print $$3}' | tr -d ' ')

# Tools directory
TOOLS_DIR := .tools

# Tool versions (pinned; override with `make VAR=...` as needed)
PROTOC_GEN_GO_VERSION ?= v1.36.11
PROTOC_GEN_GO_GRPC_VERSION ?= v1.6.0
PROTOC_GEN_GO_ERRORS_VERSION ?= c7a58ff59f80
PROTOC_GEN_GO_HTTP_VERSION ?= c7a58ff59f80
PROTOC_GEN_OPENAPI_VERSION ?= v0.7.1
PROTOC_GEN_VALIDATE_VERSION ?= v1.3.0
WIRE_VERSION ?= v0.7.0
SQL_MIGRATE_VERSION ?= v1.8.1
GENTOOL_VERSION ?= v0.0.2
PROTOC_GEN_GO := $(TOOLS_DIR)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(TOOLS_DIR)/bin/protoc-gen-go-grpc
PROTOC_GEN_GO_ERRORS := $(TOOLS_DIR)/bin/protoc-gen-go-errors
PROTOC_GEN_GO_HTTP := $(TOOLS_DIR)/bin/protoc-gen-go-http
PROTOC_GEN_OPENAPI := $(TOOLS_DIR)/bin/protoc-gen-openapi
PROTOC_GEN_VALIDATE := $(TOOLS_DIR)/bin/protoc-gen-validate
WIRE := $(TOOLS_DIR)/bin/wire
SQL_MIGRATE := $(TOOLS_DIR)/bin/sql-migrate
GENTOOL := $(TOOLS_DIR)/bin/gentool

# golangci-lint configuration
GOLANGCI_LINT_VERSION ?= v2.8.0
GOLANGCI_LINT := $(TOOLS_DIR)/bin/golangci-lint

ifeq ($(GOHOSTOS), windows)
	Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
	API_PROTO_FILES=$(shell find api -maxdepth 2 -name *.proto)
endif

.PHONY: init
# init env (install local toolchain into .tools/bin)
init:
	@mkdir -p $(TOOLS_DIR)/bin
	@echo "=== Installing toolchain to $(TOOLS_DIR)/bin ==="
	@if [ ! -x "$(PROTOC_GEN_GO)" ]; then \
		echo "  → protoc-gen-go ($(PROTOC_GEN_GO_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION); \
	fi
	@if [ ! -x "$(PROTOC_GEN_GO_GRPC)" ]; then \
		echo "  → protoc-gen-go-grpc ($(PROTOC_GEN_GO_GRPC_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION); \
	fi
	@if [ ! -x "$(PROTOC_GEN_GO_ERRORS)" ]; then \
		echo "  → protoc-gen-go-errors ($(PROTOC_GEN_GO_ERRORS_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@$(PROTOC_GEN_GO_ERRORS_VERSION); \
	fi
	@if [ ! -x "$(PROTOC_GEN_GO_HTTP)" ]; then \
		echo "  → protoc-gen-go-http ($(PROTOC_GEN_GO_HTTP_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@$(PROTOC_GEN_GO_HTTP_VERSION); \
	fi
	@if [ ! -x "$(PROTOC_GEN_OPENAPI)" ]; then \
		echo "  → protoc-gen-openapi ($(PROTOC_GEN_OPENAPI_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install github.com/google/gnostic/cmd/protoc-gen-openapi@$(PROTOC_GEN_OPENAPI_VERSION); \
	fi
	@if [ ! -x "$(PROTOC_GEN_VALIDATE)" ]; then \
		echo "  → protoc-gen-validate ($(PROTOC_GEN_VALIDATE_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install github.com/envoyproxy/protoc-gen-validate@$(PROTOC_GEN_VALIDATE_VERSION); \
	fi
	@if [ ! -x "$(WIRE)" ]; then \
		echo "  → wire ($(WIRE_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install github.com/google/wire/cmd/wire@$(WIRE_VERSION); \
	fi
	@if [ ! -x "$(SQL_MIGRATE)" ]; then \
		echo "  → sql-migrate ($(SQL_MIGRATE_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install github.com/rubenv/sql-migrate/...@$(SQL_MIGRATE_VERSION); \
	fi
	@if [ ! -x "$(GENTOOL)" ]; then \
		echo "  → gentool ($(GENTOOL_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install gorm.io/gen/tools/gentool@$(GENTOOL_VERSION); \
	fi
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "  → golangci-lint ($(GOLANGCI_LINT_VERSION))"; \
		GOBIN="$(PWD)/$(TOOLS_DIR)/bin" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi
	@echo "=== Toolchain installation complete ==="

.PHONY: config
# generate internal proto
config:
	@protoc --plugin=protoc-gen-go=$(PWD)/$(PROTOC_GEN_GO) \
	       --proto_path=./internal \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	@mkdir -p docs
	@for NAME in $(API_PROTO_FILES); do \
		echo $$NAME; \
		protoc --plugin=protoc-gen-go=$(PWD)/$(PROTOC_GEN_GO) \
			--plugin=protoc-gen-go-grpc=$(PWD)/$(PROTOC_GEN_GO_GRPC) \
			--plugin=protoc-gen-go-errors=$(PWD)/$(PROTOC_GEN_GO_ERRORS) \
			--plugin=protoc-gen-go-http=$(PWD)/$(PROTOC_GEN_GO_HTTP) \
			--plugin=protoc-gen-openapi=$(PWD)/$(PROTOC_GEN_OPENAPI) \
			--plugin=protoc-gen-validate=$(PWD)/$(PROTOC_GEN_VALIDATE) \
			--proto_path=./api \
			--proto_path=./third_party \
			--go_out=. \
			--go-errors_out=. \
			--go-http_out=. \
			--go-grpc_out=. \
			--validate_out=lang=go:. \
			--openapi_out=fq_schema_naming=true,default_response=false,output_mode=source_relative:docs \
			$$NAME; \
	done
	@echo 'You can import *.json into https://editor.swagger.io/'

.PHONY: wire
# generate wire
wire:
	cd cmd/auth && ../../$(TOOLS_DIR)/bin/wire

.PHONY: gen-createdb
# create database (only needed once)
gen-createdb:
	@echo "Creating PostgreSQL database auth..."
	@if docker ps --filter "name=pg" --format "{{.Names}}" | grep -q "^pg$$"; then \
		echo "Found running PostgreSQL container: pg"; \
		PG_USER=$$(grep 'dsn:' configs/db.yaml | sed 's|.*postgresql://\([^:]*\):.*|\1|'); \
		PG_PASSWORD=$$(grep 'dsn:' configs/db.yaml | sed 's|.*postgresql://[^:]*:\([^@]*\)@.*|\1|'); \
		echo "Using PostgreSQL user: $$PG_USER"; \
		docker exec -e PGPASSWORD=$$PG_PASSWORD pg psql -U $$PG_USER -d postgres -c "CREATE DATABASE auth;" 2>&1 | grep -v "already exists" || true; \
		echo "✓ Database 'auth' is ready"; \
	fi

.PHONY: gen-config
# generate sql-migrate config from db.yaml (internal use)
gen-config:
	@echo "development:" > /tmp/sql-migrate-$$$$.yml
	@echo "  dialect: $$(grep 'driver:' configs/db.yaml | grep -v '#' | awk '{print $$2}')" >> /tmp/sql-migrate-$$$$.yml
	@echo "  datasource: $$(grep 'dsn:' configs/db.yaml | sed 's/.*dsn: *"\(.*\)".*/\1/' | sed 's/TimeZone=UTC//')" >> /tmp/sql-migrate-$$$$.yml
	@echo "  dir: internal/db/migrations" >> /tmp/sql-migrate-$$$$.yml
	@echo "  table: schema_migrations" >> /tmp/sql-migrate-$$$$.yml
	@echo "/tmp/sql-migrate-$$$$.yml"

.PHONY: gen-up
# run database migrations up using sql-migrate
gen-up:
	@echo "Running database migrations up..."
	@TMP_CONFIG=$$(mktemp /tmp/sql-migrate.XXXXXX.yml); \
	echo "development:" > $$TMP_CONFIG; \
	echo "  dialect: $$(grep 'driver:' configs/db.yaml | grep -v '#' | awk '{print $$2}')" >> $$TMP_CONFIG; \
	echo "  datasource: $$(grep 'dsn:' configs/db.yaml | sed 's/.*dsn: *"\(.*\)".*/\1/' | sed 's/TimeZone=UTC//')" >> $$TMP_CONFIG; \
	echo "  dir: internal/db/migrations" >> $$TMP_CONFIG; \
	echo "  table: schema_migrations" >> $$TMP_CONFIG; \
	$(SQL_MIGRATE) up -config=$$TMP_CONFIG -env=development; \
	EXIT_CODE=$$?; \
	rm -f $$TMP_CONFIG; \
	if [ $$EXIT_CODE -ne 0 ]; then \
		echo ""; \
		echo "💡 Hint: If database doesn't exist, run: make gen-createdb"; \
	fi; \
	exit $$EXIT_CODE

.PHONY: gen-down
# run database migrations down using sql-migrate
gen-down:
	@echo "Running database migrations down..."
	@TMP_CONFIG=$$(mktemp /tmp/sql-migrate.XXXXXX.yml); \
	echo "development:" > $$TMP_CONFIG; \
	echo "  dialect: $$(grep 'driver:' configs/db.yaml | grep -v '#' | awk '{print $$2}')" >> $$TMP_CONFIG; \
	echo "  datasource: $$(grep 'dsn:' configs/db.yaml | sed 's/.*dsn: *"\(.*\)".*/\1/' | sed 's/TimeZone=UTC//')" >> $$TMP_CONFIG; \
	echo "  dir: internal/db/migrations" >> $$TMP_CONFIG; \
	echo "  table: schema_migrations" >> $$TMP_CONFIG; \
	$(SQL_MIGRATE) down -config=$$TMP_CONFIG -env=development; \
	rm -f $$TMP_CONFIG

.PHONY: gen-status
# show database migration status using sql-migrate
gen-status:
	@TMP_CONFIG=$$(mktemp /tmp/sql-migrate.XXXXXX.yml); \
	echo "development:" > $$TMP_CONFIG; \
	echo "  dialect: $$(grep 'driver:' configs/db.yaml | grep -v '#' | awk '{print $$2}')" >> $$TMP_CONFIG; \
	echo "  datasource: $$(grep 'dsn:' configs/db.yaml | sed 's/.*dsn: *"\(.*\)".*/\1/' | sed 's/TimeZone=UTC//')" >> $$TMP_CONFIG; \
	echo "  dir: internal/db/migrations" >> $$TMP_CONFIG; \
	echo "  table: schema_migrations" >> $$TMP_CONFIG; \
	$(SQL_MIGRATE) status -config=$$TMP_CONFIG -env=development; \
	rm -f $$TMP_CONFIG

.PHONY: gen-model
# generate GORM models from database using gentool
gen-model:
	@echo "Generating GORM models from database..."
	@DB_TYPE=$$(grep 'driver:' configs/db.yaml | grep -v '#' | awk '{print $$2}'); \
	DSN=$$(awk '/^[[:space:]]*driver:/{flag=1} flag && /^[[:space:]]*dsn:/{print; exit}' configs/db.yaml | sed 's/.*dsn: *"\(.*\)".*/\1/'); \
	echo "Database: $$DB_TYPE"; \
	echo "DSN: $$DSN"; \
	$(GENTOOL) -dsn "$$DSN" \
		-db $$DB_TYPE \
		-onlyModel \
		-outPath internal/data/model \
		-modelPkgName model \
		-fieldNullable \
		-fieldWithIndexTag \
		-fieldWithTypeTag
	@echo "Removing excluded tables..."; \
	EXCLUDE_TABLES=$$(awk '/^[[:space:]]*exclude:/{flag=1; next} /^[[:space:]]*[a-zA-Z]/{flag=0} flag && /^[[:space:]]*-/{print}' configs/db.yaml | sed 's/^ *- *//' | tr '\n' ' '); \
	if [ -n "$$EXCLUDE_TABLES" ]; then \
		echo "Excluded tables: $$EXCLUDE_TABLES"; \
		for table in $$EXCLUDE_TABLES; do \
			rm -f internal/data/model/$${table}.gen.go; \
			echo "  - Removed $${table}.gen.go"; \
		done; \
	fi
	@echo "Applying customizations..."; \
	USE_CARBON=$$(grep 'useCarbonTime:' configs/db.yaml | awk '{print $$2}'); \
	REMOVE_PTR=$$(grep 'removeNullablePointer:' configs/db.yaml | awk '{print $$2}'); \
	JSON_STRING_FIELDS=$$(awk '/^[[:space:]]*jsonStringFields:/{flag=1; next} /^[[:space:]]*[a-zA-Z]/{flag=0} flag && /^[[:space:]]*-/{print}' configs/db.yaml | sed 's/^ *- *//' | tr '\n' '|' | sed 's/|$$//'); \
	for file in internal/data/model/*.gen.go; do \
		if [ -f "$$file" ]; then \
			if [ "$$USE_CARBON" = "true" ]; then \
				sed -i.bak 's/"time"/"github.com\/golang-module\/carbon\/v2"/g' "$$file"; \
				sed -i.bak 's/\*time\.Time/*carbon.DateTime/g' "$$file"; \
				sed -i.bak 's/time\.Time/carbon.DateTime/g' "$$file"; \
				echo "  - Converted time.Time to carbon.DateTime in $$file"; \
			fi; \
			if [ -n "$$JSON_STRING_FIELDS" ]; then \
				TABLE_NAME=$$(basename "$$file" .gen.go); \
				IFS='|' read -ra FIELDS <<< "$$JSON_STRING_FIELDS"; \
				for field_spec in "$${FIELDS[@]}"; do \
					if echo "$$field_spec" | grep -q '\.'; then \
						SPEC_TABLE=$$(echo "$$field_spec" | cut -d'.' -f1); \
						SPEC_FIELD=$$(echo "$$field_spec" | cut -d'.' -f2); \
						if [ "$$TABLE_NAME" = "$$SPEC_TABLE" ]; then \
							sed -i.bak "s/\(json:\"$$SPEC_FIELD\)\"/\1,string\"/g" "$$file"; \
							echo "  - Added ,string to $$SPEC_TABLE.$$SPEC_FIELD in $$file"; \
						fi; \
					else \
						sed -i.bak "s/\(json:\"$$field_spec\)\"/\1,string\"/g" "$$file"; \
					fi; \
				done; \
			fi; \
			if [ "$$REMOVE_PTR" = "true" ]; then \
				sed -i.bak 's/\*carbon\.DateTime/carbon.DateTime/g' "$$file"; \
				sed -i.bak 's/\*time\.Time/time.Time/g' "$$file"; \
				sed -i.bak 's/\*string/string/g' "$$file"; \
				sed -i.bak 's/\*int64/int64/g' "$$file"; \
				sed -i.bak 's/\*int32/int32/g' "$$file"; \
				sed -i.bak 's/\*float64/float64/g' "$$file"; \
				sed -i.bak 's/\*bool/bool/g' "$$file"; \
				echo "  - Removed nullable pointers in $$file"; \
			fi; \
			rm -f "$$file.bak"; \
		fi; \
	done
	@echo "Adding associations..."; \
	ASSOCIATIONS=$$(awk '/^[[:space:]]*association:/{flag=1; next} /^[[:space:]]*[a-zA-Z]/{flag=0} flag && /^[[:space:]]*-/{print}' configs/db.yaml | sed "s/^ *- *'//;s/'$$//"); \
	if [ -n "$$ASSOCIATIONS" ]; then \
		echo "$$ASSOCIATIONS" | while IFS= read -r assoc; do \
			if [ -z "$$assoc" ]; then continue; fi; \
			MODEL_TABLE=$$(echo "$$assoc" | cut -d'|' -f1); \
			RELATED_TABLE=$$(echo "$$assoc" | cut -d'|' -f2); \
			FIELD_NAME=$$(echo "$$assoc" | cut -d'|' -f3); \
			RELATION_TYPE=$$(echo "$$assoc" | cut -d'|' -f4); \
			GORM_TAG=$$(echo "$$assoc" | cut -d'|' -f5); \
			MODEL_FILE="internal/data/model/$${MODEL_TABLE}.gen.go"; \
			if [ ! -f "$$MODEL_FILE" ]; then continue; fi; \
			STRUCT_NAME=$$(grep "^type.*struct" "$$MODEL_FILE" | awk '{print $$2}'); \
			IS_SLICE=false; \
			FIELD_TYPE="$$RELATED_TABLE"; \
			if echo "$$FIELD_NAME" | grep -q '^\[\]'; then \
				IS_SLICE=true; \
				FIELD_NAME=$$(echo "$$FIELD_NAME" | sed 's/^\[\]//'); \
			fi; \
			if echo "$$FIELD_TYPE" | grep -q '^\*'; then \
				FIELD_TYPE=$$(echo "$$FIELD_TYPE" | sed 's/^\*//'); \
			else \
				RELATED_MODEL=$$(grep "^type.*$$RELATED_TABLE.*struct" internal/data/model/*.gen.go 2>/dev/null | head -1 | awk '{print $$2}'); \
				if [ -n "$$RELATED_MODEL" ]; then \
					FIELD_TYPE="$$RELATED_MODEL"; \
				else \
					RELATED_PASCAL=$$(echo "$$RELATED_TABLE" | awk '{for(i=1;i<=NF;i++){sub(/./,toupper(substr($$i,1,1)),$$i);}}1' FS="_" OFS=""); \
					FIELD_TYPE="$$RELATED_PASCAL"; \
				fi; \
			fi; \
			if [ "$$IS_SLICE" = "true" ]; then \
				FIELD_TYPE="[]$$FIELD_TYPE"; \
			fi; \
			if [ "$$RELATION_TYPE" = "has_one" ] && [ "$$IS_SLICE" = "false" ]; then \
				FIELD_TYPE="*$$FIELD_TYPE"; \
			elif [ "$$RELATION_TYPE" = "has_many" ]; then \
				FIELD_TYPE="[]$$FIELD_TYPE"; \
			fi; \
			JSON_NAME=$$(echo "$$FIELD_NAME" | awk '{print tolower($$0)}'); \
			FIELD_LINE="	$$FIELD_NAME $${FIELD_TYPE} \`gorm:\\\"$$GORM_TAG\\\" json:\\\"$$JSON_NAME\\\"\`"; \
			awk -v line="$$FIELD_LINE" '/^}$$/ && !found {print line; found=1} {print}' "$$MODEL_FILE" > "$$MODEL_FILE.tmp" && mv "$$MODEL_FILE.tmp" "$$MODEL_FILE"; \
			echo "  - Added $$FIELD_NAME to $$MODEL_TABLE"; \
		done; \
	fi
	@echo "✓ Models generated in internal/data/model/"

.PHONY: lint
# run code linters (gofmt + golangci-lint)
lint:
	@echo "Running gofmt..."
	@gofmt -w -s .
	@if [ -x "$(GOLANGCI_LINT)" ]; then \
		$(GOLANGCI_LINT) run --fix ./...; \
	elif command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix ./...; \
	else \
		echo "golangci-lint not found. Run 'make init' to install $(GOLANGCI_LINT_VERSION), or install via: brew install golangci-lint"; \
		echo "Running basic checks (go vet, gofmt) instead..."; \
		go vet ./...; \
		[ -z "$$(gofmt -s -l .)" ] || (echo "Files not gofmt'd:" && gofmt -s -l . && exit 1); \
	fi

.PHONY: build
# build service binary
build:
	@mkdir -p bin/
	@go build -ldflags "-s -w -X main.Version=$(VERSION)" -o ./bin/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)

.PHONY: run
# run service (dev)
run:
	@go run ./cmd/$(SERVICE_NAME)

.PHONY: test
# run tests
test:
	@go test ./...

.PHONY: debug-grpc
# debug grpc with grpcui
debug-grpc:
	@echo "Starting gRPC UI on port $(GRPC_PORT)..."
	@API_PROTO=$$(find api -name "*.proto" | head -1); \
	if [ -n "$$API_PROTO" ]; then \
		echo "Using proto file: $$API_PROTO"; \
		grpcui -plaintext \
			-use-reflection \
			-import-path ./api \
			-import-path ./third_party \
			-proto $$API_PROTO \
			"127.0.0.1:$(GRPC_PORT)"; \
	else \
		echo "No proto files found, using reflection only..."; \
		grpcui -plaintext "127.0.0.1:$(GRPC_PORT)"; \
	fi

.PHONY: debug-http
# debug http api with swagger ui
debug-http:
	@OPENAPI_FILE=$$(find docs -name "*.openapi.yaml" | head -1); \
	if [ -z "$$OPENAPI_FILE" ]; then \
		echo "Error: No OpenAPI file found in docs directory. Please run 'make api' first."; \
		exit 1; \
	fi; \
	echo "Generating Swagger UI with embedded OpenAPI spec..."; \
	TMP_YAML="docs/.openapi-temp.yaml"; \
	awk 'BEGIN {added=0; in_info=0} \
		/^openapi:/ {print; next} \
		/^info:/ {in_info=1; print; next} \
		in_info && /^[a-z]+:/ && !added { \
			print "servers:"; \
			print "    - url: http://localhost:$(HTTP_PORT)"; \
			added=1; \
		} \
		/^servers:/ {skip=1; next} \
		skip && /^    - / {next} \
		skip && /^[a-z]+:/ {skip=0} \
		{if (!skip) print}' "$$OPENAPI_FILE" | \
		sed 's/\\/\\\\/g; s/`/\\`/g' > "$$TMP_YAML"; \
	awk '/OPENAPI_YAML_CONTENT/ { \
		print "      const openapiYaml = `"; \
		while ((getline line < "'"$$TMP_YAML"'") > 0) print line; \
		print "`"; \
		next; \
	} \
	{print}' docs/swagger-for-debug.html > docs/swagger-ui.html; \
	rm -f "$$TMP_YAML"; \
	echo "Swagger UI available at:"; \
	echo "  http://localhost:$(HTTP_PORT)/docs/swagger-ui.html"; \
	echo ""; \
	echo "Make sure your service is running with 'make run'"; \
	open "http://localhost:$(HTTP_PORT)/docs/swagger-ui.html" || \
	echo "Please open http://localhost:$(HTTP_PORT)/docs/swagger-ui.html in your browser"

.PHONY: all
# generate all
all:
	go mod tidy
	make api
	make config
	go mod tidy
	make wire
	make lint

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
