# kafka

## 1.概念

1. broker
   1. kafka服务器，负责将收到的消息存储到磁盘
   2. 消息中间件处理节点，一个kafka节点就是一个broker，多个broker组成一个kafka集群
2. topic
   1. 消息主题，是一个逻辑概念
   2. kafka中的消息以主题为单位进行归类
   3. 生产者负责将消息发送到特定的主题(发送到 Kafka 集群中的每一条消息都要指定一个主题)，而消费者负责订阅主题并进行消费
3. partition
   1. 不考虑多副本的情况，一个分区对应一个日志(Log)。为了防止 Log 过大，Kafka 又引入了日志分段(LogSegment)的概念，将 Log 切分为多个 LogSegment，相当于一个巨型文件被平均分配为多个相对较小的文件，这样也便于消息的维护和清理。
   2. topic物理上的分组，一个分区只属于一个topic
   3. 优点: 分区类似于分片条带化，支持并发读写，提高读写效率
   4. 分区策略: 决定生产者将消息发到那个分区是算法
      1. 轮询
      2. 随机
      3. 按照消息键保存
4. segment: 每个partition又由多个segment file组成
5. offset: 
   1. 消息在被追加到partition日志文件的时候都会分配一个特定的偏移量(offset)
   2. partition中的每个消息都有一个连续的序列号叫做offset，用于partition唯一标识一条消息
   3. 偏移量offset 是消息在partition中的唯一标识，是一个单调递增且不变的值
   4. Kafka 通过它来保证消息在partition内的顺序性，不过 offset 并不跨越partition，也就是说，Kafka 保证的是partition有序而不是主题有序

总结

1. 主从是以partition为单位的(不是以topic为单位的)，一个topic包含多个partition，多个partition可以横跨多个broker，每个partition包含副本，一个个partition形成自己从主从。
2. kafka的消费者采用pull模式从topic获取消息
   - push/pull对比
     1. push
        - 缺点: (1)没考虑消费者的消费能力 (2)推送完消息后设置消费成功，但是消费者挂了，推送的消息会丢失。需要复杂的逻辑来保证一致性
        - 优点: 及时性强
     2. pull
        - 缺点: 及时性差
        - 优点: (1)消费者可以根据自己的消费能力拉取消息 (2)消费成功后，修改offset，消息不会丢失

## 2.kafka的存储层(Kafka Broker 是如何持久化数据的)

1. kafka使用消息日志log保存数据，一个日志就是磁盘上一个只能追加写消息的物理文件。
2. 如果不断的向一个日志写入消息，最终会耗尽所有的磁盘空间，因此kafka必然要定期的删除消息来回收磁盘。==> kafka通过日志段机制(log segment)
   1. 在kafka底层，一个日志又进一步细分成多个日志段，消息被追加写到当前最新的日志段中
   2. 当写满了一个日志段后，kafka会自动切分成多个日志段，消息被追加写到当前最新的日志段中
   3. 当写满一个日志段后，kafka会自动切分出一个新的日志段，并将老的日志段冯存起来
   4. kafka后台任务会定期检查删除老的日志段，从而实现磁盘回收
3. 向Kafka发送数据并不是真要等数据被写入磁盘才会认为成功，而是只要数据被写入到操作系统的`页缓存(Page Cache)`上就可以了，随后操作系统根据LRU算法定期(约5s)将页缓存上的脏数据刷写物理磁盘
   1. 如果在页缓存中的数据在写入到磁盘前机器宕机了，那岂不是数据就丢失了。
   2. 的确，这种情况数据确实就丢失了，但鉴于 Kafka 在软件层面已经提供了多副本的冗余机制，因此这里稍微拉大提交间隔去换取性能还是一个合理的做法
4. kafka只对“已提交offset”的消息做有限度的持久化
   1. kafka的broker接收到生产者发送的消息写入日志文件后，它们会告诉生产者程序这条消息已经成功提交
   2. 此时，这条消息在kafka看来就正式的变成“已提交”消息了

