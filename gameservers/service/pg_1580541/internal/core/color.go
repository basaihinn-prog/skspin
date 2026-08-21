package core

type Color int

const (
	C_Invalid Color = -1 + iota
	C_BaiDa         //百搭
	C_NoColor       //占位用 1没用
	C_YuanBao       //元宝
	C_RuYi          //如意
	C_FuDai         // 福袋
	C_HongBao       // 红包
	C_BianPao       //鞭炮
	C_Orange        //橘子
	CMax
)

var itemNames = [CMax]string{
	C_BaiDa:   "百搭",
	C_YuanBao: "元宝",
	C_RuYi:    "如意",
	C_FuDai:   "福袋",
	C_HongBao: "红包",
	C_BianPao: "鞭炮",
	C_Orange:  "橘子",
}
