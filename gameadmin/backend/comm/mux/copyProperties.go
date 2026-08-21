package mux

import (
	"errors"
	"reflect"
)

func ShallowCopy(dst, src interface{}) (err error) {
	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}

// 初始化字段
// i 必须是想 *struct
// map, slice , *struct 类型的字段如果是nil的话将被初始化
func InitNilFieldsWithSelector(i interface{}) {
	t := reflect.TypeOf(i)
	Assert(t.Kind() == reflect.Ptr)
	t = t.Elem()
	Assert(t.Kind() == reflect.Struct)

	v := reflect.ValueOf(i).Elem()

	propertyNums := t.NumField()

	for i := 0; i < propertyNums; i++ {
		stField := t.Field(i)
		field := v.Field(i)

		if !field.CanSet() {
			continue
		}

		switch stField.Type.Kind() {
		case reflect.Map:
			if field.IsNil() {
				field.Set(reflect.MakeMap(stField.Type))
			}
		case reflect.Slice:
			if field.IsNil() {
				field.Set(reflect.MakeSlice(stField.Type, 0, 4))
			}
		case reflect.Ptr:
			if field.IsNil() && stField.Type.Elem().Kind() == reflect.Struct {
				field.Set(reflect.New(stField.Type.Elem()))
			}
		}
	}
}

// *map[string]int     *[]int     *struct{}
func ZeroAt(i interface{}) {
	t := reflect.TypeOf(i)
	Assert(t.Kind() == reflect.Ptr)
	tt := t.Elem()

	v := reflect.ValueOf(i).Elem()

	switch tt.Kind() {
	case reflect.Map:
		v.Set(reflect.MakeMap(tt))
	case reflect.Slice:
		v.Set(reflect.MakeSlice(tt, 0, 4))
	case reflect.Ptr:
		ttt := tt.Elem()
		v.Set(reflect.New(ttt))
	default:
		v.Set(reflect.Zero(tt))
	}
}
