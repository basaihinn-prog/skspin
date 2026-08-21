package gendata

import (
	"encoding/json"

	"serve/comm/slotsmongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MFMap struct {
	MT any `json:"mt"`
	MS any `json:"ms"`
	MI any `json:"mi"`
}

func (m MFMap) MarshalJSON() ([]byte, error) {
	mt_v, _ := json.Marshal(m.MT)
	ms_v, _ := json.Marshal(m.MS)
	mi_v, _ := json.Marshal(m.MI)
	return []byte(`{"mt":` + string(mt_v) + `,"ms":` + string(ms_v) + `,"mi":` + string(mi_v) + `}`), nil
}

func (m MFMap) UnmarshalJSON(data []byte) error {
	var mm MFMap
	if err := json.Unmarshal(data, &mm); err != nil {
		return err
	}
	return nil
}

type SimulateData = slotsmongo.SimulateData

func DealMF(pan map[string]any) map[string]any {
	mf := pan["mf"].(primitive.M)
	newMF := MFMap{
		MT: mf["mt"],
		MS: mf["ms"],
		MI: mf["mi"],
	}
	pan["mf"] = newMF
	return pan
}

const Line = 5
const BuyMul = 5
