package slotsmongo

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
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

type GatewaySend struct {
	Pid  int64  `json:"pid"`
	Data []byte `json:"data"`
}

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

// var CompressionThreshold = 1024
// var MaxMessageSize = 100000
var CompressionThreshold = 512
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
func NewFromBinaryData(data []byte) (*SFSObject, error) {
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

	// 5. 如果 flag 的 32（0x20）位被置位，则说明数据体经过压缩，需要解压
	if (flag & 32) > 0 {
		// 4. 读取数据体
		packetData := make([]byte, packetLength)
		if _, err := io.ReadFull(buffer, packetData); err != nil {
			return nil, fmt.Errorf("failed to read packet data: %v", err)
		}
		data, err := decompress(packetData)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress packet data: %v", err)
		}

		buffer = bytes.NewReader(data)
	}

	var dataType byte
	if err := binary.Read(buffer, binary.BigEndian, &dataType); err != nil {
		return nil, err
	}
	// 根据类型读取 value
	value, err := readValue(buffer, dataType)
	if err != nil {
		return nil, err
	}
	fmt.Println(value)
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

// PutByteArray 添加 byte 数组到 SFSObject 中
func (obj *SFSObject) PutByteArray(key string, value []byte) {
	obj.dataHolder.Put(key, SFSDataWrapper{Type: BYTE_ARRAY, Value: value})
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

func CreateSpin() SFSObject {
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
	entity.PutString("playerBet", "9")
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
	p.PutInt("ct", int32(CompressionThreshold))
	p.PutInt("ms", int32(MaxMessageSize))
	uid, _ := uuid.NewUUID()
	p.PutString("tk", uid.String())
	so.PutSFSObject("p", &p)
	return so
}

func C0A1P(pid int, un string) SFSObject {
	so := SFSObject{}
	so.Init()

	p := SFSObject{}
	p.Init()
	p.PutShort("rs", 0)
	p.PutString("zn", "JDB_ZONE_GAME")
	p.PutString("un", un)
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
	p.PutInt("id", int32(pid))

	so.PutSFSObject("p", &p)
	so.PutShort("a", 1)
	so.PutByte("c", 0)

	return so
}

func C0A1PFC(pid int64, un string) SFSObject {
	so := NewSFSObject()
	so.PutByte("c", 0)
	so.PutShort("a", 1)

	p := NewSFSObject()
	p.PutShort("rs", 0)
	p.PutString("zn", "FC_GAME_ZONE")
	p.PutString("un", un)
	p.PutShort("pi", 0)
	// 创建 rl 数组
	rl := NewSFSArray()

	p.PutSFSArray("rl", rl)
	p.PutInt("id", int32(pid))

	so.PutSFSObject("p", p)
	so.PutShort("a", 1)
	so.PutByte("c", 0)

	return *so
}

func A4C0FACAI() SFSObject {
	bytes, _ := base64.StdEncoding.DecodeString("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")
	got, _ := NewFromBinaryData(bytes)
	return *got
}

func A1000C0FACAI() SFSObject {
	bytes, _ := base64.StdEncoding.DecodeString("gABFEgADAAFwEgACAAFyBAAAAAIAAXURAAUEAAAHlwgAEkJFTi1iZW5fMDU1NTU2NjY2NgMAAAMARhEAAAABYQMD6AABYwIA")
	got, _ := NewFromBinaryData(bytes)
	return *got
}

// facai
func C1A13FACaiLobbyInfo(pid string) SFSObject {
	//自己做一个
	so := SFSObject{}
	so.Init()
	p := SFSObject{}
	p.Init()
	pp := SFSObject{}
	pp.Init()
	pp.PutShort("Status", 0)
	pp.PutString("TableName", "SlotRoom")
	pp.PutBool("testMode", false)
	pp.PutString("NickName", pid)
	pp.PutString("LobbyName", "SlotRoom")

	p.PutSFSObject("p", &pp)
	p.PutString("c", "LobbyInfo")
	so.PutSFSObject("p", &p)

	so.PutShort("a", 13)
	so.PutByte("c", 1)
	return so
}

// todo 跟余额有关系了
func C1A13gameLoginReturn(pid int64) SFSObject {
	//自己做一个
	so := SFSObject{}
	so.Init()
	p := SFSObject{}
	p.Init()
	pp := SFSObject{}
	pp.Init()
	pp.PutString("loginRoom", "SLOT_ROOM")
	pp.PutBool("data", true)
	balance, _ := GetBalance(pid)

	pp.PutDouble("balance", float64(balance)/10000) //在游戏里用的是分
	pp.PutBool("testMode", false)
	pp.PutString("serverId", "01")
	pp.PutLong("ts", time.Now().UnixNano())
	p.PutSFSObject("p", &pp)
	p.PutString("c", "gameLoginReturn")
	so.PutSFSObject("p", &p)
	so.PutByte("c", 1)
	so.PutShort("a", 13)
	return so
}

// todo 跟投注有关系了应该走info
func C1A13h5initResponse(byte []byte) SFSObject {
	//自己做一个
	so := SFSObject{}
	so.Init()
	p := SFSObject{}
	p.Init()
	pp := SFSObject{}
	pp.Init()
	pp.PutString("code", "initResponse")
	//pp.PutByteArray("entity", []byte(`{"maxBet":9223372036854775807,"defaultWaysBetIdx":-1,"singleBetCombinations":{"10_10_9_NoExtraBet":90,"10_1_9_NoExtraBet":9,"10_20_9_NoExtraBet":180,"10_30_9_NoExtraBet":270,"10_40_9_NoExtraBet":360,"10_50_9_NoExtraBet":450,"10_5_9_NoExtraBet":45},"minBet":0,"gambleTimes":0,"defaultLineBetIdx":0,"defaultConnectBetIdx":-1,"defaultQuantityBetIdx":-1,"gameFeatureCount":3,"executeSetting":{"settingId":"v3_14027_05_03_002","betSpecSetting":{"paymentType":"PT_001","extraBetTypeList":["NoExtraBet"],"betSpecification":{"lineBetList":[1,5,10,20,30,40,50],"betLineList":[9],"betType":"LineGame"}},"gameStateSetting":[{"gameStateType":"GS_003","frameSetting":{"screenColumn":3,"screenRow":3,"wheelUsePattern":"Dependent"},"tableSetting":{"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":55,"noWinIndex":[0],"wheelData":[3,7,5,0,7,6,2,8,6,1,9,6,9,5,0,8,0,2,6,1,9,4,7,7,2,6,0,3,7,4,8,1,9,3,8,5,0,9,8,3,7,5,1,4,8,2,7,8,0,9,6,3,8,9,1]},{"wheelLength":67,"noWinIndex":[0],"wheelData":[6,4,8,2,5,3,7,0,9,4,6,2,5,8,1,4,7,2,9,5,0,6,8,3,9,4,7,1,6,8,3,9,7,2,5,9,3,8,6,0,9,5,8,3,6,4,9,1,8,6,4,9,2,8,5,3,9,7,0,8,6,1,3,7,9,2,5]},{"wheelLength":62,"noWinIndex":[0],"wheelData":[8,4,3,8,6,4,1,5,7,2,6,4,9,0,7,5,7,2,9,4,5,9,1,6,9,4,3,9,3,0,9,4,7,2,8,5,9,1,8,4,5,2,7,7,0,5,7,3,9,4,1,9,3,6,7,2,5,9,0,5,4,8]}]]},"symbolSetting":{"symbolCount":10,"symbolAttribute":["Wild_01","FreeGame_01","M1","M2","M3","M4","M5","M6","M7","M8"],"payTable":[[0,0,100],[0,0,0],[0,0,50],[0,0,25],[0,0,15],[0,0,12],[0,0,8],[0,0,8],[0,0,3],[0,0,3]],"mixGroupCount":0,"mixGroupSetting":[]},"lineSetting":{"maxBetLine":9,"lineTable":[[1,1,1],[0,0,0],[2,2,2],[0,1,2],[2,1,0],[2,1,2],[0,1,0],[1,0,1],[1,2,1]]},"gameHitPatternSetting":{"gameHitPattern":"LineGame_LeftToRight","maxEliminateTimes":0},"specialFeatureSetting":{"specialFeatureCount":1,"specialHitInfo":[{"specialHitPattern":"HP_05","triggerEvent":"Trigger_01","basePay":0}]},"progressSetting":{"triggerLimitType":"RoundLimit","stepSetting":{"defaultStep":1,"addStep":0,"maxStep":1},"stageSetting":{"defaultStage":1,"addStage":0,"maxStage":1},"roundSetting":{"defaultRound":1,"addRound":0,"maxRound":1}},"displaySetting":{"readyHandSetting":{"readyHandLimitType":"NoReadyHandLimit","readyHandCount":1,"readyHandType":["ReadyHand_01"]}}},{"gameStateType":"GS_069","frameSetting":{"screenColumn":5,"screenRow":3,"wheelUsePattern":"Dependent"},"tableSetting":{"tableCount":1,"tableHitProbability":[1.0],"wheelData":[[{"wheelLength":48,"noWinIndex":[0],"wheelData":[5,3,7,6,5,2,8,7,9,4,8,7,9,0,9,7,8,5,9,7,3,8,7,5,8,6,9,2,8,7,4,9,7,3,6,9,5,7,9,0,6,7,4,5,9,7,8,8]},{"wheelLength":39,"noWinIndex":[0],"wheelData":[8,6,4,7,3,8,9,5,7,9,6,7,4,8,6,3,7,5,8,4,7,9,6,0,9,7,5,8,9,7,4,9,5,8,2,9,5,8,6]},{"wheelLength":44,"noWinIndex":[0],"wheelData":[7,9,2,6,1,7,7,0,6,5,8,1,7,9,4,6,5,7,8,4,9,1,7,4,8,0,9,5,1,8,6,3,7,9,3,9,5,7,9,1,8,4,9,6]},{"wheelLength":51,"noWinIndex":[0],"wheelData":[8,3,1,5,4,6,2,7,6,3,7,5,1,8,4,6,7,0,6,5,8,4,9,6,1,7,8,4,9,6,2,8,6,4,9,5,1,7,9,4,6,8,0,5,9,3,6,4,9,6,4]},{"wheelLength":58,"noWinIndex":[0],"wheelData":[6,5,9,4,0,6,5,7,1,8,4,9,6,2,8,5,3,6,9,5,8,4,9,5,8,0,6,4,7,3,9,5,6,1,8,4,9,6,2,8,6,4,7,9,1,9,6,3,8,5,9,2,6,9,3,6,4,6]}]]},"symbolSetting":{"symbolCount":10,"symbolAttribute":["Wild_01","FreeGame_01","M1","M2","M3","M4","M5","M6","M7","M8"],"payTable":[[0,0,100,200,1000],[0,0,0,0,0],[0,0,50,100,500],[0,0,25,50,200],[0,0,15,25,100],[0,0,12,25,100],[0,0,8,15,50],[0,0,8,15,50],[0,0,3,8,30],[0,0,3,5,30]],"mixGroupCount":0,"mixGroupSetting":[]},"lineSetting":{"maxBetLine":25,"lineTable":[[1,1,1,1,1],[0,0,0,0,0],[2,2,2,2,2],[0,1,2,1,0],[2,1,0,1,2],[2,1,2,2,2],[0,1,0,0,0],[1,0,1,2,2],[1,2,1,0,0],[0,1,1,1,2],[2,1,1,1,0],[1,1,0,1,1],[1,1,2,1,1],[1,0,0,0,1],[1,2,2,2,1],[0,0,1,2,1],[2,2,1,0,1],[0,0,1,2,2],[2,2,1,0,0],[0,0,0,1,2],[2,2,2,1,0],[2,1,0,0,0],[0,1,2,2,2],[0,1,2,1,2],[2,1,0,1,0]]},"gameHitPatternSetting":{"gameHitPattern":"LineGame_LeftToRight","maxEliminateTimes":0},"specialFeatureSetting":{"specialFeatureCount":1,"specialHitInfo":[{"specialHitPattern":"HP_07","triggerEvent":"ReTrigger_01","basePay":0}]},"progressSetting":{"triggerLimitType":"RoundLimit","stepSetting":{"defaultStep":1,"addStep":0,"maxStep":1},"stageSetting":{"defaultStage":1,"addStage":0,"maxStage":1},"roundSetting":{"defaultRound":5,"addRound":5,"maxRound":25}},"displaySetting":{"readyHandSetting":{"readyHandLimitType":"NoReadyHandLimit","readyHandCount":1,"readyHandType":["ReadyHand_06"]}},"extendSetting":{"roundOddsRadix":2,"startPower":0,"maxRoundOdds":64}}],"doubleGameSetting":{"doubleRoundUpperLimit":5,"doubleBetUpperLimit":1000000000,"rtp":0.96,"tieRate":0.1},"boardDisplaySetting":{"winRankSetting":{"BigWin":38,"MegaWin":246,"UltraWin":550}},"gameFlowSetting":{"conditionTableWithoutBoardEnd":[["CD_False","CD_True","CD_False"],["CD_False","CD_False","CD_01"],["CD_False","CD_False","CD_False"]]}},"denoms":[10],"defaultDenomIdx":0,"defaultBetLineIdx":0,"betCombinations":{"10_9_NoExtraBet":90,"1_9_NoExtraBet":9,"20_9_NoExtraBet":180,"30_9_NoExtraBet":270,"40_9_NoExtraBet":360,"50_9_NoExtraBet":450,"5_9_NoExtraBet":45},"gambleLimit":0,"buyFeatureLimit":2147483647,"buyFeature":true,"defaultWaysBetColumnIdx":-1}`))
	pp.PutByteArray("entity", byte)
	p.PutSFSObject("p", &pp)
	p.PutString("c", "h5.initResponse")
	so.PutSFSObject("p", &p)
	so.PutByte("c", 1)
	so.PutShort("a", 13)
	return so
}

// todo facai跟投注有关系了应该走info
func C1A13EV_SC_GAME_INIT(byte []byte) SFSObject {
	//自己做一个
	so := SFSObject{}
	so.Init()
	p := SFSObject{}
	p.Init()
	pp := SFSObject{}
	pp.Init()
	pp.PutByteArray("body", byte)
	p.PutSFSObject("p", &pp)
	p.PutString("c", "EV_SC_GAME_INIT")
	so.PutSFSObject("p", &p)
	so.PutShort("a", 13)
	so.PutByte("c", 1)
	return so
}

// todo spin数据
func C1A13h5spinResponse(byte []byte) SFSObject {
	//自己做一个
	so := SFSObject{}
	so.Init()
	p := SFSObject{}
	p.Init()
	pp := SFSObject{}
	pp.Init()
	pp.PutString("code", "spinResponse")
	pp.PutByteArray("entity", byte)
	p.PutSFSObject("p", &pp)
	p.PutString("c", "h5.spinResponse")
	so.PutSFSObject("p", &p)
	so.PutShort("a", 13)
	so.PutByte("c", 1)
	return so
}

// spin数据
func C1A13h5FaCaiSpinResponse(byte []byte) SFSObject {
	//自己做一个
	so := SFSObject{}
	so.Init()
	p := SFSObject{}
	p.Init()
	pp := SFSObject{}
	pp.Init()
	pp.PutShort("total", 1)
	pp.PutBool("merge", false)
	pp.PutByteArray("body", byte)
	pp.PutShort("seq", 1)
	p.PutSFSObject("p", &pp)
	p.PutString("c", "EV_SC_SPIN_RESULT")
	p.PutShort("r", 2)
	so.PutSFSObject("p", &p)
	so.PutShort("a", 13)
	so.PutByte("c", 1)
	return so
}

// decompressSpinData 模拟 JS 中的 ZW 函数：先 Base64 解码，再使用 gzip 解压
func decompressSpinData(s string) (string, error) {
	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	// gzip 解压
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("gzip reader creation failed: %w", err)
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("gzip decompression failed: %w", err)
	}
	return string(decompressed), nil
}

