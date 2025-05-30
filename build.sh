# shellcheck disable=SC2164
cd rpc/file && goctl rpc protoc file.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
cd ../user && goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
cd ../friend && goctl rpc protoc friend.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
cd ../message && goctl rpc protoc message.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
cd ../group && goctl rpc protoc group.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true