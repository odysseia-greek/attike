PROTO_DIRS := aristophanes
TOOLS_DIR := $(PWD)/.bin
PROTOC_GEN_DOC := $(TOOLS_DIR)/protoc-gen-doc

.PHONY: all
all: generate docs

.PHONY: generate
generate:
	@for dir in $(PROTO_DIRS); do \
		echo "Generating Protobuf files in $$dir..."; \
		buf generate --template $$dir/buf.gen.yaml $$dir; \
	done

.PHONY: tools
tools: $(PROTOC_GEN_DOC)

$(PROTOC_GEN_DOC):
	@mkdir -p $(TOOLS_DIR)
	@echo "Installing protoc-gen-doc..."
	@GOBIN=$(TOOLS_DIR) go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest

.PHONY: docs
docs: docs-grpc docs-graphql

.PHONY: docs-grpc
docs-grpc: tools
	@for dir in $(PROTO_DIRS); do \
		echo "Generating gRPC docs in $$dir..."; \
		PATH=$(TOOLS_DIR):$$PATH buf generate --template $$dir/buf.gen.docs.yaml $$dir; \
	done

.PHONY: docs-graphql
docs-graphql:
	@echo "Generating Euripides docs..."
	@spectaql -c euripides/docs/spectaql.yaml
