

LSM-Tree的设计思路是，将数据拆分为几百M大小的Segments，并是顺序写入

## 优点

1. LSM-Tree的优点是支持高吞吐的写(可认为是O(1))
   - 写入 memTable 后就返回
2. 针对读取普通的LSM-Tree结构，读取是O(N)的复杂度，在使用索引或者缓存优化后的也可以达到O(logN)的复杂度。

## 缺点

1. 写放大：整个系统需要频繁的compaction，消耗CPU和存储IO
2. 读放大

## write

1. WAL(write ahead log): 先写日志，即使宕机了，也能恢复之前写入的数据
2. “顺序追加”写内存中的memTable，memTable采用`跳表`的数据结构，因此按照key进行排序
3. 当memTable超过一定大小后，会在内存中冻结，变成不可变的memTable(immutable memTable)，同时，为了不阻塞写，会新建一个memTable继续提供服务
4. 把内存中的immutable memTable保存到SSTable层中，此步骤称为minor compaction
   - 这里需要注意
     - 在L0层的SSTable是没有进行合并的，所以这里的key range在多个SSTable中可能会出现重叠
     - 在层数大于0之后的SSTable不存在重叠key(因为>0层的SSTable会发生合并)
5. 当每层磁盘上的SSTable的体积超过一定的大小或个数，就会周期性的合并。此步骤称为major compaction。这个阶段会真正的清除掉被标记删除掉的数据以及多版本数据的合并，避免浪费空间。
   - 由于每个SSTable都是有序的，在合并的时候，可以采用merge sort进行高效合并

[Q1]: 为什么LSM不直接顺序写磁盘，而是需要在内存中缓存一下？

[A1]: ①单条写的性能肯定没有批量写来的快，批量写入来提高写入速度 ②针对新增数据，可以直接从内存中查询到，加速查询速度

## read  
1. 查询内存中的memTable
2. 依次下沉，直到把所有level层的SSTable查询一遍得到最终结果

[Q]: level 0的多个SSTable有重叠的key，它是怎么保证读取时怎么保证读到的不是老的数据？

[A]: read(key)时得到的val也一定是最新的==>是因为是按照顺序追加写入的，后写入的key，一定会落在新的SSTable上(读取时的顺序，也是按照读取最新的SSTable)

- read优化
   1. 优化原因
      - read最坏的情况，将所有level上的SSTable都扫描一遍
   2. 优化方式
      1. 压缩
         - 这种压缩并不是将整个SSTable一起压缩，而是根据locality将数据分组，每个分组分别压缩，这样的好处是当读取数据的时候，不需要解压缩整个文件，而是解压缩部分Group就可以读取
      2. 缓存
         - 因为SSTable在写入磁盘后，除了Compaction之外，是不会变化的 ==> 所以可以将Scan的Block进行缓存，从而提高检索效率
      3. 索引、BloomFilter
         1. 正常情况下，一个read操作需要读取所有SSTable将结果合并后再返回，但是对于某些key而言，有写SSTable上根本不包含该key对应的数据 ==> 所以，可以给每个SSTable添加BloomFilter ==> 用于判断某个SSTable不存在某个key    
      4. 合并时间在晚上开启，白天禁用合并
         - 这个在前面的写入流程中已经介绍过，通过定期合并瘦身，可以有效的清除无效数据，缩短读取路径，提高磁盘利用空间。但Compaction操作是非常消耗CPU和磁盘IO的，尤其是在业务高峰期，如果发生了Major Compaction，则会降低整个系统的吞吐量，这也是一些NoSQL数据库，比如Hbase里面常常会禁用Major Compaction，并在凌晨业务低峰期进行合并的原因。



---

[深入探讨LSM Compaction机制](https://zhuanlan.zhihu.com/p/141186118)
