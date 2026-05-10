.PHONY: proto

PROTO_DIR=docs/protobuff
PROTO_FILE=$(PROTO_DIR)/auth-service.proto
GEN_DIR=gen

proto:
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
		$(PROTO_FILE)