// HistoryRecord 模拟 API 返回的结构体
type HistoryRecord struct {
	GameSeqNo          string      `json:"gameseqno"`
	GameID             string      `json:"gameid"`
	PlayerID           string      `json:"playerid"`
	BeforeGameCredits  string      `json:"beforegamecredits"`
	AfterGameCredits   string      `json:"aftergamecredits"`
	PlayDenom          string      `json:"playdenom"`
	TtlBet             string      `json:"ttlbet"`
	TtlWinGame         string      `json:"ttlwingame"`
	TtlWinJackpot      string      `json:"ttlwinjackpot"`
	HasFreeGame        string      `json:"has_freegame"`
	HasBonus           string      `json:"has_bonus"`
	StartTime          string      `json:"starttime"`
	StartTimestamp     int64       `json:"starttimestamp"`
	Score              string      `json:"score"`
	SpinData           string      `json:"spin_data"`
	GamePageURL        string      `json:"game_page_url"`
	WinTypeCombination string      `json:"win_type_combination"`
	Extra              interface{} `json:"extra,omitempty"`
}

// ResponseData 结构体对应最外层的JSON结构
type ResponseData struct {
	Code string        `json:"code"`
	Data HistoryRecord `json:"data"`
}

// processHistory 对 spin_data 和 extra 字段进行转换处理
func processHistory(record *HistoryRecord) error {
	// 处理 SpinData 字段（仅当其存在并且不等于默认占位值时进行转换）
	//if record.SpinData != nil {
	// 确保 spin_data 是字符串类型（例如来自 JSON 反序列化时为 string）
	// 如果等于占位字符串则不做处理
	if record.SpinData != "H4sIAAAAAAAAAAMAAAAAAAAAAAA" {
		// 调用解码解压函数（即 JS 中的 ZW）
		decompressed, err := decompressSpinData(record.SpinData)
		if err != nil {
			return fmt.Errorf("failed to decompress spin_data: %w", err)
		}
		// 将解压后的 JSON 字符串反序列化为一个通用对象
		var data interface{}
		if err := json.Unmarshal([]byte(decompressed), &data); err != nil {
			return fmt.Errorf("failed to unmarshal spin_data: %w", err)
		}
		fmt.Println(data)
		//record.SpinData = data

	}
	//}

	// 处理 Extra 字段：如果字段存在且不为空（通常为字符串），则解析之
	if record.Extra != nil {
		// 判断 extra 是否为字符串且非空
		if s, ok := record.Extra.(string); ok && s != "" {
			var data interface{}
			// 尝试将字符串反序列化为 JSON 对象
			if err := json.Unmarshal([]byte(s), &data); err != nil {
				// 如果解析失败，则可选择保留原始字符串，或返回错误
				// 此处选择保留原始内容
				// return fmt.Errorf("failed to unmarshal extra: %w", err)
			} else {
				record.Extra = data
			}
		}
	}
	return nil
}

// --------------------------飞机
func (so *SFSObject) AddCreatePAC(p *SFSObject, c byte, a int16) {
	so.PutSFSObject("p", p)
	so.PutShort("a", a)
	so.PutByte("c", c)
}
