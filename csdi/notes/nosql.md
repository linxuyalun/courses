# NoSQL & BigTable

* [GitHub Notes](https://github.com/Emilio66/CSDI/blob/master/9_Bigtable_%E7%8E%8B%E8%8B%B1%E8%89%BA_%E5%BC%A0%E9%87%91%E7%9F%B3_%E5%A7%9C%E5%BE%B7%E6%99%BA.md)
* [谷歌技术"三宝"之BigTable](https://blog.csdn.net/OpenNaive/article/details/7532589)

## NoSQL

1. 关系型数据库遵循ACID规则
2. NoSQL，指的是非关系型的数据库。NoSQL有时也称作Not Only SQL的缩写，是对不同于传统的关系型数据库的数据库管理系统的统称。NoSQL用于超大规模数据的存储。（例如谷歌或Facebook每天为他们的用户收集万亿比特的数据）。这些类型的数据存储不需要固定的模式，无需多余操作就可以横向扩展。
3. **CAP**:  CAP指出对于一个分布式计算系统来说，不可能同时满足以下三点:
   1. 一致性(Consistency) (所有节点在同一时间具有相同的数据)
   2. 可用性(Availability) (保证每个请求不管成功或者失败都有响应)
   3. 分隔容忍(Partition tolerance) (系统中任意信息的丢失或失败不会影响系统的继续运作)
4. CAP理论的核心是：一个分布式系统不可能同时很好的满足一致性，可用性和分区容错性这三个需求，最多只能同时较好的满足两个：
   1. CA: 相当于几乎不是分布式，单点集群，满足一致性，可用性的系统，通常在可扩展性上不太强大。 例子：RDBMS
   2. CP - 满足一致性，分区容忍必的系统，通常性能不是特别高。 例子：MogoDB、BigTabHBase、Redis
   3. AP - 满足可用性，分区容忍性的系统，通常可能对一致性要求低一些。 例子：CouchDB，DynamoDB、Riak
5. NoSQL的优点/缺点
   1. 优点: 高可扩展性、分布式计算、低成本、架构的灵活性、半结构化数据、没有复杂的关系
   2. 缺点: 没有标准化（各种情况不同的tradeoff）、有限的查询功能（到目前为止）、最终一致是不直观的程序

## BigTable

### Data Model

本质上说，Bigtable是一个键值（key-value）映射。按作者的说法，Bigtable是一个稀疏的，分布式的，持久化的，多维的排序映射。

Bigtable的Data model是一个稀疏的、分布式的、持久化存储的多维度排序Map。Map由key和value组成：Map的索引key是由行关键字、列关键字以及时间戳组成；Map中的每个value都是一个未经解析的byte数组。即 `(row:string, column:string,time:int64) -> string`

先来看看多维、排序、映射。Bigtable的键有三维，分别是行键（row key）、列键（column key）和时间戳（timestamp），行键和列键都是字节串，时间戳是64位整型；而值是一个字节串。可以用 (row:string, column:string, time:int64)→string 来表示一条键值对记录。

行键可以是任意字节串，通常有10-100字节。对同一个行关键字的读或者写操作都是原子的（不管读或者写这一行里多少个不同列），这个设计决策能够使用户很容易的理解程序在对同一个行进行并发更新操作时的行为。 Bigtable通过行关键字的字典顺序来组织数据。表中的每个行都可以动态分区。每个分区叫做一个"Tablet"，Tablet是数据分布和负载均衡调整的最小单位。**这样做的结果是，当操作只读取行中很少几列的数据时效率很高，通常只需要很少几次机器间的通信即可完成。**用户可以通过选择合适的行关键字，在数据访问时有效利用数据的位置相关性，从而更好的利用这个特性。举例来说，在Webtable里，通过反转URL中主机名的方式，可以把同一个域名下的网页聚集起来组织成连续的行。具体来说，我们可以把maps.google.com/index.html的数据存放在关键字com.google.maps/index.html下。把相同的域中的网页存储在连续的区域可以让基于主机和域名的分析更加有效。

行的读写都是原子性的。Bigtable按照行键的字典序存储数据。Bigtable的表会根据行键自动划分为片（tablet），片是负载均衡的单元。最初表都只有一个片，但随着表不断增大，片会自动分裂，片的大小控制在100-200MB。行是表的第一级索引，我们可以把该行的列、时间和值看成一个整体，简化为一维键值映射，类似于：

```bash
table{
  "1" : {sth.}, # 一行
  "aaaaa" : {sth.},
  "aaaab" : {sth.},
  "xyz" : {sth.},
  "zzzzz" : {sth.}
}
```

列是第二级索引，每行拥有的列是不受限制的，可以随时增加减少。为了方便管理，列被分为多个列族（column family，是访问控制的单元），一个列族里的列一般存储相同类型的数据。一行的列族很少变化，但是列族里的列可以随意添加删除。列键按照family:qualifier格式命名的。这次我们将列拿出来，将时间和值看成一个整体，简化为二维键值映射，类似于：

```bash
table{
  # ...
  "aaaaa" : { # 一行
    "A:foo" : {sth.}, # 一列
    "A:bar" : {sth.}, # 一列
    "B:" : {sth.} # 一列，列族名为B，但是列名是空字串
  },
  "aaaab" : { # 一行
    "A:foo" : {sth.},
    "B:" : {sth.}
  },
  # ...
```

或者可以将列族当作一层新的索引，类似于：

```bash
table{
  # ...
  "aaaaa" : { # 一行
    "A" : { # 列族A
      "foo" : {sth.}, # 一列
      "bar" : {sth.}
    },
    "B" : { # 列族B
      "" : {sth.}
    }
  },
  "aaaab" : { # 一行
    "A" : {
      "foo" : {sth.},
    },
    "B" : {
      "" : "ocean"
    }
  },
  # ...
}
```

时间戳是第三级索引。Bigtable允许保存数据的多个版本，版本区分的依据就是时间戳。时间戳可以由Bigtable赋值，代表数据进入Bigtable的准确时间，也可以由客户端赋值。

数据的不同版本按照时间戳降序存储，因此先读到的是最新版本的数据。为了减轻多个版本数据的管理负担，我们对每一个列族配有两个设置参数，Bigtable通过这两个参数可以对废弃版本的数据自动进行垃圾收集。用户可以指定只保存最后n个版本的数据，或者只保存“足够新”的版本的数据（比如，只保存最近7天的内容写入的数据）。我们加入时间戳后，就得到了Bigtable的完整数据模型，类似于：

```bash
table{
  # ...
  "aaaaa" : { # 一行
    "A:foo" : { # 一列
        15 : "y", # 一个版本
        4 : "m"
      },
    "A:bar" : { # 一列
        15 : "d",
      },
    "B:" : { # 一列
        6 : "w"
        3 : "o"
        1 : "w"
      }
  },
  # ...
}
```

查询时，如果只给出行列，那么返回的是最新版本的数据；如果给出了行列时间戳，那么返回的是时间小于或等于时间戳的数据。比如，我们查询"aaaaa"/"A:foo"，返回的值是"y"；查询"aaaaa"/"A:foo"/10，返回的结果就是"m"；查询"aaaaa"/"A:foo"/2，返回的结果是空。

![](https://img-my.csdn.net/uploads/201205/04/1336139889_6039.jpg)

图1是Bigtable论文里给出的例子，Webtable表存储了大量的网页和相关信息。在Webtable，每一行存储一个网页，其反转的url作为行键，比如maps.google.com/index.html的数据存储在键为com.google.maps/index.html的行里，反转的原因是为了让同一个域名下的子域名网页能聚集在一起。图1中的列族"anchor"保存了该网页的引用站点（比如引用了CNN主页的站点），qualifier是引用站点的名称，而数据是链接文本；列族"contents"保存的是网页的内容，这个列族只有一个空列"contents:"。图1中"contents:"列下保存了网页的三个版本，我们可以用("com.cnn.www", "contents:", t5)来找到CNN主页在t5时刻的内容。

### Components of BigTable

Bigtable依赖于google的几项技术。用GFS来存储日志和数据文件；按SSTable文件格式存储数据；用Chubby管理元数据。

GFS参见[谷歌技术"三宝"之谷歌文件系统](http://blog.csdn.net/opennaive/article/details/7483523)。BigTable的数据和日志都是写入GFS的。

**SSTable**：SSTable的全称是Sorted Strings Table，是一种不可修改的有序的键值映射，提供了查询、遍历等功能。每个SSTable由一系列的块（block）组成，Bigtable将块默认设为64KB。在SSTable的尾部存储着块索引，在访问SSTable时，整个索引会被读入内存。每一个片（tablet）在GFS里都是按照SSTable的格式存储的，每个片可能对应多个SSTable。

**Scheduler**：BigTable集群往往运行在一个共享的机器池中，池中的机器还会运行其它各种各样的分布式应用程序，BigTable的进程经常要和其它应用的进程共享机器。BigTable依赖集群管理系统在共享机器上调度作业、管理资源、处理机器的故障、以及监视机器的状态。

**Chubby**：负责master的选举，保证在任意时间最多只有一个活动的Master；存储BigTable数据的引导程序的位置；发现tablet服务器，以及在Tablet服务器失效时进行善后；存储BigTable的模式信息（每张表的列族信息）；以及存储访问控制列表。Chubby是一种高可用的分布式锁服务，Chubby有五个活跃副本，同时只有一个主副本提供服务，副本之间用Paxos算法维持一致性，Chubby提供了一个命名空间（包括一些目录和文件），每个目录和文件就是一个锁，Chubby的客户端必须和Chubby保持会话，客户端的会话若过期则会丢失所有的锁。关于Chubby的详细信息可以看google的另一篇论文：The Chubby lock service for loosely-coupled distributed systems。Chubby用于片定位，片服务器的状态监控，访问控制列表存储等任务。

**Table, Tablet 和 SStable 的关系**：Tablet是从Table中若干范围的行组成的一个相当于原Table的子表，所以多个Tablet就能组成一个Table。而SSTable是Tablet在GFS文件系统中的**持久化存储的形式**，即Tablet在Bigtable中，是存在一个一个SSTable格式的文件中的。 它们的关系可以由下图来总结： 

![](https://github.com/Emilio66/CSDI/blob/master/img/9-1.png?raw=true)



### Refinement

- 局部性群组（Locality groups）
  客户程序可以将多个列族组合成一个局部性群族。对每个tablet中的每个局部性群组都会生成一个单独的SSTable。将通常不会一起访问的列族分割成单独的局部性群组使读取操作更高效。此外，可以以局部性群组为单位指定一些有用的调整参数。比如，可以把一个局部性群组设定为全部存储在内存中。设定为存入内存的局部性群组的SSTable依照惰性加载的策略装载进tablet服务器内存。加载完成之后，属于该局部性群组的列族的不用访问硬盘即可读取。
- 压缩
  客户程序可以控制一个局部性群组的SSTable是否压缩；用户指定的压缩格式应用到每个SSTable的块中（块的大小由局部性群组的调整参数操纵）。尽管为每个分别压缩浪费了少量空间，我们却受益于在只读取小部分数据SSTable的时候就不必解压整个文件了。许多客户程序使用双步（two-pass）定制压缩模式。第一步采用Bentley and McIlroy’s模式，这种模式横跨一个很大窗口压缩常见的长字符串；第二步采用快速压缩算法，即在一个16KB数据的小窗口中寻找重复数据。在Webtable的例子里，我们使用这种压缩方式来存储网页内容，在空间上达到了10:1的压缩比。
- Bloom过滤器
  一个读操作必须读取组成tablet状态的所有SSTable的数据。如果这些SSTable不在内存中，那么就需要多次访问硬盘。我们通过允许客户程序对特定局部性群组的SSTable指定Bloom过滤器，来减少硬盘访问的次数。通过bloom过滤器我们可以查询一个SSTable是否包含了特定行/列对的数据。对于某些应用程序，只使用了少量的tablet服务器内粗来存储Bloom过滤器，却大幅度减少了读操作需要的磁盘访问次数。Bloom过滤器的使用也意味着对不存在的行或列的大多数查询不需要访问硬盘。
- 通过缓存提高读操作的性能
  为了提高读操作的性能，tablet服务器使用二级缓存的策略。对tablet服务器代码而言，扫描缓存是第一级缓存，其缓存SSTable接口返回的键值对；Block缓存是二级缓存，其缓存从GFS读取的SSTable块。对于趋向于重复读取相同数据的应用程序来说，扫描缓存非常有效；对于趋向于读取刚读过的数据附近的数据的应用程序来说，Block缓存很有用。
- 提交日志的实现
  为了避免重复读取日志文件，我们首先把提交日志的条目按照关键字（table，row name，log sequence number）排序。排序之后，对一个特定tablet的修改操作连续存放，因此，随着一次询盘操作之后的顺序读取，修改操作的读取将更高效。为了并行排序，我们将日志文件分割成64MB的段，之后在不同的tablet服务器对每段进行并行排序。这个排序过程由master来协调，并且当一个tablet服务器指出它需要从一些提交日志文件中回复修改时排序被初始化。为了使修改操作免受GFS瞬时延迟的影响，每个tablet服务器实际上有两个日志写入线程，每个线程写自己的日志文件，并且同一时刻，两个线程只有其中之一是活跃的。如果写入活跃日志文件的效率很低，日志文件写入切换到另外一个线程，在提交日志队列中的修改操作就会由新的活跃日志写入线程写入。日志条目包含序列号，这使得恢复进程可以省略掉由于日志进程切换而造成的重复条目。
- Tablet恢复提速
  如果master将一个tablet从一个tablet服务器移到另外一个tablet服务器，源tablet服务器会对这个tablet做一次Minor Compaction。这个Compaction操作减少了tablet服务器日志文件中没有压缩的状态的数目，从而减少了恢复的时间。Compaction完成之后，该tablet服务器停止为该tablet提供服务。在真正卸载tablet之前，tablet服务器还会再做一次Minor Compaction，以消除tablet服务器日志中第一次minor compaction执行过程中产生的未压缩的状态残留。当第二次minor compaction完成以后，tablet就在不需要任何日志条目恢复的情况下被装载到另一个tablet服务器上了。
- 利用不变性
  因为SSTable是不变的，所以永久移除已被删除数据的问题就转换成对废弃的SSTable进行垃圾收集的问题了。每个tablet的SSTable都在注册在元数据表中。Master将废弃的SSTable作为对SSTable集合的“标记-清除”的垃圾回收而删除，元数据表则保存了root的集合。最后，SSTable的不变性使得分割tablet的操作非常快捷。与为每个子tablet生成新的SSTable集合相反，我们让子tablet共享父tablet的SSTable。

### Advantage and Limitation

由于BigTable针对数据存储进行包括压缩等各方面的优化，以及在事务的一致性上做出了让步，BigTable对于那些需要海量数据存储，高扩展性以及大量的数据处理，但又不要求强一致性的应用是十分适合的，比如Google Earth等。也因此，对于那些需要强一致性，需要同步更改多行数据的应用来说，BigTable是不合适的。