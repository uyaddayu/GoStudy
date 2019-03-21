package main

import (
	"fmt"
	"os"
)

func main() {
	maze := readMaze("./maze/maze.in") // 假设迷宫是长方形
	for _, row := range maze {
		for _, val := range row {
			fmt.Printf("%3d", val)
		}
		fmt.Println()
	}
	steps := wark(maze, point{0, 0}, point{len(maze) - 1, len(maze[0]) - 1})
	fmt.Println("走法：")
	for _, row := range steps {
		for _, val := range row {
			fmt.Printf("%3d", val)
		}
		fmt.Println()
	}
}

type point struct {
	i, j int
}

// 上左下右
var dirs = [4]point{
	{0, -1}, {-1, 0}, {0, 1}, {1, 0},
}

// 点和点的加
func (p point) add(a point) point {
	return point{p.i + a.i, p.j + a.j}
}

// 得到二维坐标系对应点的值
func (p point) at(grid [][]int) (int, bool) {
	// 判断是否越界
	if (p.i < 0 || p.i >= len(grid)) || (p.j < 0 || p.j >= len(grid[p.i])) {
		return 0, false
	}
	return grid[p.i][p.j], true
}

// 走迷宫
func wark(maze [][]int, start, end point) [][]int {
	steps := make([][]int, len(maze))
	for i := range steps {
		steps[i] = make([]int, len(maze[i]))
	}

	Q := []point{start} // 起点,队列
	for len(Q) > 0 {    // 队列为空则退出
		cur := Q[0] // 当前点
		Q = Q[1:]
		if cur == end { // 发现终点
			break
		}
		curSteps, _ := cur.at(steps) // 当前点的值
		for _, dir := range dirs {
			next := cur.add(dir)

			// 判断迷宫里对应位置必须是0，还有走过的路径里对应位置也必须是0，也不能等于起点
			if val, ok := next.at(maze); !ok || val != 0 {
				continue
			}
			if val, ok := next.at(steps); !ok || val != 0 {
				continue
			}
			if next == start {
				continue
			}

			steps[next.i][next.j] = curSteps + 1

			Q = append(Q, next) // 增加到队列去
		}
	}
	return steps
}

// 从文件中读取迷宫
func readMaze(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var row, col int
	fmt.Fscanf(file, "%d %d\n", &row, &col)
	maze := make([][]int, row)
	for i := range maze {
		maze[i] = make([]int, col)
		for j := range maze[i] {
			// fmt.Fscanf(file, "%d", &maze[i][j])
			fmt.Fscan(file, &maze[i][j])
		}
	}
	return maze
}
