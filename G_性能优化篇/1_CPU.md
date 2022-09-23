

# 1. 到底应该怎么理解“平均负载”？

1. top的S列(status): 进程状态
   1. R: running或runnable
   2. D: Disk-Sleep，不可中断状态，一般表示进程正在跟硬件交互，等待IO返回
      1. 正常情况下，D状态的进程时间都很短，一般可以忽略
      2. 当系统或硬盘发生故障，进程可能处于D状态保持很久，甚至导致系统中出现大量D进程。这时，就要注意下，系统是不是出现了IO性能
   3. Z: Zombie，僵尸进程，进程结束了，父进程没回收它的资源
   4. S: Interrupt-Sleep，可中断状态睡眠，表示进程因为等待某个时间被系统挂起。当进程等待的事件发生时，他会被唤醒进入R状态
   5. I: idel，空闲状态，用在不可中断睡眠的内核线程上。
      - 前面说了，硬件交互导致的不可终端进程用D表示，但对某些内核线程来说，他们有可能实际上并没有任何负载，用idle正式为了区分这种情况
      - 注意：D状态的进程会导致load average升高，I状态的进程却不会
   6. T/t: Stopped/Traced，进程处于暂停或跟踪状态
      1. stopped: 进程因为等待某个信号，进入暂停状态
      2. Traced: Gdb进程时，进程会进入跟踪状态
   7. X: Dead，进程已经消亡，在top/ps中看不到它
2. 当发现系统变慢时，要查看系统的负载情况（top/uptime)
    ```shell
    $ uptime 
    02:34:03 up 2 days, 20:14,  1 user,  load average: 0.63, 0.83, 0.88
    #                                               近1min,5min,15min的平均负载
    ```

## 1.1.平均负载load average - 平均活跃进程数
1. 定义：在一个时间段内，进程(R状态/D状态)个数
   1. R状态
      1. Running：正在使用CPU的进程
      2. Runable：正在等待CPU的进程
   2. D状态：不可中断状态的进程，常见的是等待硬件设备IO响应，就是ps看到的D状态(Disk-Sleep)的进程
      - 进程向磁盘读写数据时，为了保证数据的一致性，在等待磁盘回复前，它是不能被其他进程终端的，这个时候的进程处于不可中断状态
      - 不可中断进程实际上是系统对进程和硬件设备的一种保护机制
    3. 总结：造成load average高的原因，分为2种
      1. IO密集型：进程因为等待IO返回，导致负载高
      2. CPU密集型：进程
2. 分析：近1min，5min，15min
      1. 近1min >> 5min / 15min ==> 当前时刻，系统负载突然升高
      2. 近1min << 5min / 15min ==> 过去一段时间内，出现系统负载过载


## 1.2.load average取值多大才合适

1. 结论：理想情况下，load average == CPU核心个数，才认为CPU被充分利用。
   - 当load average > CPU核心个数，系统出现过载
2. 查看CPU个数: grep 'model name' /proc/cpuinfo | wc -l

## 1.3.load average / CPU使用率

CPU使用率：是单位时间内CPU繁忙情况的统计，跟load average并不一定完全对应。
1. CPU密集型：CPU使用率 == load average
   - 使用大量的CPU，导致load average升高，此时二者是一致的
2. 大量等待CPU的进程调度，也会导致load average升高，此时CPU使用率也会比较高
3. I/O密集型: load average高，但CPU使用率不一定高
   - 等待IO导致load average升高，但CPU使用率不一定高

## 1.4.load average案例: 分析load average升高的根源

分析工具：iostat、mpstat、pidstat，找出load average升高的根源

1. 机器配置：2CPU、8GB内存
2. 压测工具：stress
    ```shell
    stress --cpu 1 --timeout 600  # CPU密集型: 构造某个CPU使用率100%的场景
    stress    -i 1 --timeout 600  # I/O密集型: 构造IO压力
    stress    -c 8 --timeout 600  # 大量等待CPU的进程调度: 构造多进程使用CPU运算的场景 
    ```
3. 分析工具
   1. uptime：查看平均负载load average
      - watch -d uptime：-d参数表示高亮显示变化的区域
   2. mpstat：查看CPU使用率变化的情况
        ```shell
        $ mpstat -P ALL 5   # -P ALL 表示监控所有CPU，后面数字5表示间隔5秒输出一组数据
        Linux 4.15.0 (ubuntu) 09/22/18 _x86_64_ (2 CPU)  # 系统总共有2个CPU
        # 采集时间点  CPU号码 用户态   优先级   内核态   等待IO  硬中断   软中断                           CPU空闲   
        13:30:06     CPU    %usr   %nice    %sys %iowait    %irq   %soft  %steal  %guest  %gnice   %idle
        13:30:11     all   50.05    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00   49.95  # 所有CPU平均值
        13:30:11       0    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00  100.00  # CPU-0
        13:30:11       1  100.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00  # CPU-1
        ```
   3. pidstat(进程分析工具)：实时查看进程CPU\内存\IO\上下文切换等性能指标
        ```shell
        $ pidstat -u 5 1    # 间隔5秒后输出一组数据
                                      用户    系统            等待CPU使用率
        13:37:07      UID       PID    %usr %system  %guest   %wait    %CPU   CPU  Command
        13:37:12        0      2962  100.00    0.00    0.00    0.00  100.00     1  stress
        ```

