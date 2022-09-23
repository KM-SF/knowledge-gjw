

## [DFS](https://leetcode-cn.com/problems/number-of-islands/solution/dao-yu-lei-wen-ti-de-tong-yong-jie-fa-dfs-bian-li-/) 深度优先遍历

下面是最简单的二叉树的前序遍历递归写法：

```C++
void dfs(TreeNode root) {
    if (root == null) { // 判断 base case
        return;
    }
    
    cout << root.val << endl;  // 访问root
    
    // 访问两个相邻结点：左子结点、右子结点
    dfs(root.left);
    dfs(root.right);
}
```

----

### 1. [200. 岛屿数量](https://leetcode-cn.com/problems/number-of-islands/)

```c++
class Solution {
public:
    // 深度优先遍历: 一直将岛屿变成陆地
    void dfs(vector<vector<char>>& grid, int x, int y) {
        if (x < 0 || x >= m || y < 0 || y >= n)  // 不在岛屿内
            return;
        if (grid[x][y] != '1')  // 节点不是岛屿
            return;

        grid[x][y] = '0'; // 将岛屿变成陆地

        dfs(grid, x, y+1);
        dfs(grid, x +1, y);
        dfs(grid, x, y-1);
        dfs(grid, x-1, y);
    }
    int numIslands(vector<vector<char>>& grid) {
        m = grid.size();
        n = grid[0].size();
        int ans = 0;

        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (grid[i][j] == '1') {
                    dfs(grid, i, j);
                    ans += 1;
                }
            }
        }

        return ans;
    }

private:
    int m;
    int n;
};
```

### 2. [695. 岛屿的最大面积](https://leetcode-cn.com/problems/max-area-of-island/)

```c++
class Solution {
public:
    // 深度优先遍历: 一直将岛屿变成陆地
    // 返回值: 包含(x,y)的岛屿的面积
    int dfs(vector<vector<int>>& grid, int x, int y) {

        int ans = 0;

        if (x < 0 || x >= m || y < 0 || y >= n)  // 不在岛屿内
            return 0;
        if (grid[x][y] != 1)  // 节点不是岛屿
            return 0;

        grid[x][y] = 0; // 将岛屿变成陆地

        // ans = (x,y)的相邻的上下左右岛屿面积
        ans += dfs(grid, x, y+1);
        ans += dfs(grid, x +1, y);
        ans += dfs(grid, x, y-1);
        ans += dfs(grid, x-1, y);

        return 1 + ans; // (x,y)本身 + ans
    }
    int maxAreaOfIsland(vector<vector<int>>& grid) {
        m = grid.size();
        n = grid[0].size();
        int ans = 0;

        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (grid[i][j] == 1) {
                    ans = max( ans, dfs(grid, i, j) );
                }
            }
        }

        return ans;
    }
private:
    int m;
    int n;
};
```
