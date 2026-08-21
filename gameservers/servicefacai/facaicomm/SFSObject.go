package facaicomm

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"strings"
	"time"
)

// SFSDataType 数据类型枚举
const (
	NULL             byte = 0
	BOOL             byte = 1
	BYTE             byte = 2
	SHORT            byte = 3
	INT              byte = 4
	LONG             byte = 5
	FLOAT            byte = 6
	DOUBLE           byte = 7
	UTF_STRING       byte = 8
	BOOL_ARRAY       byte = 9
	BYTE_ARRAY       byte = 10
	SHORT_ARRAY      byte = 11
	INT_ARRAY        byte = 12
	LONG_ARRAY       byte = 13
	FLOAT_ARRAY      byte = 14
	DOUBLE_ARRAY     byte = 15
	UTF_STRING_ARRAY byte = 16
	SFS_ARRAY        byte = 17
	SFS_OBJECT       byte = 18
	TEXT             byte = 20
)

// SFSArray 结构定义
type SFSArray struct {
	data []SFSDataWrapper
}

// NewSFSArray 创建新的 SFSArray
func NewSFSArray() *SFSArray {
	return &SFSArray{
		data: make([]SFSDataWrapper, 0),
	}
}

// Add 添加元素到 SFSArray
func (arr *SFSArray) Add(value interface{}, dataType byte, forceType bool) {
	arr.data = append(arr.data, SFSDataWrapper{Type: dataType, Value: value})
}

// orderedMap 有序映射结构
type orderedMap struct {
	keys []string
	data map[string]SFSDataWrapper
}

// newOrderedMap 创建新的有序映射
func newOrderedMap() *orderedMap {
	return &orderedMap{
		keys: make([]string, 0),
		data: make(map[string]SFSDataWrapper),
	}
}

// Put 添加或更新键值对
func (m *orderedMap) Put(key string, value SFSDataWrapper) {
	if _, exists := m.data[key]; !exists {
		m.keys = append(m.keys, key)
	}
	m.data[key] = value
}

// Get 获取值
func (m *orderedMap) Get(key string) (SFSDataWrapper, bool) {
	val, ok := m.data[key]
	return val, ok
}

// Keys 返回所有键（按插入顺序）
func (m *orderedMap) Keys() []string {
	return m.keys
}

// Size 返回键的数量
func (m *orderedMap) Size() int {
	return len(m.keys)
}

// 定义包头的常量（此处规定为固定4个字节）
var headerBytes = []byte{0x80, 0x00, 0x33}
var CompressionThreshold = 2147483647
var MaxMessageSize = 500000

// SFSObject 结构定义
type SFSObject struct {
	dataHolder *orderedMap
}

// SFSDataWrapper 包装值与类型
type SFSDataWrapper struct {
	Type  byte
	Value interface{}
}

// NewSFSObject 创建新的SFSObject
func NewSFSObject() *SFSObject {
	return &SFSObject{
		dataHolder: newOrderedMap(),
	}
}

func (obj *SFSObject) Init() {
	if obj.dataHolder == nil {
		obj.dataHolder = newOrderedMap()
	}
}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// decompress 使用 Zlib 解压缩数据
func decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %v", err)
	}
	defer reader.Close()

	decompressedData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %v", err)
	}

	return decompressedData, nil
}

// decodeSFSArray 解码 SFSArray
func decodeSFSArray(e io.Reader) (*SFSArray, error) {

	//var dataType byte
	//if err := binary.Read(e, binary.BigEndian, &dataType); err != nil {
	//	return nil, err
	//}
	//if dataType != SFS_ARRAY {
	//	return nil, fmt.Errorf("Invalid SFSDataType. Expected: %d, found: %d", SFS_ARRAY, dataType)
	//}

	var size int16
	if err := binary.Read(e, binary.BigEndian, &size); err != nil {
		return nil, err
	}
	if size < 0 {
		return nil, fmt.Errorf("Can't decode SFSArray. Size is negative: %d", size)
	}

	arr := NewSFSArray()
	for i := 0; i < int(size); i++ {
		obj, err := decodeObject(e)
		if err != nil {
			return nil, fmt.Errorf("Could not decode value for index %d: %v", i, err)
		}
		arr.Add(obj.Value, obj.Type, true)
	}

	return arr, nil
}

