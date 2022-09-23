### BFS 广度优先遍历

- queue



---



#### [994. 腐烂的橘子](https://leetcode-cn.com/problems/rotting-oranges/)


```python
class Solution:
    def orangesRotting(self, grid):
        m = len(grid)
        n = len(grid[0])
        
        if grid == [[0,0,0,0]]:
            return 0
        if m==1 and n==1 and grid[0][0] != 1:
            return 0

        queue = []

        # 先把烂橘子放入队列
        for i in range(m):
            for j in range(n):
                if grid[i][j] == 2:
                    queue.append((i,j))
        
        time = -1
        # BFS
        while queue:
            size = len(queue)
            for _ in range(size):
                i, j = queue.pop(0)
                for x,y in [(1,0), (-1,0), (0,1), (0,-1)]:
                    i_next = i+x
                    j_next = j+y
                    if 0<=i_next<m and 0<=j_next<n and grid[i_next][j_next]==1:
                        grid[i_next][j_next] = 2
                        queue.append((i_next, j_next))
            time += 1
        
        for i in range(m):
            for j in range(n):
                if grid[i][j] == 1:
                    return -1

        return time
```

