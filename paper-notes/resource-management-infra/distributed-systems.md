* [Distributed systems for fun and profit](http://book.mixu.net/distsys/index.html)

# Introduction

在本文中，我试图为分布式系统提供更易于理解的介绍。 对我而言，这意味着两件事：介绍需要的关键概念，以便有更好的时间阅读更严肃的文本，并提供一个足够详细的事情的叙述，了解正在发生的事情而不会卡住在细节上。 

在我看来，许多分布式编程都是关于处理以下两个内容：

* information travels at the speed of light
* independent things fail independently*

换句话说，分布式编程的核心是处理距离（呃！）和多台计算机（呃！）。 这些约束定义了可能的系统设计空间，我希望在阅读完之后，你将更好地了解距离，时间和一致性模型如何相互作用。

本文重点介绍了解数据中心中商业系统所需的分布式编程和系统概念。 你将学习许多关键协议和算法（例如，涵盖该学科中许多被引用次数最多的论文），包括一些新的令人兴奋的方法来研究最终的一致性，这些方法还没有进入大学教科书——例如CRDT 和CALM定理。

## 1. Basics

[第一章](#1.-distributed-systems-at-a-high-level)通过介绍一些重要的术语和概念，从高层次上介绍了分布式系统。 它涵盖了高级目标，例如可伸缩性，可用性，性能，延迟和容错; 这些是如何难以实现的，以及抽象和模型以及分区和复制如何发挥作用。

## 2. Up and down the level of abstraction

第二章深入研究了抽象和不可能结果。 它以Nietzsche引用开始，然后介绍系统模型和在典型系统模型中做出的许多假设。 然后讨论了CAP定理并总结了FLP不可能性结果。 然后转向CAP定理的含义，其中之一是应该探索其他一致性模型。 然后讨论了许多一致性模型.

## 3. Time and order

理解分布式系统的一个重要部分是了解时间和顺序。 如果我们无法理解和模拟时间，我们的系统将会失败。 第三章讨论时间和顺序，时钟以及时间，顺序和时钟（如矢量时钟和故障检测器）的各种用途。

## 4. Replication: preventing divergence

第四章介绍了复制问题，以及它可以执行的两种基本方法。 事实证明，大多数相关特征可以通过这种简单的表征来讨论。 然后，从最小容错（2PC）到Paxos讨论了用于维护单拷贝一致性的复制方法。

## 5. Replication: accepting divergence

第五章讨论了弱一致性保证的复制。 它引入了一个基本的协调方案，其中分区副本尝试达成协议。 然后，它讨论了亚马逊的Dynamo作为具有弱一致性保证的系统设计的示例。 最后，讨论了关于无序编程的两个观点：CRDT和CALM定理。

# 1. Distributed systems at a high level

> Distributed programming is the art of solving the same problem that you can solve on a single computer using multiple computers.

任何计算机系统都需要完成两个基本任务：

* 存储
* 计算

分布式编程是一种艺术，它可以用多台计算机解决单台计算机上的相同问题——通过，这种问题不再适用于单台计算机。

不是说非要用分布式系统。如果给定无限的资金和无限的研发时间，我们不需要分布式系统。 所有的计算和存储都可以在一个神奇的盒子上完成 ——一个单一的，令人难以置信的快速和令人难以置信的可靠系统。

但是，谁又拥有无限的资源呢。 因此，我们必须在一些现实世界的成本效益曲线上找到合适的位置。 在小规模的时候，升级硬件是一种可行的策略。 但是，随着问题规模的增加，当你想要解决的问题无法通过升级硬件解决或者升级硬件的成本太高的时候——欢迎来到分布式系统的世界。

目前的现实情况是，采用中档商品硬件——维护成本可以通过容错软件来降低。

计算主要受益于高端硬件，它们可以用内部内存访问来代替慢速网络访问。 但是当遇到在节点之间需要大量通信的任务时，高端硬件的性能优势就会受到限制。

![figure 1](http://book.mixu.net/distsys/images/barroso_holzle.png)

如上图所示，假设所有节点都有统一的内存访问模式，高端硬件和商用硬件之间的性能差距随着cluster的增大而减小。

理想情况下，添加新机器将线性地提高系统的性能和容量。 但当然这是不可能的，因为这些单独的计算机会产生一些额外的开销。 需要复制数据，必须协调计算任务等。 这就是研究分布式算法的原因——它们为特定问题提供了有效的解决方案，并提供了可行的指导，正确实施的最低成本是什么，以及什么是不可能的。

本文的重点是分布式编程和系统中一个普通但和商业上息息相关的设置：数据中心。 例如，我不会讨论由异常网络配置或共享内存设置中出现的特殊问题。 此外，重点是探索系统设计空间而不是优化任何特定设计——后者是更专业化的主题。

## What we want to achieve: Scalability and other good things

在我看来，一切都从对处理规模（size）的需要开始。

大多数事情在小规模上都是微不足道的，但是一旦超过一定的规模，数量或其他物理上有限的事情，同样的问题就会变得更加困难。 提起一块巧克力很容易，很难举起一座山。 很容易计算一个房间里有多少人，很难计算一个国家有多少人。

所以一切都从规模开始——可扩展性。可扩展性，简单的说，就是在可扩展系统中，当事情从小变大，事情不应该越来越糟。 这是另一个定义：

>[Scalability](http://en.wikipedia.org/wiki/Scalability): is the ability of a system, network, or process, to handle a growing amount of work in a capable manner or its ability to be enlarged to accommodate that growth.

什么是增长？ 几乎可以用任何方式衡量增长（人数，用电量等）。但是有三个特别有趣的事情要看：

* 规模可扩展性：添加更多节点应该使系统线性更快; 增长数据集不应该增加延迟
* 地理可扩展性：应该可以使用多个数据中心来减少响应用户查询所需的时间，同时以一种合理的方式处理跨数据中心延迟。
* 管理可扩展性：添加更多节点不应增加系统的管理成本（例如管理员与机器的比率）。

当然，在实际系统中，增长同时发生在多个不同的轴上；每个指标都只捕获了增长的某些方面。

可扩展系统是随着规模的增加而不断满足其用户需求的系统。 有两个特别相关的方面——性能和可用性——可以通过各种方式进行衡量。

### Performance (and latency)

> [Performance](http://en.wikipedia.org/wiki/Computer_performance): is characterized by the amount of useful work accomplished by a computer system compared to the time and resources used.

根据具体情况，这可能涉及实现以下一项或多项：

* 对于给定的工作，响应时间短/低延迟
* 高吞吐量（处理工作率）
* 计算资源利用率低

优化任何这些结果都需要权衡。 例如，系统可以通过处理更大批量的工作来实现更高的吞吐量，从而减少操作开销。 但是由于批处理，权衡对于各个工作的响应时间会更长。

低延迟——实现较短的响应时间——是性能中最有趣的方面，因为它与物理（而非金钱）限制有很强的联系。 使用钱去解决延迟比执行其他方面更困难。

延迟有很多非常具体的定义，但是通过词源去理解延迟真的很棒：

> Latency: The state of being latent; delay, a period between the initiation of something and the occurrence.

那么“latent”是什么意思？

> Latent: From Latin latens, latentis, present participle of lateo ("lie hidden"). Existing or present but concealed or inactive.

这个定义非常酷，因为它突出了延迟实际上是事件发生后变得可见需要的时间。

例如，假设感染了一种将人变成僵尸的空气传播病毒。 潜伏期是被感染之后和变成僵尸之间的时间——这就是延迟：已经发生的事情被隐藏起来的时间。

让我们暂时假设我们的分布式系统只执行一个高级任务：给定一个查询，它会获取系统中的所有数据并计算单个结果。 换句话说，将分布式系统视为数据存储，能够在其当前内容上运行单个确定性计算（函数）：

```
result = query(all data in the system)
```

那么，对于延迟而言重要的不是旧数据的数量，而是新数据在系统中“生效”的速度。 例如，延迟可以这样衡量——how long it takes for a write to become visible to readers。

基于这个定义的另一个关键点是，如果无事发生，就没有“潜伏期”。 数据不变的系统不会（或不应该）存在延迟问题。

在分布式系统中，存在无法克服的最小延迟：光速限制了信息传输的速度，硬件组件每次操作都会产生的最低延迟花费（想想RAM和硬盘驱动器以及CPU）。

最小延迟对查询的影响程度取决于这些查询的性质以及信息需要传输的物理距离。

### Availability (and fault tolerance)

可扩展系统的第二个方面是可用性。

>[Availability](http://en.wikipedia.org/wiki/High_availability): the proportion of time a system is in a functioning condition. If a user cannot access the system, it is said to be unavailable.

分布式系统使我们能够实现在单个系统上难以实现的理想特性。例如，单个机器无法容忍任何故障，因为它要么故障要么正常运行。

分布式系统可以采用一堆不可靠的组件，并在它们之上构建可靠的系统。

没有冗余的系统只能作为其底层组件可用。 使用冗余构建的系统可以容忍部分故障，因此可用性更高。 值得注意的是，“冗余”可能意味着不同的东西，具体取决于您所看到的内容——组件，服务器，数据中心等。

在公式上，可用性是：`Availability = uptime / (uptime + downtime)`

从技术角度来看，可用性主要是关于容错。 因为发生故障的概率随着组件的数量的增加而增加，所以系统应该能够进行补偿，以防止随着组件数量的增加而变得不那么可靠。

比如：

| Availability %         | How much downtime is allowed per year? |
| ---------------------- | :------------------------------------- |
| 90% ("one nine")       | More than a month                      |
| 99% ("two nines")      | Less than 4 days                       |
| 99.9% ("three nines")  | Less than 9 hours                      |
| 99.99% ("four nines")  | Less than an hour                      |
| 99.999% ("five nines") | ~ 5 minutes                            |
| 99.9999% ("six nines") | ~ 31 seconds                           |

从某种意义上说，可用性是一个比正常运行时间更广泛的概念，因为服务的可用性也会受到网络中断或拥有该服务的公司的影响（这是一个与容错无关但仍会影响系统可用性的因素）。 但是，如果不了解系统的每个特定方面，我们所能做的最好的是容错设计。

什么是容错？

> Fault tolerance: ability of a system to behave in a well-defined manner once faults occur

容错归结为：定义您期望的故障，然后设计一个容忍它们的系统或算法。

## What prevents us from achieving good things?

分布式系统受两个物理因素的限制:

* 节点数（随着所需的存储和计算能力增加而增加）
* 节点之间的距离（认为以光速进行信息传播）

在这些限制范围内工作：

* 独立节点数量的增加会增加系统故障的可能性（降低可用性并增加管理成本）
* 独立节点数量的增加可能会增加节点之间通信的需求（随着规模的增加而降低性能）
* 地理距离的增加会增加远程节点之间通信的最小延迟（降低某些操作的性能）

除了这些限制——这是物理限制的结果——剩下就是系统设计选择的世界。

性能和可用性都由系统的外部保证定义。 在较高的层面上，可以将保证视为系统的SLA（service level agreement）：如果我写数据，我可以多快在其他地方访问它？ 写完数据后，我有什么保证耐用性？ 如果我要求系统运行计算，它返回结果的速度有多快？ 当组件发生故障或停止运行时，这会对系统产生什么影响？

还有另一个标准，没有明确提及但暗示：可理解性。 当然，没有简单的指标来衡量什么是可理解性。

我有点想在物理限制下加入“可理解性”。 毕竟，对于人来说，这是一个硬件限制，我们很难理解任何涉及比手指更动人的东西。 这是错误和异常之间的区别——错误是不正确的行为，而异常是意外行为。 如果你更聪明，你会发现异常发生。

## Abstractions and models

这就是抽象和模型发挥作用的地方。 通过消除与解决问题无关的现实世界方面，抽象使事情更易于管理。 模型以精确的方式描述分布式系统的关键属性。 我将在下一章讨论多种模型，例如：

- System model (asynchronous / synchronous)
- Failure model (crash-fail, partitions, Byzantine)
- Consistency model (strong, eventual)

良好的抽象使得使用系统更容易理解，同时捕获与特定目的相关的因素。

在存在许多节点的现实与我们对“像单个系统一样工作”的系统的需求之间存在着紧绷关系。 通常，最熟悉的模型（例如，在分布式系统上实现共享内存抽象）太昂贵了。

制定较弱保证的制度具有更大的行动自由，因此可能具有更高的性能，但也可能难以推理。 人们更擅长推理像单个系统一样工作的系统，而不是节点集合。

人们通常可以通过暴露有关系统内部的更多细节来获得性能。 例如，在列式存储中，用户可以（在某种程度上）推断系统内键值对的位置，从而做出影响典型查询性能的决策。 隐藏这些细节的系统更容易理解（因为它们更像是单个单元，需要考虑更少的细节），而暴露更多真实细节的系统可能有更好的性能（因为它们更接近现实）。

几种类型的故障使得编写像单个系统一样的分布式系统变得十分困难。 网络延迟和网络分区（例如某些节点之间的总网络故障）意味着系统有时需要做出艰难的选择，以确定是否更好地保持可用但丢失一些无法实施的关键保证，或者在发生这些类型的故障时保证安全并拒绝客户端。

最后，理想的系统满足程序员的需求（干净的语义）和业务需求（可用性/一致性/延迟）。

## Design techniques: partition and replicate

数据集在多个节点之间分配的方式非常重要。为了进行任何计算，我们需要定位数据然后对其进行操作。

有两种基本技术可以应用于数据集。 它可以分割为多个节点（分区）以允许更多并行处理。 它还可以复制或缓存在不同的节点上，以减少客户端和服务器之间的距离，并提高容错能力（复制）。

> Divide and conquer - I mean, partition and replicate.

下图说明了这两者之间的区别：分区数据（下面的A和B）被分成独立的集合，而复制的数据（下面的C）被复制到多个位置。

![figure2](http://book.mixu.net/distsys/images/part-repl.png)

这是解决分布式计算中任何问题的组合拳。 当然，诀窍在于为具体实施选择正确的技术; 有许多算法实现复制和分区，每个算法都有不同的限制和优点，需要根据设计目标进行评估。

### Partitioning

分区是将数据集划分为更小的不同独立集合;这用于减少数据集增长的影响，因为每个分区都是数据的子集。

* 分区通过限制要检查的数据量并通过在同一分区中分配相关数据来提高性能；
* 分区通过允许分区独立失效来提高可用性，从而增加在牺牲可用性之前失效的节点数。

分区也是特定于应用程序的，因此在不了解具体细节的情况下很难说清楚。这就是为什么大多数文章重点放在复制，包括本文。

分区主要是根据认为的主要访问模式来定义分区，并处理来自独立分区的限制（例如跨分区的低效访问，不同的增长率等）。

### Replication

复制是在多台机器上复制相同的数据;这允许更多服务器参与计算。

> To replication! The cause of, and solution to all of life's problems.

复制——复制或再现某些东西，是分布式系统抵御延迟的主要方式。

* 复制通过用额外的计算能力和带宽在新的数据副本来提高性能；
* 复制通过创建数据的其他副本来提高可用性，从而增加了在牺牲可用性之前失败的节点数。

复制是关于提供额外的带宽，以及缓存重要的内容。它还涉及根据某种一致性模型以某种方式保持一致性。

复制允许我们实现可伸缩性，性能和容错。 害怕可用性或性能降低？ 复制数据以避免瓶颈或单点故障。 计算慢？ 在多个系统上复制计算。 慢I/O？ 将数据复制到本地缓存以减少延迟，或将数据复制到多台计算机上以提高吞吐量。

复制也是许多问题的根源，因为必须在多台机器上保持同步的数据的独立副本——这意味着确保复制遵循一致性模型。

一致性模型的选择至关重要：良好的一致性模型为程序员提供了干净的语义（换句话说，它保证的属性易于推理）并满足业务/设计目标，如高可用性或强一致性。

只有一个用于复制的一致性模型——强一致性——允许你进行编程，就好像未复制基础数据一样。其他一致性模型将复制的一些内部暴露给程序员。 但是，较弱的一致性模型可以提供较低的延迟和较高的可用性 - 并且不一定难以理解，只是不同。

## Further reading

- [The Datacenter as a Computer - An Introduction to the Design of Warehouse-Scale Machines](http://www.morganclaypool.com/doi/pdf/10.2200/s00193ed1v01y200905cac006) - Barroso & Hölzle, 2008
- [Fallacies of Distributed Computing](http://en.wikipedia.org/wiki/Fallacies_of_Distributed_Computing)
- [Notes on Distributed Systems for Young Bloods](http://www.somethingsimilar.com/2013/01/14/notes-on-distributed-systems-for-young-bloods/) - Hodges, 2013

[返回Introduction](#introduction)

