# NewSQL & Distributed Transactions

* [GitHub Notes](https://github.com/Emilio66/CSDI/blob/master/10_Spanner_%E8%96%9B%E7%BF%94.md)
* [What's the difference between sharding DB tables and partitioning them?](https://www.quora.com/Whats-the-difference-between-sharding-DB-tables-and-partitioning-them)
* [[论文笔记] Google Spanner Distributed Database](https://blog.csdn.net/chen_kkw/article/details/81262140)
* [Spanner的分布式事务实现](https://zhuanlan.zhihu.com/p/20868175)

## Why NewSQL?

> * [SQL VS. NOSQL VS. NEWSQL: FINDING THE RIGHT SOLUTION](https://dataconomy.com/2015/08/sql-vs-nosql-vs-newsql-finding-the-right-solution/)

传统的SQL很好，使用广泛，标准统一，技术支持丰富，但是传统SQL不是横向扩展架构，事务吞吐量通常由单个机器的容量来控制，扩展性太差。另外，SQL被设计为"one size fits all"，它过于通用以至于性能提升受限，调试复杂。

于是NoSQL来了，最终一致性使得NoSQL的availability很好，扩展性好，半结构化数据的设计适于动态调整schema，但是NoSQL的代价是ACID的弱化，这在之前都讲过。

NewSQL并不像NoSQL那样宽泛。NewSQL系统都是从关系数据模型和SQL查询语言开始的，它们都试图解决类似NoSQL的一样相同类型的可伸缩性，不灵活性或缺乏焦点。很多提供了比NoSQL更强的一致性保证。因此NewSQL 更容易实现强一致性，事务支持，支持SQL语义和工具，使用NoSQL风格的集群架构而提供传统的数据和查询模型；通用性还是没SQL好，并不支持所有传统SQL工具。

## Spanner

### Overview

Spanner 是谷歌的可伸缩、多版本、全球分布、支持同步复制的数据库，它是第一个在全球范围内传递数据且保证外部一致的分布式事务的系统。

Google 内部 Bigtable 和 Megastore 存在一些问题，前者对于复杂可变的模式，或者需要大范围强一致性复制的场合不适用；后者的写吞吐量很差，即使是它的半关系型数据模型和对实施复制的支持很好，但是也非完美的选择。最后的结果是，Spanner 从一个与 Bigtable 相似的 KV 存储进化成了一个数据多版本的数据库。数据存储在模式化、半关系型的表中；数据有版本的区分，数据在 commit 的时候会为每个版本生成一个 timestamp，旧版本会受制于垃圾回收机制；同时应用也可以读旧版本的数据。Spanner 支持通用的事务，提供了基于 SQL 的查询语言。

简单的说，Spanner同BigTable的区别如下：

- 从简单的key-value store加强到temporal multi-version database；
- 数据以半关系型的table组织；
- 支持txn语义；
- 支持SQL查询。

另外，Spanner还有两个特性：第一是应用可以自己细粒度的控制数据复制的配置项，包括哪个 datacenter 包含什么数据、数据离用户多远（控制读延迟）、副本间距多远（控制写延迟）、有几个副本（决定可用性和容灾等级），数据也可以在 datacenter 之间流动，以平衡资源使用。第二 Spanner 实现了两个分布式数据库的难点—— 外部一致的读写操作和在一个 timestamp 下全球一致的跨数据中心的读操作。这些特性使得 Spanner 可以在全球层面上，支持一致性备份、一致的 MapReduce 任务执行，和对 schema 原子更新操作，即使是事务正在进行中也可以。简单的说，如下两个特性：

- 其上的应用程序能够动态配置replication，达到负载均衡和降低延迟；
- **External Consistency**：分布式事务系统中，txn1的commit先于txn2的start，那么txn1的提交时间应小于txn2的提交时间。即能看见txn2的时候一定要能看见txn1。

### Architecture

![](https://github.com/Emilio66/CSDI/blob/master/img/10_1.png?raw=true)

从宏观到微观看，可以类比一下k8s，非常的像。

Universe：Spanner的整个部署

- universe master：单例，维护所有zones；
- placement driver：单例，负责在zones之间迁移数据。
- zone：等同于其下BigTable的部署节点，是管理配置的单位，也是物理隔离的单位；（类比k8s的一个node）

Zone内部：

- zonemaster：每个zone都有一个，负责将数据分配给当前zone的spanserver；
- location proxy：每个zone有多个，为client提供定位到需要的SpanServer的服务；
- spanserver：每个zone有成百上千个，负责为client提供数据服务；（类比k8s的一个pod）

核心设计就在于span server。

### Span Server Architecture

每个 spanserver 负责一百至一千个数据实例 tablet，它和 Bigtable 中的概念有相似也有不同，它通过如果形势组织数据：

```
(key:string, timestamp:int 64) -> string
```

timestamp 在 key 上而不是在 data 上（Bigtable 中数据具有不同的timestamp），这也是为什么 Spanner 更像多版本数据库而不是KV数据库的原因。每个 Tablet 的状态被存储在一系列类似 B-Tree 结果的文件和一个 write-ahead（写前操作） 日志中。为了支持复制，每个 span server 在自己 tablet 的顶层实现了一个 Paxos 状态机，我们的 Paxos 实现通过基于时间的租约而支持长时间存活的 leader。**副本的集合被称为一个 Paxos Group**，写操作必须在 leader 上初始化 Paxos 协议，读操作可以在任意一个足够新的副本上进行。每个副本 leader 上 span server 实现了 lock table 来控制并发，它把 key 的范围映射到锁状态上，对于需要同步的操作例如事务读需要获取锁表中的锁，其他操作不需要。

为了支持分布式事务，每个 span server 实现了 txn mngr 一个 participant leader，当 txn 中仅有一个 Paxos Group 时可以忽略 txn mngr，因为依靠锁表和 Paxos 就可以支持事务，如果涉及多个 Paxos Group，那么这些 Group 的 leader 协商进行 Two-phase commit，其中一个 Group 被选出担任协调者（coordinator），那么这个 Group 中的 participant leader 就会变成 coordinator leader，Group 中的其他副本变成 coordinator slaves 用于容灾。

简单来说，每个span server有如下组成部分：

- **Tablets**：存数据的最小单位，概念和BigTable的tablet相近（一个tablet往往有多个备份副本，会存在其他zone的span server上）；
- **Paxos state machine**：每个span server维护一个用来选举（使用场景：当需要修改本地某个tablet时，由于需要同时修改其他span server上的副本，此时用每个span server上的Paxos状态机来处理一致性的问题(选出leader等)），tablet副本的集合组成Paxos group；写请求由Paxos leader负责，读请求由任意足够up-to-date的tablet所在span server执行都行；
- **Lock table**：（单Paxos group中选为leader可用）标示该Paxos group中对应数据的上锁情况；
- **Txn mngr**：（多group leaders中选为coordinator leader可用）当要执行的txn跨多个Paxos group时，需要对这些groups leader进行再选举，选举出来coordinator leader & non-coordinator-participant leader。前者使用txn mngr来协调这个txn（协调方法：txn mngr对其下管理的Paxos leaders执行2PL，即“尝试拿锁-修改并释放”）。

然后我们就能看懂下图了：

![](https://github.com/Emilio66/CSDI/blob/master/img/10_2.png?raw=true)

**Directory** 里是一段连续的、具有公共前缀的 Keys，数据在不同的 Paxos Group 间是一个一个 directory 移动的，Movedir 是后台移动数据的任务，也可以用来添加和删除副本，它会开启一个事务用于转移数据的最后一部分，然后借助事务更新两个 Paxos Group 的元数据。事实上，Spanner 会将大的 directory 切分成 segment，segment也会被保存在不同的分组上，转移时是转移 segment。

### TrueTime

**API**

`TT.now()`：返回一个时间段[earliest, latest]，保证被调用的一刻，所有span server的本地时间都处在这个范围内；

`TT.after(t), TT.before(t)`：检查是否所有spanserver都经历了t；或都还没有经历t。

**TrurTime 实现方式**

在底层，TrueTime API 使用的时间是 GPS 和原子钟。TrueTime 是由每个 datacenter 上面的许多 time master 机器和每台机器上的一个 timeslave daemon 来共同实现的。大多数 master 都有具备专用天线的 GPS 接收器，这些 master 在物理上是相互隔离的，这样可以减少天线失效、电磁干扰和电子欺骗的影响。剩余的 master (我们称为 Armageddon master) 则配备了原子钟。所有 master 的时间 参考值都会进行彼此校对。每个 master 也会交叉检查时间参考值和本地时间的比值，如果二者差别太大，就会把自己驱逐出去。每个 daemon 会从许多 master（可能是附近的也可能是很远的 datacenter的 ） 中收集投票，获得时间参考值，从而减少误差

### **Concurrency Control**

Clients通过location proxy来定位spanserver，将请求交由该spanserver处理。

#### 主要的txn类别

- read-write txn：普通的读写txn；
- read-only txn：一个只读事务必须事先被声明不会包含任何写操作，在一个只读事务中的读操作，在执行时会采用一个系统选择的时间戳，不包含锁机制，因此，后面到达的写操作不会被阻塞，可以到任何足够新的副本上去执行。
- snapshot reads：一个快照读操作，是针对历史数据的读取，执行过程中，不需要锁机制。

对于 read-only txn 和 snapshot reads 而言，当一个服务器失效的时候，客户端就可以使用同样的时间戳和当前的读位置，在另外一个服务器上继续执行读操作。

#### Read-write txns

1. （**spanserver执行部分**）对要读的数据，向对应的group leader拿读锁（如果拿不到就放弃，重新拿，即wound-wait）；
2. 执行本地读写（外部不可见）；
3. 修改完成，开始2PC。选择coordinator group，将修改发送给coordinator leader和non-coordinator-participant leader；
4. （**每个non-coordinator-participant leader执行部分**）收到txn修改内容后，选择本地最新成功的txn commit timestamp作为"prepare timestamp"返回给coordinator leader；
5. （**coordinator leader执行部分**）获得每个leader相应的写锁；
6. 等待所有participant leader的"*prepare timestamps*"，选择最大的"**s**"。再将**s**与 `TT.now().latest` 和本地最新成果的txn commit timestamp比较，取最大的作为commit timestamp，赋值给"**s**"；
7. 持续调用TrueTime获取interval，等待s < `TT.now().earliest`，即`TT.after(s)` 为 true，确保所有在s之前的txn都全局生效；
8. 以s为commit timestamp提交当前txn，并反馈client；
9. 释放锁。

**说明**：第六步中，实际上，就是想拿到一个timestamp **s**作为当前txn的commit timestamp。由于txn的顺序是遵循他们的commit timestamp，这个**s**就要保证大于之前所有txn commit timestamp。所以第七步等待花的时间就是用来确保所有机器的时间都“经历”了这些“之前的”txn的commit timestamp时间点，这些“之前的”txn此刻确定全局可见。

#### Read-Only Txns

首先，提取所有会被读到的key作为scope，然后分类讨论：

1. 如果scope都落在一个Paxos group：将这个RO txn发送给group leader；leader调用`LastTS()`获取本spanserver最近一次的write commit timestamp作为RO txn的timestamp并执行；
2. 如果scope跨多个Paxos groups：读取`TT.now().latest`作为当前RO txn的timestamp并执行。

说明：以上两种处理都能保证这次读在所有已全局生效的写之后