# 数据结构

### SDS 简单动态字符串
```c++
struct sdschdr{
    //记录buf数组中已使用的字节的数量
    int len;
    //记录buf数组中未使用的字节的数量
    int free;
    //字节数组，用于保存字符串
    char buf[];
}
```

* 和c字符串的区别：
1. 获取字符串长度只需要O（1）的时间复杂度，程序仅需要访问SDS的len属性即可
2. c语言会产生缓冲区溢出，在Redis中如果需要对SDS进行修改时，API回显检查SDS的空间是否满足修改所需要的要求，如果不满足的话，API会自动将SDS的空间拓展到执行修改所需要的大小，然后才执行实际的修改操作
3. 减少修改字符串时带来的内存重分配次数，SDS通过free属性记录未使用的字节，通过未使用空间，SDS实现了空间预分配和惰性空间释放两种优化策略。是的修改字符串长度N次最多需要执行N次的内存重分配。
4. c字符串只能保存文本数据，SDS可以保存文本或者二进制数据

### 链表
```C++
typedef struct listNode{
    //前置节点
    struct listNode *prev;
    //后置节点
    struct listNode *next;
    //节点值
    void *value;
}listNode;
```
Redis 提供了adlist.h/list来持有链表
```c++
typedef struct list {
    // 表头节点
    listNode *head;
    // 表尾节点
    listNode *tail;
    // 节点值复制函数
    void *(*dup)(void *ptr);
    // 节点值释放函数
    void (*free)(void *ptr);
    // 节点值对比函数
    int (*match)(void *ptr, void *key);
    // 链表所包含的节点数量
    unsigned long len;

} list;
```
特性
* 双端 ：链表节点带有prev和next指针
* 无环：表头节点的prev和表尾节点的next都指向NULL
* 带表头和表尾指针：获取表头节点和表尾节点的复杂度都是O（1）
* 带链表长度计数器
* 多态：链表节点使用void*指针来保存节点值，并且可以通过list结构的dup、free、match三个属性为节点值设置类型特定函数，所以链表可以保存不同类型的值
  
### 字典
字典是哈希键的底层实现之一
Redis 的字典使用哈希表作为底层实现之一
#### 哈希表
```C++
typedef struct dictht {
    
    // 哈希表数组
    dictEntry **table;
    // 哈希表大小
    unsigned long size;
    // 哈希表大小掩码，用于计算索引值
    // 总是等于 size - 1
    unsigned long sizemask;
    // 该哈希表已有节点的数量
    unsigned long used;

} dictht;
```
table属性是一个数组，每个元素都是一个指向dictEntry的指针，每个dictEntry保存一个键值对。
```C++
typedef struct dictEntry {
    // 键
    void *key;
    // 值
    union {
        void *val;
        uint64_t u64;
        int64_t s64;
    } v;
    // 指向下个哈希表节点，形成链表
    struct dictEntry *next;

} dictEntry;
```
next属性是指向下一个哈希表节点的指针，用来解决*键冲突问题*
#### 字典结构
```c++
typedef struct dict {
    // 类型特定函数
    dictType *type;
    // 私有数据
    void *privdata;
    // 哈希表
    dictht ht[2];
    // rehash 索引
    // 当 rehash 不在进行时，值为 -1
    int rehashidx; /* rehashing not in progress if rehashidx == -1 */
    // 目前正在运行的安全迭代器的数量

    int iterators; /* number of iterators currently running */

} dict;
```
ht是一个包含两个项的数组，每个项都是一个dictht哈希表，一般情况下字典只使用ht[0]哈希表，ht[1]哈希表只会在ht[0]哈希表进行rehash时使用

#### 哈希
Redis首先计算出键的哈希值，根据得到的哈希值以及哈希表的sizemask属性计算出索引值

冲突解决：通过*链地址法*解决冲突

