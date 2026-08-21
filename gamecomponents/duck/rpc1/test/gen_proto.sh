


# download protoc at https://github.com/protocolbuffers/protobuf/releases
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

OUTDIR=./pb/

protoc \
	--go_out=$OUTDIR  \
	--go-grpc_out=$OUTDIR \
	--interface_out=$OUTDIR \
	--go-grpc_opt=require_unimplemented_servers=false \
	pb/*.proto

# 删除 json tag omitempty标记
sed -i "" 's/,omitempty//g' pb/pb/*.go
