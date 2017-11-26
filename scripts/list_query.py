import rocksdb


if __name__ == "__main__":
    db = rocksdb.DB("rocks_data", rocksdb.Options(create_if_missing=True))
    db.get("aaa")
    it = db.iteritems()
    prefix = ''
    it.seek(prefix)
    print(list(it))
