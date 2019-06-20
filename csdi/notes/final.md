# Final Review

## NVM

NVM主要特点是NVM的诞生使得非易失性来到了RAM这一层，这使得它变化很大。

必须知道**为什么存储变了，CPU指令要变，背后的原因是什么**？原因在于原来memory的组织不考虑持久性，所以如果中间坏了，内存直接丢掉，没有关系，而内存变成一个可持久的部分后，也就意味着当数据进入内存后，状态就确定下来了，也就意味着内存中的数据在任何一个时间点数据都得是正确的，原来不需要保证。原来是disk保证，那我可以人为控制，确保什么时候写入，保证disk状态的正确，但是现在要保证内存状态的正确，而原来从cache去memory取数据的时候是不受人为控制的，由计算机自动处理，而在NVM情况下，原来的disk是现在的memory，原来的memory是现在的cache，原来的memory刷到disk的时候是人自己控制的，我batch的提交，然后按照一定的顺序写回去，这样保证数据都是一致完整的；而现在cache刷到memory是由硬件背后的自己完成，和软件没有关系，你完成不了，于是就引出了现在这个问题，所以对于NVM的问题，最重要的一个点是如何保证内存里面的数据始终是正确的一致的。

解决这个问题第一个需要硬件的支持，需要特定的指令，来完成cache往内存刷新的事情，原来不可以控制，那么现在就要提供一个控制的手段，因此，它有了新的指令，然后呢，为了从内存刷到disk中，它本来有各种各样的优化，会先在内存中 buffer 起来，然后批量的写回，在buffer的时候还可以调整顺序，现在变成 cache 刷到内存中，也要考虑这些问题，我怎么刷（第一个，我有了手段，第二个，我要考虑怎么刷，用什么机制，第三个，要考虑怎么测试，怎么验证）。

于是介绍了PMFS，主要三个工作：

* 硬件原语支持，`pm_wbarrier`
* 重新设计了系统结构，原来从内存到disk是块操作，现在cache到内存是 byte addressable
* 保证数据是正确的，原来的方式 atomic point，在 point 之前之后都是对的，关于point的寻找，有两种经典操作，一个 COW，一个是 Journaling。前者的方式牵涉到具体的修改，会牵涉到数据结构，需要case by case的做，而 Journaling 方案做之前先做 log，log 记录下来之后代表我做了，log里面可以加入 atomic 记录，最后后台把 log 上的变化更新上去就是 Journaling。现在的问题是内存里的文件系统我该选择哪种方案？在PMFS中，如果数据是分散的且小的，这种时候它适合 Journaling ，原因在于 COW 粒度很大，如果修改很小，那么很多空间都浪费了，而且如果它很分散，那么我调整的时候要调整很多指针，要很多的指针更新 atomic 完成，那么 atomic操作不知道辣么多指针同步更新，所以 COW 不适合这种场景。反过来，集中的大数据用 COW，同时说为什么  Journaling 不适合这种呢？因为 Journaling 先要写 log，然后还要更新到最终地方，意味着要**写两次**，有一个写放大的效果，当修改数据很多的时候，意味着放大的数据也很多，造成影响很大。

## RAMCloud

看分布式场景下数据可靠性的问题，在这个场景第一个它是一个内存的场景，现在很多内存计算的场景，因为要求 latency 越来越高，IO 太慢了，所以把数据都扔到内存中。我们之前已经讲了 NVM，它当然会带来一些新的问题和变化，在不考虑 NVM 的情况下，也有新问题。原来内存只是 cache，数据都在 disk，现在在 in-memory 的场景下，我们假设的场景是数据都在内存里，内存是我们操作处理数据的地方，这时候就要考虑计算过程中如果 crash，如果出现 failure 的情况，这些数据丢了怎么办，然后如果要提高数据处理速度，我们必须要在计算过程中，就要尽量避免 IO 操作，那么你没有 IO操作，又怎么保证数据不会丢这件事情呢？那么必须要有 IO 操作，那么又该怎么设计？所以这个 topic 是在 in-memory 计算下，怎么处理内存数据可靠性问题（现在 cache 是计算的第一现场）。NVM当然是一个解决方案，还没有 NVM的时候，提出了 RAMCloud。

RAMCloud 通过分布式的方式提供高可靠性，数据可以出现在多台机器内存里，从而达到高可用，当出现一部分 failure 的情况下，可以保证数据不会丢。如果整个 cluster 挂了，数据还是会丢，所以还要考虑durability，最后还要考虑当机器 crash 之后，能够尽快恢复，所以还要关注快速恢复。

