run:
	cd rpc/file && goctl rpc protoc file.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
	cd rpc/user && goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
	cd rpc/friend && goctl rpc protoc friend.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
	cd rpc/message && goctl rpc protoc message.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
	cd rpc/group && goctl rpc protoc group.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
	cd rpc/notify && goctl rpc protoc notify.proto --go_out=. --go-grpc_out=. --zrpc_out=. --client=true
	docker compose up -d --build