// decodeObject 解码单个对象
func decodeObject(e io.Reader) (SFSDataWrapper, error) {
	var dataType byte
	if err := binary.Read(e, binary.BigEndian, &dataType); err != nil {
		return SFSDataWrapper{}, err
	}

	value, err := readValue(e, dataType)
	if err != nil {
		return SFSDataWrapper{}, err
	}

	return SFSDataWrapper{Type: dataType, Value: value}, nil
}

// NewFromBinaryData2 按照 JS 代码逻辑解析数据
func NewFromBinaryData2(data []byte) (*SFSObject, error) {
	// 数据必须至少包含 1 字节 flag + 2 字节长度字段
	if len(data) < 3 {
		return nil, errors.New("binary data too short")
	}

	buffer := bytes.NewReader(data)

	// 1. 读取标志字节
	var flag uint8
	if err := binary.Read(buffer, binary.BigEndian, &flag); err != nil {
		return nil, fmt.Errorf("failed to read flag byte: %v", err)
	}

	// 2. 根据标志判断读取长度字段的宽度
	var packetLength uint32
	if (flag & 8) > 0 {
		if err := binary.Read(buffer, binary.BigEndian, &packetLength); err != nil {
			return nil, fmt.Errorf("failed to read packet length (uint32): %v", err)
		}
	} else {
		var length16 uint16
		if err := binary.Read(buffer, binary.BigEndian, &length16); err != nil {
			return nil, fmt.Errorf("failed to read packet length (uint16): %v", err)
		}
		packetLength = uint32(length16)
	}

	// 4. 读取数据体
	packetData := make([]byte, packetLength)
	if _, err := io.ReadFull(buffer, packetData); err != nil {
		return nil, fmt.Errorf("failed to read packet data: %v", err)
	}

	// 5. 如果 flag 的 32（0x20）位被置位，则说明数据体经过压缩，需要解压
	if (flag & 32) > 0 {
		data, err := decompress(packetData)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress packet data: %v", err)
		}
		packetData = data
	}

	// 6. 调用 binary2object 将数据体解析为 SFSObject
	return binary2object(bytes.NewReader(packetData))
}

// NewFromBinaryData 从包含包头的二进制数据创建SFSObject
func NewFromBinaryData(data []byte) (*SFSObject, error) {
	// 检查数据至少长度要大于包头长度
	if len(data) < len(headerBytes) {
		return nil, errors.New("binary data too short, missing header")
	}
	t := bytes.NewReader(data)
	// 读取一个 uint8
	var n uint8
	binary.Read(t, binary.BigEndian, &n)

	// 检查第6位（32）是否被设置
	r := (n & 32) > 0
	_ = r

	// 检查第4位（8）是否被设置，决定读取 uint32 还是 uint16
	var i uint32
	if (n & 8) > 0 {
		binary.Read(t, binary.BigEndian, &i)
	} else {
		var i16 uint16
		binary.Read(t, binary.BigEndian, &i16)
		i = uint32(i16)
	}

	// 读取长度为 i 的字节数组
	// o := make([]byte, i)
	// binary.Read(t, binary.BigEndian, &o)
	// 读取 type
	var dataType byte
	if err := binary.Read(t, binary.BigEndian, &dataType); err != nil {
		return nil, err
	}

	// 根据类型读取 value
	value, err := readValue(t, dataType)
	if err != nil {
		return nil, err
	}

	// 剥离包头后，剩余部分作为 SFSObject 的数据
	//buffer := bytes.NewReader(data[len(headerBytes):])
	//return binary2object(buffer)
	return value.(*SFSObject), nil
}

// ToBinary 将 SFSObject 转换为包含包头的二进制数据
func (obj *SFSObject) ToBinary() ([]byte, error) {
	binaryData, err := obj.toBinary()
	//fmt.Printf("%x\n", binaryData)
	if err != nil {
		return nil, err
	}
	byteLength := len(binaryData)
	var header byte = 128

	if byteLength > CompressionThreshold {
		header += 32
		var err error
		binaryData, err = compress(binaryData)
		if err != nil {
			return nil, err
		}
		byteLength = len(binaryData)
	}

	if byteLength > 65335 {
		header += 8
	}

	var buf bytes.Buffer
	buf.WriteByte(header)

	if byteLength > 65335 {
		binary.Write(&buf, binary.BigEndian, uint32(byteLength))
	} else {
		binary.Write(&buf, binary.BigEndian, uint16(byteLength))
	}

	buf.Write(binaryData)

	result := buf.Bytes()
	if len(result) > MaxMessageSize {
		return nil, errors.New("MaxMessageSize to long")
	}
	return result, nil
}