## 3.kafka中的ZooKeeper是做什么的呢？

### 3.1.zookeeper在kafka中的作用

它是一个分布式协调框架，负责协调管理并保存 Kafka 集群的所有元数据信息，比如: 集群都有哪些Broker在运行、创建了哪些Topic、每个Topic都有多少partition以及这些partition的Leader副本都在哪些机器上等信息。

1. /brokers/ids: 临时节点，保存所有broker节点信息，存储broker的物理地址、版本信息、启动时间等，节点名称为brokerID(broker定时发送心跳到zk，如果断开，则删除该brokerID对应的节点)
2. /broker/topics: 临时节点，节点保存broker节点下所有topic信息，每一个topic节点下包含一个固定的partition节点，partition的子节点就是topic的partition，每个partition下保存了一个state节点，保存着当前leaderpartition和ISR的brokerID，state节点由leader创建，若leader宕机该节点会被删除，知道有新的leader选举产生，重新生成state节点
3. /consumer/[组id]/offset/[topic]/[broker_id-partition_id]: partition消息的消费进度offset
4. /consumer/[组id]/owners/[topic]/[broker_id-partition_id]: 维护消费者和partition的注册关系

### 3.2.kafka的控制器Controller

1. 介绍: 在kafka集群中，每个broker都有一个控制器，但是只有一个控制器运行(即：所有的控制器，只有一个控制器提供服务；当该控制器挂掉后，会有其他控制器来提供服务)
2. 控制器的功能: 主题分配、分区重分配、集群broker成员管理
   1. 获取某个broker上的所有分区
   2. 某组broker上的所有副本
   3. 当前存活的副本
   4. 正在进行rebalance的分区列表
   5. 某组分区下的所有副本
   6. 当前存活、正在关闭的的broker列表
   7. 正在进行preferred leader选举的分区
   8. 分配给每个分区的副本列表
   9. topic列表
   10. 每个分区的leader和ISR信息
   7. 移除某个topic的所有信息
3. 控制器故障转移(failover)
   1. 场景: 在kafka集群中，只有一台broker上的控制器提供服务，那么就会存在单点失效问题==>控制器故障转移功能，就是所谓的failover
   2. 故障转移: 指当运行中的控制器突然宕机或意外终止时，Kafka能够快速地感知到，并立即启用备用控制器来替代之前失败的控制器
   3. 过程
      1. 最开始，broker0的控制器提供服务，当broker0宕机后，zookeeper通过watch机制感知到/controller节点被删除了
      2. 之后，所有存活的broker开始竞选新的控制器执行权
      3. 若broker3的控制器赢得了选举，(成功的在zookeeper上重建了/controller节点)
      4. 之后，broker3的控制器从zookeeper中读取集群元数据信息，并初始化到自己的缓存，至此，控制器的failover完成

### 3.3.offset topic 位移主题

1. 为什么引入offset topic？
   - consumer的位移管理之前是由Zookeeper保存的，需要对zookeeper高频的写操作
2. 新版本的consumer的offset管理机制
   - consumer的offset数据，统一写入__consumer_offset__topic__中，消息的格式是一个key-value键值对，key={consumer_group_id, topic_id, partition_id}，value={offset, }

## 4.分区多副本机制

### 4.1.特点

1. 副本之间采用一主多从
2. leader 副本负责处理读写请求，follower 副本只负责与 leader 副本的消息同步
3. 分区副本处于不同的broker中

- 问题: 与sql/redis不同，kafka的从节点为什么不提供读能力
  1. kafka并不是“读多写少”的读写分离场景，它通常涉及到生产/消费msg
  2. kafka副本机制采用的是异步消息拉取，因此存在主从数据不一致性问题。如果从节点提供读，数据不一致
  3. 主写从读无非就是为了减轻leader节点的压力，将读请求的负载均衡到follower节点，如果Kafka的分区相对均匀地分散到各个broker上，同样可以达到负载均衡的效果，没必要刻意实现主写从读增加代码实现的复杂程度

