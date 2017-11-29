package rocks

import (
	"bytes"
	"errors"
	"sync"

	"github.com/rfyiamcool/rocks_queue/libs/counter"
	"github.com/tecbot/gorocksdb"
)

// list
// +keyName,l = ""
// l[keyName]0 = "a"
// l[keyName]1 = "b"
// l[keyName]2 = "c"
type ListElement struct {
	db       *DB
	key      []byte
	mu       sync.RWMutex
	isScaned bool
	left     counter.Counter
	right    counter.Counter
}

func NewListElement(db *DB, key []byte) *ListElement {
	dao := &ListElement{db: db, key: key}
	dao.left = counter.Counter(dao.leftIndex())
	dao.right = counter.Counter(dao.rightIndex())
	dao.isScaned = true
	return dao
}

func (l *ListElement) Ltrim(start, stop int) error {
	return nil
}

func (l *ListElement) Range(start, stop int, fn func(i int, value []byte, quit *bool)) error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if start < 0 || (stop != -1 && start > stop) {
		return errors.New("bad start/stop index")
	}

	min := l.indexKey(l.leftIndex() + int64(start))
	max := []byte{MAXBYTE} // use rightIndex is better
	prefix := l.keyPrefix()
	l.db.RangeEnumerate(min, max, IterForward, func(i int, key, value []byte, quit *bool) {
		if !bytes.HasPrefix(key, prefix) {
			*quit = true
			return
		}
		fn(start+i, value, quit)
		if stop != -1 && (i >= (stop - start)) {
			*quit = true
		}
	})

	return nil
}

func (l *ListElement) Index(i int64) ([]byte, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	x := l.leftIndex()
	idxkey := l.indexKey(x + i)
	return l.db.RawGet(idxkey)
}

func (l *ListElement) RPush(vals ...[]byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	batch := gorocksdb.NewWriteBatch()
	defer batch.Destroy()

	x, y := l.leftIndex(), l.rightIndex()
	if x == 0 && y == -1 {
		// empty
		batch.Put(l.rawKey(), KEY_EXIST)
	}

	for _, val := range vals {
		batch.Put(l.indexKey(l.right.Incr(1)), val)
		// batch.Put(l.indexKey(y+int64(i)+1), val)
	}
	return l.db.WriteBatch(batch)
}

func (l *ListElement) LPush(vals ...[]byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	batch := gorocksdb.NewWriteBatch()
	defer batch.Destroy()

	x, y := l.leftIndex(), l.rightIndex()
	if x == 0 && y == -1 {
		// empty
		batch.Put(l.rawKey(), nil)
	}

	for _, val := range vals {
		// batch.Put(l.indexKey(x-int64(i)-1), val)
		batch.Put(l.indexKey(l.left.Decr(1)), val)
	}
	return l.db.WriteBatch(batch)
}

func (l *ListElement) RPop() ([]byte, error) {
	return l.pop(false)
}

func (l *ListElement) LPop() ([]byte, error) {
	return l.pop(true)
}

// true for LPop(), false for RPop()
func (l *ListElement) pop(left bool) ([]byte, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	x, y := l.leftIndex(), l.rightIndex()
	size := y - x + 1
	if size == 0 {
		return nil, nil
	}

	var idxkey []byte
	if left {
		idxkey = l.indexKey(x)
	} else {
		idxkey = l.indexKey(y)
	}

	val, err := l.db.RawGet(idxkey)
	if err != nil {
		return nil, err
	}

	if size > 1 {
		return val, l.db.RawDelete(idxkey)
	} else if size == 1 {
		batch := gorocksdb.NewWriteBatch()
		defer batch.Destroy()
		batch.Delete(l.rawKey())
		batch.Delete(idxkey)
		return val, l.db.WriteBatch(batch)
	} else {
		return nil, errors.New("size less than 0")
	}
}

func (l *ListElement) Drop() error {
	return l.drop()
}

// move out side
func (l *ListElement) Exist() (bool, error) {
	res, err := l.db.GetList(l.key)

	if Bytes2str(res) == "" || err != nil {
		return false, err
	}
	return true, err
}

func (l *ListElement) drop() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var err error
	var ok bool

	ok, err = l.Exist()
	if !ok {
		return errors.New("not found")
	}
	if err != nil {
		return err
	}

	batch := gorocksdb.NewWriteBatch()
	defer batch.Destroy()

	l.db.PrefixEnumerate(l.keyPrefix(), IterForward, func(i int, key, value []byte, quit *bool) {
		batch.Delete(copyBytes(key))
	})
	batch.Delete(l.rawKey())

	err = l.db.WriteBatch(batch)
	if err == nil {
		l.db = nil
	}
	return err
}

func (l *ListElement) Len() int64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.len()
}

func (l *ListElement) len() int64 {
	x, y := l.leftIndex(), l.rightIndex()
	return y - x + 1
}

func (l *ListElement) leftIndex() int64 {
	if l.isScaned {
		return l.left.Count()
	}
	idx := int64(0) // default 0
	l.db.PrefixEnumerate(l.keyPrefix(), IterForward, func(i int, key, value []byte, quit *bool) {
		idx = l.indexInKey(key)
		*quit = true
	})
	return idx
}

func (l *ListElement) rightIndex() int64 {
	if l.isScaned {
		return l.right.Count()
	}
	idx := int64(-1) // default -1
	l.db.PrefixEnumerate(l.keyPrefix(), IterBackward, func(i int, key, value []byte, quit *bool) {
		idx = l.indexInKey(key)
		*quit = true
	})
	return idx
}

// +keyName,l = ""
func (l *ListElement) rawKey() []byte {
	return rawKey(l.key, LIST)
}

// l[keyName]
func (l *ListElement) keyPrefix() []byte {
	return bytes.Join([][]byte{[]byte{LIST}, SOK, l.key, EOK}, nil)
}

// l[keyName]0 = "a"
func (l *ListElement) indexKey(i int64) []byte {
	sign := []byte{0}
	if i >= 0 {
		sign = []byte{1}
	}
	return bytes.Join([][]byte{l.keyPrefix(), sign, Int64ToBytes(i)}, nil)
}

// split l[key]index into index
func (l *ListElement) indexInKey(key []byte) int64 {
	idxbuf := bytes.TrimPrefix(key, l.keyPrefix())
	return BytesToInt64(idxbuf[1:]) // skip sign "0/1"
}
