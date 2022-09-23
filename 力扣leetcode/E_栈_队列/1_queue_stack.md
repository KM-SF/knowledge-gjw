

# 1. 栈 <==>队列

:artificial_satellite: 栈、队列（相互实现）

没有什么技巧可言

- 添加时，都是向data中push元素
- 写法注意：因为C++的pop函数返回值为void，因此，要先top/front，所以，在书写代码pop的时候，要转换为两行（为了防止漏写，将下面2行写在一行）
  - int top = help.top();     help.pop();
  - int ret = data.front();   data.pop();





### 1.1. 两个栈实现一个队列

```c++
class MyQueue {
public:
    /** Initialize your data structure here. */
    stack<int> data;
    stack<int> help;
    MyQueue() {

    }
    
    void push(int x) {
        data.push(x);  // 一直进入data
    }
    
    int pop() {
        if (empty())
            return -1;

        if(!help.empty()){  // 返回时，如果help有值，就返回help栈顶
            int top = help.top(); help.pop();
            return top;
        }
        // 如果help没有值，就先将data中的数据[“全部”“倒入”help]
        while(!data.empty()){
            int top = data.top(); data.pop();
            help.push(top);
        }
        // 最后，弹出&&返回help的栈顶
        int top = help.top(); help.pop();
        return top;  
    }
    
    int peek() {
        if (empty())
            return -1;
        if(!help.empty()){
            int top = help.top();
            return top;
        }
        while(!data.empty()){
            int top = data.top(); data.pop();
            help.push(top);
        }
        // 最后，只是返回help的栈顶
        return help.top();
    }
    
    bool empty() {
        if (help.empty() && data.empty())
            return true;
        return false;
    }
};
```

### 1.2. 两个队列实现一个栈

```c++
class MyStack {
public:
    queue<int> data;
    queue<int> help;

    MyStack() {
    }
    
    void push(int x) {
        data.push(x);
    }
    
    int pop() {
        if (empty())
            return -1;
        while(data.size() > 1){
            int f = data.front(); data.pop();
            help.push(f);
        }
        int ret = data.front(); data.pop();
        swap(data, help);
        return ret;
    }
    
    int top() {
        if (empty())
            return -1;
        while(data.size() > 1){
            int f = data.front(); data.pop();
            help.push(f);
        }
        int ret = data.front(); data.pop();
        help.push(ret);
        swap(data, help);
        return ret;
    }
    
    bool empty() {
        return data.empty();
    }
};
```

### 1.3. 最小栈  TODO

---

# 2. 栈: 基本题型


**python的list模拟栈**

- 获取栈顶：stack[-1]
- 弹出栈顶：stack.pop(-1)  或  stack.pop()

### 2.1. [71. 简化路径](https://leetcode-cn.com/problems/simplify-path/)

```python
class Solution(object):
    def simplifyPath(self, path):
        slist = path.split('/')
        stack = []
        for ch in slist:
            if ch == '.' or ch == '':
                continue
            if ch == '..':
                if stack != []:
                    stack.pop(-1)
            else:
                stack.append(ch)
        return '/' + '/'.join(stack)
```

### 2.2. [20. 有效的括号](https://leetcode-cn.com/problems/valid-parentheses/)

```python
class Solution(object):
    def isValid(self, s):
        hash = {
            ')' : '(', 
            '}' : '{', 
            ']' : '['
        }

        stack = []

        for ch in s:
            if ch in hash.values():  # ch是左括弧
                stack.append(ch)
            else:  # ch是右括弧
                if stack == [] or hash[ch] != stack[-1]: # 当前右括弧对应的左括弧hash[ch] != 栈顶左括弧
                    return False
                else:
                    stack.pop(-1)
        return len(stack) == 0
```

### 2.3. [844. 比较含退格的字符串](https://leetcode-cn.com/problems/backspace-string-compare/)

```
输入：S = "ab#c", T = "ad#c"
输出：true
解释：S 和 T 都会变成 “ac”。
```

```python
class Solution(object):
    def backspaceCompare(self, S, T):
        def getStr(s):
            stack = []
            for ch in s:
                if ch != '#':
                    stack.append(ch)
                else:
                    if stack != []:
                        stack.pop(-1)
            return "".join(stack)
        return getStr(S) == getStr(T)
```



### 2.4. [1544. 整理字符串](https://leetcode-cn.com/problems/make-the-string-great/)

```python
class Solution(object):
    def makeGood(self, s):
        stack = []
        for ch in s:
            if stack != [] and abs(ord(stack[-1]) - ord(ch)) == ord('a') - ord('A'):
                stack.pop(-1)
            else:
                stack.append(ch)
        return "".join(stack)
```

### 2.5. [1047. 删除字符串中的所有相邻重复项](https://leetcode-cn.com/problems/remove-all-adjacent-duplicates-in-string/)


