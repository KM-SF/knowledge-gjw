:kissing_smiling_eyes: 递归回溯 、 dfs深度优先遍历


1. 先画出二叉树的结果图
2. 采用dfs执行深度优先遍历 （dfs: 不撞南墙不回头）


```python
def xxx(nums):
    ans = []
    land = []

    def dfs():
        if 满足递归退出条件
            ans.append(land)
            return
        for 元素 in range(可选集合):
            land.append(元素)  # 选中一个元素
            dfs(下一个可选集合)
            land.pop()        # 弹出刚才选中的元素

    dfs()
    return ans
```

# 1. 子集\组合

### 1.1. [子集](https://leetcode-cn.com/problems/subsets/)

**题78：无重复元素**

```python
输入: nums = [1,2,3]
输出:
[
  [3],
  [1],
  [2],
  [1,2,3],
  [1,3],
  [2,3],
  [1,2],
  []
]
```

```python
class Solution(object):
    def subsets(self, nums):
        ans = []
        land = []

        size = len(nums)

        def dfs(start):
            ans.append(land[:])
            for i in range(start, size):
                land.append(nums[i])
                dfs(i+1)
                land.pop()
        
        dfs(0)
        return ans
```

### 1.2. [题90.子集Ⅱ](https://leetcode-cn.com/problems/subsets-ii/)

给定一个可能包含`重复元素`的整数数组 ***nums***，返回该数组所有可能的子集（幂集）。

**说明：**解集不能包含重复的子集。

```python
输入: [1,2,2]
输出:
[
  [2],
  [1],
  [1,2,2],
  [2,2],
  [1,2],
  []
]
```

```python
class Solution(object):
    def subsetsWithDup(self, nums):
        ans = []
        land = []

        size = len(nums)
        nums.sort() # 剪枝前，必须排序
        
        def dfs(start):
            ans.append(land[:])
            for i in range(start, size):
                if i>start and nums[i] == nums[i-1]: # 剪枝: 解集不能包含重复的子集
                    continue
                land.append(nums[i])
                dfs(i+1)
                land.pop()
        dfs(0)
        return ans
```




### 1.3. [77. 组合](https://leetcode-cn.com/problems/combinations/)

给定两个整数 n 和 k，返回范围 [1, n] 中所有可能的 k 个数的组合

```python
输入: n = 4, k = 2
输出:
[
  [2,4],
  [3,4],
  [2,3],
  [1,2],
  [1,3],
  [1,4],
]
```

```python
class Solution(object):
    def combine(self, n, k):
        nums = [i+1 for i in range(n)]

        ans = []
        land = []

        size = len(nums)

        def dfs(start):
            if len(land) == k:
                ans.append(land[:])
                return
            for i in range(start, size):
                land.append(nums[i])
                dfs(i+1)
                land.pop()
        
        dfs(0)
        return ans
```

### 1.4. [39.组合总和](https://leetcode-cn.com/problems/combination-sum/)

给你一个 无重复元素 的整数数组 candidates 和一个目标整数 target ，找出 candidates 中可以使数字和为目标数 target 的 所有不同组合 ，并以列表形式返回。你可以按 任意顺序 返回这些组合。

candidates 中的 同一个数字可以`无限制重复`被选取。如果至少一个数字的被选数量不同，则两种组合是不同的。 


```python
class Solution(object):
    def combinationSum(self, candidates, target):
        ans = []
        land = []

        size = len(candidates)

        def dfs(start, target):
            if target < 0:
                return
            if target == 0:
                ans.append(land[:])
                return
            for i in range(start, size):
                land.append(candidates[i])
                dfs(i, target-candidates[i]) # i: 相同位置的数可以被重复使用多次
                land.pop()
        
        dfs(0, target)
        return ans 
```

### 1.5. [40. 组合总和 II](https://leetcode-cn.com/problems/combination-sum-ii/)

给你一个由候选元素组成的集合 candidates 和一个目标数 target ，找出 candidates 中所有可以使数字和为 target 的组合。

candidates 中的每个元素在每个组合中只能使用一次 。

注意：解集不能包含重复的组合。

```python
class Solution(object):
    def combinationSum2(self, candidates, target):
        ans = []
        land = []
        size = len(candidates)
        candidates.sort()
        def backtrace(start, target):
            if target < 0:
                return
            if target == 0:
                ans.append(land[:])
                return
            for i in range(start, size):
                if i>start and candidates[i] == candidates[i-1]: # 剪枝: 解集不包含重复组合
                    continue
                land.append(candidates[i])
                backtrace(i+1, target-candidates[i]) # i+1: 不能重复使用
                land.pop()
        
        backtrace(0, target)
        return ans
```

### 1.6.[组合总和3](https://leetcode-cn.com/problems/combination-sum-iii/)

找出所有相加之和为 n 的 k 个数的组合。组合中只允许含有 1 - 9 的正整数，并且每种组合中不存在重复的数字。

说明：
1. 所有数字都是正整数。
2. 解集不能包含重复的组合。

```python
class Solution(object):
    def combinationSum3(self, k, n):
        ans = []
        land = []

        def dfs(start, n):
            if len(land) == k and n==0:
                ans.append(land[:])
                return 
            for i in range(start, 10):
                land.append(i)
                dfs(i+1, n-i)
                land.pop()
        
        dfs(1, n)
        return ans
```

### 1.6.[组合总和4](https://leetcode-cn.com/problems/combination-sum-iv/)

说明: dp，不是dfs

---


# 2. 全排列

### 2.1. [题46. 全排列](https://leetcode-cn.com/problems/permutations/submissions/)

给定一个 没有重复 数字的序列，返回其所有可能的全排列。

