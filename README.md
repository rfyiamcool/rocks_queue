# rocks_queue

使用rocksdb实现的高性能队列, 可引用到项目中使用.

### 怎么用

安装rocksdb存储引擎

```
https://github.com/facebook/rocksdb/blob/master/INSTALL.md
```

运行

```
git clone git@github.com:rfyiamcool/rocks_queue.git
go run main.go
```

### 改进

* 增加更多的redis指令.
* 收敛代码, 去除其他的数据结构.

### 数据结构

表明类型为list队列

```
    key                               value

+queue_name,l                           1
```

list数据存储格式

```
			key                                         value

l[queue_name]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 945
l[queue_name]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 946
l[queue_name]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 947
l[queue_name]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 948
l[queue_name]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 949
l[queue_name]\x01\x00\x00\x00\x00\x00\x00\x03   xiaorui.cc index: 950

```

### benchmark

push queue

```
190000 qps
```

pop queue

```
188000 qps
```

### ChangeLog

* v1.0 2017-11-28

计数改为内存计数，只在启动初始化时会扫描一次left、right

*参考了约炮神器陌陌 { GoRedis半成品代码 }*
