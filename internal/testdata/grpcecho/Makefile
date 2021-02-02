PROTO_DIR = .

protos: $(PROTO_DIR)
	for dir in $^ ; do protoc \
		-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I ${GOPATH}/src/github.com/golang/protobuf/ptypes/struct \
		--proto_path=. \
		--go_out=plugins=grpc,paths=source_relative:./ \
		$${dir}/*.proto \
		--grpc-gateway_out=paths=source_relative:./ \
	; done