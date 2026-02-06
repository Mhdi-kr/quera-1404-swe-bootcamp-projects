protoc \
  --proto_path=protos \
  --go_out=protos-gen/ --go_opt=paths=source_relative \
  --go-grpc_out=protos-gen/ --go-grpc_opt=paths=source_relative \
  protos/post/v1/post.proto