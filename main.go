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

	fmt.Println("start...")

	l := db.List([]byte("list"))

	for index := 0; index < 10000; index++ {
		l.RPush([]byte(fmt.Sprintf("xiaorui.cc index: %d", index)))
	}

	for index := 0; index < 1; index++ {
		res, err = l.LPop()
		fmt.Println(string(res), err)
	}

	fmt.Printf("queue length: %d \n", l.Len())
}
