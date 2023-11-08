protoc -I./protos --go_out=./server \
 --go_opt=paths=source_relative --go-grpc_out=./server \
  --go-grpc_opt=paths=source_relative server.proto