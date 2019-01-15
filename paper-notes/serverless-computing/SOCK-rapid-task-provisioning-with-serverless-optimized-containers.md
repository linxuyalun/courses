# SOCK: Rapid Task Provisioning with Serverless-Optimized Containers

* [SOCK: Rapid Task Provisioning with Serverless-Optimized Containers](https://www.usenix.org/system/files/conference/atc18/atc18-oakes.pdf)

这篇文章谈了一个serverless平台容器冷启动的问题并给出了他们的解决方案，SOCK，一种专门针对serverless冷启动进行优化的容器。这篇文章看完感觉没什么意思= =

serverless的三种策略（在更高抽象级别进行编程，重用库以及将应用程序分解为自动扩展的lambda）可提高开发人员的速度，但它们也会产生新的基础结构问题。 具体而言，这些技术使得冷启动过程更加昂贵和频繁。

* 像Python和JavaScript这样的语言需要大量的运行时间，使启动速度比启动等效的C程序慢10倍；
* 重用代码引入了库加载和初始化的进一步启动延迟。 Serverless 计算放大了这些成本：如果将整个应用程序分解为N个无服务器的lambda，则启动频率也会被放大；
* Lambda通常通过容器彼此隔离，这需要进一步的沙箱开销。

为了更好地理解干扰有效冷启动的沙箱和应用特性，论文先做了很多实验。

论文先解构了一下Docker，分析Docker的实现在serverless中的性能瓶颈。首先，在serverless环境中，所有处理程序都运行在少数基本映像之一上，因此联合文件系统的灵活堆叠可能比不上bind mounting的性能成本。 创建根位置后，依赖于复制装入命名空间的文件系统树转换的规模很大。当不需要灵活的文件系统树构造时，可以使用更便宜的`chroot`调用来access drop。 其次，network namespace是一个主要的可扩展性瓶颈; 虽然静态端口分配在基于服务器的环境中可能很有用，但serverless平台（如AWS Lambda）会在网络地址转换器后面执行处理程序，从而使network namespace几乎没有价值。 第三，重用cgroup的速度是创建新cgroup的两倍，这表明维护一组初始化的cgroup可能会减少启动延迟并提高整体吞吐量。

论文进一步分析了Python初始化的一些问题，实验结果表明download和install package需要花费几秒钟，而97%的可安装package可以在单个文件中共存，因此在磁盘上构建一个本地存储大型包存储库是可行的。import package 需要超过100毫秒，但是包与包间存在巨大的流行度偏差，36％的import只有20个package。因此，可以将一部分包预先import到解释器内存中。

于是SOCK就登场了，SOCK的实现机制和评估具体就看论文吧。这篇论文真的挺没意思的，也可能是因为今天整个人蔫蔫的，看着才觉着无聊吧。说不定之后会补补这块内容。



[返回目录](../README.md)



