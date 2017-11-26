# rocks_queue

使用rocksdb实现的高性能队列.

### 改进

* 计数改为内存计数，只在启动初始化时会扫描一次left、right
* 增加更多的redis指令.

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

push queue

```
75000 qps
```

pop queue

```
68000 qps
```

参考了约炮神器陌陌 { GoRedis半成品代码 }
