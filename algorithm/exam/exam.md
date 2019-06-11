# Examination

# 2014

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

动态规划方法。时间复杂度 $O(n^2)$，空间复杂度 $O(n^2)$。令$L[i][j]$为 $x[i]$ 到 $x[j]$ 的最大回文字符串

因此可得到动态规划的递归方程：
$$
L[i][j]=\left\{
\begin{array}{lcl}
i && if\;i=j \\
2&& if\;j=i+1\;and\;x[i]=x[j] \\
max(L[i+1][j],L[i][j-1]) && if\;x[i]=x[j] \\
max(L[i+1][j],L[i][j-1], 2+L[i+1][j+1]) && o.w.
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

题目再描述，问题就是对于一个无向图$G=(V,E)$，有$|V|=n$ 和 $|E|=m$ ，对于每个节点 $x_i\in V$，它都有一个各自的权重值 $c(i)$。问题问的就是是否能找到一个子集 $C \subseteq V$ ，使得对于原图中的每条边 $e \in E$ 都有至少一个点属于集合 $C$，并且权重值 $\sum_{i\in C}c_i$ 为最小。

# 2015

## 02

### Question

Give an $O(n\cdot t)$ algorithm for the following task.

**Input**: A list of $n$ positive integers $a_1,a_2,…a_n$ and a positive integer $t$.

**Question**: Does some subset of the $a_i$'s add up to $t$ ? (You can use each $a_i$ at most once)

### Solution

定义：
$$
f(i,j)=\left\{
\begin{array}{lcl}
true && if\;some\;subset\;of\;a_1,a_2,...,a_i\;add\;up\;to\;j\\ 
false && otherwise\\
\end{array}\right.
$$
故得到动态规划方程组：
$$
f(i,j)=\left\{
\begin{array}{lcl}
true && if\;i=0\;and\;j=0\\ 
true && if\;f(i-1,j)=true\\ 
true && if\;j\geq a_i\;and\;f(i-1,j-a_i)=true\\ 
false && otherwise\\
\end{array}\right.
$$
得到算法：

![](img/2.png)

很明显，该算法时间复杂度为 $O(nt)$

## 03

### Question

In a **facility location problem**, there is a set of facilities and a set of cities, and our job is to choose a subset of facilities to open, and to connect every city to some one of the open facilities. There
is a nonnegative cost $f_j$ for opening facility $j$, and a nonnegative connection cost $c_{i,j}$ for connecting city $i$ to facility $j$. Given these as input, we look for a solution that minimizes the total cost. Formulate this
facility location problem as an integer programming problem.

### Solution

该问题的整数方程为：
$$
Min \quad \sum f_jx_j+\sum c_{i,j}x_jy_{i,j}
$$

$$
s.t.\quad\sum x_jy_{i,j} \geq 1 \; for \; any \; i \\
x_i\geq0\;for\;any\;i \\
\quad y_{i,j}\geq0\;for\;any\;i,j
$$

其中，$x_i$表示设施 $i$ 是否开着；$y_{i,j}$ 表示城市 $i$ 到设施 $j$ 是否连通

## 05

### Question

You are given an undirected graph. The problem is to remove a minimum number of edges such that the residual graph contains no triangle. i.e., there is no three vertices $a,b,c$ such that edges $(a,b), (b,c), (c,a)$ are all in the residual graph. Give a factor 3 approximation algorithm that runs in polynomial time.

### Solution

**Summary**

基本思想是不断重复的寻找图G中的三角形，找到一个就把对应的三条边全部移除。因为一个对应的最优算法至少需要移除一条边，所以可以说该算法是一个近似比为3的算法。

**Algorithm description**

首先，需要有一个小算法，对于给定的图 $G=(V,E)$，该算法找到一个三角形。假设我们的图 $G$ 存储在邻接矩阵中，并且这个算法会检查任何三个节点 $(u,v,w)$ 是否构成一个三角形。显然，这个小算法的时间复杂度是 $O(n^3)$。当然，这个小算法存在进一步优化。

接下来，整个算法的处理流程如下，当图中有任何三角形时，执行以下步骤：

* 这个三角形是 $(u,v,w)$
* 删除 $G$ 中的 $(u,v),(v,w),(u,w)$

**Explain why it works**

首先很容易看出来，当 $G$ 中没有三角形时，算法结束。接下来要证明这个近似算法是一个近似比为3的近似算法。定义一个 $G$ 的 *triangle-set* ，该集合中是 $G$ 中的一堆无关联的三角形。所谓无关联指的是，该集合中不会有两个三角形有共享边。所以该算法会生成一个  *triangle-set*，每次算法的一个循环，都会移除一个三角形的三条边。进一步的，如果这个图构成的 *triangle-set* 为S，那么整个算法移除的边数必然是 $3|S|$。

让 $OPT$ 是本题的最优解，$OPT$ 里包含的是所有要移除的三角形的边。因为最优算法每个三角形移除的边数至少有一条边，因此 $|OPT| \geq |S|$。而近似算法移除的总边数是 $3|S|$，即 $3|OPT| \geq 3|S|$，即 $3|OPT| \geq 近似算法移除的边数$，因此该算法是一个近似比为3的近似算法。

## 08

### Question

Suppose in a given network, all edges are undirected (or think of every edge as bi-directional with the same capacity), and the length of the longest simple path from the source s to sink t is at most L. Show that the running time of **Edmond‐Karp algorithm** on such networks can be improved to be $O(L \cdot |E|^2)$

### Solution

# 2016

## 01

### Question

Write the dual to the following linear program.
$$
Max \quad \frac{1}{2}x+2y-4 \\
x+y\leq3 \\
\frac{x}{2}+3y\leq5 \\
x,y \geq 0
$$
Find the optimal solutions to both primal and dual LPs.

### Solution

$$
Min \quad 3x+5y-4 \\
s.t. \quad x + \frac{y}{2} \geq \frac{1}{2} \\
\qquad  x + 3y \geq 2 \\
\quad x,y \geq 0
$$

> 用画图的方式得到最优解，还有其他方法吗= =

原问题：(8/5,7/5) (max = ‐2/5)

对偶问题：(1/5,3/5)(min=‐2/5)

## 02

### Question

Find a polynomial time 4/3‐approximation algorithm for instance of **Metric TSP** where the distances are either 1 or 2.

*Hint:* The 2‐match problem (find a minimum weight subset of edges 􏰀 such that each node is adjacent to exactly 2 edges in 􏰀) can be solved in polynomial time.

### Solution

## 03

### Question

Consider the following algorithm for **Cardinality Vertex Cover**: In each connect component of the input graph execute a **depth first search (DFS)**. Output the nodes that are not leaves in the DFS tree. Show that the output is indeed vertex cover, and that it approximates the minimum vertex cover within a factor of 2.

### Solution

## 04

### Question

Prove that the **Graph‐Isomorphism problem** is a NP problem.

### Solution

## 05

### Question

**Set Multicover problem** is defined as an extension of **Set Cover problem**, such that each element, $e$, needs to be covered a specified integer number, $r_e$, of times. The objective again is to cover all elements up to their coverage requirements at minimum cost. We will assume that the cost of picking a set $S$ $k$ times is $k \cdot cost(S)$. Represent the problem as an integer linear program. Then relax to a LP and work out its dual.

### Solution

## 06

### Question

Given two strings $x=x_1x_2…x_n$ and $y=y_1y_2…y_m$, we wish to find the length of their longest common substring, that is, the largest $k$ for which there are indices $i$ and $j$ with $x_ix_{i+1}…x_{i+k-1}=y_jy_{j+1}…y_{j+k-1}$. Show how to do this in time $O(mn)$.

### Solution

使用动态规划的思想，如下示例：

|       | F    | I    | S    | H    |
| ----- | ---- | ---- | ---- | ---- |
| **F** | 0    | 0    | 0    | 0    |
| **I** | 0    | 1    | 0    | 0    |
| **S** | 0    | 0    | 2    | 0    |
| **H** | 0    | 0    | 0    | 3    |

$ L[i][j]$ 是 $x_1x_2…x_i$ 和 $y_1y_2…y_j$ 的最长公共子串，满足如下条件：
$$
L[i][j]=\left\{
\begin{array}{lcl}
0 && if \;i=0\;or\;j=0 \\
L[i-1][j-1]+1&& x[i]==y[j] \\
o && o.w.
\end{array}\right.
$$
对应算法为：

```latex
For i = 1 to n
	L[i][0] = 0
For i = 1 to m
	L[0][i] = 0
For i = 1 to n
	For j = 1 to m
		Update L[i][j]

Return max L[i][j]
```

显然，该算法复杂度为 $O(mn)$