### 1.4.1.(场景一)CPU密集型进程

1. 模拟一个CPU使用率100%的场景：`stress --cpu 1 --timeout 600`
2. 查看load average变化情况：`watch -d uptime`，输出`...,  load average: 1.00, 0.75, 0.39`
   - 可以看到，load average最近1min负载=1.00，表示有一个CPU跑满了
3. 查看CPU使用率变化情况：mpstat -P ALL 5
    ```shell
    $ mpstat -P ALL 5   
    Linux 4.15.0 (ubuntu) 09/22/18 _x86_64_ (2 CPU) 
    13:30:06     CPU    %usr   %nice    %sys %iowait    %irq   %soft  %steal  %guest  %gnice   %idle
    13:30:11     all   50.05    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00   49.95  # 所有CPU平均值
    13:30:11       0    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00  100.00  # CPU-0
    13:30:11       1  100.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00    0.00  # CPU-1
    ```
    - 现象：CPU[0]的指标全是0，CPU[1]的usr%=100
    - 分析：CPU[1]使用率为100%，但是它的iowait=0 ==> 说明，load average升高是由于CPU使用率为100%
4. 查看到底是哪个进程导致了CPU使用率为100%：pidstat
    ```shell
    $ pidstat -u 5 1   # 间隔5秒后输出一组数据
    13:37:07      UID       PID    %usr %system  %guest   %wait    %CPU   CPU  Command
    13:37:12        0      2962  100.00    0.00    0.00    0.00  100.00     1  stress
    ```
    - 可以看到，stress进程的CPU使用率为100%

### 1.4.4.(场景二)IO密集型进程

1. 使用stress，模拟IO压力，即不停地执行sync：`stress -i 1 --timeout 600`
2. 查看load average变化情况：`watch -d uptime`，输出`...,  load average: 1.00, 0.75, 0.39`
   - 可以看到，load average最近1min负载=1.06，表示有一个CPU跑满了
3. 查看CPU使用率变化情况：mpstat -P ALL 5
    ```shell
    $ mpstat -P ALL 5 1
    Linux 4.15.0 (ubuntu)     09/22/18     _x86_64_    (2 CPU)
    13:41:28     CPU    %usr   %nice    %sys %iowait    %irq   %soft  %steal  %guest  %gnice   %idle
    13:41:33     all    0.21    0.00   12.07   32.67    0.00    0.21    0.00    0.00    0.00   54.84
    13:41:33       0    0.43    0.00   23.87   67.53    0.00    0.43    0.00    0.00    0.00    7.74
    13:41:33       1    0.00    0.00    0.81    0.20    0.00    0.00    0.00    0.00    0.00   98.99
    ```
    - 现象：CPU[1]的指标几乎≈0，CPU[0]：{%usr=43%，%sys=23.87，%iowait=67.53，%idle=7.74}
    - 分析：CPU[1]被跑满了，%idle=7.74%，空闲率≈0，系统CPU(%sys)升高到23.87%，而iowait高达67.53% ==> 说明，load average升高由于iowait的升高，阻塞在IO上了（负载的升高，不是因为CPU使用率高，而是因为等待IO导致的）
4. 查看到底是哪个进程导致了CPU使用率为100%：pidstat
    ```shell
    $ pidstat -u 5 1
    Linux 4.15.0 (ubuntu)     09/22/18     _x86_64_    (2 CPU)
    13:42:08      UID       PID    %usr %system  %guest   %wait    %CPU   CPU  Command
    13:42:13        0       104    0.00    3.39    0.00    0.00    3.39     1  kworker/1:1H
    13:42:13        0       109    0.00    0.40    0.00    0.00    0.40     0  kworker/0:1H
    13:42:13        0      2997    2.00   35.53    0.00    3.99   37.52     1  stress
    13:42:13        0      3057    0.00    0.40    0.00    0.00    0.40     0  pidstat
    ```
    - 可以看到，还是stress进程导致的load average升高


### 1.4.3.(场景三)CPU密集型进程


1. 构造 (进程数>>CPU数) 的场景：`stress -c 8 --timeout 600`，系统CPU个数为2，构造了8个进程
2. 查看load average变化情况：`watch -d uptime`，输出`...,  load average: 7.97, 5.93, 3.02`
   - 可以看到，由于进程数8远远大于2，导致CPU处于严重过载状态，平均负载高达7.97