// // ToBinary 将 SFSObject 转换为包含包头的二进制数据
func (obj *SFSObject) toBinary() ([]byte, error) {
	buffer := new(bytes.Buffer)
	//buffer.WriteByte(SFS_OBJECT)
	if err := binary.Write(buffer, binary.BigEndian, SFS_OBJECT); err != nil {
		return nil, err
	}
	//buffer.Write([]byte(SFS_OBJECT))
	// 写入对象自身的序列化数据
	if err := object2binary(obj, buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// binary2object 反序列化时按顺序插入键
func binary2object(buffer io.Reader) (*SFSObject, error) {
	obj := NewSFSObject()

	// 读取 size
	var size int16
	if err := binary.Read(buffer, binary.BigEndian, &size); err != nil {
		return nil, err
	}

	// 读取所有键值对
	for i := 0; i < int(size); i++ {
		// 读取 key
		var keySize int16
		if err := binary.Read(buffer, binary.BigEndian, &keySize); err != nil {
			return nil, err
		}

		keyBytes := make([]byte, keySize)
		if _, err := buffer.Read(keyBytes); err != nil {
			return nil, err
		}
		key := string(keyBytes)

		// 读取 type
		var dataType byte
		if err := binary.Read(buffer, binary.BigEndian, &dataType); err != nil {
			return nil, err
		}

		// 根据类型读取 value
		value, err := readValue(buffer, dataType)
		if err != nil {
			return nil, err
		}

		obj.dataHolder.Put(key, SFSDataWrapper{Type: dataType, Value: value})
	}

	return obj, nil
}

// object2binary 序列化时按顺序写入键
func object2binary(obj *SFSObject, buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, int16(obj.dataHolder.Size())); err != nil {
		return err
	}

	for _, key := range obj.dataHolder.Keys() {
		wrapper, _ := obj.dataHolder.Get(key)
		keyBytes := []byte(key)

		// 写入键长度和键内容
		if err := binary.Write(buffer, binary.BigEndian, int16(len(keyBytes))); err != nil {
			return err
		}
		if _, err := buffer.Write(keyBytes); err != nil {
			return err
		}

		// 写入类型和值
		if err := binary.Write(buffer, binary.BigEndian, wrapper.Type); err != nil {
			return err
		}
		if err := writeValue(buffer, wrapper.Type, wrapper.Value); err != nil {
			return err
		}
	}
	return nil
}

// readValue 根据类型读取值
func readValue(buffer io.Reader, dataType byte) (interface{}, error) {
	switch dataType {
	case NULL:
		return nil, nil
	case BOOL:
		var val byte
		if err := binary.Read(buffer, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val != 0, nil
	case BYTE:
		var val byte
		if err := binary.Read(buffer, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case SHORT:
		var val int16
		if err := binary.Read(buffer, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case INT:
		var val int32
		if err := binary.Read(buffer, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case LONG:
		var val int64
		if err := binary.Read(buffer, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case FLOAT:
		var val float32
		if err := binary.Read(buffer, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case DOUBLE:
		var val float64
		if err := binary.Read(buffer, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case UTF_STRING:
		var size int16
		if err := binary.Read(buffer, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		bytes := make([]byte, size)
		if _, err := buffer.Read(bytes); err != nil {
			return nil, err
		}
		return string(bytes), nil
	case SFS_ARRAY:
		return decodeSFSArray(buffer)
	case SFS_OBJECT:
		return binary2object(buffer)
	case BYTE_ARRAY:
		// 读取数组长度 (int32)
		var size int32
		if err := binary.Read(buffer, binary.BigEndian, &size); err != nil {
			return nil, fmt.Errorf("failed to read BYTE_ARRAY length: %v", err)
		}

		// 检查长度是否为负值
		if size < 0 {
			return nil, fmt.Errorf("error decoding typed array size. Negative size: %d", size)
		}

		// 创建字节数组
		byteArray := make([]byte, size)

		// 逐个读取字节
		for i := int32(0); i < size; i++ {
			var byteVal byte
			if err := binary.Read(buffer, binary.BigEndian, &byteVal); err != nil {
				return nil, fmt.Errorf("failed to read BYTE_ARRAY content at index %d: %v", i, err)
			}
			byteArray[i] = byteVal
		}

		return byteArray, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %d", dataType)
	}
}

// writeValue 根据类型写入值
func writeValue(buffer *bytes.Buffer, dataType byte, value interface{}) error {
	switch dataType {
	case NULL:
		return nil
	case BOOL:
		boolVal := value.(bool)
		var byteVal byte = 0
		if boolVal {
			byteVal = 1
		}
		return binary.Write(buffer, binary.BigEndian, byteVal)
	case BYTE:
		return binary.Write(buffer, binary.BigEndian, value.(byte))
	case SHORT:
		return binary.Write(buffer, binary.BigEndian, value.(int16))
	case INT:
		return binary.Write(buffer, binary.BigEndian, value.(int32))
	case LONG:
		return binary.Write(buffer, binary.BigEndian, value.(int64))
	case FLOAT:
		return binary.Write(buffer, binary.BigEndian, value.(float32))
	case DOUBLE:
		return binary.Write(buffer, binary.BigEndian, value.(float64))
	case UTF_STRING:
		strVal := value.(string)
		strBytes := []byte(strVal)
		if err := binary.Write(buffer, binary.BigEndian, int16(len(strBytes))); err != nil {
			return err
		}
		_, err := buffer.Write(strBytes)
		return err
	case SFS_ARRAY:
		array := value.(*SFSArray)
		if err := binary.Write(buffer, binary.BigEndian, int16(len(array.data))); err != nil {
			return err
		}
		for _, item := range array.data {
			if err := binary.Write(buffer, binary.BigEndian, item.Type); err != nil {
				return err
			}
			if err := writeValue(buffer, item.Type, item.Value); err != nil {
				return err
			}
		}
		return nil
	case SFS_OBJECT:
		return object2binary(value.(*SFSObject), buffer)
	case BYTE_ARRAY:
		byteArray := value.([]byte)
		if err := binary.Write(buffer, binary.BigEndian, int32(len(byteArray))); err != nil {
			return fmt.Errorf("failed to write BYTE_ARRAY length: %v", err)
		}
		for _, byteVal := range byteArray {
			if err := binary.Write(buffer, binary.BigEndian, byteVal); err != nil {
				return fmt.Errorf("failed to write BYTE_ARRAY content: %v", err)
			}
		}
		return nil
	default:
		return errors.New("unsupported data type")
	}
}

// 以下是访问器方法
func (obj *SFSObject) PutNull(key string) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: NULL, Value: nil})
}

func (obj *SFSObject) PutBool(key string, value bool) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: BOOL, Value: value})
}

func (obj *SFSObject) PutByte(key string, value byte) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: BYTE, Value: value})
}

