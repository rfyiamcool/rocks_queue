package main

import (
	"fmt"

	"github.com/rfyiamcool/rocks_queue/rocks"
)

func main() {
	var err error
	var res []byte

	rocksdb := rocks.NewRocksDB("./rocks_data")
	db := rocks.New(rocksdb)
	defer db.Close()

	l := db.List([]byte("qq"))

	for index := 0; index < 10; index++ {
		l.RPush([]byte(fmt.Sprintf("xiaorui.cc index: %d", index)))
	}

	fmt.Printf("queue length: %d \n", l.Len())

	for index := 0; index < 1; index++ {
		res, err = l.LPop()
		fmt.Printf("LPOP value: %v, err: %v \n", string(res), err)
	}

	fmt.Printf("queue length: %d \n", l.Len())
	// l.Drop()
}
