package gendata

//func init() {
//	mux.RegRpc(fmt.Sprintf("/%s/AdminSlotsRpc/combine/gen", internal.GameID), "生成组合数据", "gendata", genCombineData, GetCombine()).SetOnlyDev()
//}
//
//func genCombineData(ps jdbcomm.Combine, ret *mux.EmptyResult) (err error) {
//	err = syncCombine(ps)
//	if err != nil {
//		return
//	}
//
//	err = LoadCombineData()
//	return err
//}
//
//func syncCombine(ps jdbcomm.Combine) error {
//	coll := db.Collection("combine")
//	var models []mongo.WriteModel
//	var md mongo.WriteModel
//	for _, c := range ps {
//		md = mongo.NewUpdateOneModel().SetFilter(db.D("_id", c.ID)).SetUpdate(db.D("$set", db.D("count", c.Count)))
//		models = append(models, md)
//	}
//
//	if len(models) > 0 {
//		_, err := coll.BulkWrite(context.TODO(), models)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