func (obj *SFSObject) PutShort(key string, value int16) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: SHORT, Value: value})
}

func (obj *SFSObject) PutInt(key string, value int32) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: INT, Value: value})
}

func (obj *SFSObject) PutLong(key string, value int64) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: LONG, Value: value})
}

func (obj *SFSObject) PutFloat(key string, value float32) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: FLOAT, Value: value})
}

func (obj *SFSObject) PutDouble(key string, value float64) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: DOUBLE, Value: value})
}

func (obj *SFSObject) PutString(key string, value string) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: UTF_STRING, Value: value})
}

func (obj *SFSObject) PutSFSObject(key string, value *SFSObject) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: SFS_OBJECT, Value: value})
}

func (obj *SFSObject) PutSFSArray(key string, value *SFSArray) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: SFS_ARRAY, Value: value})
}

// 获取值的方法

func (obj *SFSObject) IsNull(key string) bool {
	wrapper, exists := obj.dataHolder.Get(key)
	return !exists || wrapper.Type == NULL
}

func (obj *SFSObject) GetBool(key string) (bool, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != BOOL {
		return false, false
	}
	return wrapper.Value.(bool), true
}

func (obj *SFSObject) GetByte(key string) (byte, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != BYTE {
		return 0, false
	}
	return wrapper.Value.(byte), true
}