### 4.2.ISR、HW、LEO

1. ISR(In Sync Replicas) 存活节点
   - 存活条件
      1. 节点必须和ZK保持会话
      2. follower的leo落后leader的leo不超过阈值
2. HW 高水位
   - 概念: 表示了一个特定的消息偏移量，消费者只能拉取到这个offset之前的消息
      - 作用
         1. 定义消息的可见性
         2. 帮助kafka完成副本同步
3. LEO(log end offset)
   1. 概念: 它表示了当前日志文件中下一条待写入消息的offset
4. LEO和HW的关系
   1. ISR集合中的每个副本都会维护自身的LEO
   2. ISR集合中最小的LEO，就是分区副本的HW(对于消费者而言，只能消费HW之前的消息)

### 4.2.1.HW/LEO更新过程(主从数据同步过程)

[背景知识]

1. leader副本: 既保存了自己的HW、LEO，还保存了所有follower的HW、LEO
2. follower副本: 只保存了自己的HW、LEO


[leader副本和follower副本的HW/LEO更新过程] 以单分区一主一从举例介绍

1. 初始状态，mq中没有消息
   - leader
     - HW=0, LEO=0
     - remote_LEO=0
   - follower
     - HW=0, LEO=0
2. 生产者向topic分区发送一条msg，此时，leader副本成功将msg写入了本地磁盘，故`leader的LEO`的值更新为1(写入消息的位移fetchOffset=0)
   - leader
     - HW=0, `LEO=1`
     - remote_LEO=0
   - follower
     - HW=0, LEO=0
3. follower尝试从leader拉取消息(fetchOffset=0，拉取位移为0的消息)，此时，有数据拉取了，`follower的leo`也更新为1
   - leader
     - HW=0, LEO=1
     - remote_LEO=0
   - follower
     - HW=0, `LEO=1`
4. 此时，leader和follower的leo都是1，但是各自的HW依然是0，还没有被更新。**他们需要在下一轮的拉取中被更新**
   1. 在新一轮的拉取请求中，由于fetchOffset=0的消息已经拉取成功，所以follower副本这次请求拉取的fetchOffset=1
   2. leader副本接收到拉取fetchOffset=1的请求后，会:①先更新remote_LEO=1 ②再`更新自己的HW=1` ③最后，将当前已经更新过HW=1的消息，告诉follower副本
      - leader
        - `HW=1`, LEO=1
        - remote_LEO=1
      - follower
        - HW=0, LEO=1
   3. `follower副本`接收到leader更新HW=1消息后，也会将自己的`HW更新为1`(至此，一次完整的消息同步周期就结束了)
      - leader
        - HW=1, LEO=1
        - remote_LEO=1
      - follower
        - `HW=1`, LEO=1

可以看到，follower副本的HW更新需要一轮额外的拉取请求才能实现。也就是说，leader的HW的更新、follower的HW的更新，存在时间上的偏差 ==> 这会导致“数据丢失”或“数据不一致”(发生日志截断) ==> 引入leader epoch来解决这个问题

- 发生数据丢失的场景
   1. 当follower还没来得及更新自己的HW时，即使follower's LEO和leader's LEO一致
   2. 当A宕机了，B变成了leader，那么，此时B会执行“日志截断” ==> 数据丢失

### 4.2.2.leader epoch

引入leader epoch机制，解决数据截断导致丢失数据的问题

1. 所谓leader epoch，可以认为是leader版本，它有两部分数据组成
   1. epoch: 每当副本领导权发生变化，epoch+1
   2. start offset(起始位移): leader副本在该epoch值上写入的首条消息的位移
