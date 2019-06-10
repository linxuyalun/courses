# Examination

## 01

### Question

Write the dual to the following linear program. 
$$
Max \quad 6x-4z+7 \\
3x-y \leq 1	\\
4y-z \leq 2 \\
x, y, z \geq 0
$$
Is the solution $(x, y, z)=(\frac{1}{2},\frac{1}{2},0)$ optimal? Write the dual program of the given linear program and 

find out its optimal solution. 

### Solution

为了让结果更清楚，将原问题改写为：
$$
Max \quad 6x+y-4z+7
$$

$$
3x - y + 0z \leq 1 \\
0x + 4y - z \leq 2 \\
x,y,z \geq 0
$$

从而写出我们的dual program：
$$
Min \quad x + 2y + 7
$$

$$
3x + 0y \geq 6 \\
-x + 4y \geq 1 \\
0x-y \geq -4 \\
x,y \geq 0
$$

**如何确定这个$(x, y, z)=(\frac{1}{2},\frac{1}{2},0)$是最优解？？？**

## 02

### Question

A **Minimum Makespan Scheduling** problem is as follows:

**Input** processing times for 􏰊 jobs, 􏰋􏰌, 􏰋􏰍, . . . , 􏰋􏰎, and an integer 􏰏. 

**Output** an assignment of the job to m􏰏 identical machines so that the completion time is minimized. 

We know that by a greedy approach on the problem, the approximation factor 2. Give a tight example to show the approximation guarantee. 

### Solution

Minimun Makespan Scheduling 算法描述：

**输入**：所需计算时间分别为$t_1$，$t_2$，…，$t_n$的n个任务；m台相同的机器

**输出**：最少的完成时间。

使用了一个贪心算法解决这个问题，思路是：

1. 随机排列这些job获得一个队列；
2. 依次分配任务1到n，每次选择当前队列时间最短的机器，放入该任务，更新该机器的队列时间。

一个 tight example 可以是这样一个队列：有 $m^2$ 个计算时间为1的 job，后面跟着一个计算时间为m的job。在这种情况下，使用贪心算法所花费的时间是2m，而最优解所花费的时间是m+1。

## 03

### Question

**Steiner Forest Problem** is defined as follows,

**Input** an undirected graph $G=(V, E)$, nonnegative costs $c_e \geq 0$ for all edges $e∈E$ and $􏰭k$ pairs of vertices $(s_i,t_i)∈V$􏰇􏰮􏰐􏰒*.* 

**Output** a minimum cost subset of edges $F$􏰰 such that every $(s_i,t_i)$􏰒 pair is connected in the set of selected edges. 

Represent the problem as an integer program. 

### Solution

看不懂题目！

## 04

### Question

Given a reduction from the **Clique Problem** to the following problem, which you also need to show to be a search problem.

**Input** a undirected graph *G* and a positive integer $k$.

**Output** a Clique of size $k$ as well as an Independent Set of size $k$, provided both exist.

### Solution

1. 给定一个团的大小k和独立集的大小k，我们可以在多项式时间内验证是否能在图G里找到对应尺寸的团和独立集，所以这是一个搜索问题；

2. 规约过程如下，对于一个图$G=(V,E)$，构造图$G'=(V\cup V',E)$，其中$|V'|=k$。此时，如果 $G'$ 有一个尺寸$k>1$的团，那么这个团中的所有节点一定在集合$V$中。

   **P的输入转化为Q的输入**：P的输入自然就是一个图G和正整数k，Q的输入是图 $G'$ 和正整数k

   **Q的输出转化为P的输出**：如果$G'$有且只有一个尺寸为k的团和尺寸为k的独立集，那么在G中一定也有一个尺寸为k的团，否则$G'$无法构成尺寸为k的团，因为剩下的节点要么属于G，要么是独立集中的节点。

## 05

### Question

Given an undirected graph $G=(V, E)$  in which each node has $degree\leq d$, find an approximation algorithm for maximal independent set with the factor $\frac{1}{d+1}$.

### Solution

![](img/1.png)