rehash：
根据负载因子进行相应的可拓展或者收缩 大小为ht[0].used*2的2^n或者第一个大于等于ht[0].used的2^n
渐进式rehash：
1. 为ht[1]分配空间，让字典同时持有ht[0]和ht[1]两个哈希表
2. 在字典中维持一个索引计数器变量rehashidx，并将他的值设置为0，表示rehash工作正式开始
3. 在rehash期间，每次对字典之星添加、删除等操作时，除了执行制定的操作外，还会顺带将ht[0]哈希表在rehashidx索引上的所有键值对rehash到ht[1]上，当rehash工作完成后，程序将rehashidx值增一
4. 当所有键值对都被rehash到ht[1]时，将rehashodx的值设为-1，表示操作完成

### 跳跃表
跳跃表（skiplist）是一种有序数据结构，通过每个节点维持多个指向其它节点的指针，从达到快速访问的目的。
*跳跃表支持平均O(logN)、最坏O(N)复杂度的查找*，还可以通过顺序性操作来批处理结点。

Redis使用skiplist作为zset的底层实现之一，如果一个有序集合包含的元素数量比较多，又或者有序集合中的元素的成员是比较长的字符串时，Redis就会使用skiplist来作为zset的底层实现。

####  实现
Redis的跳跃表的实现由在skiplistNode和在skiplist两个数据结构定义,其中zskiplistNode结构用于表示跳跃表结点，而在zskiplist结构则用于保存跳跃表结点的相关信息。

```c++
typedef struct zskiplist {

    // 表头节点和表尾节点
    struct zskiplistNode *header, *tail;
    // 表中节点的数量
    unsigned long length;
    // 表中层数最大的节点的层数
    int level;

} zskiplist;
```
跳跃表的节点
```c++
typedef struct zskiplistNode {
    // 成员对象
    robj *obj;
    // 分值
    double score;
    // 后退指针
    struct zskiplistNode *backward;
    // 层
    struct zskiplistLevel {
        // 前进指针
        struct zskiplistNode *forward;
        // 跨度
        unsigned int span;
    } level[];
} zskiplistNode;
```
* 层：跳表节点的level数组每个元素都包含一个指向其他节点的指针，每次创建一个新条约标的接电的时候，程序根据*幂次定律*随机生成一个介于1和32之间的值作为level数组的大小
* 跨度 span 用于记录两个节点之间的距离；跨度石基商用来计算排位（rank），在查找某个节点的过程中，将沿途访问过的所有层的跨度累计起来，得到的结果就是目标节点在跳跃表中的排位
* 分值和成员 跳跃表中的所有节点按照分值score从小到大进行排序，结点的成员对象obj是一个指针，指向一个字符串对象，*在跳跃表中各个节点保存的成员对象必须是唯一的，多个节点保存的分值可以是相同的*分值相同的节点按照成员对象在字典序中的大小来进行排序。

### 整数集合
intset是作为set键的底层实现之一，当一个集合只包含整数值元素或者作者个集合元素不多的情况下，Redis就会使用整数集合作为set的底层实现。  
<font color=red>整数集合保证集合中不会出现重复元素</font>

```c++
typedef struct intset {
    // 编码方式
    uint32_t encoding;
    // 集合包含的元素数量
    uint32_t length;
    // 保存元素的数组
    int8_t contents[];
} intset;
```
contents数组从小到大保存每个元素，类型由encoding决定
#### 升级
每当我么要将一个新元素添加到整数集合里面，并且新元素的类型比整数集合现在所有元素的类型都要长时，整数集合需要先进行升级，然后才能将新元素添加到整数集合里面去。
1. 根据新元素的类型，拓展整数集合底层数组的空间大小，并为新元素分配空间
2. 将底层数组现有的所有元素都转换成与新元素相同的类型，并将类型转换后的元素放置到正确的位上，并保持底层数组的有序性
3. 将新元素添加到底层数组里面
向整数集合添加新元素的时间复杂度为O（n)
好处：提升整数集合的灵活性 尽可能地节约内存
*整数集合不支持降级*