输入："abbaca"

输出："ca"

解释：在 "abbaca" 中，我们可以删除 "bb" 由于两字母相邻且相同，这是此时唯一可以执行删除操作的重复项。之后我们得到字符串 "aaca"，其中又只有 "aa" 可以执行重复项删除操作，所以最后的字符串为 "ca"。


```python
class Solution(object):
    def removeDuplicates(self, s):
        stack = []
        for ch in s:
            if stack != [] and stack[-1] == ch:
                stack.pop(-1)
            else:
                stack.append(ch)
        return "".join(stack)
```

### 2.6. [1209. 删除字符串中的所有相邻重复项 II](https://leetcode-cn.com/problems/remove-all-adjacent-duplicates-in-string-ii/)


输入：s = "deeedbbcccbdaa", k = 3

输出："aa"
解释： 
1. 先删除 "eee" 和 "ccc"，得到 "ddbbbdaa"
2. 再删除 "bbb"，得到 "dddaa"
3. 最后删除 "ddd"，得到 "aa"


- 解题关键: **栈中存放的是 [ch, ch出现次数]**

```python
class Solution(object):
    def removeDuplicates(self, s, k):
        stack = [] # 栈中元素是[], 即 [ch, ch出现次数]
        for ch in s:
            if stack != [] and stack[-1][0] == ch:
                if stack[-1][1] == k-1:
                    stack.pop(-1)
                else:
                    stack[-1][1] = stack[-1][1]+1
            else:
                stack.append([ch, 1])
        
        ans = ""
        for [ch,cnt] in stack:
            ans += ch * cnt
        return ans
```

### 2.7. [394. 字符串解码](https://leetcode-cn.com/problems/decode-string/) :kissing_smiling_eyes:

给定一个经过编码的字符串，返回它解码后的字符串。

编码规则为: k[encoded_string]，表示其中方括号内部的 encoded_string 正好重复 k 次。注意 k 保证为正整数。

```c
示例 1：
输入：s = "3[a]2[bc]"
输出："aaabcbc"
示例 2：
输入：s = "3[a2[c]]"
输出："accaccacc"
```