算法分析：首先使用这个算法，S肯定是G的一个独立集。接下来考虑一下 $V/S$ 中的节点，通过这种算法，每次从G中选取一个节点 $v$ 到 $S$，最多会删除d个相邻节点（因为每个节点 $degree\leq d$ ），所以：
$$
|V/S| \leq d|S|
$$
又因为：
$$
|S| + |V/S|=|V|
$$
即：
$$
|V/S|=|V| - |S|
$$
联立，得：
$$
|V|-|S| \leq d|S| \\
\therefore |S| \geq \frac{1}{d+1}|V| \geq \frac{1}{d+1}opt
$$

## 06

### Question

A subsequence is **palindromic** if it is the same whether read left to right or right to left. For instance, the sequence 
$$
A,C,G,T,G,T,C,A,A,A,A,T,C,G
$$
has many palindromic subsequences, including $A,C,G,C,A$ and $A,A,A,A$. Devise an algorithm that takes a sequence 􏰀􏰳$x[1,…,n]$􏰊􏰴, and returns the (length of the) longest palindromic subsequence. Its running time should be $O(n^2)$.

### Solution

动态规划方法。时间复杂度 $O(n^2)$，空间复杂度 $O(n^2)$，如果用 $f[i][j]$ 保存子串从 $i$ 到 $j$ 是否是回文子串，那么在求 $f[i][j]$ 的时候如果 $j-i>=2$ 时，如果 $f[i][j]$ 为回文，那么 $f[i+1][j-1]$，也一定为回文。否则 $f[i][j]$ 不为回文。如下图：

![](https://img-blog.csdn.net/20170425171617065?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQv/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)

因此可得到动态规划的递归方程：
$$
f(i,j)=\left\{
\begin{array}{lcl}
true && i=j \\
S[i]=S[j]&& j=i+1 \\
S[i]=S[j]\quad and \quad f(i+1,j-1) && j>i+1
\end{array}\right.
$$
对应的伪代码为：

```
for i = 1 to n
	update f(i,i) and f(i,i+1)
for i = (n - 2) to 1
	for j = (i + 2) to n
		update f(i,j)
```

从该伪代码易知，算法复杂度为 $O(n^2)$

## 07

### Question

The **Maximum Cut** problem is defined as follows: Given an undirected graph $G=(V,E)$􏰇􏰒,􏰓􏰈 along with a non negative weight $w_{ij} \geq 0$ for each $(i,j) \in E$. The goal is to partition the vertex set into two parts, $U$ and $W=V-U$, so as to maximize the weight of the edges whose two endpoints are in different parts. Give a 2‐approximation randomized algorithm for maximum cut problem.

### Solution

**使用如下算法：**

对于 $V$ 中的每个节点 $v \in V$，它有$\frac{1}{2}$ 的概率被分割至集合$U$

**算法分析：**

令$X_{ij}$为：
$$
X_{ij}=\left\{
\begin{array}{lcl}
0 && 若边不被分割 \\ 
1 && 若边被分割\\
\end{array}\right.
$$
令Z为该分割状态下的权重和，则有：
$$
Z= \sum _{(i,j)\in E}w_{ij}X_{ij}
$$
由于每个节点有$\frac{1}{2}$ 的概率被分割至集合$U$，因此我们这里求的应该是Z的期望，所以：
$$
E(Z)=\sum _{(i,j)\in E}w_{ij}E(X_{ij})=\sum _{(i,j)\in E}w_{ij}Probability[Edge(i,j)\quad in\quad the\quad cut]=\frac{1}{2}\sum _{(i,j)\in E}w_{ij}
$$
令OPT为该图下的最大分割权的最优解。显然，OPT至多为$\sum _{(i,j)\in E}w_{ij}$，故有：
$$
E(Z)=\frac{1}{2}\sum _{(i,j)\in E}w_{ij} \geq \frac{1}{2}OPT
$$
所以，$\frac{OPT}{E(Z)} \leq 2$。

## 08

### Question

The **Weighted Vertex Cover** problem is defined as follows: Given an undirected graph $G=(V,E)$, where $|V|=n$ and $|E|=m$ and a cost function on vertices $c:V \rightarrow Q^+$, find a subset $C \subseteq V$ such that every edge $e \in E$ has at least one endpoint in $C$ and $C$ has a minimum cost. Use the primal‐dual method to give a 2‐approximation algorithm for this problem, and prove the guarantee.

### Solution

