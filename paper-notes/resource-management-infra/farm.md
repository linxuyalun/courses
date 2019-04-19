# No compromises: distributed transactions with consistency, availability, and performance

* [No compromises: distributed transactions with consistency, availability, and performance](https://pdos.csail.mit.edu/6.824/papers/farm-2015.pdf)

## Prequirement

- [分布式事务](https://juejin.im/post/5b5a0bf9f265da0f6523913b)
- [分布式系统的CAP理论](https://www.hollischuang.com/archives/666)
- [深入浅出全面解析RDMA](https://zhuanlan.zhihu.com/p/37669618)
- [深入理解乐观锁与悲观锁](https://www.hollischuang.com/archives/934)

## Reference

- [vedio](https://www.youtube.com/watch?v=fYrDPK_t6J8)
- [slides](https://pdfs.semanticscholar.org/adff/d5d50baa1c98b3562c995cef35c4e3256092.pdf)
- [translation](farm-translation.md)
- [Slides](farm.pptx)

## Background

分布式事务是一个非常强大的primitive，它提供了一个抽象"在一台**永不故障**的机器上，事务的执行都是**严格串行**的"，从而隐藏并行和分布式系统的失效，具有强一致性和高可用性的事务简化了分布式系统的构建和推理。

尽管如此，在如今的数据中心中并没有被广泛部署，一个主要原因是，大家相信，事务本身很慢，因此大部分现在的系统或许完全不支持它，或许提供了一种弱一致性形式的事务，只能保证单机事务，跨机器的则无能为力。

这篇论文想要展示的是，我们完全没有必要做这种妥协，分布式事务可以用强一致性，因此做到严格序列化，并且于此同时提高高可用性，要达到这一点，需要一些硬件上的支持。

* 首先用了大量的DRAM，使用256GB的内存，并且会持续增加（因为摩尔定律）
* 我们使用了非易失性内存（ssd），论文使用了一种性价比很高的方法，将锂离子电池与机架内每个24机器机箱中的电源单元集成在一起，分布式UPS有效地使DRAM耐用。当发生电源故障时，分布式UPS利用电池的能量将内存内容保存到商品SSD。这不仅通过避免同步写入到SSD提高了常见情况下的性能，而且还通过仅在发生故障时写入到SSD来保留SSD的生命周期。
* 第三我们使用了支持RDMA的网络，它拥有很高的吞吐量，很低的延迟

如果把这三个结合在一起，会发现它移除了访问存储的开销（除了失效情况），移除大量网络开销，这意味着CPU成为了主要的瓶颈。

ok，现在假设已经有了使用这些硬件的系统，它当然能比没有这些硬件条件下性能更好，但是为了完全利用这些硬件，需要设计一个新的协议，这个协议是为我刚刚之前讲的那些hardware trends量身定做的。因此，这篇论文的主要工作是"design transaction and recovery protocols"，那之前说了，在这种硬件条件下，CPU会成为瓶颈，因此在设计这个协议的同时，论文需要遵循以下三个principle：

* 使用单边的RDMA操作，因为它不仅提供了高吞吐量，低延迟，还做到了CPU efficient
* 第二个准则是在我们的协议中减少消息数量，因为处理消息需要占用CPU，比如论文中设计了一个4PC的事务设计而不是传统的2PC，因为它更加CPU efficient，这个会在之后再详细解释
* 那当然我们可以肯定也要好好的利用CPU，所以它使用了高效的并行，一个例子是，论文设计的协议在系统恢复时允许新事务的执行，从而提高可用性

那接下来我就要介绍论文如果利用上面这三个准则去构建一个高效的事务和恢复协议。

先讲一下背景，论文构建的系统叫做FaRM，是一个分布式计算平台，它的目标是简化构建分布式系统的任务，并能够以简单的方式利用我们刚刚提及的硬件。所以这是一个通用的平台，它的目标应用往往有一些不规则的访问模式啦或者延迟敏感，所以比如一个key-value存储，graph存储或者是OLTP数据库。FaRM有一个一个集群，它把所有数据都存储在内存中，每个机器都要存储数据并且也要执行客户端代码和服务端代码，所以它是一个对称模型，允许充分利用本地的过剩资源，FaRM把所有集群的内存作为一个共享的内存，集群中所有机器中的内存组成统一地址空间，应用就可以在这个上面进行一些事务的操作。ok，那接下来讲的事务就是在这个上下文。

FaRM中的每个机器都会把它大部分内存贡献给共享内存，因为每台机器会把数据存在内存中，每个数据块大小是2gb，论文称之为region，这些数据块可以通过RDMA访问。所以FaRM通过RDMA read去读取对象，比如机器B发起一个RDMA read，通过NIC网卡，到内存中去fetch这些数据，然后直接把数据传回给机器B，而不需要中断CPU，（因为不需要打断CPU，所以它是cpu efficient），即使事务更新了数据，它同样可以保证了数据一致性，这些内容是在14年他们发表的论文讲述的，在这里就不展开了。对于RDMA写，FaRM同样实现了一个非常高效的messaging，为了做到这一点，每台机器会分配一个循环缓冲区并暴露给RDMA，所以每个sender都有一个专用的循环buffer，因此machine B会发送一个write request，然后通过RDMA放入circular buffer，这些信息肯定需要cpu处理，所以当cpu检测到circular buffer有数据的时候，cpu就会拉取这些数据，但是一旦sender知道这些数据存入远程机器后，可能cpu都还没开始处理这些数据，它就会返回一个ack给machine B，这是一个硬件级别的通知，这种处理就是为了优化后面我要介绍的论文的协议。还有一个要注意的点是这些在buffer里的数据会被持久化到ssd，而在后面它在实现事务性的log的时候也用到了相同的方案。

## Transaction

正如之前所说的，为了实现高性能，协议广泛使用了RDMA 单边操作来优化系统，协议同样尽可能的减少信息数来降低cpu的开销。比如论文使用primary-backup replication来减少消息数；使用乐观锁，它可以在一些情况下避免必须进行messageing的锁对象；第三点，协调器只读primaries，来减少锁，减少messaging

ok，那么一个事务执行是怎么work的呢？