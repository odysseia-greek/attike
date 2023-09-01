generate_docs:
	docker run --rm \
		-v $(CURDIR)/docs:/out \
		-v $(CURDIR)/proto:/protos \
		pseudomuto/protoc-gen-doc --doc_opt=markdown,docs.md