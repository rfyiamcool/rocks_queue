# rocks_queue

使用rocksdb实现的高性能队列.

### 数据结构

```
			key                                 value
l[list]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 945
l[list]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 946
l[list]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 947
l[list]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 948
l[list]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 949
l[list]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 950

```


### benchmark

insert (push queue)

```
62000 / s
```

参考 陌陌 GoRedis代码...