2. 假设现在有两个leader epoch <0,0> 和 <1,120>
   - 第一个 Leader Epoch 表示版本号是 0，这个版本的 Leader 从位移 0 开始保存消息，一共保存了 120 条消息
   - 之后，Leader 发生了变更，版本号增加到 1，新版本的起始位移是 120

过程

1. 初始时
   - A是leader:   (epoch,offset) = <0,0>    HW=1, LEO=1, remote_LEO=1
   - B是follower: (epoch,offset) = <0,0>    HW=0, LEO=1
2. 之后，A宕机，B升级为主，之后A又重新上线
   1. 此时B会向A发送一个特殊的请求，去获取leader的LEO值(该值当前=1)
   2. B发现该LEO值不比它自己的LEO值小，而且缓存中也没有保存任何起始位置大于1的epoch条目 ==> 因此，B无需执行任何日志截断操作(即：副本机制是否执行日志截断不再依赖于高水位)
   3. 当A重启回来后，执行与B相同的逻辑判断，发现也不用执行日志截断，至此，offset=1的那条消息就在两个副本中均得到保留
   4. 后面当生产者程序向B写入新消息时，副本B所在的Broker缓存中，会生成新的Leader Epoch，条目：[Epoch=1, Offset=1]。之后，副本 B 会使用这个条目帮助判断后续是否执行日志截断操作。这样，通过 Leader Epoch 机制，Kafka 完美地规避了这种数据丢失场景。

## 5.消费者组

### 5.1.特点

1. 消费者/分区的关系
   1. 一个消费者，可以消费多个分区
   2. 一个分区，只能分配给一个消费者
2. topic、partition、consumer_group、consumer
   1. 一个topic可以有多个分区，每个分区可以指定给一个consumer消费
   2. 一个消息，只能被某个消费组内的某个消费者消费一次

### 5.2.consumer与partition的分配机制

1. group中的consumer数 < topic中的partition数 ==> group中的consumer就会消费多个partition
2. group中的consumer数 = topic中的partition数 ==> group中的一个consumer就会消费topic中的一个partition
3. group中的consumer数 > topic中的partition数 ==> group中就会有一部分的consumer处于空闲状态

### 5.3.Coordinator协调者

功能: 每个broker上都有一个Coordinator进程(所有broker启动时，都会创建和开启Coordinator进程)，它专门为consumer-group服务，负责消费者组执行rebalance、提交位移管理、组成员管理等
1. rebalance
2. 位移管理
   - consumer提交位移时，其实是想Coordinator所在的broker提交位移
3. 组成员管理
   1. consumer启动时，向Coordinator所在的broker发送各种请求
   2. Coordinator负责执行consumer-group的注册、成员管理记录等元数据管理操作

### 5.4.[rebalance重平衡](https://www.bilibili.com/video/BV1HP4y157tx?p=29)

先说明，之前理解错了，rebalance并不会造成数据的丢失！是主从复制的时候，才会导致丢失！

#### 5.4.1.概念

1. **什么是重平衡rebalance**
   - 消费者组，多个分区与多个消费者，重新匹配
   - rebalance过程需要借助Coordinator组件
2. **发生rebalance的条件**
   1. 消费者组成员个数变化 
   2. 订阅主题的partition数发生变更 
   3. 订阅主题数发生变化(正则topic)
3. **<u>危害</u>**
   - 消费暂停
   - 消费突增
   - 消费重复
4. **消费者partition分配策略**: 范围、轮询、Sticky

#### 5.4.2.重平衡过程