### 压缩列表
压缩列表是列表键和哈希键的底层实现之一。压缩列表是Redis为了节约内存而开发的，是由一系列特殊编码的连续内存快组成的顺序性数据结构。
```
-------------------------------------------
|zlbytes|zltail|zllen|entry1|entry2|...|zd|
-------------------------------------------
```
* zlbytes 记录整个ziplist占用的内存字节数
* zltail 记录ziplist表尾节点距离压缩列表的起始地址有多少字节；通过这个偏移量，程序无需遍历整个压缩列表就可以确定表尾节点的地址
* zllen 记录了ziplist包含的节点数量；当这个属性的值小于UINT16_MAX(65535)时，表示正确的数量；当这个值等于时，节点的真实数量需要遍历才可以得出
* zlend 用于标记ziplist的末端
#### 节点
`|previous_entry_lengrh|encoding|content|`
previous_entry_lengrh 记录压缩列表中前一个节点的长度
encoding 记录了节点的content属性所保存的数据类型及其长度
content 保存结点的值 可以为一个字节数组或者整数
# 对象
Redis有五种类型的对象，分别是*字符串对象、列表对象、哈希对象、集合对象、有序集合对象*
```c
typedef struct redisObject{
    //类型 
    unsigned type:4;
    //编码
    unsigned encoding :4;
    //指针
    void *ptr;
}
```
### 字符串对象
REDIS_STRING对应三种不同的编码方式，分别对应

| 编码方式              | 解释                                |
| --------------------- | ----------------------------------- |
| REDIS_ENCODING_INT    | 使用整数值实现的字符串对象          |
| REDIS_ENCODING_RAW    | 使用SDS实现的字符串对象             |
| REDIS_ENCODING_EMBSTR | 使用embstr编码的SDS实现的字符串对象 |

当字符串值长度大于32字节，那么字符串对象将使用一个SDS来保存字符串值
如果字符串长度值小于32字节，那么字符串对象将使用embstr编码方式来保存这个字符串值
* embstr编码是一种专门用于保存短字符串的一种优化编码方式，这种编码方式都适用redisObject结构和sdshdr结构表示字符串对象，raw调用两次内存分配函数来分别创界，*embstr则通过调用一次内存分配函数来分配一块连续的空间。*  
好处：1.将内存分配次数降低为一次，也只需要调用一次内存释放函数  
2.所有数据都保存在一块连续的内存里面，能更好地利用缓存带来的优势


### 列表对象
| 编码方式                  | 解释                       |
| ------------------------- | -------------------------- |
| REDIS_ENCODING_ZIPLIST    | 使用压缩列表实现的列表对象 |
| REDIS_ENCODING_LINKEDLIST | 使用双端链表实现的列表对象 |
当一个列表键只包含少量(元素数量少于512个)的列表项，并且每个列表想要么是小整数值，要么是长度比较短的字符串（长度小于64字节），那么Redis就会使用压缩列表来做列表键的底层实现。

当一个条件不满足时，对象的编码转换操作就会被执行，原本保存在ziplist中的所有元素会被转移并保存在双端列表里。

### 哈希对象
| 编码方式               | 解释                       |
| ---------------------- | -------------------------- |
| REDIS_ENCODING_ZIPLIST | 使用压缩列表实现的哈希对象 |
| REDIS_ENCODING_HT      | 使用字典实现的哈希对象     |

ziplist编码的哈希对象使用ziplist作为底层实现，每当有新的键值对要加入到哈希对象时，程序会先将保存了键的接地那推入到ziplist结尾，再将保存了值得节点推入到表尾。
1. 保存同一键值对的两个节点始终紧挨，键在前值在后
2. 先添加的键值对会被放在压缩列表的表头方向
   
hashtable编码的哈希对象使用字典作为底层实现，哈希对象的每一个键值都是用一个字典键值对来保存