func (obj *SFSObject) GetShort(key string) (int16, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != SHORT {
		return 0, false
	}
	return wrapper.Value.(int16), true
}

func (obj *SFSObject) GetInt(key string) (int32, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != INT {
		return 0, false
	}
	return wrapper.Value.(int32), true
}

func (obj *SFSObject) GetLong(key string) (int64, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != LONG {
		return 0, false
	}
	return wrapper.Value.(int64), true
}

func (obj *SFSObject) GetFloat(key string) (float32, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != FLOAT {
		return 0, false
	}
	return wrapper.Value.(float32), true
}

func (obj *SFSObject) GetDouble(key string) (float64, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != DOUBLE {
		return 0, false
	}
	return wrapper.Value.(float64), true
}

func (obj *SFSObject) GetString(key string) (string, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != UTF_STRING {
		return "", false
	}
	return wrapper.Value.(string), true
}

func (obj *SFSObject) GetSFSObject(key string) (*SFSObject, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != SFS_OBJECT {
		return nil, false
	}
	return wrapper.Value.(*SFSObject), true
}

func (obj *SFSObject) GetSFSArray(key string) (*SFSArray, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != SFS_ARRAY {
		return nil, false
	}
	return wrapper.Value.(*SFSArray), true
}

func (obj *SFSObject) GetByteArray(key string) ([]byte, bool) {
	wrapper, exists := obj.dataHolder.Get(key)
	if !exists || wrapper.Type != BYTE_ARRAY {
		return nil, false
	}
	return wrapper.Value.([]byte), true
}

func (obj *SFSObject) String() string {
	var builder strings.Builder
	builder.WriteString("SFSObject{")

	for i, key := range obj.dataHolder.Keys() {
		if i > 0 {
			builder.WriteString(", ")
		}
		value, _ := obj.dataHolder.Get(key)
		builder.WriteString(fmt.Sprintf("%s: %v", key, valueToString(value)))
	}

	builder.WriteString("}")
	return builder.String()
}

