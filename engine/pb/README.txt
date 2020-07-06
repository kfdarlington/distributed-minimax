To use protobufs:
(1) Ensure that protobuf is installed on your machine. For MacOS, use "brew install protobuf"

To use the protoc gRPC compiler to generate client and server code:
(1) Ensure that "protoc-gen-go-grpc" is installed in your Go bin with "go install google.golang.org/grpc/cmd/protoc-gen-go-grpc"
(2) Ensure that the Go bin directory is in your PATH
(3) In the root directory of this project, enter "protoc -I ./engine/pb/ --go-grpc_out=./engine/pb engine/pb/minimax.proto"

To use the protoc compiler to generate protobufs:
(1) Ensure that "" is installed in your Go bin with "go install google.golang.org/protobuf/cmd/protoc-gen-go"
(2) Ensure that the Go bin is in your PATH
(3) In the root directory of this project, enter "protoc -I ./engine/pb/ --go_out=./engine/pb ./engine/pb/minimax.proto"
