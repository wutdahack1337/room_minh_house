DOCKER=docker
protoVer=0.12.1
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR)/:/workspace --workdir /workspace $(protoImageName)

.PHONY: proto
proto:
	@$(protoImage) sh ./scripts/protoc-gen-go.sh