// valueToString 将 SFSDataWrapper 的值转换为字符串
func valueToString(value SFSDataWrapper) string {
	switch value.Type {
	case NULL:
		return "null"
	case BOOL:
		return fmt.Sprintf("%v", value.Value.(bool))
	case BYTE:
		return fmt.Sprintf("%v", value.Value.(byte))
	case SHORT:
		return fmt.Sprintf("%v", value.Value.(int16))
	case INT:
		return fmt.Sprintf("%v", value.Value.(int32))
	case LONG:
		return fmt.Sprintf("%v", value.Value.(int64))
	case FLOAT:
		return fmt.Sprintf("%v", value.Value.(float32))
	case DOUBLE:
		return fmt.Sprintf("%v", value.Value.(float64))
	case UTF_STRING:
		return fmt.Sprintf("\"%s\"", value.Value.(string))
	case BOOL_ARRAY:
		return fmt.Sprintf("%v", value.Value.([]bool))
	case BYTE_ARRAY:
		return fmt.Sprintf("%v", value.Value.([]byte))
	case SHORT_ARRAY:
		return fmt.Sprintf("%v", value.Value.([]int16))
	case INT_ARRAY:
		return fmt.Sprintf("%v", value.Value.([]int32))
	case LONG_ARRAY:
		return fmt.Sprintf("%v", value.Value.([]int64))
	case FLOAT_ARRAY:
		return fmt.Sprintf("%v", value.Value.([]float32))
	case DOUBLE_ARRAY:
		return fmt.Sprintf("%v", value.Value.([]float64))
	case UTF_STRING_ARRAY:
		return fmt.Sprintf("%v", value.Value.([]string))
	case SFS_ARRAY:
		arr := value.Value.(*SFSArray)
		var elements []string
		for _, elem := range arr.data {
			elements = append(elements, valueToString(elem))
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
	case SFS_OBJECT:
		return value.Value.(*SFSObject).String()
	case TEXT:
		return fmt.Sprintf("\"%s\"", value.Value.(string))
	default:
		return "unknown"
	}
}
func CreateSFSObject(c byte, a int32) SFSObject {
	so := SFSObject{}
	so.Init()
	so.PutByte("c", c)
	so.PutInt("a", a)
	return so
}
func CreateC(val byte) SFSObject {
	so := SFSObject{}
	so.Init()
	so.PutByte("c", val)
	return so
}

func CreateA(val int32) SFSObject {
	so := SFSObject{}
	so.Init()
	so.PutInt("a", val)
	return so
}

func CreateHeartBeat() SFSObject {
	so := SFSObject{}
	so.Init()
	so.PutString("c", "GEN_HEARTBEAT")
	so.PutInt("r", -1)
	so2 := SFSObject{}
	so2.Init()
	so.PutSFSObject("p", &so2)
	return so
}

func CreateH5Init() SFSObject {
	so := SFSObject{}
	so.Init()
	so.PutString("c", "h5.init")
	so.PutInt("r", -1)
	so2 := SFSObject{}
	so2.Init()
	so2.PutString("code", "h5.init")
	so3 := SFSObject{}
	so3.Init()
	so2.PutSFSObject("entity", &so3)
	so.PutSFSObject("p", &so2)
	return so
}

func CreateSpinTest(bet string) SFSObject {
	sfsTop := SFSObject{}
	sfsTop.Init()
	sfsTop.PutByte("c", 1)
	sfsTop.PutShort("a", 13)

	so := SFSObject{}
	so.Init()
	so.PutString("c", "h5.spin")
	so.PutInt("r", -1)
	entity := SFSObject{}
	entity.Init()
	entity.PutString("denom", "10")
	entity.PutString("extraBetType", "NoExtraBet")
	entity.PutString("gameStateId", "0")
	entity.PutString("playerBet", bet)
	entity.PutString("buyFeatureType", "null")
	betRequest := SFSObject{}
	betRequest.Init()
	betRequest.PutString("betType", "LineGame")
	betRequest.PutInt("betLine", 9)
	betRequest.PutInt("lineBet", 1)
	entity.PutSFSObject("betRequest", &betRequest)
	p := SFSObject{}
	p.Init()
	p.PutSFSObject("entity", &entity)
	so.PutSFSObject("p", &p)
	sfsTop.PutSFSObject("p", &so)
	return sfsTop
}

// FindKeyInNestedMap 在嵌套的 map 中查找键
func FindKeyInNestedMap(data map[string]any, key string) (any, bool) {
	return findKeyRecursive(data, key)
}

// 递归函数，用于在嵌套的 map 中查找键
func findKeyRecursive(data any, key string) (any, bool) {
	if data == nil {
		return nil, false
	}

	switch v := data.(type) {
	case map[string]any:
		// 首先检查当前 map 是否包含该键
		if val, exists := v[key]; exists {
			return val, true
		}
		// 如果不在当前 map 中，递归检查嵌套的 map 和 slice
		for _, val := range v {
			if result, found := findKeyRecursive(val, key); found {
				return result, true
			}
		}
		//case []any:
		//	// 如果是数组，递归检查每个元素
		//	for _, item := range v {
		//		if result, found := findKeyRecursive(item, key); found {
		//			return result, true
		//		}
		//	}
	}
	return nil, false
}

type GateWayData struct {
	Pid  int64  `json:"pid"`
	Data []byte `json:"data"`
}

func ParsePidSfsObject(data []byte) (int64, *SFSObject, error) {
	var gd GateWayData
	err := json.Unmarshal(data, &gd)
	if err != nil {
		return 0, nil, err
	}
	so, err := NewFromBinaryData(gd.Data)
	if err != nil {
		return 0, nil, err
	}
	return gd.Pid, so, nil
}

func C0A0Pctmstk() SFSObject {
	//自己做一个
	so := SFSObject{}
	so.Init()
	so.PutByte("c", 0)
	so.PutShort("a", 0)
	p := SFSObject{}
	p.Init()
	p.PutInt("ct", 2147483647)
	p.PutInt("ms", 500000)
	uid, _ := uuid.NewUUID()
	p.PutString("tk", uid.String())
	so.PutSFSObject("p", &p)
	return so
}

func C0A1P() SFSObject {
	so := SFSObject{}
	so.Init()

	p := SFSObject{}
	p.Init()
	p.PutShort("rs", 0)
	p.PutString("zn", "JDB_ZONE_GAME")
	p.PutString("un", "demo001027@XX")
	p.PutShort("pi", 0)
	// 创建 rl 数组
	rl := NewSFSArray()
	{
		// Inner array 1
		innerArray1 := NewSFSArray()
		innerArray1.Add(int32(2), INT, true)
		innerArray1.Add("SLOT_ROOM", UTF_STRING, true)
		innerArray1.Add("default", UTF_STRING, true)
		innerArray1.Add(true, BOOL, true)
		innerArray1.Add(false, BOOL, true)
		innerArray1.Add(false, BOOL, true)
		innerArray1.Add(int16(2224), SHORT, true)
		innerArray1.Add(int16(5000), SHORT, true)
		innerArray1.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		innerArray1.Add(int16(0), SHORT, true)
		innerArray1.Add(int16(0), SHORT, true)
		rl.Add(innerArray1, SFS_ARRAY, true)

		// Inner array 2
		innerArray2 := NewSFSArray()
		innerArray2.Add(int32(3), INT, true)
		innerArray2.Add("PUSOYS_LOBBY", UTF_STRING, true)
		innerArray2.Add("default", UTF_STRING, true)
		innerArray2.Add(false, BOOL, true)
		innerArray2.Add(false, BOOL, true)
		innerArray2.Add(false, BOOL, true)
		innerArray2.Add(int16(22), SHORT, true)
		innerArray2.Add(int16(5000), SHORT, true)
		innerArray2.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray2, SFS_ARRAY, true)

		// Inner array 3
		innerArray3 := NewSFSArray()
		innerArray3.Add(int32(4), INT, true)
		innerArray3.Add("TONGITS_LOBBY", UTF_STRING, true)
		innerArray3.Add("default", UTF_STRING, true)
		innerArray3.Add(false, BOOL, true)
		innerArray3.Add(false, BOOL, true)
		innerArray3.Add(false, BOOL, true)
		innerArray3.Add(int16(23), SHORT, true)
		innerArray3.Add(int16(5000), SHORT, true)
		innerArray3.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray3, SFS_ARRAY, true)

		// Inner array 4
		innerArray4 := NewSFSArray()
		innerArray4.Add(int32(5), INT, true)
		innerArray4.Add("RUMMY_LOBBY", UTF_STRING, true)
		innerArray4.Add("default", UTF_STRING, true)
		innerArray4.Add(false, BOOL, true)
		innerArray4.Add(false, BOOL, true)
		innerArray4.Add(false, BOOL, true)
		innerArray4.Add(int16(0), SHORT, true)
		innerArray4.Add(int16(5000), SHORT, true)
		innerArray4.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray4, SFS_ARRAY, true)

		// Inner array 5
		innerArray5 := NewSFSArray()
		innerArray5.Add(int32(6), INT, true)
		innerArray5.Add("RUNNING_GAME", UTF_STRING, true)
		innerArray5.Add("default", UTF_STRING, true)
		innerArray5.Add(false, BOOL, true)
		innerArray5.Add(false, BOOL, true)
		innerArray5.Add(false, BOOL, true)
		innerArray5.Add(int16(23), SHORT, true)
		innerArray5.Add(int16(5000), SHORT, true)
		innerArray5.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray5, SFS_ARRAY, true)

		// Inner array 6
		innerArray6 := NewSFSArray()
		innerArray6.Add(int32(231), INT, true)
		innerArray6.Add("18020", UTF_STRING, true)
		innerArray6.Add("default", UTF_STRING, true)
		innerArray6.Add(false, BOOL, true)
		innerArray6.Add(false, BOOL, true)
		innerArray6.Add(false, BOOL, true)
		innerArray6.Add(int16(0), SHORT, true)
		innerArray6.Add(int16(5000), SHORT, true)
		innerArray6.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray6, SFS_ARRAY, true)

		// Inner array 7
		innerArray7 := NewSFSArray()
		innerArray7.Add(int32(232), INT, true)
		innerArray7.Add("18021", UTF_STRING, true)
		innerArray7.Add("default", UTF_STRING, true)
		innerArray7.Add(false, BOOL, true)
		innerArray7.Add(false, BOOL, true)
		innerArray7.Add(false, BOOL, true)
		innerArray7.Add(int16(0), SHORT, true)
		innerArray7.Add(int16(5000), SHORT, true)
		innerArray7.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray7, SFS_ARRAY, true)

		// Inner array 8
		innerArray8 := NewSFSArray()
		innerArray8.Add(int32(233), INT, true)
		innerArray8.Add("SINGLE_SPIN", UTF_STRING, true)
		innerArray8.Add("default", UTF_STRING, true)
		innerArray8.Add(false, BOOL, true)
		innerArray8.Add(false, BOOL, true)
		innerArray8.Add(false, BOOL, true)
		innerArray8.Add(int16(6), SHORT, true)
		innerArray8.Add(int16(5000), SHORT, true)
		innerArray8.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray8, SFS_ARRAY, true)

		// Inner array 9
		innerArray9 := NewSFSArray()
		innerArray9.Add(int32(234), INT, true)
		innerArray9.Add("18026", UTF_STRING, true)
		innerArray9.Add("default", UTF_STRING, true)
		innerArray9.Add(false, BOOL, true)
		innerArray9.Add(false, BOOL, true)
		innerArray9.Add(false, BOOL, true)
		innerArray9.Add(int16(6), SHORT, true)
		innerArray9.Add(int16(5000), SHORT, true)
		innerArray9.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray9, SFS_ARRAY, true)

		// Inner array 10
		innerArray10 := NewSFSArray()
		innerArray10.Add(int32(235), INT, true)
		innerArray10.Add("MINES", UTF_STRING, true)
		innerArray10.Add("default", UTF_STRING, true)
		innerArray10.Add(false, BOOL, true)
		innerArray10.Add(false, BOOL, true)
		innerArray10.Add(false, BOOL, true)
		innerArray10.Add(int16(133), SHORT, true)
		innerArray10.Add(int16(5000), SHORT, true)
		innerArray10.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray10, SFS_ARRAY, true)

		// Inner array 11
		innerArray11 := NewSFSArray()
		innerArray11.Add(int32(236), INT, true)
		innerArray11.Add("CASINO_ROOM", UTF_STRING, true)
		innerArray11.Add("default", UTF_STRING, true)
		innerArray11.Add(false, BOOL, true)
		innerArray11.Add(false, BOOL, true)
		innerArray11.Add(false, BOOL, true)
		innerArray11.Add(int16(0), SHORT, true)
		innerArray11.Add(int16(5000), SHORT, true)
		innerArray11.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray11, SFS_ARRAY, true)

		// Inner array 12
		innerArray12 := NewSFSArray()
		innerArray12.Add(int32(237), INT, true)
		innerArray12.Add("18022", UTF_STRING, true)
		innerArray12.Add("default", UTF_STRING, true)
		innerArray12.Add(false, BOOL, true)
		innerArray12.Add(false, BOOL, true)
		innerArray12.Add(false, BOOL, true)
		innerArray12.Add(int16(0), SHORT, true)
		innerArray12.Add(int16(5000), SHORT, true)
		innerArray12.Add(NewSFSArray(), SFS_ARRAY, true) // Empty SFSArray
		rl.Add(innerArray12, SFS_ARRAY, true)
	}
	p.PutSFSArray("rl", rl)
	p.PutInt("id", 636958)

	so.PutSFSObject("p", &p)
	so.PutShort("a", 1)
	so.PutByte("c", 0)

	return so
}

// todo 跟余额有关系了
func C1A13gameLoginReturn() SFSObject {
	//自己做一个
	so := SFSObject{}
	so.Init()
	p := SFSObject{}
	p.Init()
	pp := SFSObject{}
	pp.Init()
	pp.PutString("loginRoom", "SLOT_ROOM")
	pp.PutBool("data", true)
	pp.PutDouble("balance", 2000)
	pp.PutBool("testMode", false)
	pp.PutString("serverId", "01")
	pp.PutLong("ts", int64(time.Now().Second()))
	p.PutSFSObject("p", &pp)
	p.PutString("c", "gameLoginReturn")
	so.PutSFSObject("p", &p)
	so.PutByte("c", 1)
	so.PutShort("a", 13)
	return so
}
