package main

// 商户类型
const (
	GENERALCONTROL = iota + 1 //总控
	LINEMERCHANTS             //线路商
	OPERATOR                  //运营商
)

// 钱包类型
const (
	SINGLEWALLET   = iota + 1 //单一钱包
	TRANSFERWALLET            //转账钱包
)

// 开关类
const (
	BucketOff = 1 //打开
	BucketOn  = 0 //关闭
)

// 游戏模式
const (
	BUCKET_HEARTBEAT = iota + 1
	BUCKET_FLUCTUATING
	BUCKET_SIMULATED
	BUCKET_MIXED
	BUCKET_STABLE
	BUCKET_HIGH
	BUCKET_EXTREMELY
)