[解答](https://leetcode-cn.com/problems/decode-string/solution/decode-string-fu-zhu-zhan-fa-di-gui-fa-by-jyd/)

- 阶梯关键: 栈存放的元素 (当前元素出现的次数, 上一个字符串)

```python
class Solution(object):
    def decodeString(self, s):
        stack = []
        res, multi = "", 0
        for c in s:
            if '0' <= c <= '9': # [0-9]
                multi = multi * 10 + int(c)            
            elif c == '[':  # 入栈
                stack.append([multi, res])
                res, multi = "", 0
            elif c == ']':  # 出栈
                top = stack.pop(-1)
                res = top[1] + top[0] * res
            else: # [A-Z,a-z]
                res += c
        return res
```

--- 


# 3. 单调队列、单调栈

## 单调队列、单调栈

解题套路：先分析题意，是否要使用单调栈、单调队列，并知道数据的变化过程（代码的书写是有规律的，按照下面的写法，代码十分清晰）

```python
1. 从前向后，遍历每一个元素 (该元素一定会进入栈/队列)
	for r in range(size):
2. while循环, 保持单调性
	2.1. while 循环条件 and 题目条件
    		从[]中移除不满足单调性的栈顶/队尾，即:pop(-1)
	2.2. 当出 单调[] 后，可能需要(保存结果 + 更新条件)
3. ...
4. 元素进入 单调[]
	4.1. 当进入 单调[] 后，可能需要(更新结果 + 更新条件)
```

## 3.1. 单调队列

### 3.1.1. [剑指 Offer 59 - I. 滑动窗口的最大值](https://leetcode-cn.com/problems/hua-dong-chuang-kou-de-zui-da-zhi-lcof/)

:small_airplane: 单调递减队列

```python
class Solution(object):
    def maxSlidingWindow(self, nums, k):
        size = len(nums)
        ans = []
        queue = []  # 单调队列 (由大到小), 存放了"下标"
        for r in range(size):
            while queue and nums[r] > nums[queue[-1]]:
                queue.pop(-1)
            queue.append(r)
            if r-queue[0]+1 > k:
                queue.pop(0)
            if r >= k-1:
                ans.append(nums[queue[0]])
        return ans
```



### 3.1.2. [LeetCode 面试题59 - II. 队列的最大值](https://segmentfault.com/a/1190000021962984):slightly_smiling_face:

:small_airplane: 单调递减队列

请定义一个队列并实现函数 max_value 得到队列里的最大值，要求函数max_value、push_back 和 pop_front 的均摊时间复杂度都是O(1)。

若队列为空，pop_front 和 max_value 需要返回 -1


```python

class MaxQueue:

    def __init__(self):
        queue = []     # 原队列
        max_queue = [] # 单调递减队列

    def max_value(self) -> int:
        return self.max_queue[0] if len(max_queue) > 0 else -1

    def push_back(self, value: int) -> None:
        # 保证单调递减队列的性质
		while len(max_queue) > 0 and value > max_queue[-1]:
            max_queue.pop(-1)
        max_queue.append(value)
        # 原始队列
		queue.append(value)

    def pop_front(self) -> int:
        if not self.queue:
            return -1
        popVal = queue.pop(0) # 待弹出元素
        if popVal == max_queue[0]: # max_queue的首元素 == 待弹出元素，将其弹出
            max_queue.pop(0)
        return res
```


## 3.2. 单调栈

### 3.2.1. [739. 每日温度](https://leetcode-cn.com/problems/daily-temperatures/):slightly_smiling_face:

```python
请根据每日气温列表，重新生成一个列表。对应位置的输出为：要想观测到更高的气温，至少需要等待的天数。如果气温在这之后都不会升高，请在该位置用 0 来代替。
例如，
	给定一个列表 temperatures = [73, 74, 75, 71, 69, 72, 76, 73]
	你的输出应该是 [1, 1, 4, 2, 1, 1, 0, 0]。
```

- 栈中存放的元素，是下标，不是nums[i]

```python
class Solution(object):
    def dailyTemperatures(self, T):
        ans = [0] * len(T)  # 保存结果[0, 0, 0, ... ...]
        stack = []  # 单调栈, 存放“下标”
        for i in range(len(T)):
            while stack and T[i] > T[stack[-1]]:  # 当前温度 > 栈顶温度
                ans[stack[-1]] = i - stack[-1]    # 出栈前, 得到栈顶(下标位置)的结果值
                stack.pop(-1)
            stack.append(i)
        return ans
```

### 3.2.2. [402. 移掉K位数字](https://leetcode-cn.com/problems/remove-k-digits/)

给定一个以字符串表示的非负整数 *num*，移除这个数中的 *k* 位数字，使得剩下的数字最小。

```
输入: num = "1432219", k = 3
输出: "1219"
解释: 移除掉三个数字 4, 3, 和 2 形成一个新的最小的数字 1219	。
```

栈中存放的元素是 num[i]，不是下标

```python
class Solution(object):
    def removeKdigits(self, num, k):
        stack = []
        for e in num:
            while stack and stack[-1] > e and k: # 保持栈的单调性 && k>0
                stack.pop()
                k -= 1
            stack.append(e)
        # 别忘记(k>0)时，要处理下喔~
        if k > 0:
            stack = stack[:-k]
        # 去掉左边的0
        ret = "".join(stack).lstrip('0')
        return '0' if len(ret)==0 else ret
```



### 3.2.3. [496. 下一个更大元素 I](https://leetcode-cn.com/problems/next-greater-element-i/) TODO

给定两个 `没有重复元素` 的数组 nums1 和 nums2 ，其中nums1 是 nums2 的子集。找到 nums1 中每个元素在 nums2 中的下一个比其大的值。

nums1 中数字 x 的下一个更大元素是指 x 在 nums2 中对应位置的右边的第一个比 x 大的元素。如果不存在，对应位置输出 -1 。

```
输入: nums1 = [4,1,2], nums2 = [1,3,4,2].
输出: [-1,3,-1]
解释:
    对于num1中的数字4，你无法在第二个数组中找到下一个更大的数字，因此输出 -1。
    对于num1中的数字1，第二个数组中数字1右边的下一个较大数字是 3。
    对于num1中的数字2，第二个数组中没有下一个更大的数字，因此输出 -1。
```



```python
class Solution(object):
    def nextGreaterElement(self, nums1, nums2):
        stack = []
        hash = {}
        for e in nums2:
            while stack != [] and stack[-1] < e:  # 单调栈
                hash[stack[-1]] = e  # 将<元素,第一个比元素大的值>保存到哈希表
                stack.pop(-1)
            stack.append(e)
        ret = []
        for e in nums1: # 遍历每个元素，从哈希表中找其对应的第一个较大值
            if e in hash.keys():
                ret.append(hash[e])
            else:
                ret.append(-1)
        return ret
```



### 3.2.4. [503. 下一个更大元素 II](https://leetcode-cn.com/problems/next-greater-element-ii/) TODO

给定一个循环数组（最后一个元素的下一个元素是数组的第一个元素），输出每个元素的下一个更大元素。数字 x 的下一个更大的元素是按数组遍历顺序，这个数字之后的第一个比它更大的数，这意味着你应该循环地搜索它的下一个更大的数。如果不存在，则输出 -1。

```shell
输入: [1,2,1]
输出: [2,-1,2]
解释: 第一个 1 的下一个更大的数是 2；
数字 2 找不到下一个更大的数； 
第二个 1 的下一个最大的数需要循环搜索，结果也是 2。
```


