# 1.语法知识

### 1.1. lambda 表达式

lambda表达式，返回的是匿名函数指针


```python
def add(a, b):
    return a+b

add_func = lambda a,b : (a+b) 

add_func(1,2)
```

```python
<函数名> = lambda <参数> : <表达式>
```



### 1.2. 函数 list.sort 


```python
list.sort(cmp=None,key=None,reverse=False)
```

1. cmp
   1. cmp是个函数对象，它的函数声明 = `def func(iter1,iter2) --> bool`
   2. 形参: 任意两个迭代器
   3. 返回值: bool
2. key
   1. key是个函数对象，它的函数声明 = `def func(iter iterator) --> int`
   2. 形参: list的迭代器对象
   3. 返回值: int，作为比较的依据
3. reverse 排序规则 
   1. True 降序 
   2. False 升序 (默认)

例子1: cmp

```python
nums = [[1,3],[2,6],[8,10],[15,18]]
# key的形参 = 迭代器
nums.sort(key = lambda iterator : iterator[0])  
nums.sort(key = lambda (x,y) : x) 
# cmp的形参 = 迭代器
nums.sort(cmp = lambda iter1, iter2: iter1[0]-iter2[0])
```


例子2: key

传递给key的是函数名，它指定`可迭代对象中的每个元素`按照`该函数`进行排序


```python
nums = [3, 30, 34, 5, 9]
nums.sort(cmp = lambda x,y : cmp(str(x)+str(y), str(y)+str(x)))
```

---

# 2.题目

### 2.1. [56.合并区间](https://leetcode-cn.com/problems/merge-intervals/)

```python
输入: intervals = [[1,3],[2,6],[8,10],[15,18]]
输出: [[1,6],[8,10],[15,18]]
解释: 区间 [1,3] 和 [2,6] 重叠, 将它们合并为 [1,6].
```

```python
class Solution(object):
    def merge(self, intervals):
        
        # 按照第一个元素排序
        # intervals.sort(key = lambda iter : iter[0])  
        # intervals.sort(key = lambda (x,y) : x) 
        intervals.sort(cmp = lambda iter1, iter2: iter1[0]-iter2[0])

        ans = []
        for iter in intervals:
            # (结果集中没有区间) or (当前区间的左边界 > 结果集中的最后一个区间的右边界)
            if ans == [] or iter[0] > ans[-1][1]:
                ans.append(iter) # 将该区间加入结果集中
            else: # 结果集合中有区间 && 当前区间的左边界在结果集中最后一个区间内部
                # 当前区间的右边界 > 结果集的最后一个区间的右边界
                if iter[1] > ans[-1][1]:
                    ans[-1][1] = iter[1]  # 更新结果集的最后一个区间的右边界 = 当前区间的右边界
        return ans
```


### 2.2 [179.最大组合数](https://leetcode-cn.com/problems/largest-number/)

```shell
输入：nums = [3,30,34,5,9]
输出："9534330"
```

```python
class Solution(object):
    def largestNumber(self, nums):
        
        nums.sort(cmp = lambda x,y : cmp(str(y)+str(x), str(x)+str(y)))  # 按照 y+x, x+y 的字典序排序
        
        ans = ""
        for elem in nums:
            ans += str(elem)

        ans = ans.lstrip('0')
        return '0' if ans == "" else ans
```

