package main

//
//func init() {
//	//RegMsgProc("/AdminInfo/GetOperatorGameList", "获取运营商报表数据", "AdminInfo", getOperatorReportData, getOperatorGameListParams{
//	//	AppID:     "",
//	//	Operator:  1,
//	//	PageIndex: 0,
//	//	PageSize:  0,
//	//})
//}
//
//type getOperatorGameListParams struct {
//	AppID     string `json:"AppID"`
//	Operator  int64
//	PageIndex int64
//	PageSize  int64
//	Type      string
//}
//
//func GetGameCofig(ctx *Context, ps getOperatorGameListParams, ret *respGameCofig) (err error) {
//	var operator comm.Operator_V2
//	err = CollAdminOperator.FindOne(bson.M{"AppID": ps.AppID}, &operator)
//	if err != nil {
//		return
//	}
//
//	// 获取游戏配置信息是否存在，如果不存在就创建，并设置默认值
//	var gconfig *comm.GameConfig
//	err = CollGameConfig.FindOne(bson.M{"AppID": ps.AppID}, &gconfig)
//	//创建默认数据
//	if errors.Is(err, mongo.ErrNoDocuments) {
//		var list_OperatorGameConfig []*operatorGameConfig
//		err = NewOtherDB("game").Collection("Games").FindAll(bson.M{}, &list_OperatorGameConfig)
//		if err != nil {
//			return err
//		}
//		for _, config := range list_OperatorGameConfig {
//			operatorConfig := bson.M{}
//			operatorConfig["AppID"] = ps.AppID
//			operatorConfig["RewardPercent"] = -900       //todo:有算法 根据RTP生成
//			operatorConfig["NoAwardPercent"] = 200       //todo:有算法 根据RTP生成
//			operatorConfig["GameId"] = config.GameId     //游戏编号
//			operatorConfig["ConfigPath"] = "config/path" //配置文件
//			operatorConfig["RTP"] = 96                   //RTP设置
//			operatorConfig["StopLoss"] = 0               //止盈止损开关
//			operatorConfig["MaxMultipleOff"] = 0         //赢取最高押注倍数
//			operatorConfig["MaxMultiple"] = 80.00
//			operatorConfig["BetBase"] = "0.1,0.2,0.3" //游戏投注
//			operatorConfig["GamePattern"] = 1         //游戏类型
//			operatorConfig["Preset"] = 10             //预设面额
//			operatorConfig["GameOn"] = 0              //游戏状态  0 开启  1 关闭
//			operatorConfig["FreeGameOff"] = 0         //免费游戏开关
//			err = CollGameConfig.InsertOne(operatorConfig)
//			if err != nil {
//				return err
//			}
//		}
//	}
//
//	//查询游戏列表
//	if ps.PageSize == 0 {
//		ps.PageSize = 2000
//	}
//	game_filter := mongodb.FindPageOpt{
//		Page:     ps.PageIndex,
//		PageSize: ps.PageSize,
//		Sort:     bson.M{"_id": 1},
//		Query:    bson.M{},
//	}
//	var list_OperatorGameConfig []*operatorGameConfig
//	count, err := NewOtherDB("game").Collection("Games").FindPage(game_filter, &list_OperatorGameConfig)
//	if err != nil {
//		return
//	}
//
//	//查询商户游戏配置
//	for _, gameConfig := range list_OperatorGameConfig {
//		gameid := gameConfig.GameId
//		filter := bson.M{"GameId": gameid, "AppID": ps.AppID}
//		sel := bson.M{"_id": 0}
//		//_ = NewOtherDB("GameConfig").Collection(operator.AppID).Coll().FindOne(context.T/ODO(), filter, options.FindOne().SetProjection(sel)).Decode(&gameConfig)
//		_ = CollGameConfig.Coll().FindOne(context.TODO(), filter, options.FindOne().SetProjection(sel)).Decode(&gameConfig)
//		gameConfig.StopLossOff = operator.StopLoss
//		gameConfig.RTPOff = operator.RTPOff
//		gameConfig.MaxMultipleOff = operator.MaxMultipleOff
//	}
//
//	ret.List = list_OperatorGameConfig
//	ret.CountAll = count
//	return
//}
