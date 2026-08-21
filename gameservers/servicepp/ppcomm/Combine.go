package ppcomm

type CombineMeta struct {
	BucketIds []int
	Count     int
}

type CombineItem struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"-"`
	Count int
	Meta  []CombineMeta `bson:"-"`
}

type Combine []*CombineItem

func NewCombine(combine Combine) Combine {
	for i, v := range combine {
		v.ID = i
	}
	return combine
}
