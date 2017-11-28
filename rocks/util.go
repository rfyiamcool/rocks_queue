package rocks

import (
	"bytes"
	"encoding/binary"
	"math"
	"unsafe"
)

var (
	SEP = []byte{','}
	KEY = []byte{'+'} // Key Prefix
	SOK = []byte{'['} // Start of Key
	EOK = []byte{']'} // End of Key
)

type ElementType byte

const (
	STRING    ElementType = 's'
	HASH                  = 'h'
	LIST                  = 'l'
	SORTEDSET             = 'z'
	NONE                  = '0'
)

func (e ElementType) String() string {
	switch byte(e) {
	case 's':
		return "string"
	case 'h':
		return "hash"
	case 'l':
		return "list"
	case 'z':
		return "sortedset"
	case 'e':
		return "set" // not design
	default:
		return "none"
	}
}

type IterDirection int

const (
	IterForward IterDirection = iota
	IterBackward
)

// 字节范围
const (
	MINBYTE byte = 0
	MAXBYTE byte = math.MaxUint8
)

func rawKey(key []byte, t ElementType) []byte {
	return bytes.Join([][]byte{KEY, key, SEP, []byte{byte(t)}}, nil)
}

// 范围判断 min <= v <= max
func between(v, min, max []byte) bool {
	return bytes.Compare(v, min) >= 0 && bytes.Compare(v, max) <= 0
}

// 复制数组
func copyBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

// 使用二进制存储整形
func Int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func SplitKeyName(key []byte) (string, string) {
	k := string(key)
	length := len(key)
	okString := string(k[1 : length-2])
	ttype := string(k[length-1 : length])
	return okString, ttype
}

func Str2bytes(s string) []byte {
	ptr := (*[2]uintptr)(unsafe.Pointer(&s))
	btr := [3]uintptr{ptr[0], ptr[1], ptr[1]}
	return *(*[]byte)(unsafe.Pointer(&btr))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
