* [Distributed systems for fun and profit](http://book.mixu.net/distsys/index.html)

# Introduction

在本文中，我试图为分布式系统提供更易于理解的介绍。 对我而言，这意味着两件事：介绍需要的关键概念，以便有更好的时间阅读更严肃的文本，并提供一个足够详细的事情的叙述，了解正在发生的事情而不会卡住在细节上。 

在我看来，许多分布式编程都是关于处理以下两个内容：

* information travels at the speed of light
* independent things fail independently*

换句话说，分布式编程的核心是处理距离（呃！）和多台计算机（呃！）。 这些约束定义了可能的系统设计空间，我希望在阅读完之后，你将更好地了解距离，时间和一致性模型如何相互作用。

本文重点介绍了解数据中心中商业系统所需的分布式编程和系统概念。 你将学习许多关键协议和算法（例如，涵盖该学科中许多被引用次数最多的论文），包括一些新的令人兴奋的方法来研究最终的一致性，这些方法还没有进入大学教科书——例如CRDT 和CALM定理。

## Basics

[第一章](#1-distributed-systems-at-a-high-level)通过介绍一些重要的术语和概念，从高层次上介绍了分布式系统。 它涵盖了高级目标，例如可伸缩性，可用性，性能，延迟和容错; 这些是如何难以实现的，以及抽象和模型以及分区和复制如何发挥作用。

## Up and down the level of abstraction

[第二章](#2-up-and-down-the-level-of-abstraction)深入研究了抽象和不可能的结果。 它以Nietzsche引用开始，然后介绍系统模型和在典型系统模型中做出的许多假设。 然后讨论了CAP定理并总结了FLP不可能性结果。 然后转向CAP定理的含义，其中之一是应该探索其他一致性模型。 然后讨论了许多一致性模型.

## Time and order

理解分布式系统的一个重要部分是了解时间和顺序。 如果我们无法理解和模拟时间，我们的系统将会失败。 第三章讨论时间和顺序，时钟以及时间，顺序和时钟（如矢量时钟和故障检测器）的各种用途。

## Replication: preventing divergence

第四章介绍了复制问题，以及它可以执行的两种基本方法。 事实证明，大多数相关特征可以通过这种简单的表征来讨论。 然后，从最小容错（2PC）到Paxos讨论了用于维护单拷贝一致性的复制方法。

## Replication: accepting divergence

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

延迟有很多非常具体的定义，但是通过词源去理解延迟真的很有意思：

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

没有冗余的系统只能作为其底层组件可用。 使用冗余构建的系统可以容忍部分故障，因此可用性更高。 值得注意的是，“冗余”可能意味着不同的东西，具体取决于所看到的内容——组件，服务器，数据中心等。

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

容错归结为：定义一个期望的故障，然后设计一个容忍它们的系统或算法。

## What prevents us from achieving good things?

分布式系统受两个物理因素的限制:

* 节点数（随着所需的存储和计算能力增加而增加）
* 节点之间的距离（认为以光速进行信息传播）

在这些限制范围内工作：

* 独立节点数量的增加会增加系统故障的可能性（降低可用性并增加管理成本）
* 独立节点数量的增加可能会增加节点之间通信的需求（随着规模的增加而降低性能）
* 地理距离的增加会增加远程节点之间通信的最小延迟（降低某些操作的性能）

除了这些限制——这是物理限制的结果——剩下就是系统设计选择的世界。

性能和可用性都由系统的外部保证定义。 在较高的层面上，可以将保证视为系统的SLA（service level agreement）：如果我写数据，我可以多快在其他地方访问它？ 写完数据后，我有什么保证这些数据的耐用性？ 如果我要求系统运行计算，它返回结果的速度有多快？ 当组件发生故障或停止运行时，这会对系统产生什么影响？

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

# 2. Up and down the level of abstraction

在本章中，我们将在抽象层次上下移动，查看一些不可能的结果（CAP和FLP），然后为了性能而向下移动。

只要写过任何编程，那么抽象级别的概念就很熟悉。 我们始终在某种抽象级别工作，通过某些API与较低级别层接口，并可能为用户提供一些更高级别的API或用户界面。 计算机网络的七层OSI模型就是一个很好的例子。

我断言，分布式编程在很大程度上处理了分布的后果。 也就是说，我们对“像单一系统一样工作”的分布式系统的需求和现实情况多节点的分布式系统之间存在一个紧绷关系。 这意味着要找到一个好的抽象，以平衡可理解性和性能。

存在许多节点的现实与我们对“像单一系统一样工作”的系统的需求之间存在着紧张关系。 这意味着要找到一个好的抽象，以平衡可能与可理解和高效的东西。

当说X比Y更抽象时，到底意味着什么？ 首先，X不会引入任何新的或与Y基本不同的东西。事实上，X可能会删除Y的某些方面或以一种使它们更易于管理的方式呈现它们。 其次，假设从Y中移除的X对于手头的事情并不重要，X在某种意义上比Y更容易掌握。

正如[Nietzsche](http://oregonstate.edu/instruct/phl201/modules/Philosophers/Nietzsche/Truth_and_Lie_in_an_Extra-Moral_Sense.htm)写得那样：

>Every concept originates through our equating what is unequal. No leaf ever wholly equals another, and the concept "leaf" is formed through an arbitrary abstraction from these individual differences, through forgetting the distinctions; and now it gives rise to the idea that in nature there might be something besides the leaves which would be "leaf" - some kind of original form after which all leaves have been woven, marked, copied, colored, curled, and painted, but by unskilled hands, so that no copy turned out to be a correct, reliable, and faithful image of the original form.
>
>每个概念都源于我们将不平等的东西等同起来。 没有叶子完全等于另一个叶子，“叶子”的概念是通过对这些个体差异的任意抽象，通过忘记区别而形成的。现在，产生了这样一种想法：在自然界中，除了叶子之外，可能还有其他什么东西是“叶子”——一种最初的形式，在这种形式之后，所有的叶子都被编织、标记、复制、着色、卷曲和上色，但都是用不熟练的手完成的，所以没有一个复制品是原来形式的正确、可靠和忠实的形象。

从根本上说，抽象是假的。 每种情况都是独特的，每个节点都是如此。 但是抽象使得世界变得易于管理：更简单的问题陈述——没有现实世界——更易于分析，只要我们不忽视任何必要的东西，解决方案就可以广泛应用。

实际上，如果我们保留的东西是必不可少的，那么我们可以得出的结果将是广泛适用的。 这就是之前提及的不可能的结果如此重要的原因：它们采用最简单的问题表达方式，并证明在某些约束或假设中无法解决。

所有的抽象都忽略了一些有利于将现实中独一无二的事物等同起来的东西。 诀窍是摆脱一切不重要的东西。 你怎么知道什么是必要的？ 先验。

每次我们从系统规范中排除系统的某些方面时，我们都会冒险引入错误源和/或性能问题。 这就是为什么有时我们需要走向另一个方向，并有选择地介绍真实硬件和现实世界问题的某些方面。 重新引入一些特定的硬件特性（例如物理顺序性）或其他物理特性以获得性能足够好的系统。

考虑到这一点，当我们仍在处理一些仍然可以识别为分布式系统的东西时，我们可以保留的最少的“现实量”是多少？系统模型是我们认为重要的特性的规范；指定了一个特性之后，我们就可以看看一些不可能的结果和挑战。

## A system model

分布式系统的关键属性是分布。更具体地说，分布式系统中的程序：

* 在独立节点上并发运行；
* 网络连接可能引入不确定性和消息丢失的；
* 没有共享内存或共享时钟。

这里有很多含义：

* 每个节点同时执行一个程序；
* 知识是本地的：节点只能快速访问其本地状态，任何有关全局状态的信息都可能超时；
* 节点可能会故障并独立地从故障中恢复；
* 消息可能会延迟或丢失（这与节点故障无关——区分网络故障和节点故障并不容易）；
* 时钟不跨节点同步（本地时间戳与全局顺序不对应，无法轻易观察到）。

系统模型列举了与特定系统设计相关的许多假设。

> System model: a set of assumptions about the environment and facilities on which a distributed system is implemented

系统模型对环境和设施的假设各不相同。这些假设包括：

* 节点具有哪些功能以及它们可能如何失效；
* 通信链路如何运作以及它们如何失效；
* 整个系统的属性，例如关于时间和顺序的假设。

一个鲁棒性最强的系统模型是做出最弱假设的模型：为这样的系统编写的任何算法 能容忍不同的环境，因为它做出非常少且非常弱的假设。

换言之，我们可以通过作出强有力的假设来创建一个易于推理的系统模型。例如，假设节点没有失败意味着我们的算法不需要处理节点失败。显然，这样的系统模型是不现实的，因此很难应用到实践中。

让我们更详细地看一下节点，链接以及时间和顺序的属性。

### Nodes in our system model

节点充当计算和存储的主机。它们有：

* 执行程序的能力
* 能够将数据存储到易失性存储器（可能在发生故障时丢失）并进入稳定状态（可在故障后读取）
* 时钟（可能会或可能不会被认为是准确的）

节点执行确定性算法：本地计算，计算后的本地状态以及发送的消息由接收到的消息和接收消息时的本地状态唯一确定。

有许多可能的失效模型描述了节点失效的方式。在实践中，大多数系统都假设一个崩溃-恢复的失败模型：即节点只能通过崩溃来失败，并且在稍后的某个时刻崩溃后可以（可能）恢复。

另一种选择是假设节点可以通过任意方式的错误行为而失败。这被称为拜占庭容错。拜占庭式错误在现实世界的商业系统中很少被处理，因为对任意错误有弹性的算法运行起来更昂贵，实现起来也更复杂。我不会在这里讨论它们。

### Communication links in our system model

通信链路将各个节点相互连接，并允许消息以任意方向发送。许多讨论分布式算法的书籍假定每对节点之间都有单独的链接，这些链接为消息提供FIFO顺序，它们只能传递已发送的消息，并且发送的消息可能会丢失。

一些算法假定网络是可靠的：消息不会丢失，也不会无限期延迟。对于某些实际设置来说，这可能是一个合理的假设，但一般来说，最好考虑网络不可靠，并且会受到消息丢失和延迟的影响。

当网络出现故障，而节点本身仍在运行时，就会出现网络分区。出现这种情况时，消息可能会丢失或延迟，直到网络分区修复。分区节点可以被一些客户机访问，因此必须区别于崩溃的节点。下图说明了节点故障与网络分区的关系：

![](http://book.mixu.net/distsys/images/system-of-2.png)

很少对通信链路做进一步的假设。我们可以假设链接只在一个方向上工作，或者我们可以为不同的链接引入不同的通信成本（例如，由于物理距离造成的延迟）。但是，在商业环境中，除了长距离链路（广域网延迟）之外，这些问题很少被关注，因此我不会在这里讨论它们；更详细的成本和拓扑模型允许以复杂性为代价进行更好的优化。

### Timing / ordering assumptions

物理分布的结果之一是每个节点以独特的方式体验世界。这是不可避免的，因为信息只能以光速传播。如果节点之间的距离不同，那么从一个节点发送到另一个节点的任何消息将在不同的时间到达，并且可能在其他节点以不同的顺序到达。

时间假设是一种方便的简写，用于捕捉关于我们将这一现实考虑在内的程度的假设。两个主要的替代方案是：

> Synchronous system model: Processes execute in lock-step; there is a known upper bound on message transmission delay; each process has an accurate clock
>
> Asynchronous system model: No timing assumptions - e.g. processes execute at independent rates; there is no bound on message transmission delay; useful clocks do not exist

同步系统模型对时间和顺序施加了许多约束。它基本上假定节点具有相同的经历：发送的消息总是在特定的最大传输延迟内接收，并且进程在锁定步骤中执行。这很方便，因为它允许您作为系统设计者对时间和顺序进行假设，而异步系统模型则不是如此。

异步性是一种非假设：它只是假设你不能依赖于时间（或“时间传感器”）。

同步系统模型中的问题更容易解决，因为关于执行速度、最大消息传输延迟和时钟精度的假设都有助于解决问题，因为您可以根据这些假设进行推断，并通过假设不发生不方便的故障场景来排除它们。

当然，假设同步系统模型不太现实。现实世界中的网络容易出现故障，并且消息延迟没有严格的限制。现实世界中的系统充其量只是部分同步的：它们可能偶尔会正常工作并提供一些上限，但有时消息会无限期延迟，时钟也会不同步。我不会在这里真正讨论同步系统的算法，但您可能会在许多其他介绍性书籍中遇到它们，因为它们在分析上更容易（但不现实）。

### The consensus problem

在本文的其余部分，我们将改变系统模型的参数。接下来，我们将看看如何改变两个系统属性：

* 网络分区是否包含在故障模型中，以及
* 同步与异步时间假设

通过讨论两个不可能的结果（FLP和CAP）来影响系统设计选择。

当然，为了进行讨论，我们还需要引入一个问题来解决。我要讨论的问题是共识问题。

如果几个计算机（或节点）都同意一些价值，那么它们就会达成共识。更正式地说：

* Agreement: Every correct process must agree on the same value.
* Integrity: Every correct process decides at most one value, and if it decides some value, then it must have been proposed by some process.
* Termination: All processes eventually reach a decision.
* Validity: If all correct processes propose the same value V, then all correct processes decide V.

共识问题是许多商业分布式系统的核心。毕竟，我们希望一个分布式系统的可靠性和性能不必处理分布式带来的一些不好的后果（例如节点之间的分歧），解决共识问题可以解决几个相关的、更高级的问题，例如原子广播和原子提交。

### Two impossibility results

第一个不可能结果称为FLP不可能结果，主要和设计分布式算法的人特别相关。第二个——CAP定理——是一个与实践者更相关的相关结果——那些需要在不同系统设计之间进行选择，但不直接关注算法设计的人。

## The FLP impossibility result

虽然在学术界被认为是更重要的，但我只会简单地总结一下FLP不可能的结果。FLP不可能的结果（以作者Fischer、Lynch和Patterson的名字命名）检验了异步系统模型下的共识问题（从技术上讲，是共识问题的一种非常弱的形式）。假设节点只能通过崩溃而失败；网络可靠；异步系统模型的典型时间假设成立：例如，消息延迟没有限制。

在这些假设下，FLP结果表明，“在容易发生故障的异步系统中，不存在共识问题的（确定性）算法，即使消息永远不会丢失，至多一个进程可能会失败，并且只能通过崩溃（停止执行）来失败。”

这个结果意味着，**在一个极小的系统模型下，没有办法以一种永远不能延迟的方式来解决共识问题**。论证是，如果存在这样的算法，那么可以设计一种算法的执行，在这种算法中，通过延迟消息传递（异步系统模型中允许），它将在任意时间内保持不确定（“二价”）。因此，这种算法不可能存在。

这种不可能的结果很重要，因为它强调了假设异步系统模型会导致一种权衡：解决共识问题的算法必须在消息传递边界的保证不成立时放弃安全性或活跃性。

这种见解与设计算法的人尤其相关，因为它对异步系统模型中我们知道可以解决的问题施加了一个硬约束。CAP定理是一个与实践者更相关的相关定理：它做出了稍微不同的假设（网络故障而不是节点故障），并对实践者在系统设计之间进行选择具有更明确的含义。

## The CAP theorem

CAP定理最初是由计算机科学家Eric Brewer提出的猜想。 在系统设计的保证中考虑权衡是一种流行且相当有用的方法。该定理说明了这三个属性：

* 一致性：所有节点同时看到相同的数据；
* 可用性：节点故障不会阻止幸存者继续运行；
* 分区容差：尽管由于网络和/或节点故障导致消息丢失，系统仍继续运行。

只有两个可以同时满足。我们甚至可以将它绘制成一个漂亮的图表，从三个中选择两个属性为我们提供了三种类型的系统，它们对应于不同的交叉点：

![](http://book.mixu.net/distsys/images/CAP.png)

注意，该定理表明中间件（具有所有三个属性）是不可实现的。然后我们得到三种不同的系统类型：

* CA（一致性+可用性）。示例包括完全严格的仲裁协议，例如两阶段提交；
* CP（一致性+分区容差）。示例包括多数分区协议，例如Paxos；
* AP（可用性+分区容差）。示例包括使用冲突解决的协议，例如Dynamo。

CA和CP系统设计都提供了相同的一致性模型：强一致性。唯一的区别是，CA系统不能容忍任何节点故障；在非拜占庭故障模型中，给定2f+1节点，CP系统可以容忍最多f个的故障（换句话说，只要大多数f+1保持不变，它可以容忍少数f个节点的故障）。原因很简单：

* CA系统不能区分节点故障和网络故障，因此为了避免引起分歧（多个副本），CA系统必须在所有地方都停止写操作。它无法判断远程节点是否关闭，或者只是网络连接是否关闭：所以唯一安全的事情就是停止接受写操作；
* CP系统通过强制分区两侧的非对称行为来防止分歧（例如保持单个拷贝的一致性）。它只保留大多数分区，并要求少数分区变为不可用（例如停止接受写入），从而保持一定程度的可用性（多数分区），并且仍然确保单个副本的一致性。

当我讨论paxos时，我将在关于复制的那一章中更详细地讨论这个问题。重要的是，CP系统将网络分区合并到其故障模型中，并使用诸如paxos、raft或viewstamped复制之类的算法区分大多数分区和少数分区。CA系统不支持分区，而且在历史上更为常见：它们通常使用两阶段提交算法，并且在传统的分布式关系数据库中很常见。

假设发生了分区，该定理简化为可用性和一致性之间的二元选择。

![](http://book.mixu.net/distsys/images/CAP_choice.png)

我认为应该从CAP定理中得出四个结论：

首先，**在早期的分布式关系数据库系统中使用的许多系统设计没有考虑分区容差**（例如，它们是CA设计）。分区容差是现代系统的一个重要属性，因为如果系统是地理分布的（和许多大型系统一样），那么网络分区的可能性就大得多。

其次，**在网络分区期间，强一致性和高可用性之间存在紧绷关系**。 CAP定理说明了强保证和分布式计算之间的权衡。

在某种意义上，承诺由不可预知网络连接的独立节点组成的分布式系统“以与非分布式系统不可区分的方式运行”是非常疯狂的。

强一致性保证要求我们在分区期间放弃可用性。这是因为在继续接受分区两侧的写操作的同时，不能防止两个无法相互通信的副本之间的分歧。

我们如何解决这个问题？通过加强假设（假设没有分区）或削弱担保。一致性可以与可用性（以及离线可访问性和低延迟的相关功能）进行权衡。如果“一致性”被定义为小于“所有节点同时看到相同的数据”，那么我们可以同时拥有可用性和一些（较弱的）一致性保证。

第三，**在正常操作中强一致性和性能之间存在紧绷关系**。

强一致性/单拷贝一致性要求节点在每个操作上进行通信并达成一致。这导致正常操作期间的高延迟。

如果可以使用传统一致性模型以外的一致性模型（允许副本延迟或分散的一致性模型），那么可以减少正常操作期间的延迟，并在存在分区的情况下保持可用性。

当涉及的消息更少且节点更少时，操作可以更快地完成。但实现这一目标的唯一方法是放宽保证：让一些节点的联系频率降低，这意味着节点可以包含旧数据。

这也使异常发生成为可能——不再保证获得最新价值。根据所提供的保证类型，可能会读取比预期更早的值，甚至会丢失一些更新。

第四，**如果我们不想在网络分区期间放弃可用性，那么我们需要探索除了强一致性之外的一致性模型是否可用于我们的目的**。

例如，即使将用户数据地理复制到多个数据中心，并且这两个数据中心之间的链接暂时不正常，在许多情况下，我们仍然希望允许用户使用网站/服务。这意味着稍后要协调两组不同的数据，这既是一个技术挑战，也是一个业务风险。但通常技术挑战和业务风险都是可管理的，因此最好提供高可用性。

一致性和可用性并不是真正的二元选择，除非您将自己限制为强一致性。但是，强一致性只是一个一致性模型：在这个模型中，您必须放弃可用性，以防止多个数据副本处于活动状态。

简而言之：**“一致性”不是一个单一的，明确的属性**：

> [ACID](http://en.wikipedia.org/wiki/ACID) consistency != [CAP](http://en.wikipedia.org/wiki/CAP_theorem) consistency != [Oatmeal](http://en.wikipedia.org/wiki/Oatmeal) consistency

相反，一致性模型是数据存储向使用它的程序提供的保证——任何保证。

> Consistency model: a contract between programmer and system, wherein the system guarantees that if the programmer follows some specific rules, the results of operations on the data store will be predictable

CAP中的“C”是“强一致性”，但“consistency”不是“strong consistency”的同义词。我们来看看一些替代的一致性模型。

## Strong consistency vs. other consistency models

一致性模型可以分为两类：强一致性模型和弱一致性模型：

* 强一致性模型（能够维护单个副本）
  * Linearizable consistency
  * Sequential consistency
* 弱一致性模型（不强）
  * Client-centric consistency models
  * Causal consistency: strongest model available
  * Eventual consistency models

强一致性模型保证更新的明显顺序和可见性等同于非复制系统。另一方面，弱一致性模型并不能做出这样的保证。

请注意，这绝不是一个详尽的清单。同样，一致性模型只是程序员和系统之间的一个契约，所以它们几乎可以是任何东西。

### Strong consistency models

强一致性模型可以进一步划分为两个相似但略有不同的一致性模型：

* Linearizable consistency：在Linearizable consistency下，所有操作似乎都按照与全局实时操作顺序一致的顺序原子执行。 
* Sequential consistency：在Sequential consistency下，所有操作似乎都以某种顺序原子执行，这与在各个节点上看到的顺序一致，并且在所有节点上都是相等的。

关键的区别在于，前者要求操作生效的顺序等于操作的实际实时顺序。后者允许操作重新排序，只要在每个节点上观察到的顺序保持一致。如果他们能够观察到进入系统的所有输入和时间，可以区分这两者——如果光从客户机与节点交互的角度来看，这两者是等效的。

差异似乎并不重要，但值得注意的是equential consistency并不构成。

强一致性模型允许您作为程序员用分布式节点集群替换单个服务器，而不会遇到任何问题。

所有其他一致性模型都有异常（与保证强一致性的系统相比），因为它们的行为方式可以与非复制系统区分开来。但这些异常现象通常是可以接受的，要么是因为我们不关心偶尔出现的问题，要么是因为我们编写了代码来处理在以某种方式发生不一致之后出现的问题。

请注意，弱一致性模型确实没有任何通用类型，因为“不是一个强一致性模型”（例如“在某种程度上可以与非复制系统区分开”）几乎可以是任何东西。

### Client-centric consistency models

以客户机为中心的一致性模型是以某种方式涉及客户机或会话概念的一致性模型。例如，以客户机为中心的一致性模型可以保证客户机永远不会看到数据项的旧版本。这通常是通过在客户端库中构建额外的缓存来实现的，这样，如果客户端移动到包含旧数据的副本节点，那么客户端库将返回其缓存值，而不是从副本返回旧值。

客户端仍然可以看到旧版本的数据，如果它们所在的副本节点不包含最新版本，但是它们永远不会看到旧版本的值重新出现的异常（例如，因为它们连接到不同的副本）。 请注意，有许多种以客户为中心的一致性模型。

### Eventual consistency

最终一致性模型表明，如果停止更改值，则在一些未定义的时间后，所有副本将同意相同的值。 暗示在此之前副本之间的结果在某种不确定的方式上是不一致的。 

说一些事情是最终一致的就像说“人们最终死了”。这是一个非常弱的约束，我们可能想要至少对两件事进行更具体的表征：

首先，“最终”有多长？具有严格的下限，或者至少对于系统收敛到相同值通常需要多长时间。

其次，副本如何就值达成一致？ 始终返回“42”的系统最终是一致的：所有副本都同意相同的值。 它只是没有收敛到有用的值，因为它只是保持返回相同的固定值。 相反，我们希望更好地了解该方法。 例如，一种决定方法是使具有最大时间戳的值始终获胜。

因此，当供应商说“最终一致”时，他们的意思是一些更精确的术语，例如“最终last-writer-wins，同时read-the-latest-observed-value”的一致性。 “how”是很重要的，因为糟糕的方法会导致写入丢失——例如，如果一个节点上的时钟设置不正确并且使用了时间戳。

我将在弱一致性模型的复制方法一章中更详细地研究这两个问题。

## Further reading

- [Brewer's Conjecture and the Feasibility of Consistent, Available, Partition-Tolerant Web Services](http://lpd.epfl.ch/sgilbert/pubs/BrewersConjecture-SigAct.pdf) - Gilbert & Lynch, 2002
- [Impossibility of distributed consensus with one faulty process](http://scholar.google.com/scholar?q=Impossibility+of+distributed+consensus+with+one+faulty+process) - Fischer, Lynch and Patterson, 1985
- [Perspectives on the CAP Theorem](http://scholar.google.com/scholar?q=Perspectives+on+the+CAP+Theorem) - Gilbert & Lynch, 2012
- [CAP Twelve Years Later: How the "Rules" Have Changed](http://www.infoq.com/articles/cap-twelve-years-later-how-the-rules-have-changed) - Brewer, 2012
- [Uniform consensus is harder than consensus](http://scholar.google.com/scholar?q=Uniform+consensus+is+harder+than+consensus) - Charron-Bost & Schiper, 2000
- [Replicated Data Consistency Explained Through Baseball](http://pages.cs.wisc.edu/~remzi/Classes/739/Papers/Bart/ConsistencyAndBaseballReport.pdf) - Terry, 2011
- [Life Beyond Distributed Transactions: an Apostate's Opinion](http://scholar.google.com/scholar?q=Life+Beyond+Distributed+Transactions%3A+an+Apostate%27s+Opinion) - Helland, 2007
- [If you have too much data, then 'good enough' is good enough](http://dl.acm.org/citation.cfm?id=1953140) - Helland, 2011
- [Building on Quicksand](http://scholar.google.com/scholar?q=Building+on+Quicksand) - Helland & Campbell, 2009

# 3. Time and order

What is order and why is it important?

What do you mean "what is order"?

我的意思是，为什么我们一开始就如此关注于顺序？为什么我们关心A是否发生在B之前？为什么我们不关心其他的属性，比如“颜色”？

正如你可能记得的，我将分布式编程描述为：

> The art of solving the same problem that you can solve on a single computer using multiple computers.

事实上，这就是为什么分布式编程如此关注顺序。任何一次只能做一件事的系统都将创建一个操作的总顺序。就像人们通过一扇门一样，每一个行动都将有一个明确的前驱和后继者。这基本上就是我们努力维护的编程模型。

传统的模式是：在一个CPU上运行一个程序、一个进程、一个内存空间。操作系统抽象出这样一个事实：可能有多个CPU和多个程序，并且计算机上的内存实际上在许多程序之间共享。我不是说线程编程和面向事件的编程不存在；只是它们是“one/one/one”模型之上的特殊抽象。程序是按顺序执行的：从顶部开始，然后向下进入底部。

顺序作为一种属性受到了如此多的关注，因为定义“正确性”的最简单方法是说“它像在一台机器上那样工作”。这通常意味着a）我们运行相同的操作，b）我们以相同的顺序运行它们——即使有多台机器。

保持顺序的分布式系统（定义为单个系统）的好处在于它们是通用的。不需要关心这些操作是什么，因为它们将像在一台机器上一样执行。这很好，因为这样无论操作是什么，都可以使用相同的系统。

实际上，一个分布式程序在多个节点上运行；有多个CPU和多个操作流进入。你仍然可以分配一个总订单，但它需要精确的时钟或某种形式的通信。可以使用一个完全精确的时钟为每个操作设置时间戳，然后使用它计算出总顺序。或者，可能有某种通信系统，它可以按总顺序分配序列号。

## Total and partial order

分布式系统的自然状态是偏序的。无论是网络节点还是独立节点，都不能保证相对顺序；但是在每个节点上，都可以观察到本地顺序。

全序是一个二进制关系，它定义了某个集合中每个元素的顺序。当其中一个元素大于另一个元素时，两个截然不同的元素是可比较的。在偏序的集合中，一些元素对是不可比较的，因此偏序并不指定每个项的确切顺序。

全序和偏序都具有传递性和反对称性。对于X中的所有A、B和C，以下语句同时具有全序和偏序：

```
If a ≤ b and b ≤ a then a = b (反对称性);
If a ≤ b and b ≤ c then a ≤ c (传递性);
```

但是，全序符合整体性：

```
a ≤ b or b ≤ a (整体性) for all a, b in X
```

偏序只符合自反性：

```
a ≤ a (自反性) for all a in X
```

请注意，整体性意味着自反性；因此部分顺序是整体顺序的较弱变体。对于偏序的某些元素，totality属性不成立——换句话说，有些元素是不可比的。

Git Branch就是偏序的一个例子：

```
[ branch A (1,2,0)]  [ master (3,0,0) ]  [ branch B (1,0,2) ]
[ branch A (1,1,0)]  [ master (2,0,0) ]  [ branch B (1,0,1) ]
                  \  [ master (1,0,0) ]  /
```

分支A和分支B源于一个共同的祖先，但它们之间没有明确的顺序：它们代表不同的历史，如果没有额外的工作（合并），就不能简化为单一的线性历史。

在一个由一个节点组成的系统中，全序是必要的：指令被执行，消息在一个程序中以一个特定的、可观察的顺序被处理。于是我们开始依赖于这个全序——它使程序的执行具有可预测性。这种顺序可以在分布式系统上维护，但代价是：通信成本高昂，时间同步困难且脆弱。

## What is time?

时间是顺序之源——它允许我们定义操作的顺序——巧合的是，它也有一种人们可以理解的解释（一秒钟、一分钟、一天等等）。

在某种意义上，时间就像其他整数计数器一样。它恰好足够重要，大多数计算机都有一个专用的时间传感器，也就是时钟。这是非常重要的，我们已经找到了如何用一些不完善的物理系统（从蜡烛到铯原子）合成相同计数器的近似值。通过“综合”，我的意思是我们可以通过一些物理性质，在物理上遥远的地方近似整数计数器的值，而不需要直接通信。

时间戳实际上是表示从宇宙开始到当前时刻的世界状态的一个速记值——如果某个事件发生在某个特定的时间戳上，那么它可能受到之前发生的所有事情的影响。这个想法可以概括为一个因果时钟，它明确地跟踪原因（依赖性），而不是简单地假设时间戳之前的所有内容都是相关的。当然，通常的假设是，我们只应该担心特定系统的状态，而不是整个世界。

假设时间在任何地方都以相同的速度进行——这是一个很大的假设，我稍后将回到这个假设——时间和时间戳在程序中使用时有几个有用的解释。这三种解释是：

* 顺序
* 持续时间
* 解释

顺序。当我说时间是顺序之源，我的意思是：

* 我们可以将时间戳附加到无序事件从而使它们有序；
* 我们可以使用时间戳来强制执行特定的操作顺序或消息的传递（例如，如果操作无序到达则延迟操作）；
* 我们可以使用时间戳的值来确定某事物是否在其他事物之前按时间顺序发生。

解释。时间是一个普遍可比的价值。时间戳的绝对值可以解释为日期，这对人们很有用。

持续时间。持续时间与现实世界有一定的关系。算法通常不关心时钟的绝对值或它作为日期的解释，但它们可能会使用持续时间来作出一些判断。特别是，等待所花费的时间量可以提供有关线索来判断系统是分区的还是仅仅经历高延迟。

就其性质而言，分布式系统的组件的行为不可预测。它们不保证任何特定的顺序、增长率或延迟。每个节点都有一些本地命令——因为执行是（大致）连续的——但是这些本地命令彼此独立。

当事情可以以任何顺序发生时，人类很难对事情进行推理——只是有太多的排列需要考虑。

## Does time progress at the same rate everywhere?

我们都有一个直观的时间概念，基于我们个人的经验。不幸的是，直观的时间概念使我们更容易描绘出总序而不是偏序。更容易想象事情发生的顺序，一个接一个，而不是同时发生。对一个消息顺序进行推理要比对以不同顺序和不同延迟到达的消息进行推理容易得多。

然而，在实施分布式系统时，我们希望避免对时间和顺序做出强有力的假设，因为假设越强，系统就越容易受到“时间传感器”或车载时钟的问题的影响。此外，执行命令也会带来成本。我们越能容忍时间上的不确定性，就越能利用分布式计算。

“每个地方的时间增长的速度都相同吗？”这个问题有3个答案。这些是：

- "Global clock": yes
- "Local clock": no, but
- "No clock": no!

这些大致符合我在第二章中提到的三个计时假设：同步系统模型有一个全局时钟，部分同步模型有一个本地时钟，在异步系统模型中，根本不能使用时钟。让我们更详细地看看这些。

### Time with a "global-clock" assumption

全球时钟的假设是有一个完全准确的全球时钟，并且每个人都可以使用该时钟。这是我们考虑时间的方式，因为在人类互动中，时间上的微小差异并不真正重要。

![](http://book.mixu.net/distsys/images/global-clock.png)

全局时钟基本上是全序的（所有节点上每个操作的确切顺序，即使这些节点从未通信过）。

然而，这是一个理想化的世界观：在现实中，时钟同步只能在有限的精度范围内实现。这一点受到商品计算机时钟精度不足、使用NTP等时钟同步协议时的延迟以及时空本质的限制。

假设分布式节点上的时钟是完全同步的，这意味着假设时钟以相同的值开始，并且永不分离。这是一个很好的假设，因为可以自由地使用时间戳来确定一个由时钟漂移而不是延迟约束的全局总顺序，但这是一个非常重要的操作挑战，也是一个潜在的异常源。有许多不同的场景，其中一个简单的故障——例如用户意外地更改了机器上的本地时间，或者过时的机器加入了一个集群，或者同步时钟以稍微不同的速率漂移，等等，这可能导致难以跟踪的异常。

然而，有一些现实世界的系统做出了这个假设。Facebook的[Cassandra](http://en.wikipedia.org/wiki/Apache_Cassandra)就是一个假设时钟是同步的系统的例子。它使用时间戳来解决写入之间的冲突——使用新时间戳的写入操作将获胜。这意味着如果时钟漂移，新数据可能会被旧数据忽略或覆盖；同样，这是一个操作上的挑战（从我所听到的，人们敏锐地意识到的）。另一个有趣的例子是谷歌的Spanner：[本文](http://research.google.com/archive/spanner.html)描述了他们的TrueTimeAPI，它可以同步时间，但也可以估计最坏情况下的时钟漂移。

### Time with a "Local-clock" assumption

第二种可能更合理的假设是，每台机器都有自己的时钟，但没有全球时钟。这意味着不能使用本地时钟来确定远程时间戳是在本地时间戳之前还是之后发生的；换句话说，不能有意义地比较来自两台不同机器的时间戳。

![](http://book.mixu.net/distsys/images/local-clock.png)

本地时钟假设更接近于现实世界。它分配了一个偏序：每个系统上的事件都是有序的，但是不能只用一个时钟来跨系统对事件进行排序。

但是，可以使用时间戳在单台计算机上订购事件。当然，在终端用户控制的机器上，假设太多了：例如，用户可能在使用操作系统的日期控件查找日期时意外地将其日期更改为其他值。

### Time with a "No-clock" assumption

最后，还有逻辑时间的概念。在这里，我们根本不使用时钟，而是以其他方式跟踪因果关系。记住，时间戳只是到世界状态某一瞬间的简写，因此我们可以使用计数器和通信来确定是否发生了什么事情，是之前、之后还是同时发生了什么事情。

通过这种方式，我们可以确定不同机器之间事件的顺序，但不能谈论间隔，也不能使用超时（因为我们假设没有“时间传感器”）。这是偏序的：事件可以使用计数器在单个系统上进行排序，而无需通信，但跨系统排序事件需要消息交换。

在分布式系统中引用最多的论文之一是Lamport关于时间、时钟和事件顺序的论文。矢量时钟，这一概念的概括（我将更详细地介绍），是一种不用时钟跟踪因果关系的方法。Cassandra的堂兄弟Riak（Basho）和Voldemort（LinkedIn）使用矢量时钟，而不是假设节点可以访问一个完全精确的全局时钟。这使得这些系统可以避免前面提到的时钟精度问题。

当不使用时钟时，可以在远程机器上对事件进行排序的最大精度受通信延迟的限制。

## How is time used in a distributed system?

时间有什么好处？

* 时间可以定义整个系统的顺序（无通信）
* 时间可以定义算法的边界条件

事件顺序在分布式系统中很重要，因为分布式系统的许多属性是根据操作/事件的顺序定义的：

* 正确性取决于（协议）正确的事件排序，例如分布式数据库中的可序列化
* 当资源争用发生时，可以将时间顺序用作判断，例如，如果窗口小部件有两个订单，则执行第一个订单并取消第二个订单

全球时钟将允许在两台不同的机器上进行操作，而不需要两台机器直接通信。如果没有全球时钟，我们需要通信以确定顺序。

时间还可以用来定义算法的边界条件——特别是区分“高延迟”和“服务器或网络链路关闭”。这是一个非常重要的用例；在大多数实际系统中，超时用于确定远程计算机是否发生故障，或者它是否只是经历了高网络延迟。做出这一决定的算法称为故障检测器；我将很快讨论它们。

## Vector clocks (time for causal order)

前面，我们讨论了关于分布式系统中时间进度的不同假设。假设我们不能实现精确的时钟同步——或者从我们的系统不应该对时间同步问题敏感的目标开始，我们如何排序？

LAMPORT时钟和矢量时钟是物理时钟的替代品，它们依靠计数器和通信来确定分布式系统中事件的顺序。这些时钟提供了一个在不同节点之间可比较的计数器。

Lamport时钟很简单。每个进程使用以下规则维护计数器：

* 只要一个进程工作，就增加计数器
* 每当进程发送消息时，计数器也会跟着发送
* 收到消息时，将计数器设置为 `max(local_counter, received_counter) + 1`

表达为代码：

```
function LamportClock() {
  this.value = 1;
}

LamportClock.prototype.get = function() {
  return this.value;
}

LamportClock.prototype.increment = function() {
  this.value++;
}

LamportClock.prototype.merge = function(other) {
  this.value = Math.max(this.value, other.value) + 1;
}
```

一个lamport时钟允许在系统间比较计数器，但要注意：lamport时钟定义了一个偏序。如果 `timestamp(a) < timestamp(b)`：

* a可能发生在b之前或
* a可能与b无法比较

这被称为时钟一致性条件：如果一个事件先于另一个事件，那么该事件的逻辑时钟先于其他事件。如果a和b来自同一因果历史，例如，两个时间戳值都是在同一进程上生成的；或者b是对a中发送的消息的响应，那么我们知道a发生在b之前。

直观地说，这是因为lamport时钟只能携带一个时间线/历史的信息；因此，比较从不相互通信的系统的lamport时间戳可能会导致并发事件在不进行通信时看起来是有序的。

想象一下这样一个系统，它在一个初始阶段之后分成两个独立的子系统，这些子系统从不相互通信。

对于每个独立系统中的所有事件，如果a发生在b之前，则 `ts(a) < ts(b)`；但如果从不同独立系统中选取两个事件（例如，与因果关系无关的事件），则不能对它们的相对顺序说任何有意义的话。虽然系统的每个部分都为事件分配了时间戳，但这些时间戳彼此之间没有关系。两个事件似乎是有序的，即使它们是无关的。	

然而，从一台机器的角度来看，这仍然是一个有用的属性，用 `ts(a)` 发送的任何消息都将收到一个用 `ts(b)` 发送的响应，该响应`> ts(a)`。

矢量时钟是lamport时钟的一个扩展，它保持一个数组 `[ t1, t2, ... ]` 共n个逻辑时钟——每个节点一个。每个节点在每个内部事件上将向量中自己的逻辑时钟增加一个，而不是增加一个公共计数器。因此，更新规则是：

* 只要进程工作，就增加向量中节点的逻辑时钟值
* 每当进程发送消息时，逻辑时钟的完整向量也会跟着发送
* 收到消息时：
  * 将向量中的每个元素更新为 `max(local, received)`
  * 增加表示矢量中当前节点的逻辑时钟值

表达为代码：

```
function VectorClock(value) {
  // expressed as a hash keyed by node id: e.g. { node1: 1, node2: 3 }
  this.value = value || {};
}

VectorClock.prototype.get = function() {
  return this.value;
};

VectorClock.prototype.increment = function(nodeId) {
  if(typeof this.value[nodeId] == 'undefined') {
    this.value[nodeId] = 1;
  } else {
    this.value[nodeId]++;
  }
};

VectorClock.prototype.merge = function(other) {
  var result = {}, last,
      a = this.value,
      b = other.value;
  // This filters out duplicate keys in the hash
  (Object.keys(a)
    .concat(b))
    .sort()
    .filter(function(key) {
      var isDuplicate = (key == last);
      last = key;
      return !isDuplicate;
    }).forEach(function(key) {
      result[key] = Math.max(a[key] || 0, b[key] || 0);
    });
  this.value = result;
};
```

下图展示了一个矢量时钟的逻辑：

![](http://book.mixu.net/distsys/images/vector_clock.svg.png)

三个节点（A、B、C）中的每一个都跟踪矢量时钟。当事件发生时，它们用矢量时钟的当前值进行时间戳。通过检查向量时钟，如 `{ A: 2, B: 4, C: 1 }` ，我们可以准确地识别（可能）影响该事件的消息。

矢量时钟的问题主要是它们需要每个节点一个条目，这意味着对于大型系统来说，它们可能会变得非常大。已经应用了各种技术来减小矢量时钟的大小（通过执行定期垃圾收集，或者通过限制大小来降低精度）。

我们已经研究了如何在没有物理时钟的情况下跟踪顺序和因果关系。现在，让我们看看如何使用持续时间进行截止。

## Failure detectors (time for cutoff)

如我前面所述，等待所花费的时间量可以提供一个系统是分区的还是仅仅经历高延迟的线索。在这种情况下，我们不需要假设一个完全精确的全球时钟——只要有一个足够可靠的本地时钟就足够了。

给定一个程序在一个节点上运行，它怎么能告诉远程节点失败了？在缺乏准确信息的情况下，我们可以推断在经过一段合理的时间后，一个没有响应的远程节点已经失败。

但什么是“合理的数额”？这取决于本地和远程节点之间的延迟。与其显式地指定具有特定值的算法（在某些情况下不可避免地是错误的），不如处理适当的抽象。

故障检测器是一种提取准确时间假设的方法。故障检测器使用心跳消息和计时器实现。程交换心跳消息。如果在超时发生之前没有收到消息响应，那么进程会怀疑另一个进程。

基于超时的故障检测器将承担过度激进（声明节点发生故障）或过度保守（检测崩溃需要很长时间）的风险。故障检测器需要多精确才能使其可用？

钱德拉等人（1996）在解决共识的背景下讨论故障检测器——这是一个特别相关的问题，因为它是大多数复制问题的基础，其中复制副本需要在具有延迟和网络分区的环境中一致。

它们使用两个特性来描述故障检测器，即完整性和准确性：

> Strong completeness: Every crashed process is eventually suspected by every correct process.
>
> Weak completeness: Every crashed process is eventually suspected by some correct process.
>
> Strong accuracy: No correct process is suspected ever.
>
> Weak accuracy: Some correct process is never suspected.

完整性比准确性更容易实现；事实上，所有重要的故障检测器都能做到这一点——所需要做的就是不要永远怀疑某个进程。注意，具有弱完整性的故障检测器可以转换为具有强完整性的故障检测器（通过广播有关可疑过程的信息），从而使我们能够集中关注准确性的频谱。

除非能够假设消息延迟有一个固定最大值，否则很难避免错误地怀疑正确的进程。这种假设可以在同步系统模型中进行，因此故障检测器在这种系统中可以非常精确。在不对消息延迟施加硬限制的系统模型下，故障检测最好的情况下最终会是准确的。

钱德拉等人表明即使是非常弱的故障检测器——最终弱故障检测器W（准确性弱+完整性弱）也可以用来解决共识问题。下图说明了系统模型与问题解决能力之间的关系：

![](http://book.mixu.net/distsys/images/chandra_failure_detectors.png)

如您所见，在异步系统中，没有故障检测器，某些问题是无法解决的。这是因为如果没有故障检测器（或对时间界限的强假设，例如同步系统模型），就无法判断远程节点是否崩溃，或者只是经历了高延迟。这种崩溃或高延迟的区别对于任何旨在实现单一副本一致性的系统都很重要：失败的节点可以被忽略，因为它们不能导致分歧，但是分区节点不能被安全地忽略。

如何实现故障检测器？从概念上讲，对于一个简单的故障检测器来说并没有什么，它只在超时结束时检测故障。最有趣的部分是如何判断远程节点是否失败。

理想情况下，我们希望故障检测器能够适应不断变化的网络条件，并避免将超时值硬编码到其中。例如，Cassandra使用应计故障检测器，它是一个故障检测器，输出怀疑级别（介于0和1之间的值），而不是二进制“向上”或“向下”判断。这允许使用故障检测器的应用程序自行决定准确检测和早期检测之间的权衡。

## Time, order and performance

早些时候，我暗示必须支付顺序费用。什么意思？

如果你正在编写分布式系统，那么你可能拥有多台计算机。世界的自然（和现实）观是一个偏序，而不是一个全序。可以将一个部分顺序转换为一个全序，但这需要通信、等待并施加限制，以限制有多少计算机可以在任何特定时间点工作。

所有时钟都只是网络延迟（逻辑时间）或物理限制的近似值。即使在多个节点之间保持一个简单的整数计数器同步也是一个挑战。

虽然时间和顺序经常在一起讨论，但时间本身并不是一个有用的属性。算法并不真正关心时间，而是关心更抽象的属性：

* 事件的因果排序
* 故障检测（例如，消息传递的上限近似值）
* 一致的快照（例如，在某个时间点检查系统状态的能力;此处未讨论）

实施一个全面的命令是可能的，但代价昂贵。它要求你以共同的（最低的）速度前进。通常，确保以某种定义的顺序传递事件的最简单方法是指定一个（瓶颈）节点，通过该节点传递所有操作。

时间/顺序/同步性真的有必要吗？这要看情况而定。在某些用例中，我们希望每个中间操作将系统从一个一致状态移动到另一个一致状态。例如，在许多情况下，我们希望来自数据库的响应表示所有可用信息，并且希望避免处理如果系统返回不一致的结果可能发生的问题。

但在其他情况下，我们可能不需要那么多时间/顺序/同步。例如，如果你运行的是长时间运行的计算，并且直到最后才真正关心系统的工作，那么只要能够保证答案是正确的，就不需要太多的同步。

同步通常作为一种钝性工具应用于所有的操作，当只有一个子集的情况实际上对最终结果很重要时。何时需要顺序来保证正确性？我将在最后一章讨论的CALM定理提供一个答案。

在其他情况下，给出一个估计的最好答案是可以接受的——也就是说，只基于系统中包含的全部信息的一个子集。特别是，在网络分区期间，可能需要在系统的一部分可访问的情况下回答查询。在其他用例中，最终用户甚至无法真正区分可以便宜获得的相对较新的答案和保证正确且计算成本较高的答案之间的区别。例如，某个用户X或X 1的Twitter关注者数量是多少？ 或者电影A，B和C是一些查询的绝对最佳答案？ 做一个更便宜，最正确的“尽力而为”是可以接受的。

在接下来的两章中，我们将研究容错、强一致性系统的复制，这些系统在不断增强对故障的恢复能力的同时提供了强大的保证。当你需要保证正确性并愿意为此付出代价时，这些系统为第一种情况提供了解决方案。然后，我们将讨论具有弱一致性保证的系统，这些保证在分区面前仍然可用，但这只能给你一个“尽力而为”的答案。

## Further reading

### Lamport clocks, vector clocks

- [Time, Clocks and Ordering of Events in a Distributed System](http://research.microsoft.com/users/lamport/pubs/time-clocks.pdf) - Leslie Lamport, 1978

### Failure detection

- [Unreliable failure detectors and reliable distributed systems](http://scholar.google.com/scholar?q=Unreliable+Failure+Detectors+for+Reliable+Distributed+Systems) - Chandra and Toueg
- [Latency- and Bandwidth-Minimizing Optimal Failure Detectors](http://www.cs.cornell.edu/people/egs/sqrt-s/doc/TR2006-2025.pdf) - So & Sirer, 2007
- [The failure detector abstraction](http://scholar.google.com/scholar?q=The+failure+detector+abstraction), Freiling, Guerraoui & Kuznetsov, 2011

### Snapshots

- [Consistent global states of distributed systems: Fundamental concepts and mechanisms](http://scholar.google.com/scholar?q=Consistent+global+states+of+distributed+systems%3A+Fundamental+concepts+and+mechanisms), Ozalp Babaogly and Keith Marzullo, 1993
- [Distributed snapshots: Determining global states of distributed systems](http://scholar.google.com/scholar?q=Distributed+snapshots%3A+Determining+global+states+of+distributed+systems), K. Mani Chandy and Leslie Lamport, 1985

### Causality

- [Detecting Causal Relationships in Distributed Computations: In Search of the Holy Grail](http://www.vs.inf.ethz.ch/publ/papers/holygrail.pdf) - Schwarz & Mattern, 1994
- [Understanding the Limitations of Causally and Totally Ordered Communication](http://scholar.google.com/scholar?q=Understanding+the+limitations+of+causally+and+totally+ordered+communication) - Cheriton & Skeen, 1993