```python
class Solution(object):
    def permute(self, nums):
        ans = []
        land = []

        size = len(nums)
        used = [0]*size

        def dfs():
            if size == len(land):
                ans.append(land[:])
                return
            for i in range(size):
                if used[i]==1:
                    continue
                used[i]=1;   land.append(nums[i])
                dfs()
                used[i]=0;   land.pop()
        
        dfs()
        return ans
```

### 2.2.[题47. 全排列Ⅱ](https://leetcode-cn.com/problems/permutations-ii/submissions/)

给定一个可包含重复数字的序列 `nums` ，**按任意顺序** 返回所有不重复的全排列。

```python
class Solution(object):
    def permuteUnique(self, nums):
        ans = []
        land = []

        nums.sort()  # 剪枝: 有序

        size = len(nums)
        used = [0]*size

        def dfs():
            if size == len(land):
                ans.append(land[:])
                return
            for i in range(size):
                if used[i-1]==1 and i>0 and nums[i]==nums[i-1]: # 剪枝: 同层之前的元素已被使用
                    continue
                if used[i]==1:
                    continue
                used[i]=1;   land.append(nums[i])
                dfs()
                used[i]=0;   land.pop()
        
        dfs()
        return ans
```

---



## [31. 下一个排列](https://leetcode-cn.com/problems/next-permutation/)

- 此题的最优解法不是“递归回溯”

```python
class Solution(object):
    def nextPermutation(self, nums):
        # 一开始，做一个边缘条件的判断
        # 如果整个数组是降序排列的，直接返回值是原数组的逆序
        if sorted(nums, reverse=True) == nums:  # 整个数组是降序排列的
            nums[:] = nums[::-1]
            return 

        size  = len(nums)

        # 从右往左看，找寻第一个位置i, 满足(右边数 > 左边数)
        # i就是第一个破坏了降序的那个数字
        for i in range(size-1)[::-1]:
            if nums[i] < nums[i+1]:
                break

        # 再在[i+1:]内，找最小的比nums[i]大的数字
        # 通过j+1找j，如果nums[j+1] <= nums[i]，那么nums[j]就是最小的比nums[i]大的数字
        # 这是因为在nums[i]之后的数，都是降序排列的
        for j in range(i+1, size):
            if j+1 == size or nums[j+1] <= nums[i]:  # 找到之后，将nums[i], nums[j]做一次交换
                nums[i], nums[j] = nums[j], nums[i]  # 没有找到，也是将nums[i], nums[j]做一次交换
                break
        # 交换完位置之后，只需要将[i+1:]区间的字符串做一次逆序，就ok了
        nums[i+1:] = nums[i+1:][::-1]
        return
```



---

# 3.应用题

### 3.1. [22. 括号生成](https://leetcode-cn.com/problems/generate-parentheses/)

数字 *n* 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且 **有效的** 括号组合。

```shell
输入：n = 3
输出：[
       "((()))",
       "(()())",
       "(())()",
       "()(())",
       "()()()"
     ]
```

生成n对括弧，那么结果就是2n，设置left/right两个变量，分别记录左/右括弧出现的次数

> 每次，选中左/右括弧，相应的left/right都+1。当满足left+right=2n时，得到一个括弧

> 「注」在选择左/右括弧时，是有条件的！

> 「注」每次选中，都要对left/right +1，然后，将括弧加入land中，再次backtrace

> 1. 选中左括弧条件：left < n
> 2. 选中右括弧条件：left > right && right < n

```python
class Solution(object):
    # 1.左括弧+右括弧，共有2n个括弧时，一共能构成哪些
    # 2.不满足括弧匹配的进行剪枝
    def generateParenthesis(self, n):
        ans = []
        land = ""

        def backtrace(land, left, right): # left,right 左括弧数量, 右括弧的数量
            if left + right == 2*n:  # 递归终止条件: 左右括弧总数量 == 2n
                ans.append(land[:])
                return

            # 添加‘(’的条件: 左括弧 < n
            if left < n: 
                backtrace(land+'(', left+1, right)
                
            # 添加‘)’的条件: 左括弧>右括弧 && 右括弧 < n
            if left > right and right < n:  
                backtrace(land+")", left, right+1)

        backtrace(land, 0, 0)
        return ans
```

### 3.2. [17. 电话号码的字母组合](https://leetcode-cn.com/problems/letter-combinations-of-a-phone-number/)

给定一个仅包含数字 `2-9` 的字符串，返回所有它能表示的字母组合。

给出数字到字母的映射如下（与电话按键相同）。注意 1 不对应任何字母。

```c
输入："23"
输出：["ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"].
```

```python
class Solution(object):
    def letterCombinations(self, digits):
        hash = {
            '2': "abc",
            '3': "def",
            '4': "ghi",
            '5': "jkl",
            '6': "mno",
            '7': "pqrs",
            '8': "tuv",
            '9': "wxyz"
        }
        if digits == "":
            return []
        size = len(digits)
        ans = []
        land = ""

        def backtrace(land, start):
            if len(land) == size:
                ans.append(land)
                return
            for idx in range(start, size):
                for ch in hash[digits[idx]]:
                    backtrace(land + ch, idx + 1)

        backtrace(land, 0)
        return ans
```



### 3.3. [93. 复原IP地址](https://leetcode-cn.com/problems/restore-ip-addresses/)

```
输入：s = "25525511135"
输出：["255.255.11.135","255.255.111.35"]

输入：s = "0000"
输出：["0.0.0.0"]

输入：s = "1111"
输出：["1.1.1.1"]
```

---



[参考链接TODO](https://space.bilibili.com/15965981/channel/seriesdetail?sid=1215317)
