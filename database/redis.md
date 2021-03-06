<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [数据结构](#%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84)
    - [SDS 简单动态字符串](#sds-%E7%AE%80%E5%8D%95%E5%8A%A8%E6%80%81%E5%AD%97%E7%AC%A6%E4%B8%B2)
    - [链表](#%E9%93%BE%E8%A1%A8)
    - [字典](#%E5%AD%97%E5%85%B8)
      - [哈希表](#%E5%93%88%E5%B8%8C%E8%A1%A8)
      - [字典结构](#%E5%AD%97%E5%85%B8%E7%BB%93%E6%9E%84)
      - [哈希](#%E5%93%88%E5%B8%8C)
    - [跳跃表](#%E8%B7%B3%E8%B7%83%E8%A1%A8)
      - [实现](#%E5%AE%9E%E7%8E%B0)
    - [整数集合](#%E6%95%B4%E6%95%B0%E9%9B%86%E5%90%88)
      - [升级](#%E5%8D%87%E7%BA%A7)
    - [压缩列表](#%E5%8E%8B%E7%BC%A9%E5%88%97%E8%A1%A8)
      - [节点](#%E8%8A%82%E7%82%B9)
- [对象](#%E5%AF%B9%E8%B1%A1)
    - [字符串对象](#%E5%AD%97%E7%AC%A6%E4%B8%B2%E5%AF%B9%E8%B1%A1)
    - [列表对象](#%E5%88%97%E8%A1%A8%E5%AF%B9%E8%B1%A1)
    - [哈希对象](#%E5%93%88%E5%B8%8C%E5%AF%B9%E8%B1%A1)
    - [集合对象](#%E9%9B%86%E5%90%88%E5%AF%B9%E8%B1%A1)
    - [有序集合对象](#%E6%9C%89%E5%BA%8F%E9%9B%86%E5%90%88%E5%AF%B9%E8%B1%A1)
- [数据库](#%E6%95%B0%E6%8D%AE%E5%BA%93)
    - [Redis过期键的删除策略](#redis%E8%BF%87%E6%9C%9F%E9%94%AE%E7%9A%84%E5%88%A0%E9%99%A4%E7%AD%96%E7%95%A5)
      - [惰性删除](#%E6%83%B0%E6%80%A7%E5%88%A0%E9%99%A4)
      - [定期删除](#%E5%AE%9A%E6%9C%9F%E5%88%A0%E9%99%A4)
    - [持久化](#%E6%8C%81%E4%B9%85%E5%8C%96)
      - [RDB持久化](#rdb%E6%8C%81%E4%B9%85%E5%8C%96)
        - [创建与载入](#%E5%88%9B%E5%BB%BA%E4%B8%8E%E8%BD%BD%E5%85%A5)
        - [自动间隔性保存](#%E8%87%AA%E5%8A%A8%E9%97%B4%E9%9A%94%E6%80%A7%E4%BF%9D%E5%AD%98)
      - [AOF持久化](#aof%E6%8C%81%E4%B9%85%E5%8C%96)
        - [AOF重写](#aof%E9%87%8D%E5%86%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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
zset结构中的zsl跳跃表按分值从小到大保存了所有的集合元素，每个跳跃表节点保存一个集合元素；字典的键保存了元素的成员，而字典的值保存了元素的分值。通过这个字典，程序可以O(1)复杂度超找给定元素的分值。*值得一提的是，虽然zset结构通式使用跳跃表和字典来保存有序集合元素，但是这两种数据结构都会通过指针来共享相同元素的成员和分值，所以同时使用跳跃表和字典来保存集合元素不会产生任何重复成员或者分值，也不会浪费额外的内存*

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
国期间的定期删除策略由activeExpireCycle函数实现，每当Redis的服务器周期性操作serverCron函数执行时，activeExpireCycle函数就会被调用，他在规定时间内，分多次遍历服务器中的各个数据库，从数据库的expires字典中随机检查一部分键的过期时间，并删除其中的过期键。

### 持久化
#### RDB持久化
Redis提供了RDB持久化功能，这个功能可以将Redis在内存中的数据库状态保存到磁盘里
RDB持久化功能所生成的RDB文件是一个经过压缩的二进制文件，通过该文件可以还原生成RDB文件时的数据状态。

##### 创建与载入
* SAVE SAVE命令会阻塞Redis服务器进程，直到RDB文件创建完毕
* BGSAVE BGSAVE命令会派生出一个子进程，由子进程负责创建RDB文件，服务器进程继续处理命令请求
在BGSAVE期间，SAVE命令会被拒绝，避免父进程和紫禁城同事之星两个rdbSave调用，防止产生竞争条件
同时，BGSAVE也会被拒绝，因为同事之星两个BGSAVE命令也会产生竞争条件。

Redis服务器在启动时检测到RDB文件存在，就会自动载入RDB文件的

##### 自动间隔性保存
Redis允许用户配置save选项，让服务器每隔一段时间自动执行一次BGSAVE命令
1. 当服务器启动时，用户可以通过配置文件或者传入参数的方式设置save选项，如果用户没有主动设置，那么服务器会设置默认条件
2. 根据设置的保存条件，设置服务器状态redisServer结构的saveparams属性
3. 服务器状态维持一个dirty计数器，以及一个lastsave属性 dirty计数器记录距离上次成功执行SAVE或者BGSAVE命令之后，服务器对数据库状态进行了多少次修改  lastsave是一个时间戳
4. 服务器周期性操作函数serverCron默认每个100毫秒执行一次，该函数用于对正在运行的服务器进行维护，其中会检查save的选项是否被满足，满足则执行BGSAVE

#### AOF持久化
AOF持久化通过保存Redis服务器所执行的写命令来记录数据库状态的。

实现：
1. 命令追加 服务器执行完一个写命令之后，会以协议格式将被执行的谢明令最佳到服务器状态的aof_buf缓冲区的末尾
2. 文件写入与同步 服务器在每次结束一个事件loop之前，都会调用flushAppendOnlyFile函数，考虑是否将aof_buf缓冲区中的内容写入和保存到AOF文件中

| appendfsync选项的值 | flushAppendOnlyFile的行为                                                  |
| ------------------- | -------------------------------------------------------------------------- |
| always              | 将aof_buf缓冲区中所有内容写入并同步AOF文件                                 |
| everysec            | 每秒同步一次，显式将多个命令写入磁盘，这样即使系统崩溃，也只会丢失1s的数据 |
| no                  | 将aof_buf缓冲区中所有内容写入AOF文件，由os决定何时同步                     |

还原：通过伪客户端执行写命令

##### AOF重写

*AOF持久化是通过保存被执行的写命令来记录数据库状态的，随着服务器运行，文件体积会越来越大，可能会对服务器造成影响*

解决方法，AOF重写
通过子进程进行AOF重写aof_rewrite，（仅记录当前状态必须的写命令，所以文件体积会小很多啊）
子进程带有服务器进程的数据副本，适用子进程而不是线程，可以在避免使用锁的情况下，保证数据的安全性

问题：子进程重写期间，服务器还需要处理命令请求，新的请求可能会修改数据库状态，造成不一致问题
解决：设置了一个AOF重写缓冲区，在服务器执行完一个命令请求后，会将这个写命令追加到AOF缓冲区和AOF重写缓冲区
当重写工作完成后，会向父进程发送一个信号，父进程将执行重写缓冲区的内ring写入到新AOF文件中，然后原子的覆盖现有的AOF文件

#### 优缺点
AOF文件比RDB更新频率高，优先使用AOF还原数据。

AOF比RDB更安全也更大

RDB性能比AOF好

如果两个都配了优先加载AOF

### 线程模型
Redis基于Reactor模式开发了网络事件处理器，这个处理器被称为文件事件处理器（file event handler）。它的组成结构为4部分：多个套接字、IO多路复用程序、文件事件分派器、事件处理器。因为文件事件分派器队列的消费是单线程的，所以Redis才叫单线程模型。

文件事件处理器使用 I/O 多路复用（multiplexing）程序来同时监听多个套接字， 并根据套接字目前执行的任务来为套接字关联不同的事件处理器。
当被监听的套接字准备好执行连接应答（accept）、读取（read）、写入（write）、关闭（close）等操作时， 与操作相对应的文件事件就会产生， 这时文件事件处理器就会调用套接字之前关联好的事件处理器来处理这些事件。
虽然文件事件处理器以单线程方式运行， 但通过使用 I/O 多路复用程序来监听多个套接字， 文件事件处理器既实现了高性能的网络通信模型， 又可以很好地与 redis 服务器中其他同样以单线程方式运行的模块进行对接， 这保持了 Redis 内部单线程设计的简单性。
# 多机数据库
## 主从复制SLAVEOF
复制功能分为*同步*和*命令传播*两个操作
Redis通过PSYNC命令代替SYNC来执行复制使得同步操作
PSYNC具有完整重同步和部分重同步两种模式：
完整重同步：用于处理初次复制
的情况，通过让主服务器创建并发送RDB文件，以及向从服务器发送保存在缓冲区里面的写命令来进行同步

部分重同步：用于处理断线后重复制情况：主服务器可以将主从服务器连接断开的期间执行的写命令发送给从服务器，从服务器接受并执行这些写命令，就可以将数据库更新至主服务器所处的状态

部分重同步的实现：PSYNC得调用方法有两种：
1. 如果从服务器没有复制过主服务器，那么从服务器再开始一次新的复制时向主服务器发送PSYNC ? -1命令，请求进行完整的重同步
2. 如果从服务器已经复制过某个主服务器，从服务器会发送PSYNC <runid> <offset>命令，runid是上一次的主服务器运行id，offset时当前的复制偏移量。
  
主服务器接收到PSYNC命令后，会返回以下三种回复之一：
1. 返回FULLRESYNC<runid><offset>，表示将执行完整的重同步操作
2. 如果返回CONTINUE，表示将执行部分重同步
3. 如果回复-ERR回复，表示当前版本过低，执行SYNC完整重同步

## Sentinel
有一个或者多个Sentinel实例组成的Sentinel系统可以监视任意多个主服务器