all: pb

pb: *.proto
	protoc --proto_path=$(GOPATH)/src:. --go_out=plugins=grpc:. token.proto

.PHONY: pb
