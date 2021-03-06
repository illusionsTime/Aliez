* b树  有序数组+平衡多叉树
* b+树  有序数组链表+平衡多叉树

Mysql索引主要有两种结构：B+Tree索引和Hash索引 

B+索引

### 局部性原理与磁盘预读：

由于存储介质的特性，磁盘本身存取就比主存慢很多，再加上机械运动耗费，磁盘的存取速度往往是主存的几百分分之一，因此为了提高效率，要尽量减少磁盘I/O。为了达到这个目的，磁盘往往不是严格按需读取，而是每次都会预读，即使只需要一个字节，磁盘也会从这个位置开始，顺序向后读取一定长度的数据放入内存。这样做的理论依据是计算机科学中著名的局部性原理： 
当一个数据被用到时，其附近的数据也通常会马上被使用。 
程序运行期间所需要的数据通常比较集中。 
由于磁盘顺序读取的效率很高（不需要寻道时间，只需很少的旋转时间），因此对于具有局部性的程序来说，预读可以提高I/O效率。

#### 为什么说红黑树没能充分利用磁盘预读功能
红黑树这种结构，h明显要深的多。由于逻辑上很近的节点（父子）物理上可能很远，无法利用局部性，所以红黑树的I/O渐进复杂度也为O(h)，效率明显比B-Tree差很多。

#### b+比b更适合作索引

B+树的关键字全部存放在叶子节点中，非叶子节点用来做索引，而叶子节点中有一个指针指向一下个叶子节点。做这个优化的目的是为了提高区间访问的性能。而正是这个特性决定了B+树更适合用来存储外部数据。

B+树只要遍历叶子节点就可以实现整棵树的遍历。而且在数据库中基于范围的查询是非常频繁的，而B树不支持这样的操作（或者说效率太低）。

#### 一棵m阶的B+树和m阶的B树的差异：
（1）有n个子结点的结点中含有n个关键字，
（2）所有叶子结点包含了全部关键字信息及指向含这些关键字记录的指针，且叶子结点本身依关键字的大小自小而大顺序连接。
（3）所有非叶子结点可以看成是索引部分，结点中仅含有其子树中的最大（或最小）关键字。
通常在B+树上有两个头指针，一个指向根结点，另一个指向关键字最小的叶子结点。因此可以对B+树进行两种查找运算：一种是从最小关键字开始进行顺序查找，另一种是从根结点开始进行随机查找。

