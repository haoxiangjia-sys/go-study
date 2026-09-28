package main

func calculateNextState(p golParams, world [][]byte) [][]byte {

	//1. 确定网格尺寸
	rows := len(world)
	cols := len(world[0])
	//2. 创建下一代网格
	next := make([][]byte, rows) //分配初始化 make is used to create and initialize slices, maps, and channels
	for i, _ := range next {     //range is used to iterate over arrays
		next[i] = make([]byte, cols)
	}
	//3. 定义 8 个邻居方向
	dirs := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	//4. 遍历每个细胞
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			//5. 统计活邻居
			liveNeighbours := 0
			for _, d := range dirs {
				ny := (y + d[0] + rows) % rows
				nx := (x + d[1] + cols) % cols
				if world[ny][nx] == 255 {
					liveNeighbours++
				}
			}


	//6. 应用生命游戏规则
	//live ceil if neighbour 2,3 -> 255   else 0
	//dead ceil if neighbour 3 -> 255  else 0
	if world[y][x] == 255 {
		if liveNeighbours == 2 || liveNeighbours == 3 {
			next[y][x] = 255
		} else {
			next[y][x] = 0
		}
	} else {
		if liveNeighbours == 3 {
			next[y][x] = 255
		} else {
			next[y][x] = 0
		}
	}


		}
	}

	//7. 返回新世界
	return next
}

func calculateAliveCells(p golParams, world [][]byte) []cell {
	//1. 初始化结果切片
	var alive []cell
	rows := len(world)
	cols := len(world[0])

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if world[y][x] == 255 {
				alive = append(alive, cell{x, y}) //把一个新的活细胞坐标追加到 alive 切片末尾。
			}
		}
	}
	return alive
}