3. 用pidstat查看进程情况：可以看到，8个进程在争抢2个CPU，每个进程等待CPU的时间(%wait)高达75% ==> 这些超出CPU计算能力的进程，导致CPU过载
    ```shell
    $ pidstat -u 5 1
    14:23:25      UID       PID    %usr %system  %guest   %wait    %CPU   CPU  Command
    14:23:30        0      3190   25.00    0.00    0.00   74.80   25.00     0  stress
    14:23:30        0      3191   25.00    0.00    0.00   75.20   25.00     0  stress
    14:23:30        0      3192   25.00    0.00    0.00   74.80   25.00     1  stress
    14:23:30        0      3193   25.00    0.00    0.00   75.00   25.00     1  stress
    14:23:30        0      3194   24.80    0.00    0.00   74.60   24.80     0  stress
    14:23:30        0      3195   24.80    0.00    0.00   75.00   24.80     0  stress
    14:23:30        0      3196   24.80    0.00    0.00   74.60   24.80     1  stress
    14:23:30        0      3197   24.80    0.00    0.00   74.80   24.80     1  stress
    14:23:30        0      3200    0.00    0.20    0.00    0.20    0.20     0  pidstat
    ```

## 1.5.小结

- load average过高的原因，有以下3点
   1. CPU密集型任务多，执行频繁计算，此时CPU使用率高
   2. I/O密集型任务多，在IO上等待，此时CPU使用率不高
   3. 执行CPU计算的线程/进程多，(上下文切换频繁)，导致大量线程/进程等待CPU

# 2.经常说的CPU上下文切换是什么意思？

多进程在竞争一个CPU时，进程没有执行，为什么会导致系统load average升高呢？ ==> 进程上下文切换是罪魁祸首

1. 为什么上下文切换慢
   1. CPU上下文切换：寄存器、程序计数器
   2. 根据任务不同，分为
      1. 进程上下文切换：一个进程切换到另一个进程运行，每次切换平均耗时(几十纳秒~几微秒)
         1. 用户空间资源：虚拟内存、栈、全局变量等
         2. 内核空间状态：内核堆栈、寄存器等
         3. TLB失效(虚拟内存到物理内存的映射关系)：当虚拟内存更新后，TLB也要刷新
      2. 线程上下文切换
         1. 不同进程的线程间切换
            - 切换方式与“进程上下文切换”一样
         2. 相同进程的线程间切换
            1. 虚拟内存不用切换，TLB不失效
            2. 只需要切换线程私有数据、寄存器等不共享的数据
      3. 中断上下文切换
         1. 中断会打断进程的正常调度和执行，转而调用中断处理程序，响应设备文件。而在打断其他进程时，就需要将进程当前的状态保存下来，之后，以便恢复中断现场
         2. 与进程上下文切换不同，
            1. 中断处理比进程具有更高的优先级 ==> 中断上下文切换不会与进程上下文切换同时发生
            2. 中断上下文切换不涉及到进程的用户态，所以中断过程打断一个正处在用户态的进程，也不需要保存和恢复这个进程的虚拟内存、全局变量等用户态资源
            3. 中断上下文只包括内核态中断服务程序执行所必需的状态，包括CPU寄存器、内核堆栈、硬件中断参数等
2. vmstat：查看“系统”的内存使用情况，也常用来分析CPU上下文切换和中断的次数
    ```shell
    $ vmstat 5  # 每隔5秒输出1组数据
    procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
    r  b   swpd   free   buff   cache   si   so    bi    bo   in   cs us sy id wa st
    0  0      0 7005360  91564 818900    0    0     0     0   25   33  0  0 100  0  0
    # 上下文切换cs=33次，系统中断次数in=25次，就绪队列长度r和不可中断状态进程数b都是0
    ```
    - cs (context switch) 每秒上下文切换的次数
    - in (interrupt) 每秒中断的次数
    - r (running + runnable) 就绪队列的长度，(正在运行+等待运行)CPU进程数
    - b (blocked) 处于不可中断睡眠状态的进程数
3. pidstat -w：查看每个进程的详细情况
    ```shell
    $ pidstat -w 5 
    Linux 4.15.0 (ubuntu)  09/23/18  _x86_64_  (2 CPU)
    08:18:26      UID       PID   cswch/s nvcswch/s  Command
    08:18:31        0         1      0.20      0.00  systemd
    08:18:31        0         8      5.40      0.00  rcu_sched
    ```
    - cswch/s：每秒自愿上下文切换次数
      - 进程无法获取所需自愿，导致的上下文切换，比如，IO、内存等系统自愿不足时
    - nvcswch/s：每秒非自愿上下文切换次数
      - 进程由于时间片到达，被系统强制调度，比如，大量进程都在争抢CPU

