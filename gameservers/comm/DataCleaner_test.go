package comm

import (
	"serve/comm/db"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// 针对主键有序的数据库
func TestClearData(t *testing.T) {
	type args struct {
		gname      string
		mongoAddr  string
		tableName  string
		filterRule bson.D
		bsonParam  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "修复fg数据",
			args: args{
				gname:     "pp_vs20piggybank",
				mongoAddr: "mongodb://myAdmin:myAdminPassword1@18.61.185.51:27017/?authSource=admin",
				//tableName: "simulate",
				tableName:  "simulate",
				filterRule: db.D("selected", true, "bucketid", -1, "times", bson.D{{"$gte", 0}}), //todo 需要根据不同的游戏修改 最终处理时使用的mongo语句
				bsonParam:  "rs_m",                                                               //todo 需要根据不同游戏修改这个bson中的 判断 是否为fg的字段
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ClearData(tt.args.gname, tt.args.mongoAddr, tt.args.tableName, tt.args.filterRule, tt.args.bsonParam)
		})
	}
}

// 针对主键无序的数据库
func TestClearData2(t *testing.T) {
	type args struct {
		gname      string
		mongoAddr  string
		tableName  string
		filterRule bson.D
		bsonParam  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "修复fg数据",
			args: args{
				gname:     "pp_vs20daydead",
				mongoAddr: "mongodb://myAdmin:myAdminPassword1@18.61.185.51:27017/?authSource=admin",
				tableName: "simulate",
				//tableName:  "simulate_copy1",
				filterRule: db.D("selected", true, "bucketid", -1, "times", bson.D{{"$gte", 0}}, "type", 1), //todo 需要根据不同的游戏修改
				bsonParam:  "rs_m",                                                                          //todo 需要根据不同游戏修改这个bson中的 判断 是否为fg的字段
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ClearData2(tt.args.gname, tt.args.mongoAddr, tt.args.tableName, tt.args.filterRule, tt.args.bsonParam)
		})
	}
}

//夺的方法

func Test_test1(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "111",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test1()
		})
	}
}