1. “消费者端”的重平衡场景剖析
   1. 重平衡时
      1. 所有的消费都停止
      2. consumer与分区绑定关系都解除
   2. 加入组
      1. consumer会向Coordinator发送JoinGroup请求消息(该消息包含topic等信息)
      2. Coordinator一旦收到了所有consumer的JoinGroup请求消息后，Coordinator会从consumer中选择一个consumer-leader(注意：这里的consumer-leader和分区副本的leader完全不一样，consumer-leader的作用主要是①收集所有consumer的订阅消息，根据这些消息，指定具体的分区消费分配方案)
      3. 选出consumer-leader后，Coordinator会把消费者组的所有订阅消息，封装进JoinGroup请求的响应体中，然后发给consumer-leader
   3. 等待接收leader-consumer下发分区消费分配方案(尽可能使得之前绑定的consumer和partition还分在一组)
      1. 之后，consumer-leader统一作出分配方案，然后，向Coordinator发送SyncGroup请求(里面包含了分区分配方案)
      2. Coordinator接收到分配方案后，将SyncGroup发送给所有consumer，这样consumer就知道自己该消费哪个分区了
2. rebalance存在的问题: 消息重复提交
   1. 场景
      1. Consumer1消费消息，还未提交offset，此时发生rebalance
      2. rebalance完成后，因为offset未被提交，消息会发到其他的Consumer2
      3. consumer2消费完消息后，提交了offset
      4. 之后，Consumer1消费消息完成，再次提交了offset==>就会出现错误(一个消息被提交了两次offset)
   2. 解决方案: Coordinator与Generation
      - Coordinator每次rebalance，会标记一个Generation给到consumer，每次rebalance时Generation都会+1，consumer提交offset时，Coordinator会比对Generation是否一致，不一致就拒绝提交offset(即: consumer1超时完成时，此时Generation已经+1，对比发现Generation不一致，提交offset会被拒绝; consumer2允许提交offset ==> 保证了1条消息只能提交一次offset)

## 6.kafka高性能的原因

1. 日志采用: 顺序追加写+log_segment
   1. 日志采用顺序追加写
   2. 每个partition对应一个log，又将该log划分成多个log_segment
2. page-cache 
   1. 向 Kafka 发送数据并不是真要等数据被写入磁盘才会认为成功，而是只要数据被写入到操作系统的`页缓存`(Page Cache)上就可以了，随后操作系统根据 LRU 算法会定期将页缓存上的“脏”数据落盘到物理磁盘上。这个定期就是由提交时间来确定的，默认是 5 秒。一般情况下我们会认为这个时间太频繁了，可以适当地增加提交间隔来降低物理磁盘的写操作。
   2. 当然你可能会有这样的疑问: 如果在页缓存中的数据在写入到磁盘前机器宕机了，那岂不是数据就丢失了。的确，这种情况数据确实就丢失了，但鉴于 Kafka 在软件层面已经提供了多副本的冗余机制，因此这里稍微拉大提交间隔去换取性能还是一个合理的做法。
3. 发送消息: 批量+压缩，降低带宽
4. 文件分段: log_segment
5. 零拷贝
    1. 为什么kafka能使用零拷贝呢？
       - 结合Kafka的消息有“多个订阅者”的使用场景，生产者发布的消息一般会被“不同的多个消费者”消费多次
    2. 对比: 
         1. 传统一次读写(4次拷贝，4次用户态内核态的切换)
            - 磁盘文件==>内核空间的读取缓冲区page cache==>应用程序==>socket缓冲区==>网卡接口==>消费者进程
         2. 零拷贝(2次拷贝，2次用户态内核态的切换)
            - 磁盘文件==>内核空间的读取缓冲区page cache==>网卡接口==>消费者进程
    3. 什么情况下使用了零拷贝
       1. 基于mmap的索引
       2. 日志文件读写所用的transportLayer 
         

## 7.kafka常见问题

### 7.1.如何保证消息不丢失

1. 生产者: 消息发送+回调
   1. producer.send(msg, callback)——捕获失败的消息，保存到db中，重试
   2. 重试次数 > 1
2. 消费者: 手动提交消息
3. 配置
   1. ack
      - 0 - 不需要任何的broker收到消息，就立即返回ack给生产者
      - 1 - Leader收到消息，消息写入到log，才返回ack给生产者
      - -1或all，min.insync.replicas>1
   2. unclean，leader，election，enable配置为false(不允许选择OSR中的从节点作为主节点)
   3. 减小broker刷盘间隔(page cache刷到磁盘的时间间隔)

