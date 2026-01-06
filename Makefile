PROTO_DIRS := aristophanes

.PHONY: all
all: generate docs

.PHONY: generate
generate:
	@for dir in $(PROTO_DIRS); do \
		echo "Generating Protobuf files in $$dir..."; \
		buf generate --template $$dir/buf.gen.yaml $$dir; \
	done

.PHONY: docs
docs:
	@for dir in $(PROTO_DIRS); do \
		echo "Generating docs in $$dir..."; \
		docker run --rm \
			-v $$PWD/$$dir/docs:/out \
			-v $$PWD/$$dir/proto:/protos \
			pseudomuto/protoc-gen-doc --doc_opt=html,docs.html; \
		docker run --rm \
			-v $$PWD/$$dir/docs:/out \
			-v $$PWD/$$dir/proto:/protos \
			pseudomuto/protoc-gen-doc --doc_opt=markdown,docs.md; \
	done
