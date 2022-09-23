
# rocketMQ

## 1.基本概念 

1. 消息
2. product
3. consumer
4. topic
5. queue
6. Bookie
7. tag消息标识: 每个消息拥有唯一的messageID，且可以携带具有业务表示的key
8. NameServer

## 2.NameServer功能

作为broker与topic路由的注册中心，支持broker的动态注册与发现

1. 无状态，动态列表
2. 这也是与Zookeeper的重要区别之一，Zookeeper是有状态的

### 2.1.功能
1. broker管理
   1. 接收broker集群的注册信息&&保存下来作为路由与信息的基本数据
   2. 提供心跳检测机制，检查broker是否存活
2. 路由信息管理: 保存了broker集群的整个路由信息和用于客户端查询的队列信息。product和consumer通过NameServer可以获取整个broker集群的路由信息，从而进行消息的投递与消费
    1. 路由注册: broker维持和NameServer的心跳，每30s发送一次
    2. 路由踢出
    3. 路由发现: rocketMQ的路由发现采用pull模型。客户端每隔30s主动拉取topic路由信息

### 2.2.工作流程

1. 启动NameServer，NameServer启动后开始监听端口，等待broker、producter、consumer连接
2. 启动broker时，broker会与所有的NameServer建立并保持长连接，每隔30s向NameServer定时发送心跳包
3. 收发消息前，可以先创建topic，创建topic时需要指定该topic要存储在那些broker上(broker与topic的绑定关系写入到NameServer中)
4. producter发送消息: 从NameServer获取topic路由信息，根据分配算法选择queue发送
5. consumer接收消息: 从NameServer获取订阅的topic的路由信息，然后根据分配算法选择queue数据消费
6. Producer、Consumer连接到一个Broker后，查询topic的Partition对应的broker，建立连接。
7. 生产: 写入Bookie，后在broker有消息的缓存
8. 消费: 先读取broker的缓存，没有再去bk读取消息
9. Bookie: bk的物理节点
10. Broker启动会去zk抢锁，抢到的会定时去拉取其他Broker的负载情况，进行负载均衡操作。
11. broker无状态，消费进度等都存储在bk中，可以直接进行扩容
12. bk的Bookie也是无状态，其存储的segment按照单调递增的ID，固定时段或者固定大小分段。扩容Bookie后，新segment数据会自动写入到新机器中。

## 3.BookKeeper

BookKeeper: http://matt33.com/2019/01/28/bk-store-realize/


1. Ledger就是segment的具体实现。Ledger3 有两个Fragment，一个Fragment会有多个Entry，一个Entry代表写入的一个批次的数据。
2. Bookie启动会先zk注册节点信息，bk是属于slave-slave架构，没有leader和follow的区别。
3. 系统创建Ledger的时候，会分配一定数量的Bookie，构成一个Ensemble
4. bk client创建Ledger，带上写入的配置: openLedger(5,3,2) ，Ensemble内节点数量，用于控制Ensemble的整个的带宽、数据备份数目、等待刷盘节点数目–成功这么多就返回。
    - bk默认是3,3,3。在组内节点通过round robin的方式循环写入目标节点中，比如1,2,3、 2,3,4、 4,5,1。写入这三个Bookie中，用消息ID对组内节点数目取模，得到一个消息存储在哪个Bookie中。反过来，通过BookieID也可以得到存储了这个Ledger的哪些消息。
5. 通过Fencing机制保证同时只有一个broker写入ledger，分配自增的EntryID，保证EntryID的唯一，每个Bookie都有Auditor。线程检查自身的Entry是否有缺失，发现缺失会从其他的副本复制数据。


## 4.rocketmq事务消息:smile:

https://blog.csdn.net/hosaos/article/details/90050276

**场景: **同时保证(本地事务+发送消息到MQ)都成功，例如，生成订单(插入到订单表)，增加积分(发送消息到mq)

### 4.1.提前理解几个概念

1. 本地事务
2. 生产者
3. broker
4. 消费者
5. 两个topic
   1. 半消息队列: 此时，消息不能被consumer消费
   2. 半消息op队列: 执行commit/rollback的消息，能被consumer消费

### 4.2.[Rocketmq执行过程(类似2阶段提交)](https://blog.csdn.net/weixin_34452850/article/details/88851419?utm_source=app&app_version=4.14.0)

1. Producer向Broker端发送Half Message
2. Broker ACK，Half Message发送成功
3. Producer执行本地事务。本地事务完毕，根据事务的状态，Producer向Broker发送二次确认消息，确认该Half Message的Commit或者Rollback状态。
4. Broker收到二次确认消息后
   1. 对于Commit状态，则直接发送到Consumer端执行消费逻辑
   2. 对于Rollback，则直接标记为失败，一段时间后清除，并不会发给Consumer
5. 正常情况下，到此分布式事务已经完成

剩下要处理的就是超时问题，即一段时间后Broker仍没有收到Producer的二次确认消息

1. 针对超时状态，Broker主动向Producer发起消息回查
2. Producer处理回查消息，返回对应的本地事务的执行结果
3. Broker针对回查消息的结果，执行Commit或Rollback操作，同Step4

**broker定时回查事务状态**

1. 场景: 发送commit/rollback到broker中的<u>半消息op队列</u>可能会丢失
2. 解决方案: broker定时去查询本地事务的执行结果，查询到commit/rollback后，执行消息的提交/丢弃


## 5.rocketmq常见问题

### 5.1.如何保证不丢失消息

1. broker: 刷盘策略=同步刷盘
2. 生产者: send(msg, handler)
3. 消费者: 手动commit offset

### 5.2.如何保证顺序消费

1. 多个queue只能保证单个queue里的顺序，queue是典型的FIFO，天然顺序
2. 多个queue同时消费是无法绝对保证消息的有序性的