为什么使用分布式内存计算是可行的，因为现在网络速度变快了，由不需要disk访问，因此它的速度可以达到单机水准。问题主要在于 crash，RAMCloud 提出了一个 bufferd logging。

原来 log 是放在 disk 上的，原因是 disk 访问很慢，而 log 是 append的，因此 log可以充分sequential 的操作，产生一个 log structure。那 log structure 原来在 disk 上，现在把 log structure 放到内存中，因为这是 in-memory 计算，那么在计算的时候内存里面直接采用 log structure，执行的时候把这些 log 分布式出去，到其他机器上，使得数据在多台机器上存在备份。然后通过**异步的方式去写disk**，因为本身内存是log结构，那么现在写入disk也是log结构。相当于有一个 disk 操作，多个内存相关的 IO操作，剩下都是异步操作。因此整体的设计是，怎么让平时（备份）尽量快，当出现 crash 的时候要尽快的 recover。原来是从 disk 去 recover，而现在它利用数据分布在多台机器上，因此可以并行的去恢复一台机器，集所有包含该数据的机器的资源去恢复机器，加速 recover。

所以我们思考一下内存计算场景带来的区别：内存场景存储大量数据，因此需要大规模的集群，带来的问题是 crash，但是同样好处是可以集合所有机器去恢复一台机器。原来单点的瓶颈可能在 disk 上，因为要从 disk 上去读取大量数据，现在因为有上百台机器，每台机器只需要恢复1%左右的数据，disk 就不会成为瓶颈，但是 IO 可能会成为瓶颈；另外，最后恢复也可能成为瓶颈，因为要把拿来的数据重新变回log结构，这件事是单点完成。

为了使得最后一步的恢复不成为瓶颈，RAMCloud 的解决方案是 on demand 的完成，一边恢复服务一边恢复数据。它的假设是这些数据不是马上被访问的，访问的数据先替你恢复，更快的去恢复服务，后台再慢慢把数据恢复出来，相当于把数据 recovery 和 service 重新运行 overlay 在一块。当然还有其他方案，比如先把 workload 分到其他所有机器上，有这些机器先分担，然后等挂掉的这台机器恢复后，再把任务拿回来。

## Scalability Lock

所谓 scalability 就是你有 n 倍资源时，能不能把性能提升 n 倍。Scalability 决定于你程序中多少部分是并行的，多少部分不是并行的。因此，要让你的代码尽量做到大部分代码都可以是并行的。串行会很大的影响系统性能，比如一个 IO 操作 FD，FD 是由操作系统统一管理，因此它会使得无关的应用程序在操作系统这个层次上面产生关联，这个操作于是变成一个不能并行操作，于是整体会使程序之间变得不并行起来，自己写的程序会被其他程序影响，因为共用了一个操作系统，用的资源越多，性能反而越差，背后的原因就是硬件上的限制——CPU都是并行的，但是会有共享资源，数据会贡献，背后的结构导致代码可扩展性很差。

举个例子，Lock。一个人拿锁了后，需要通知别人，需要 watch 这个锁，等这个人释放锁以后，再把这个锁拿过来。这个看起来很简单的串行操作，但是为什么争抢的人越多，性能反而越差呢？这个性能虽然不会更好，但至少应该不变。背后的原因在很多人争抢的时候，硬件层面存在串行通讯的情况。因为所有人都抢，只有一个人能抢到，那么你每次都要参与抢，参与抢的人越多，你的效率就越低。因此争抢是造成 scalability 性能的问题，锁通知是串行的通知，人越多，锁通知的越多，一次 lock 的交换时间就会很长，因为和很多无关的人联系，我的锁只能给一个人，很多人拿不到，那我也要通知你你拿不到锁了。

针对这个问题，期望的应该是不管人多少，拿锁的行为的是串行的，性能应该是直线，不会随着人多变差。因此可以考虑后面能够拿到锁的人打交道，把所有人串起来，而不是通过广播，通过串起来的方式一个通知一个，没有通知到你的时候进行排队，从而避免性能骤降。因此 MSC Lock 就是每个后面的人呢都等待前面那个，每个人只和后面的人有关系，避免性能骤降。

看看MSC的实现（PPT），实现不是重点，重点是怎么设计出来的，为什么这么设计。