### 7.2.消息顺序性保证(rocketMQ实现了该机制)

说明: MQ只能保证partition内的局部有序，不能保证全局有序

结论: kafka只能一个分区，一个生产者，一个消费者，才能保证有序

多个分区，一个消费者组(多个消费者)，采用key的取模的方式发送消息，当发生rebalance时，依然会导致消息乱序

1. 生产者: 需要有序的一组消息，通过指定partition发送到同一个partition中
2. 消费者: 注册有序的监听

### 7.3.幂等

1. 数据库/缓存
2. 全局唯一ID: 带业务表示的ID，来进行幂等判断

### 7.4.消息堆积

1. 消息pull时间间隔过大
2. 消费耗时
3. 消费并发度
4. 单线程计算
5. 如果你使用的是消费者组，确保没有频繁地发生rebalance

### 7.5.三种消息传递机制

at most once: 丢失消息

at least once: 消息重复

[exactly once: 不丢消息，消息不重复](http://www.jasongj.com/kafka/transaction/)

1. producer发送到broker保证幂等性：producer_ID 和 meg_seq_num
   1. 消息的格式<producer_ID, topic, partition, meg_seq_num>，每次commit一条消息时，将对meg_seq_num增加1
   2. broker接收每条消息后，若其序号比broker维护的序号大1，则broker接收消息，否则丢弃消息
      1. 若meg_seq_num比broker维护的序号大1以上，说明中间有数据尚未写入，也即乱序，此时broker拒绝该消息
      2. 若meg_seq_num <= broker维护的序号，说明该消息已经被保存，即为重复消息，此时broker直接丢弃消息
2. 事务性保证
   1. 分析：上述幂等只能保证单个producer对于同一个<topic, partition>的exactly once语义。存在一下问题：
      1. 并不能保证写操作的原子性，即：多个写操作，要么全部被commit，要么不被commit。
      2. 不能保证多个读写操作的原子性
   2. 事务：保证了应用程序将生产数据和消费数据当做一个原子单元来处理，要么全部成功，要么全部失败
      1. 应用程序必须维护一个事务ID
      2. 

### 7.6.死信队列

1. 用来存放消费失败超过设置次数的消息，通常用来作为消息重试
2. 特征
   1. 消息不会被消费者正常消费
   2. 有效期与正常消息相同，均为3天
   3. 死信队列就是一个特殊的topic
   4. 如果一个消费这组未产生死信消息，则不会为其创建相应的死信队列


### 7.7.延时队列

存放在指定时间后被处理的消息，通常用来处理一些具有过期性操作的业务(如10min内未支付则取消订单)

## 8.常见问题

### 8.1.leader总是-1，怎么破？

- 在生产环境中，你一定碰到过“某个topic的分区不能工作了”的情形，使用命令查看的时候，会发现leader是-1，于是你使用各种命令都无济于事，最后只能用“重启大法”。
- 解决方案: 删除Zookeeper节点/controller，触发controller重新选举。controller重选举能够为所有分区重刷分区状态，可以有效解决因不一致导致的leader不可用问题。


### 8.2.[kafka事务](http://matt33.com/2018/11/04/kafka-transaction/)


---

参考链接

https://www.bilibili.com/video/BV1Xf4y1u7uD?p=38

https://www.bilibili.com/video/BV1cf4y157sz?p=102

尚硅谷rocketMq: https://www.bilibili.com/video/BV1cf4y157sz?p=1&share_medium=android&share_plat=android&share_session_id=18c0028f-b7b6-4fb6-b2ea-34a5b100ccaf&share_source=WEIXIN&share_tag=s_i&timestamp=1639821947&unique_k=zcZy8VO

rocketMQ事务消息:  https://time.geekbang.org/column/article/111269

