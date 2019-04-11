# No compromises: distributed transactions with consistency, availability, and performance

* [No compromises: distributed transactions with consistency, availability, and performance](https://pdos.csail.mit.edu/6.824/papers/farm-2015.pdf)

## Prequirement

- [分布式事务](https://juejin.im/post/5b5a0bf9f265da0f6523913b)
- [分布式系统的CAP理论](https://www.hollischuang.com/archives/666)
- [深入浅出全面解析RDMA](https://zhuanlan.zhihu.com/p/37669618)
- [深入理解乐观锁与悲观锁](https://www.hollischuang.com/archives/934)

## Abstract

具有强一致性和高可用性的事务简化了分布式系统的构建和推理。但是，以前的实现性能表现不佳。这迫使系统设计者完全避免事务，削弱一致性保证，或者提供需要程序员对其数据进行分区的单机器事务。本文表明，现代数据中心不需要折衷。研究表明，主内存分布式计算平台farm能够提供具有严格的**可序列化性、高性能、持久性和高可用性的**分布式事务。Farm在具有4.9 TB数据库的90台机器上每秒实现1.4亿个TATP事务的峰值吞吐量，并在不到50毫秒的时间内从故障中恢复。实现这些结果的关键是根据第一原则设计新的事务、复制和恢复协议，以利用RDMA和新的、不成熟的商品网络。提供非易失性DRAM的方法。