package models

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Player struct {
	Account    string  //账号 (手机格式)
	PlayerId   string  //用户id
	Password   string  //用户密码
	PlayerGold float64 //用户金币数量
	UId        string  //用户Uid
	UniqueId   string  //用户唯一id
	AvatarUrl  string  //用户的头像
	VipLevel   int32   //用户等级
	UserName   string  //用户名
	Type       string
}

func (p Player) ProtoReflect() protoreflect.Message {
	//TODO implement me
	panic("implement me")
}