# 3.某个进程的CPU使用率达到100%，该怎么排查

1. top 查看系统情况
    ```shell
    $ top # 默认每3秒刷新一次
    top - 11:58:59 up 9 days, 22:47,  1 user,  load average: 0.03, 0.02, 0.00
    Tasks: 123 total,   1 running,  72 sleeping,   0 stopped,   0 zombie
    %Cpu(s):  0.3 us,  0.3 sy,  0.0 ni, 99.3 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
    KiB Mem :  8169348 total,  5606884 free,   334640 used,  2227824 buff/cache
    KiB Swap:        0 total,        0 free,        0 used.  7497908 avail Mem

    PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
        1 root      20   0   78088   9288   6696 S   0.0  0.1   0:16.83 systemd
        2 root      20   0       0      0      0 S   0.0  0.0   0:00.05 kthreadd
        4 root       0 -20       0      0      0 I   0.0  0.0   0:00.00 kworker/0:0H
    ```
    - 用户CPU和Nice CPU高 ==> 说明用户态进程占用了较多的CPU，所以应该着重排查“进程”的性能问题
    - 系统CPU高 ==> 说明内核态占用了较多的CPU，所以应该着重排查“内核线程or系统调用”的性能问题
    - IO等待CPU高 ==> 等待IO的时间比较长，所以应该着重排查“系统存储”是不是出现了问题
    - 软中断和硬中断高 ==> 说明软中断/硬中断处理程序占用了较多的CPU，所以应该着重排查内核中的中断服务程序
2. 怎么占用CPU的到底是代码里的那个函数呢？ ==> perf
   1. perf
      1. 分析系统的各种事件和内核性能
      2. 分析定位应用程序的性能问题
   2. 用法
      1. perf top 实时显示占用CPU时钟最多的函数或指令，用于查找热点函数
            ```shell
            $ perf top
            # 采样数         事件类型(CPU时钟事件)       事件总数
            Samples: 833  of event 'cpu-clock', Event count (approx.): 97742399
            Overhead  Shared    Object     Symbol
            7.28%     perf         [.]     0x00000000001f78a4
            4.72%     [kernel]     [k]     vsnprintf
            4.32%     [kernel]     [k]     module_get_kallsym
            3.65%     [kernel]     [k]     _raw_spin_unlock_irqrestore
            ```
           - Overhead 该服务的性能事件在所有采样中的比例，用百分比表示
           - Shared 该函数或指令所在的动态共享对象，如内核、进程名、动态链接库、内核模块等
           - Object 动态共享对象的类型，比如[.]表示用户空间的可执行程序、或者动态链接库，而[k]表示内核空间
           - Symbol 函数名。当函数名未知时，用16进制地址来表示
      2. perf record/report 保存到文件中

# 3.系统CPU使用率很高，但为啥找不到高CPU的应用

碰到常规问题无法解释的CPU使用率情况时，首先要想到有可能是“短时”应用导致的问题，比如有可能是下面这2中情况：
1. 应用里直接调用了[其他二进制程序，这些程序通常运行时间比较短]，通过top等工具也不容易发现
2. 应用本身在不停地[崩溃重启]，而启动过程的资源初始化，很可能会占用相当多的CPU

# 4.总结: 如何分析出系统CPU的瓶颈在哪里

1. CPU性能指标
   - 系统整体CPU使用率：uptime\vmstat\mpstat\/proc/stat
   - 进程CPU使用率：top\pidstat\ps
   1. CPU使用率：描述了非空闲时间占总CPU时间的百分比，根据CPU上运行任务的不同，又分为
      1. 用户CPU
         1. 包含：用户态CPU使用率 和 低优先级用户态CPU 使用率
         2. 用户CPU使用率高，说明应用程序比较繁忙 
      2. 系统CPU
         1. CPU在内核运行的时间百分比（不包含中断）
         2. 系统CPU使用率高，说明内核比较繁忙
      3. 等待I/O CPU
         1. 通常被称为iowait，表示等待IO的时间百分比
         2. iowait高，说明系统与硬盘设备的IO交互时间比较长
      4. 软中断和硬中断 
         1. 内核调用软中断/硬中断处理时间百分比
         2. 他们的使用率高，说明系统发生了大量的中断
2. load average 平均负载
   - 工具：uptime/top
   1. 近1min、近5min、近15min
   2. 理想情况下，load average == CPU核心数，表示每个CPU恰好被充分利用
   3. 负载升高的3个原因
      1. CPU密集型：正在执行的进程
      2. IO密集型
      3. 多进程密集型：等待执行的进程
3. 上下文切换
   - 系统上下文切换：vsstat
   - 进程上下文切换：pidstat 
