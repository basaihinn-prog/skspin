package slotsdb

import (
	"sync"
)

type Coll struct {
	db       *DB
	collID   int
	Name     string
	items    map[int][]byte
	topRowID int

	mutex sync.Mutex
}

func NewColl(collId int, name string, db *DB) *Coll {
	sf := &Coll{
		collID:   collId,
		Name:     name,
		db:       db,
		items:    map[int][]byte{},
		topRowID: 0,
	}

	return sf
}

func (sf *Coll) load(root *DataRoot, data []byte) error {
	switch root.Type {
	case SubDateType_Record:
		sf.items[root.RowID] = data
		if sf.topRowID < root.RowID {
			sf.topRowID = root.RowID
		}
	case SubDateType_Delete:
		sf.items[root.RowID] = nil
	case SubDateType_Truncate:
		sf.items = map[int][]byte{}
		sf.topRowID = 0
	}

	return nil
}

func (sf *Coll) Insert(body []byte) error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	sf.topRowID++

	sf.items[sf.topRowID] = body

	return sf.db.writeData(sf.collID, sf.topRowID, SubDateType_Record, body)
}

func (sf *Coll) Delete(rowID int) error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	sf.items[rowID] = nil

	return sf.db.writeData(sf.collID, rowID, SubDateType_Delete, nil)
}

func (sf *Coll) Truncate() error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	sf.items = map[int][]byte{}
	sf.topRowID = 0

	return sf.db.writeData(sf.collID, sf.topRowID, SubDateType_Truncate, nil)
}

func (sf *Coll) Find(rowID int) []byte {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	return sf.items[rowID]
}

func (sf *Coll) Count() int {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	return len(sf.items)
}

func (sf *Coll) GetAllIds() []int {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	var arr = []int{}

	for id := range sf.items {
		arr = append(arr, id)
	}

	return arr
}

func (sf *Coll) Each(template any, fn func(rowid int, body []byte) error) error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	for id, body := range sf.items {
		err := fn(id, body)
		if err != nil {
			return err
		}
	}

	return nil
}
