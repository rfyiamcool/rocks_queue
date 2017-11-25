package main

import (
	"fmt"
	"log"

	"github.com/rfyiamcool/rocks_queue/rocks"
	"github.com/tecbot/gorocksdb"
)

func newRocksDB(dir string) *rocks.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	rdb, err := gorocksdb.OpenDb(opts, dir)
	if err != nil {
		panic(err)
	}

	return rocks.New(rdb)
}

func main() {
	var err error
	var res []byte

	db := newRocksDB("./rocks_6380")
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

	l.Range(28, 35, func(i int, value []byte, quit *bool) {
		log.Fatal(i, string(value))
	})
}
