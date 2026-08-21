package main

import (
	"context"
	"errors"
	"fmt"
	"game/duck/lazy"
	"game/duck/mongodb"
	"game/pb/_gen/pb"
	"game/pb/_gen/pb/adminpb"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type AdminConfigFile struct {
}

func (AdminConfigFile) SaveConfig(ctx context.Context, req *adminpb.FileContent) (*pb.Empty, error) {
	ss := GetSession(ctx)

	watcher := lazy.ConfigManager.Watcher

	oldhash := ""
	if watcher.IsExists(req.FileName) {
		bytes, err := watcher.ReadFile(req.FileName)
		if err != nil {
			return nil, err
		}

		oldhash = Md5(string(bytes))
	}

	hash := Md5(req.Content)

	// 提交的文件和当前文件md5一样，则没必要保存。
	if oldhash == hash {
		return &pb.Empty{}, nil
	}

	dbFileLog := adminpb.DBConfigFileLog{
		Id:       mongodb.NewObjectID(),
		Content:  req.Content,
		FileName: req.FileName,
		CreateAt: mongodb.NewTimeStamp(time.Now()),
		User:     ss.Pid,
	}

	err := CollAdminConfigFileLog.InsertOne(&dbFileLog)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(req.Content, "\r\n") {
		req.Content = fmt.Sprintf("%s\r\n", req.Content)
	}
	err = watcher.WriteFile(req.FileName, []byte(req.Content))
	if err != nil {
		return nil, err
	}

	// 等配置监听者处理完成
	time.Sleep(300 * time.Millisecond)

	errStr := lazy.ConfigManager.GetErr(req.FileName, hash)
	if errStr != "" {
		return nil, errors.New(errStr)
	}

	return &pb.Empty{}, nil
}

func (AdminConfigFile) LoadConfig(ctx context.Context, req *adminpb.FileContent) (*adminpb.FileContent, error) {

	watcher := lazy.ConfigManager.Watcher

	if !watcher.IsExists(req.FileName) {
		return &adminpb.FileContent{FileName: req.FileName}, nil
	}

	content, err := watcher.ReadFile(req.FileName)
	if err != nil {
		return nil, err
	}

	return &adminpb.FileContent{FileName: req.FileName, Content: string(content)}, nil
}

// 获取文件历史记录
func (AdminConfigFile) FileHistory(ctx context.Context, req *adminpb.FileHistoryReq) (*adminpb.FileHistoryResp, error) {
	var list []*adminpb.DBConfigFileLog
	filter := mongodb.FindPageOpt{
		Page:     1,
		PageSize: int64(req.Count),
		Sort:     bson.M{"_id": -1}, //可以为空
		Query:    bson.M{"FileName": req.FileName},
	}
	_, err := CollAdminConfigFileLog.FindPage(filter, &list)
	if err != nil {
		return nil, err
	}
	return &adminpb.FileHistoryResp{List: list}, nil
}
