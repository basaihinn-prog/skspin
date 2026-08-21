
# 管理proto文件的工具
# download buf at https://github.com/bufbuild/buf/releases

# 核心工具，后面的都是 protoc的插件
# download protoc at https://github.com/protocolbuffers/protobuf/releases

# 生成 proto 协议解析
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

# 生成 grpc service文件
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# 生成 ts文件
# sudo npm install ts-proto -g


# - proto 中使用 mongodb.ObjectID
# - proto 中使用 时间格式： TimeStamp

# - @gotags: bson:"_id" // 将在go结构体中添加tag
# - @ts: ignore()       // 此service 不生成到 ts
# - @ts: prefix(xxx)    // 此service ts请求路径前加上 xxx

rm -rf pb/gen

buf generate

rm -rf pb/gen/pb/duck
rm -rf pb/gen/ts1/pb/duck

node ./replace.js
