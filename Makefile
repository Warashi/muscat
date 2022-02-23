.PHONY:
protoc:
	protoc \
  	--go_out=. \
  	--go_opt=module=github.com/Warashi/muscat \
  	--go-grpc_out=. \
  	--go-grpc_opt=module=github.com/Warashi/muscat \
  	--proto_path=proto \
  	proto/*.proto
