package gendata

type CombineMeta struct {
	BucketIds []int
	Count     int
}

type CombineItem struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"-"`
	Count int
	Meta  []CombineMeta `bson:"-"`
	Type  int           `bson:"-"`
}

type Combine []*CombineItem

func newCombine() Combine {
	combine := Combine{}
	for i, v := range combine {
		v.ID = i
	}
	return combine
}