## Tx Memory

事务内存能够在内存级别上支持有具有ACID的特性的操作。其实内存上的 tx 操作就是一个 COW —— 在一个位置先写，旧的 tx 相关的数据不动，当最后 commit 时，新的内容替代旧内容。实际上，解决 atomic 的问题，general 方案上来说无非就是 COW。在 memory 层面比较适合的方案也是COW，如果用 COW ，需要有一些 policy 去管理内存，比如写在哪里；然后呢，要能够去 track 这些数据；当数据产生变化，被调度的时候要去 instead 旧的数据。所以，Tx Memory 被实现成 cache 级的 memory，因为在 cache 上，本身就要解决数据之间的 consistency 的问题，本身就要去 track 每一个数据谁访问过了，这个信息就可以帮助你去维护相应的事件，相当于这种实现可以减少新的硬件。

Tx Memory 有一个很长的发展历史，最后终于现在 CPU 里都有 tx memory 的支持。当然还没有一个定论关于这部分代码如何去写的一个很好的支持。但是最后的测试，实际上会发现所谓 tx memory 提供的性能实际上和 fine grained lock 的性能差不多，fined grained lock 就是每一个数据都有一个 lock 去管理，它也是细粒度的，无关的操作也不会形成冲突。tx memory 无关的操作也不会出现冲突，大家都是 fine grained 进行管理，tx memory 的好处是由硬件去 track 操作的变化，软件是由软件去完成，所以它有软件管理的 overhead；但是反过来呢，fined grained lock 是一个对象一个 lock，并不是真的一个地址一个lock的概念，因此对于对象可以比较好的去保护，而 tx memory 是硬件上的保护，虽然 track 开销比较小，但是它的问题是所有的数据都要去 track，所有它会有很多无用的 track 操作，所以最后平衡下来，fined grained lock 和 tx memory 的性能是差不多的。

## NoSQL

NoSQL 是 not only SQL，原来 SQL系统为了支持 tx，支持结构化的 table，所以它的 scalability 比较差，因为table 不好切分，tx 如果有conflict的话在分布式场景下去 track tx 之前的关系效率比较低。因此原来传统的 SQL 系统在分布式场景表现很差。所以引出 NoSQL。

NoSQL 只是从设计的角度提供一些 tradeoff，让你选择你需要的，不需要的剔除出去，不让它影响性能。所以，优化的方向就可能有数据的划分啦，replication啦，直接不支持join啦，或者是 in-memory 的数据库。NoSQL 更多的关注数据规模大，非结构化，并发情况高，可扩展性高。

NoSQL 有大量的产品，因为有很多优化选项，组合不同几个就会是新的产品。**所以 NoSQL 的一个问题也是怎么去选择合适的 NoSQL 系统，怎么支持业务需求的变化；NoSQL也会造成开发一定的困难，程序需要解决一些冲突，对系统维护开销很大**。

传统 SQL 系统的特征，NoSQL 可能不具备其中的某一条：

* Relational model，用 table，用 tree，用 view；
* Powerful query language，有各种 SQL 语句，支持 join，filter，sort等操作；
* Transactional semantics；
* Predefined schemas，预先设计好模型；
* Strong consistency between replicas。

## NewSQL

NewSQL 又是针对 NoSQL 的缺点推出的，它是从 NoSQL 回退到 SQL 系统，但是它继承了 NoSQL 的一些优点，比如数据怎么去 Model 的，数据的划分，副本等，同时，它做到了一些重要的特性仍然需要保留，这些特性就是和上层用户打交道的接口层次，比如，它需要支持 TX。所谓 NewSQl系统就是既需要 NoSQL 系统的扩展性，也需要 SQL 系统 ACID 这种比较容易编程的接口。

NewSQL 的设计就是 SQL 作为主要的 interface，而底层采用新的设计。当然 NewSQL 也面临一个问题，它当然不好解决，如果好解决 SQL 系统就已经解决了。NewSQL 最大的问题在于为了做一件事，额外做的事情开销很大，有用的操作可能只占总消耗的 10%，但是为了支持这一块，需要提供ACID，需要有 logging，需要有locking和latching，为了支持不同的查询需要有index，还要有 buffer 去支持 isolation，这些开销很高。

所以 NewSQL 一个主要的研究在于怎么降低额外操作的性能的开销，比如新的算法，新的数据结构，很重要的一点是新的硬件。比如利用新的硬件解决分布式下事务可扩展性差，介绍了一个DrTM系统。

