package store

/*
import (
	"errors"
	"game/duck/mongodb"
	"game/duck/ut2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 把mongo当作 KV存储。

type element = any

var ErrNotFound = errors.New("not found")

type value[T element] struct {
	V T
}

type IStore interface {
	get(k string, v any) error
	set(k string, v any) error
	del(k string) error
}

func Get[T element](s IStore, k string) (T, error) {
	one := value[T]{}
	err := s.get(k, &one)
	return one.V, err
}

func GetTo[T element](s IStore, k string, to T) error {
	one := value[T]{V: to}
	err := s.get(k, &one)
	return err
}

func Set(s IStore, k string, v any) error {
	return s.set(k, &value[any]{V: v})
}

func Del(s IStore, k string) error {
	return s.del(k)
}

type mongoStore struct {
	coll *mongodb.Collection
}

func MongoStore(c *mongodb.Collection) *mongoStore {
	return &mongoStore{coll: c}
}

func (ms *mongoStore) get(k string, v any) error {
	err := ms.coll.FindId(k, v)
	if err == mongo.ErrNoDocuments {
		return ErrNotFound
	}
	return err
}

func (ms *mongoStore) set(k string, e any) error {
	return ms.coll.UpsertId(k, bson.M{"$set": e})
}

func (ms *mongoStore) del(k string) error {
	return ms.coll.DeleteId(k)
}

///

type memStore struct {
	mp ut2.IMap[string, any]
}

func MemStore(mp ut2.IMap[string, any]) *memStore {
	return &memStore{mp: mp}
}

func (ms *memStore) get(k string, v any) error {
	e, ok := ms.mp.Load(k)
	if ok {
		return bson.UnmarshalWithRegistry(mongodb.DefaultRegistry, e.([]byte), v)
	}

	return ErrNotFound
}

func (ms *memStore) set(k string, e any) error {
	buf, err := bson.MarshalWithRegistry(mongodb.DefaultRegistry, e)
	if err != nil {
		return err
	}
	ms.mp.Store(k, buf)
	return nil
}

func (ms *memStore) del(k string) error {
	ms.mp.Delete(k)
	return nil
}

///

type rwStore struct {
	r IStore
	w IStore
}

func RWStore(r, w IStore) *rwStore {
	return &rwStore{r: r, w: w}
}

func (ms *rwStore) get(k string, v any) error {
	err := ms.r.get(k, v)

	if err == ErrNotFound {
		err = ms.w.get(k, v)
		if err == ErrNotFound {
			return ErrNotFound
		}
		err = ms.r.set(k, v)
		if err != nil {
			return err
		}
		return nil
	}

	return err
}

func (ms *rwStore) Set(k string, e any) error {
	err := ms.r.set(k, e)
	if err != nil {
		return err
	}
	return ms.w.set(k, e)
}

func (ms *rwStore) Del(k string) error {
	err := ms.r.del(k)
	if err != nil {
		return err
	}
	return ms.w.del(k)
}

*/
