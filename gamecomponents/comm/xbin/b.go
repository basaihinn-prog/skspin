package xbin

import (
	"bytes"
	"errors"
	"io"
	"math"
	"reflect"

	"google.golang.org/protobuf/encoding/protowire"
)

type decoder struct {
	buf    []byte
	offset int
}

type encoder struct {
	writer io.Writer
}

func (d *decoder) bool() bool {
	x := d.buf[d.offset]
	d.offset++
	return x != 0
}

func (e *encoder) bool(x bool) {
	if x {
		e.uint8(1)
	} else {
		e.uint8(0)
	}
}

func (e *encoder) uint8(x uint8) { e.writer.Write([]byte{x}) }
func (d *decoder) uint8() uint8 {
	x := d.buf[d.offset]
	d.offset++
	return x
}

func (e *encoder) int8(x int8) { e.uint8(uint8(x)) }
func (d *decoder) int8() int8  { return int8(d.uint8()) }

func (e *encoder) int(i uint64) {
	v := protowire.AppendVarint(nil, i)
	e.writer.Write(v)
}

func (d *decoder) int() uint64 {
	v, n := protowire.ConsumeVarint(d.buf[d.offset:])
	d.offset += n
	return v
}

func (d *decoder) value(v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		l := v.Len()
		for i := 0; i < l; i++ {
			d.value(v.Index(i))
		}
	case reflect.Slice:
		l := int(d.int())

		newv := reflect.MakeSlice(v.Type(), l, l)
		v.Set(newv)

		for i := 0; i < l; i++ {
			d.value(v.Index(i))
		}
	case reflect.String:
		l := int(d.int())

		bytes := make([]byte, l)
		for i := 0; i < l; i++ {
			bytes[i] = byte(d.int8())
		}
		v.SetString(string(bytes))

	case reflect.Struct:
		l := v.NumField()
		st := v.Type()
		for i := 0; i < l; i++ {
			field := indirect(v.Field(i))
			tag := st.Field(i).Tag.Get("xbin")
			if tag == "-" {
				continue
			}

			if field.CanSet() {
				d.value(field)
			} else {
				// e.skip(v)
			}
		}

	case reflect.Bool:
		v.SetBool(d.bool())

	case reflect.Int8:
		v.SetInt(int64(d.int8()))
	case reflect.Uint8:
		v.SetUint(uint64(d.uint8()))

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(d.int())
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(d.int()))

	case reflect.Float32, reflect.Float64:
		i := d.int()
		v.SetFloat(math.Float64frombits(i))

	case reflect.Pointer:
		d.value(indirect(v))
	}
}

func (e *encoder) value(v reflect.Value) {
	kind := v.Kind()
	switch kind {
	case reflect.Array, reflect.Slice, reflect.String:
		l := v.Len()
		e.int(uint64(l))
		for i := 0; i < l; i++ {
			e.value(v.Index(i))
		}

	case reflect.Struct:
		l := v.NumField()
		st := v.Type()
		for i := 0; i < l; i++ {
			tag := st.Field(i).Tag.Get("xbin")
			if tag == "-" {
				continue
			}
			// see comment for corresponding code in decoder.value()
			field := indirect(v.Field(i))

			if field.CanSet() {
				e.value(field)
			} else {
				// e.skip(v)
			}
		}

	case reflect.Bool:
		e.bool(v.Bool())

	case reflect.Int8:
		e.int8(int8(v.Int()))
	case reflect.Uint8:
		e.uint8(uint8(v.Uint()))

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		e.int(uint64(v.Int()))
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		e.int(v.Uint())

	case reflect.Float32, reflect.Float64:
		e.int(math.Float64bits(v.Float()))

	case reflect.Pointer:
		e.value(indirect(v))

	default:
		panic("unsupport")
	}
}

func indirect(v reflect.Value) reflect.Value {
	v0 := v
	haveAddr := false
	decodingNull := false

	if v.Kind() != reflect.Pointer && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			// if e.Kind() == reflect.Pointer && !e.IsNil() && (e.Elem().Kind() == reflect.Pointer) {
			if e.Kind() == reflect.Pointer && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Pointer) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Pointer {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return v
}

func Unmarshal(buf []byte, data any) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Pointer || data == nil {
		return errors.New("must pointer")
	}

	d := &decoder{buf: buf}
	d.value(indirect(v))
	return nil
}

func Marshal(data any) ([]byte, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Pointer || data == nil {
		return nil, errors.New("must pointer")
	}

	var buf bytes.Buffer
	e := &encoder{writer: &buf}
	e.value(v)
	return buf.Bytes(), nil
}
