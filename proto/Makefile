protoc-user:
	protoc \
	  --proto_path=. \
	  --proto_path=./google/api \
	  --proto_path=$(go list -f '{{ .Dir }}' -m github.com/grpc-ecosystem/grpc-gateway/v2) \
	  --go_out=paths=source_relative:. \
	  --go_opt=Mprotoc-gen-openapiv2/options/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options \
	  --go-grpc_out=paths=source_relative:. \
	  --go-grpc_opt=Mprotoc-gen-openapiv2/options/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options \
	  --grpc-gateway_out=paths=source_relative:. \
	  --grpc-gateway_opt=Mprotoc-gen-openapiv2/options/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options \
	  --openapiv2_out=. \
	  --openapiv2_opt=allow_merge=true,merge_file_name=api,\
	Mprotoc-gen-openapiv2/options/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options \
	  user.proto