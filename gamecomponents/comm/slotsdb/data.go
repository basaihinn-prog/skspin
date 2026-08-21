package slotsdb

type SubDataType int

const (
	SubDateType_Invalid     SubDataType = iota
	SubDateType_TableDefine             // 记录表名所属ID。（记录里用到表的ID）
	SubDateType_Record                  // 有效的记录
	SubDateType_Delete                  // 标记被删除的哪张表的第几条记录。
	SubDateType_Drop                    // 标记哪张表被删除了。
	SubDateType_Truncate                // 清空指定表
)

type DataRoot struct {
	Type     SubDataType
	TableID  int
	RowID    int
	DataSize int
}

// type SubDataTableDefine struct {
// 	TableName string
// }

// // 各个游戏自定义
// type SubDataRecord struct {
// 	RowID int
// 	Body  []byte // 盘的数据
// }

// type SubDataDelete struct {
// 	RowID int // 指定表的第几条记录
// }

// type SubDataDrop struct {
// 	Nothing int
// }

// type SubDataTruncate struct {
// 	Nothing int
// }
