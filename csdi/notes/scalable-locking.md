# Multiprocessor: Scalable Locking

* [GitHub notes](https://github.com/Emilio66/CSDI/blob/master/6_Scalable-lock_%E6%9D%8E%E6%99%BA.md)
* [Non-scalable locks are dangerous](https://blog.csdn.net/he11o_liu/article/details/80386839)
* [深入理解自旋锁](https://zhuanlan.zhihu.com/p/40729293)
* [缓存一致性（Cache Coherency）入门](https://www.infoq.cn/article/cache-coherency-primer)

## Why non-scalable lock is dangerous?

不可扩展的锁是很慢的，如果使用`spin lock`或者`ticket lock`，其实都是对单一缓存行的竞争。当出现高竞争时，不可扩展锁的性能不佳，甚至会导致性能断崖式下坠。non-scalable lock哪怕在N个核的时候表现不错，也很有可能在N+1或者N+2个核的时候突然collapse。

以`ticket lock`为例子，在 `ticket lock` 中，会需要记录两个变量，一个是 `now_serving` 表示正在使用lock的ticket（一个整数），另一个是 `next_ticket` 记录着当前最后一张ticket（就是拿票等待的核的号码）。当任何一个核去拿锁的时候，都会需要读取 `now_serving` 这个变量判断跟自己的ticket是否相等。这样一来，每个核都会对now_serving做cache，一旦这个锁被释放，`ticket lock` 中的 `now_serving` 就会增加1，这个操作会invalidate所有核的cache里的 `now_serving` ，这会触发所有的核来重新读取 `now_serving` 这个值所在的cacheline，论文说明了在现有的架构中，**这个read会被串行化处理**，一个一个来，这就导致消耗的时间与等待锁的核的数量呈线性增长关系O(N)。

当很多核想得到同一个锁，并且进入了contend状态之后，一个锁从一个核转移到另一个核的时间是与等待锁的核的数量呈线性关系的，但是这个时间会极大地增加串行部分代码的长度(the length of the serial section)，所以当某个锁积累到一定量的waiter，就会突然collapse。