编码转换：当哈希对象同时满足以下两个条件时，哈希对象使用ziplist编码：
* 哈希对象所保存的键值对的字符串长度都小于64字节
* 哈希对象保存的键值对数量小于512个
不能同时满足则使用hashtable编码

### 集合对象
| 编码方式              | 解释                       |
| --------------------- | -------------------------- |
| REDIS_ENCODING_INTSET | 使用整数集合实现的集合对象 |
| REDIS_ENCODING_HT     | 使用字典实现的集合对象     |

* set和list的区别
  1. 底层实现： list可以有压缩列表，双端链表实现，set由字典和整数集合实现

### 有序集合对象
| 编码方式                | 解释                               |
| ----------------------- | ---------------------------------- |
| REDIS_ENCODING_ZIPLIST  | 使用压缩列表实现的有序集合对象     |
| REDIS_ENCODING_SKIPLIST | 使用跳跃表和字典实现的有序集合对象 |

  REDIS_ZSET类型对应两种编码REDIS_ENCODING_SKIPLIST以及REDIS_ENCODING_ZIPLIST，REDIS_ENCODING_ZIPLIST对应的底层数据结构是压缩列表，REDIS_ENCODING_SKIPLIST使用的数据结构是跳跃表和字典。

  当有序集合对象可以同时满足两个条件时，对象使用ziplist编码
   1. 有序集合保存的元素数量小于128个
   2. 有序集合保存的所有元素成员的长度都小于64字节
  不能满足以上两个条件的有序集合对象将使用skiplist编码

ziplist编码的压缩列表对象使用压缩列表作为底层实现，每个集合元素使用两个紧挨在一起的压缩列表节点来保存，第一个节点保存元素的成员，第二个节点保存元素的分值。压缩列表内的集合元素按分值从小到大进行排序，分值较小的元素被放置在表头方向。

如果是skiplist编码的有序集合对象使用zset作为底层实现，一个zset结构是包含一个字段和一个跳跃表
```c
typedef struct zset{
	zskiplist *zsl;
	dict *dict;
}zset
```
zset结构中的zsl跳跃表按分值从小到大保存了所有的几何元素，每个跳跃表节点保存一个集合元素；字典的键保存了元素的成员，而字典的值保存了元素的分值。通过这个字典，程序可以O(1)复杂度超找给定元素的分值。*值得一提的是，虽然zset结构通式使用跳跃表和字典来保存有序集合元素，但是这两种数据结构都会通过指针来共享相同元素的成员和分值，所以同时使用跳跃表和字典来保存集合元素不会产生任何重复成员或者分值，也不会浪费额外的内存*

* 为什么需要同时使用跳跃表和字典来实现？
因为单独使用跳跃表或者字典在性能上比起使用跳跃表和字典都会有所降低，如果只是用字典，虽然能以O（1）复杂度查找成员的分值这一特性会被保留，但是因为字典以无序的方式来保存集合元素，所以在执行范围型操作时，需要现对字典中的元素进行排序，完成这种排序需要至少O(NlogN)时间复杂度，以及额外的O(N)内存空间。  
而如果单独使用跳跃表，那么跳跃表执行范围型操作的优点会被保留，但是根据成员查找分值这一操作的复杂度将上升为O(logN)

# 数据库
### Redis过期键的删除策略
Redis使用惰性删除和定期删除两种策略，达到一个在CPU使用时间和内存空间的平衡
#### 惰性删除
所有读写数据库的Redis命令在执行之前都会调用expireIfNeeded函数对输入键做一个检查：
* 如果输入键已经过期，那么expireIfNeeded函数将输入键从数据库中删除
* 如果输入键未过期，那么expireIfNeeded函数不做动作
#### 定期删除
国期间的定期删除策略由activeExpireCycle函数实现，每当Redis的服务器周期性操作serverCron函数执行时，activeExpireCycle函数就会被调用，他在规定时间内，分多次遍历服务器中的各个数据库，从数据库的expires字典中随机检查一部分间的过期时间，并删除其中的过期键。

### 持久化
