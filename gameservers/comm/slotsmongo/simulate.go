package slotsmongo

import (
	"fmt"
	"reflect"

	"serve/comm/ut"

	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SimulateData struct {
	Id       primitive.ObjectID `bson:"_id"`
	DropPan  []bson.M           `bson:"droppan"`
	HasGame  bool               `bson:"hasgame"`
	Times    float64            `bson:"times"`
	BucketId int                `json:"bucketid"`
	Type     int                `bson:"type"`
	Selected bool               `bson:"selected"`
}

// 最小下注是1
func (sd *SimulateData) Deal(num int, balance, cs, ml, line, buyMul float64, isBuy bool) map[string]any {
	bet := cs * ml * line
	pan := sd.DropPan[0]

	delete(pan, "orignid")
	pan["ml"] = ml
	pan["cs"] = cs
	mul := float64(1)
	if isBuy {
		mul = buyMul
	}

	ut.MKMul(pan, "aw", bet)
	ut.MKMul(pan, "actw", bet)
	ut.MKMul(pan, "saw", bet)
	ut.MKMul(pan, "tlw", bet)
	ut.MKMul(pan, "lwa", bet)
	ut.MKMul(pan, "afw", bet)
	ut.MKMul(pan, "acw", bet)
	ut.MKMul(pan, "ctw", bet)
	ut.MKMul(pan, "btw", bet)
	ut.MKMul(pan, "cwc", bet)
	ut.MKMul(pan, "tw", bet)
	ut.MKMul(pan, "np", bet)
	ut.MKMul(pan, "tbb", bet)
	ut.MKMul(pan, "ssaw", bet)
	ut.MKMul(pan, "stw", bet)

	ut.MKMul(pan, "ptw", bet) //

	ut.MKMul(pan, "esw", bet)

	ut.MKMul(pan, "slw", bet)
	ut.MKMul(pan, "slwm", bet)
	ut.MKMul(pan, "stwm", bet)

	ut.MKMul(pan, "otw", bet)

	ut.MKMul(pan, "trtw", bet)
	ut.MKMul(pan, "trtwbm", bet)

	ut.MKMul(pan, "rtw", bet)
	ut.MKMul(pan, "rtwbm", bet)

	ut.MKMul(pan, "twbm", bet)
	ut.MKMul(pan, "bmtw", bet)
	ut.MKMul(pan, "twb", bet)
	ut.MK2Mul(pan, "fs", "aw", bet)
	ut.MK2Mul(pan, "fs", "caw", bet)
	ut.MK2Mul(pan, "fs", "lwa", bet)
	ut.MK2Mul(pan, "fs", "twbm", bet)
	ut.MK2Mul(pan, "fs", "ssaw", bet)
	ut.MK2Mul(pan, "fs", "wa", bet)

	//
	ut.MK2Mul(pan, "fs", "slwm", bet)
	//

	ut.MK2Mul(pan, "fs", "fsaw", bet)

	ut.MK2Mul(pan, "fs", "tb", bet)

	ut.MK2Mul(pan, "rs", "aw", bet)
	ut.MK2Mul(pan, "rs", "bw", bet)

	ut.MK2Mul(pan, "rs", "wa", bet)

	ut.MK2Mul(pan, "bns", "twbm", bet)

	ut.MK2Mul(pan, "bns", "aw", bet)

	ut.MK2Mul(pan, "fs", "slw", bet)

	ut.MK2Mul(pan, "fs", "stw", bet)
	ut.MK2Mul(pan, "fs", "stwm", bet)
	cptw := pan["cptw"]
	if cptw != nil {
		rt := reflect.TypeOf(cptw)
		if rt.Kind() == reflect.Slice {
			ra := cptw.(primitive.A)
			for i := range ra {
				ra[i] = ut.Round6(ut.GetFloat(ra[i]) * bet)
			}
			pan["cptw"] = ra
		} else {
			ut.MKMul(pan, "cptw", bet)
		}
	}

	lw := pan["lw"]
	if lw != nil {
		lwm := lw.(bson.M)
		for k, v := range lwm {
			t := reflect.TypeOf(v).Kind()
			if t == reflect.Slice {
				av := v.(primitive.A)
				for i := range av {
					av[i] = ut.Round6(ut.GetFloat(av[i]) * bet)
				}
				lwm[k] = av
			} else {
				lwm[k] = ut.Round6(ut.GetFloat(v) * bet)
			}
		}
		pan["lw"] = lwm
	}
	lwm := pan["lwm"]
	if lwm != nil {
		lwmm := lwm.(bson.M)
		for k, v := range lwmm {
			lwmm[k] = ut.Round6(ut.GetFloat(v) * bet)
		}
		pan["lwm"] = lwmm
	}

	lwbm := pan["lwbm"]
	if lwbm != nil {
		lwbmm := lwbm.(bson.M)
		for k, v := range lwbmm {
			lwbmm[k] = ut.Round6(ut.GetFloat(v) * bet)
		}
		pan["lwbm"] = lwbmm
	}

	sw := pan["sw"]
	if sw != nil {
		a := reflect.TypeOf(sw).Kind()
		if a == reflect.Float64 {
			ut.MKMul(pan, "sw", bet)
		} else {
			swm := sw.(bson.M)
			for k, v := range swm {
				t := reflect.TypeOf(v).Kind()
				if t == reflect.Float64 {
					swm[k] = ut.Round6(ut.GetFloat(v) * bet)
				} else if t == reflect.Slice {
					av := v.(primitive.A)
					for i := range av {
						av[i] = ut.Round6(ut.GetFloat(av[i]) * bet)
					}
					swm[k] = av
				} else if t == reflect.Map {
					vm := v.(bson.M)
					ut.MKMul(vm, "wa", bet)
					swm[k] = vm
				}
			}
			pan["sw"] = swm
		}

	}

	//
	fss := pan["fs"]
	if fss != nil {
		fssm, ok := fss.(bson.M)
		if ok {
			for k, v := range fssm {
				if k == "lw" {
					if v != nil {
						vm := v.(bson.M)
						for k, v := range vm {
							vm[k] = ut.Round6(ut.GetFloat(v) * bet)
						}
						//pan["lw"] = vm
					}
				}
			}
			pan["fs"] = fssm
		}
	}

	//

	cpf := pan["cpf"]
	if cpf != nil {
		cpfm := cpf.(bson.M)
		for k, v := range cpfm {
			vm := v.(bson.M)
			ut.MKMul(vm, "bv", bet)
			cpfm[k] = vm
		}
	}

	cp := pan["cp"]
	if cp != nil {
		cpm := cp.(bson.M)
		for k, v := range cpm {
			vk := reflect.TypeOf(v).Kind()
			if vk == reflect.Slice {
				av := v.(primitive.A)
				for kk, vv := range av {
					avv := vv.(primitive.A)
					for kkk := range avv {
						avv[kkk] = ut.Round6((ut.GetFloat(avv[kkk]) * bet))
					}
					av[kk] = avv
				}
				cpm[k] = av
			} else if vk == reflect.Map {
				ctwm := v.(bson.M)
				for kk, vv := range ctwm {
					if _, ok := vv.(primitive.M); ok {
						vvm := vv.(primitive.M)
						for _, vvv := range vvm {
							avv := vvv.(primitive.A)
							for kkk := range avv {
								avv[kkk] = ut.Round6((ut.GetFloat(avv[kkk]) * bet))
							}
						}
						ctwm[kk] = vvm
					}
				}
			}
		}
	}

	trlw := pan["trlw"]
	if trlw != nil {
		trlww := trlw.(bson.M)
		for k, v := range trlww {
			trlww[k] = ut.Round6(ut.GetFloat(v) * bet)
		}
		pan["trlw"] = trlww
	}

	//for pg_70
	pft := pan["pft"]
	if pft != nil {
		t := reflect.TypeOf(pft).Kind()
		if t == reflect.Slice {

		} else {
			pftm := pft.(bson.M)

			wa := pftm["wa"]
			wa_arr := wa.(primitive.A)
			for k, v := range wa_arr {
				wa_arr[k] = ut.Round6(ut.GetFloat(v) * bet)
			}
			pftm["wa"] = wa_arr
			pan["pft"] = pftm
		}
	}

	//for pg_1529867
	mlw := pan["mlw"]
	if mlw != nil {
		mlww := mlw.(bson.M)
		for k, v := range mlww {
			mlww[k] = ut.Round6(ut.GetFloat(v) * bet)
		}
		pan["mlw"] = mlww
	}

	//--
	rv := pan["rv"]
	if rv != nil {
		rvv := rv.(bson.A)
		for k, v := range rvv {
			if k >= 3 && k <= 5 {
				rvv[k] = v
			} else {
				rvv[k] = ut.GetFloat(v) * bet
			}
		}
		pan["rv"] = rvv
	}

	orv := pan["orv"]
	if orv != nil {
		orvv := orv.(bson.A)
		for k, v := range orvv {
			if k >= 3 && k <= 5 {
				orvv[k] = v
			} else {
				orvv[k] = ut.GetFloat(v) * bet
			}
		}
		pan["orv"] = orvv
	}

	rsrv := pan["rsrv"]
	if rsrv != nil {
		rsrvv := rsrv.(bson.A)
		for k, v := range rsrvv {
			if k >= 3 && k <= 5 {
				rsrvv[k] = v
			} else {
				rsrvv[k] = ut.GetFloat(v) * bet
			}
		}
		pan["rsrv"] = rsrvv
	}

	// for pg_24
	bn := pan["bn"]
	if bn != nil {
		mul = ut.GetFloat(bn)
	}

	if num == 0 {
		pan["tb"] = ut.Round6(bet * mul)
		pan["blb"] = balance + ut.Round6(bet*mul) //下注前的金额
		pan["blab"] = balance                     //下注后的金额
	} else {
		pan["tb"] = 0
		pan["blb"] = balance
		pan["blab"] = balance
	}
	pan["bl"] = balance + ut.Round6(ut.GetFloat(pan["tw"])) //转动后的余额

	if pan["hashr"] != nil {
		hashr := pan["hashr"].(string)
		pan["hashr"] = caclHashr(hashr, cs, ml, line)
	}

	key := "rid"
	_, ok := pan[key]
	if ok {
		delete(pan, key)
	}
	return pan
}

// 针对最小下注不是1的
func (sd *SimulateData) Deal2(num int, balance, cs, ml, line, buyMul, originMin float64, isBuy bool) map[string]any {
	bet := cs * ml * line
	pan := sd.DropPan[0]

	delete(pan, "orignid")
	pan["ml"] = ml
	pan["cs"] = cs
	mul := float64(1)
	if isBuy {
		mul = buyMul
	}

	ut.MKMul(pan, "atw", bet)
	ut.MKMul(pan, "otw", bet)
	ut.MKMul(pan, "aw", bet)
	ut.MKMul(pan, "ctw", bet)
	ut.MKMul(pan, "afw", bet)
	ut.MKMul(pan, "cwc", bet)
	ut.MKMul(pan, "saw", bet)
	ut.MKMul(pan, "tlw", bet)
	ut.MKMul(pan, "actw", bet)
	ut.MKMul(pan, "tw", bet)
	ut.MKMul(pan, "np", bet)
	ut.MKMul(pan, "tbb", bet)
	ut.MKMul(pan, "ssaw", bet)
	ut.MKMul(pan, "twbm", bet)
	ut.MKMul(pan, "ltw", bet)

	ut.MKMul(pan, "ptw", 0)

	ut.MK2Mul(pan, "fs", "aw", bet)
	ut.MK2Mul(pan, "fs", "caw", bet)
	ut.MK2Mul(pan, "fs", "lwa", bet)
	ut.MK2Mul(pan, "fs", "ssaw", bet)
	ut.MK2Mul(pan, "fs", "twbm", bet)
	ut.MK2Mul(pan, "bns", "bac", bet)
	ut.MK2Mul(pan, "bns", "aw", bet)
	ut.MK2Mul(pan, "cf", "acp", bet)
	ut.MK2Mul(pan, "cf", "scp", bet)
	ut.MK2Mul(pan, "cf", "tcp", bet)
	ut.MK2Mul(pan, "bf", "bsp", bet)
	ut.MK2Mul(pan, "lf", "lsp", bet)
	ut.MK3Mul(pan, "fs", "wf", "wa", bet)
	ut.MK3Mul(pan, "bf", "bpz", "0", bet)
	ut.MK3Mul(pan, "bf", "bpz", "1", bet)
	ut.MK3Mul(pan, "lf", "lpz", "0", bet)
	ut.MK3Mul(pan, "lf", "lpz", "1", bet)

	ut.MK1Div(pan, "atw", originMin)
	ut.MK1Div(pan, "otw", originMin)
	ut.MK1Div(pan, "aw", originMin)
	ut.MK1Div(pan, "ctw", originMin)
	ut.MK1Div(pan, "afw", originMin)
	ut.MK1Div(pan, "cwc", originMin)
	ut.MK1Div(pan, "saw", originMin)
	ut.MK1Div(pan, "tlw", originMin)
	ut.MK1Div(pan, "actw", originMin)
	ut.MK1Div(pan, "tw", originMin)
	ut.MK1Div(pan, "np", originMin)
	ut.MK1Div(pan, "tbb", originMin)
	ut.MK1Div(pan, "ssaw", originMin)
	ut.MK1Div(pan, "twbm", originMin)
	ut.MK1Div(pan, "ltw", originMin)
	ut.MK2Div(pan, "fs", "aw", originMin)
	ut.MK2Div(pan, "fs", "caw", originMin)
	ut.MK2Div(pan, "fs", "lwa", originMin)
	ut.MK2Div(pan, "fs", "ssaw", originMin)
	ut.MK2Div(pan, "fs", "twbm", originMin)
	ut.MK2Div(pan, "bns", "bac", originMin)
	ut.MK2Div(pan, "bns", "aw", originMin)
	ut.MK2Div(pan, "cf", "acp", originMin)
	ut.MK2Div(pan, "cf", "scp", originMin)
	ut.MK2Div(pan, "cf", "tcp", originMin)
	ut.MK2Div(pan, "bf", "bsp", originMin)
	ut.MK2Div(pan, "lf", "lsp", originMin)
	ut.MK3Div(pan, "fs", "wf", "wa", originMin)
	ut.MK3Div(pan, "bf", "bpz", "0", originMin)
	ut.MK3Div(pan, "bf", "bpz", "1", originMin)
	ut.MK3Div(pan, "lf", "lpz", "0", originMin)
	ut.MK3Div(pan, "lf", "lpz", "1", originMin)

	lw := pan["lw"]
	if lw != nil {
		lwm := lw.(bson.M)
		for k, v := range lwm {
			vt := reflect.TypeOf(v)
			if vt.Kind() == reflect.Slice {
				varr := v.(primitive.A)
				for i := range varr {
					varr[i] = ut.Round6(ut.GetFloat(varr[i]) * bet / originMin)
				}
				lwm[k] = v
			} else {
				lwm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
			}
		}
		pan["lw"] = lwm
	}

	gaw := pan["gaw"]
	if gaw != nil {
		gawm := gaw.(bson.M)
		for k, v := range gawm {
			gawm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)

		}
		pan["gaw"] = gawm
	}

	bm := pan["bm"]
	if bm != nil {
		bmm := bm.(bson.M)
		bmw := bmm["bmw"]
		if bmw != nil {
			bmwm := bmw.(bson.M)
			for k, v := range bmwm {
				bmwm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
			}
			bmm["bmw"] = bmwm
		}
	}

	sw := pan["sw"]
	if sw != nil {
		swm := sw.(bson.M)
		for k, v := range swm {
			swm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
		}
		pan["sw"] = swm
	}
	cpf := pan["cpf"]
	if cpf != nil {
		vt := reflect.TypeOf(cpf)
		if vt.Kind() == reflect.Slice {
			varr := cpf.(primitive.A)
			for i := range varr {
				sfc := varr[i].(bson.M)
				ut.MKMul(sfc, "cp", bet)
				ut.MKMul(sfc, "aw", bet)
				ut.MK1Div(sfc, "aw", originMin)
				ut.MK1Div(sfc, "cp", originMin)
			}
		} else {
			cpfm := cpf.(bson.M)
			for k, v := range cpfm {
				vm := v.(bson.M)
				ut.MKMul(vm, "bv", bet)
				cpfm[k] = vm
			}
		}
	}

	bns := pan["bns"]
	if bns != nil {
		bnsm := bns.(bson.M)

		pwd := bnsm["pwd"]
		if pwd != nil {
			pwdm := pwd.(bson.M)
			ut.MKMul(pwdm, "paw", bet/originMin)

			plw := pwdm["plw"]
			if plw != nil {
				plwm := plw.(bson.M)
				for k, v := range plwm {
					varr := v.(primitive.A)
					for i := range varr {
						varr[i] = ut.Round6(ut.GetFloat(varr[i]) * bet / originMin)
					}
					plwm[k] = v
				}
			}

		}

		pl := bnsm["pl"]
		if pl != nil {
			plm := pl.(bson.M)
			for k, v := range plm {
				varr := v.(primitive.A)
				for i := range varr {
					varr[i] = ut.Round6(ut.GetFloat(varr[i]) * bet / originMin)
				}

				plm[k] = v
			}
		}
	}

	fs := pan["fs"]
	if fs != nil {
		fsm := fs.(bson.M)
		fs_lw := fsm["lw"]
		if fs_lw != nil {
			fs_lwm := fs_lw.(bson.M)
			for k, v := range fs_lwm {
				fs_lwm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
			}
			fsm["lw"] = fs_lwm
		}

		fs_slw := fsm["slw"]
		if fs_slw != nil {
			slw_arr := fs_slw.(primitive.A)
			for k, v := range slw_arr {
				slw_arr[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
			}
			fsm["slw"] = slw_arr
		}
		pan["fs"] = fsm
	}

	slw := pan["slw"]
	if slw != nil {
		slw_arr := slw.(primitive.A)
		for k, v := range slw_arr {
			slw_arr[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
		}
		pan["slw"] = slw_arr
	}

	if num == 0 {
		pan["tb"] = ut.Round6(bet * mul)
		pan["blb"] = balance + ut.Round6(bet*mul) //下注前的金额
		pan["blab"] = balance                     //下注后的金额
	} else {
		pan["tb"] = 0
		pan["blb"] = balance
		pan["blab"] = balance
	}
	pan["bl"] = balance + ut.Round6(ut.GetFloat(pan["tw"]))
	if pan["hashr"] != nil {
		hashr := pan["hashr"].(string)
		pan["hashr"] = caclHashr(hashr, cs, ml, line)
	}
	key := "rid"
	_, ok := pan[key]
	if ok {
		delete(pan, key)
	}
	return pan
}

func caclHashr(hashr string, cs, ml, line float64) string {
	args := strings.Split(hashr, "#")
	for i := 0; i < len(args); i++ {
		if (args[i] == "MV" || args[i] == "MG") && i < len(args)-1 {
			data, _ := strconv.ParseFloat(args[i+1], 64)
			data = data * cs * ml * line
			args[i+1] = fmt.Sprintf("%.1f", data)
		}
	}
	return strings.Join(args, "#")
}

// pg_1717688 swm特殊处理
func (sd *SimulateData) Deal3(num int, balance, cs, ml, line, buyMul, originMin float64, isBuy bool) map[string]any {
	bet := cs * ml * line
	pan := sd.DropPan[0]

	delete(pan, "orignid")
	pan["ml"] = ml
	pan["cs"] = cs
	mul := float64(1)
	if isBuy {
		mul = buyMul
	}

	ut.MKMul(pan, "aw", bet)
	ut.MKMul(pan, "ctw", bet)
	ut.MKMul(pan, "afw", bet)
	ut.MKMul(pan, "cwc", bet)
	ut.MKMul(pan, "saw", bet)
	ut.MKMul(pan, "tlw", bet)
	ut.MKMul(pan, "actw", bet)
	ut.MKMul(pan, "tw", bet)
	ut.MKMul(pan, "np", bet)
	ut.MKMul(pan, "tbb", bet)
	ut.MKMul(pan, "ssaw", bet)
	ut.MKMul(pan, "twbm", bet)
	ut.MKMul(pan, "ltw", bet)

	ut.MKMul(pan, "ptw", 0)

	ut.MK2Mul(pan, "fs", "aw", bet)
	ut.MK2Mul(pan, "fs", "caw", bet)
	ut.MK2Mul(pan, "fs", "lwa", bet)
	ut.MK2Mul(pan, "fs", "ssaw", bet)
	ut.MK2Mul(pan, "fs", "twbm", bet)
	ut.MK2Mul(pan, "bns", "bac", bet)
	ut.MK2Mul(pan, "bns", "aw", bet)
	ut.MK2Mul(pan, "cf", "acp", bet)
	ut.MK2Mul(pan, "cf", "scp", bet)
	ut.MK2Mul(pan, "cf", "tcp", bet)
	ut.MK2Mul(pan, "bf", "bsp", bet)
	ut.MK2Mul(pan, "lf", "lsp", bet)
	ut.MK3Mul(pan, "fs", "wf", "wa", bet)
	ut.MK3Mul(pan, "bf", "bpz", "0", bet)
	ut.MK3Mul(pan, "bf", "bpz", "1", bet)
	ut.MK3Mul(pan, "lf", "lpz", "0", bet)
	ut.MK3Mul(pan, "lf", "lpz", "1", bet)

	ut.MK1Div(pan, "aw", originMin)
	ut.MK1Div(pan, "ctw", originMin)
	ut.MK1Div(pan, "afw", originMin)
	ut.MK1Div(pan, "cwc", originMin)
	ut.MK1Div(pan, "saw", originMin)
	ut.MK1Div(pan, "tlw", originMin)
	ut.MK1Div(pan, "actw", originMin)
	ut.MK1Div(pan, "tw", originMin)
	ut.MK1Div(pan, "np", originMin)
	ut.MK1Div(pan, "tbb", originMin)
	ut.MK1Div(pan, "ssaw", originMin)
	ut.MK1Div(pan, "twbm", originMin)
	ut.MK1Div(pan, "ltw", originMin)
	ut.MK2Div(pan, "fs", "aw", originMin)
	ut.MK2Div(pan, "fs", "caw", originMin)
	ut.MK2Div(pan, "fs", "lwa", originMin)
	ut.MK2Div(pan, "fs", "ssaw", originMin)
	ut.MK2Div(pan, "fs", "twbm", originMin)
	ut.MK2Div(pan, "bns", "bac", originMin)
	ut.MK2Div(pan, "bns", "aw", originMin)
	ut.MK2Div(pan, "cf", "acp", originMin)
	ut.MK2Div(pan, "cf", "scp", originMin)
	ut.MK2Div(pan, "cf", "tcp", originMin)
	ut.MK2Div(pan, "bf", "bsp", originMin)
	ut.MK2Div(pan, "lf", "lsp", originMin)
	ut.MK3Div(pan, "fs", "wf", "wa", originMin)
	ut.MK3Div(pan, "bf", "bpz", "0", originMin)
	ut.MK3Div(pan, "bf", "bpz", "1", originMin)
	ut.MK3Div(pan, "lf", "lpz", "0", originMin)
	ut.MK3Div(pan, "lf", "lpz", "1", originMin)

	lw := pan["lw"]
	if lw != nil {
		lwm := lw.(bson.M)
		for k, v := range lwm {
			vt := reflect.TypeOf(v)
			if vt.Kind() == reflect.Slice {
				varr := v.(primitive.A)
				for i := range varr {
					varr[i] = ut.Round6(ut.GetFloat(varr[i]) * bet / originMin)
				}
				lwm[k] = v
			} else {
				lwm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
			}
		}
		pan["lw"] = lwm
	}

	gaw := pan["gaw"]
	if gaw != nil {
		gawm := gaw.(bson.M)
		for k, v := range gawm {
			gawm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)

		}
		pan["gaw"] = gawm
	}

	bm := pan["bm"]
	if bm != nil {
		bmm := bm.(bson.M)
		bmw := bmm["bmw"]
		if bmw != nil {
			bmwm := bmw.(bson.M)
			for k, v := range bmwm {
				bmwm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
			}
			bmm["bmw"] = bmwm
		}
	}

	sw := pan["sw"]
	if sw != nil {
		// swm := sw.(bson.M)
		// for k, v := range swm {
		// 	t := reflect.TypeOf(v).Kind()
		// 	if t == reflect.Float64 {
		// 		swm[k] = ut.Round6(ut.GetFloat(v) * bet)
		// 	} else if t == reflect.Slice {
		// 		av := v.(primitive.A)
		// 		for i := range av {
		// 			av[i] = ut.Round6(ut.GetFloat(av[i]) * bet)
		// 		}
		// 		swm[k] = av
		// 	} else if t == reflect.Map {
		// 		vm := v.(bson.M)
		// 		ut.MKMul(vm, "wa", bet)
		// 		swm[k] = vm
		// 	}
		// }
		// pan["sw"] = swm
	}
	cpf := pan["cpf"]
	if cpf != nil {
		vt := reflect.TypeOf(cpf)
		if vt.Kind() == reflect.Slice {
			varr := cpf.(primitive.A)
			for i := range varr {
				sfc := varr[i].(bson.M)
				ut.MKMul(sfc, "cp", bet)
				ut.MKMul(sfc, "aw", bet)
				ut.MK1Div(sfc, "aw", originMin)
				ut.MK1Div(sfc, "cp", originMin)
			}
		} else {
			cpfm := cpf.(bson.M)
			for k, v := range cpfm {
				vm := v.(bson.M)
				ut.MKMul(vm, "bv", bet)
				cpfm[k] = vm
			}
		}
	}

	bns := pan["bns"]
	if bns != nil {
		bnsm := bns.(bson.M)

		pwd := bnsm["pwd"]
		if pwd != nil {
			pwdm := pwd.(bson.M)
			ut.MKMul(pwdm, "paw", bet/originMin)

			plw := pwdm["plw"]
			if plw != nil {
				plwm := plw.(bson.M)
				for k, v := range plwm {
					varr := v.(primitive.A)
					for i := range varr {
						varr[i] = ut.Round6(ut.GetFloat(varr[i]) * bet / originMin)
					}
					plwm[k] = v
				}
			}

		}

		pl := bnsm["pl"]
		if pl != nil {
			plm := pl.(bson.M)
			for k, v := range plm {
				varr := v.(primitive.A)
				for i := range varr {
					varr[i] = ut.Round6(ut.GetFloat(varr[i]) * bet / originMin)
				}

				plm[k] = v
			}
		}
	}

	fs := pan["fs"]
	if fs != nil {
		fsm := fs.(bson.M)
		fs_lw := fsm["lw"]
		if fs_lw != nil {
			fs_lwm := fs_lw.(bson.M)
			for k, v := range fs_lwm {
				fs_lwm[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
			}
			fsm["lw"] = fs_lwm
		}

		fs_slw := fsm["slw"]
		if fs_slw != nil {
			slw_arr := fs_slw.(primitive.A)
			for k, v := range slw_arr {
				slw_arr[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
			}
			fsm["slw"] = slw_arr
		}
		pan["fs"] = fsm
	}

	slw := pan["slw"]
	if slw != nil {
		slw_arr := slw.(primitive.A)
		for k, v := range slw_arr {
			slw_arr[k] = ut.Round6(ut.GetFloat(v) * bet / originMin)
		}
		pan["slw"] = slw_arr
	}

	if num == 0 {
		pan["tb"] = ut.Round6(bet * mul)
		pan["blb"] = balance + ut.Round6(bet*mul) //下注前的金额
		pan["blab"] = balance                     //下注后的金额
	} else {
		pan["tb"] = 0
		pan["blb"] = balance
		pan["blab"] = balance
	}
	pan["bl"] = balance + ut.Round6(ut.GetFloat(pan["tw"]))
	if pan["hashr"] != nil {
		hashr := pan["hashr"].(string)
		pan["hashr"] = caclHashr(hashr, cs, ml, line)
	}
	key := "rid"
	_, ok := pan[key]
	if ok {
		delete(pan, key)
	}
	return pan
}
