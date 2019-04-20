# No compromises: distributed transactions with consistency, availability, and performance

* [No compromises: distributed transactions with consistency, availability, and performance](https://pdos.csail.mit.edu/6.824/papers/farm-2015.pdf)

## Prequirement

- [分布式事务](https://juejin.im/post/5b5a0bf9f265da0f6523913b)
- [分布式系统的CAP理论](https://www.hollischuang.com/archives/666)
- [深入浅出全面解析RDMA](https://zhuanlan.zhihu.com/p/37669618)
- [深入理解乐观锁与悲观锁](https://www.hollischuang.com/archives/934)

## Reference

- [vedio](https://www.youtube.com/watch?v=fYrDPK_t6J8)
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

ok，那么一个事务执行是怎么work的呢？在FaRM中，每个事务都是由单个线程区执行的

* 应用线程开启一个事务，同时变为协调者（coordinator）；
* C可以使用单边RDMA从集群中的其他机器读对象，C显然可以从多台机器读多个对象来执行一些操作；
* Write 都是乐观的，这些write都会被buffer在coordinator的内存中，直到commit time； 
* 在commit时，我们会去提交这些乐观的事务，如果提交成功，所有更新都会自动应用到共享内存中，否则事务就会被abort，然后就要重试。

那么如何进行commit？一个非常直观的想法是使用2PC，它需要2轮去commit transaction。然而，它需要一些messages，然后在transaction commit期间，它同样需要cpu的开销。因此，这篇文章没有使用2pc，而是设计了下面这一种思路：

* Lock：协调器将lock record写入到primary中。record包含所有已写入对象的版本和新值，以及具有已写入对象的所有区域的列表。primary将对象锁定在指定的版本来处理这些记录，并返回一条消息，报告是否成功获取了所有锁。如果任何对象版本在事务读取后发生更改，或者对象当前被另一个事务锁定，则锁定会失败。在这种情况下，协调器中止事务。它将一个中止记录写入所有主计算机，并向应用程序返回一个错误。
* Validate：因为使用了乐观锁的机制，所以validate是必要的。协调器通过从primary中读取对象的版本来执行读验证。如果任何对象已更改，验证将失败，事务将中止。论文使用RDMA进行验证，因此这一步操作不需要额外的开销。
* Replicate：验证成功的话，就会把这些数据写入到备份中，同样的，使用单边RDMA，因此只会有网卡返回一个ACK
* Update and unlock：在所有提交备份写入都被确认之后，协调器将Commit primaries log写入每个primary里。Primaries处理这些log的方法是：在适当的位置更新对象，增加它们的版本，然后解锁它们，这样就公开了事务提交的写操作。

锁定确保了写入对象的安全，验证确保了只读取对象的安全。在没有失败的情况下，这相当于在序列化点原子地执行和提交整个事务。

这张图和前面那个对比，首先，messages数量减少了，在lock阶段，事务执行期间我们用到了cpu，而在剩余阶段cpu都用不到。剩余三个阶段延迟都很低，因为我们用的是单边RDMA。

## Failure Recovery

ok，这是没有failure的情况，那么如果发生了failure会怎么样。

由于使用了单边操作，failure recovery会比较复杂。一个直觉是说，基本上我们没法拥有CPU rejecting messages，一个类似的系统，对于一致性，它可能会采取的措施就是利用租约，系统保证它存储的对象在其租约到期之前不会发生变化。但是FaRM使用RDMA，nic网卡不支持租约识别，因此不管怎么样nic一定会返回一个response。从这个角度而言，论文需要去重新设计一个故障恢复。对于这一点，论文实现了一个精确的成员关系，失败后，新配置中的所有机器必须在允许对象突变之前就其成员资格达成一致。

再比如，单边RDMA写入也会影响事务恢复。跨配置一致性的一般方法是拒绝来自旧配置的消息。FaRM无法使用此方法，我们通过Drain logs来解决这个问题，以确保在恢复过程中处理所有相关的记录。

论文需要保证高可用，对于数据高可用，这也就意味着当出现失效时，在继续处理数据前必须要重新配置系统。一个主要的技术是使用主备份，当primary失效的时候，backup会瞬间晋升成primary，

为了迅速的错误恢复，发现错误的时间也要很快。论文使用了好几种技术来达成这一点，比如使用Dedicated network queues避免网络间的竞争，使用dedicated prioritized thread来避免CPU的竞争，使用memory pre-allocation避免内存分配的干扰。

为了高可用，论文也尽可能的使用并行。

ok，那接下来就是它到底怎么work，这是一个从一个high level的角度来说。

### detecting failures

在FaRM集群中有一个特殊的角色，是configuration manager，简称CM。CM的目标就是侦查其他机器的故障，要达成这一点，它会向其他机器发送心跳信息。其他机器也可以因此知道CM是否失效。

现在假设一台机器挂了，CM就会知道，CM就要启动一个配置改变协议（显然，就是主备份），第一步是将新配置写入agreement service，FaRM用的就是ZooKeeper。

### configuration change

那么configuration change是怎么工作的，正如一开始提及的那样，先将新配置写入ZooKeeper。

同时，CM会remap regions，因此cm会重新分配先前映射到故障机器的区域，对于失败的primary，它总是将backup晋升为primary并且重新分配。

当CM从ZooKeeper收到一个response并且重新映射区域后，CM向配置中的所有计算机发送Config-New的消息，在收到以后，所有机器会更新自身的配置并停止访问所有受失败影响的region，之后这些机器向CM发送一个Config-Ack，当CM收到所有的Config-Ack后，它再向所有机器发送Config-Commit，在这之后，就可以进行下一步的Recovery。

我们之所以需要一个三步message是使用了单边RDMA，因此这种方法给了我们精确的成员，配置中的机器不会向不在配置中的机器发出RDMA请求，并且会忽略对不在配置中的机器的RDMA读操作和RDMA写操作的应答。

### Transaction recovery

修改配置后呢，集群将要恢复那些需要恢复的事务。

正如我之前提及的那样，在恢复region前，需要drain transaction logs，意味着需要处理所有事务logs中的所有记录。所以目标就是在新的配置上要接管那些在旧的配置上还未执行的messages。

draining logs完成后，就会开始transaction recovery。它会在每一个受影响的region都执行，那每个region数据都很大，所以这些恢复操作都被设计为并行的。第一步，所有backup都要根据log去找那些它们需要恢复的事务，它们会把对应region的log发送给primary，在primary拿到所有message后，primary会重新去获取这些region的对应log。要注意的是primary在内存中有一个数据的备份，因为primary是从backup晋升过来的。所以这些region看起来应该和失效前是一致的。

在确定了要恢复的事务后，之前受失败影响的region就可以被重新访问了。后续恢复步骤并行地读取对象并提交对region的更新，从而提升效率。

下一步是复制primary上的log，这是需要的，需要应对未来的失效。然后primary给协调器发送一个投票请求，根据事务更新的每个区域的投票决定是提交还是中止事务，最后coordinator决定最终结果。

所以这就是事务的恢复，我知道它涉及很多步骤，而且还隐藏了一些细节，比如如何确定每个region的投票。我们需要记住的是，上面这些步骤的直觉是恢复保留了先前已提交或中止的事务的结果。我们说只有primary公开事务修改或coordinator通知应用程序提交事务时，事务被提交了。 当coordinator发送中止消息或通知应用程序事务已中止时，事务将中止。 对于尚未确定结果的交易，恢复可能会提交或中止，但它确保从其他故障中恢复可以保留之前结果。另外一点是，FaRM通过backup晋升为primary机制，确定恢复事务后相关region可以被重新访问以及大量的并行操作去尽可能的减少系统的停机时间。

### Data recovery

最后一步是数据恢复，对于FaRM来说，就是要对那些丢失了一个backup的region重新创建一个新的backup，这些是并行完成的，当在进行数据恢复的时候，新的事务仍然可以执行。因为它不会干扰到新的事务，因此在恢复数据的时候可以慢一点，放在后台，尽量不去影响前台的操作。

## Evaluation

最后，我花一两分钟讲一下FaRM的评估结果。

首先是性能，在正常操作下，没有failures。TATP是用来测试高性能内存数据库的，它主要都是读操作。