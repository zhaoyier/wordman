clean:
	@rm -rf rpc/*

start:
	@go run main

gen:
	@protoc -I apidoc apidoc/wordman/*.proto --go_out=plugins=grpc:rpc/