对于这个问题的解决，DrTM其实从多个角度解决：

* 对于 recovery，它有 NVM，用 in-memory 的 logging 来提高；
* 对于 locking，使用 RDMA 去提升；
* 对于 latching 这种单机里面的锁，使用 tx memory 去提升；
* 对于 buffer，使用完全内存的存储，这样就不需要 buffer 了。

因此当有新的硬件时，可以用它尝试来解决原来软件的问题，提高性能；也就意味着需要重新设计系统，但是为了不影响应用，应当保持应用程序和系统之间的接口保持一致。

当然利用新的硬件，需要大量重新设计软件，把不同的硬件特性很好的 group 起来。

key value store 使用的方式是什么？数据结构是什么？读写方式是怎么样的？

## Latency

现在 latency 要求越来越高，因为它直接关联用户体验，lantency 受到多种因素的影响，这种影响超出我们原来的理解，在过去人们关注的是latency是由于程序造成的，然而实际上 latency 实际上会受到大量外围的因素，这就跟我们之前讲系统中的 scalability 一样，程序可能没毛病，毛病在系统。

外围因素有：

* 竞争：当发生竞争，只有争过别人才能拿到资源，开销就很大；
* 倾斜的访问模式：访问模式可能不是均匀的，而是skewed的，这样会造成一些集中，会带来竞争问题，资源不够问题，不平衡问题，造成性能下降；
* 队列延迟：当service进行服务的时候，并发数很高时，时间花费可能主要在等待被服务的时间；
* 后台行为：比如碰到 gc，碰到数据压缩的时候，整个性能就会出现波动，波动的结果就是一些机器性能会受到影响。

所以 Latency 指的并不是所有机器的latency，而是只是一小部分，称为 tail latency。但是对于用户的角度而言，受到 tail latency 影响的用户就会成为一个很大的问题。tail latency 的解决方案往往是使得整体 latency 变差，但是防止高 latency 的出现。

具体方案可以见[此](tail-latency.md)。

## NFV

NFV，网络功能虚拟化，就是去虚拟化Firewall，IDS和DNS这些网络功能，也就是通过软件模拟实现这些网络功能。NFV 实际上是软件替代硬件，有些工作是要不惜一切代价用新的硬件，有些工作是努力省钱。网络硬件太贵，要用软件，省钱。NFV 主要关注的是软件换硬件后，性能尽量下降的不多。介绍了 Click 系统。

NFV 一开始先解决性能，解决 package 处理速度，当性能处理好后，再去替换原来的硬件。在性能处理完后，NFV 不满足这个，还要做 Virtualization，不但要一台机器模拟一个 network function，我要模拟各种各样 network function，要一个对多个，省更多的钱。

## Graph Query

Graph Query 和 Graph Analytics 不太一样，它主要关注的是不同的 workloader，两者的数据虽然很相似，但 workloader 不一样。Graph Analytics 是分析类的，访问整个数据，去分析整个图上面的特征，行为，哪些重要，哪些不重要，哪些相互之间有关系，我从这点到那点要多少时间。Graph Query 解决的是具体的个人问题，单独的问题，它访问的是一小部分相关的数据，然后做计算。Graph Analytics 是一个个的去做任务，因为时间很长，并行反而会抢占资源，CPU还要互相协调，只要每个任务做的时候资源全用上，一个个做没关系；Graph Query 每个都是小的任务，资源都用不完，因此当然是并行。所以，Graph Query 的优化方向最根本无非是把资源都用上，提高硬件资源效率，不要有空闲，不要做无用的事情。

Wukong 是一个分布式内存图存储系统，具体技术见[此](graph-query.md)。

## Bug

写程序总有 Bug，在这里的重点是用系统的方法去寻找 bug 然后 fix 它。所以要找到这些 bug 的共性，从共性出发，找到一类 bug。

先理解问题，哪来的 bug？Bug 有哪些特征和行为？Bug 原因是什么？对应找出解决 bug 方法，最后再生成 patch。

用系统的方法找 bug 不单单是可以找到更多的bug，还可以带来程序性能的提高，还可以帮助编译器和其他工具完善它们工具的能力。比如如果 bug 和编译器相关，可以指导编译器去避免这些 bug；可以从程序设计角度修改 bug；可以让硬件去修改 bug，只要需求足够重要，带来效果足够好。